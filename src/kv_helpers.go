package main

// ripoff: https://raw.githubusercontent.com/hashicorp/vault/main/command/kv_helpers.go

import (
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/api"
)

func kvReadRequest(client *api.Client, path string, params map[string]string) (*api.Secret, error) {
	r := client.NewRequest("GET", "/v1/"+path)
	for k, v := range params {
		r.Params.Set(k, v)
	}
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := api.ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return api.ParseSecret(resp.Body)
}

func kvPreflightVersionRequest(client *api.Client, path string) (string, int, error) {
	// We don't want to use a wrapping call here so save any custom value and
	// restore after
	currentWrappingLookupFunc := client.CurrentWrappingLookupFunc()
	client.SetWrappingLookupFunc(nil)
	defer client.SetWrappingLookupFunc(currentWrappingLookupFunc)
	currentOutputCurlString := client.OutputCurlString()
	client.SetOutputCurlString(false)
	defer client.SetOutputCurlString(currentOutputCurlString)

	r := client.NewRequest("GET", "/v1/sys/internal/ui/mounts/"+path)
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		// If we get a 404 we are using an older version of vault, default to
		// version 1
		if resp != nil && resp.StatusCode == 404 {
			return "", 1, nil
		}

		return "", 0, err
	}

	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return "", 0, err
	}
	if secret == nil {
		return "", 0, errors.New("nil response from pre-flight request")
	}
	var mountPath string
	if mountPathRaw, ok := secret.Data["path"]; ok {
		mountPath = mountPathRaw.(string)
	}
	options := secret.Data["options"]
	if options == nil {
		return mountPath, 1, nil
	}
	versionRaw := options.(map[string]interface{})["version"]
	if versionRaw == nil {
		return mountPath, 1, nil
	}
	version := versionRaw.(string)
	switch version {
	case "", "1":
		return mountPath, 1, nil
	case "2":
		return mountPath, 2, nil
	}

	return mountPath, 1, nil
}

func isKVv2(path string, client *api.Client) (string, bool, error) {
	mountPath, version, err := kvPreflightVersionRequest(client, path)
	if err != nil {
		return "", false, err
	}

	return mountPath, version == 2, nil
}

func addPrefixToKVPath(p, mountPath, apiPrefix string) string {
	if p == mountPath || p == strings.TrimSuffix(mountPath, "/") {
		return path.Join(mountPath, apiPrefix)
	}

	tp := strings.TrimPrefix(p, mountPath)
	for {
		// If the entire mountPath is included in the path, we are done
		if tp != p {
			break
		}
		// Trim the parts of the mountPath that are not included in the
		// path, for example, in cases where the mountPath contains
		// namespaces which are not included in the path.
		partialMountPath := strings.SplitN(mountPath, "/", 2)
		if len(partialMountPath) <= 1 || partialMountPath[1] == "" {
			break
		}
		mountPath = strings.TrimSuffix(partialMountPath[1], "/")
		tp = strings.TrimPrefix(tp, mountPath)
	}

	return path.Join(mountPath, apiPrefix, tp)
}

func getHeaderForMap(header string, data map[string]interface{}) string {
	maxKey := 0
	for k := range data {
		if len(k) > maxKey {
			maxKey = len(k)
		}
	}

	// 4 for the column spaces and 5 for the len("value")
	totalLen := maxKey + 4 + 5

	equalSigns := totalLen - (len(header) + 2)

	// If we have zero or fewer equal signs bump it back up to two on either
	// side of the header.
	if equalSigns <= 0 {
		equalSigns = 4
	}

	// If the number of equal signs is not divisible by two add a sign.
	if equalSigns%2 != 0 {
		equalSigns = equalSigns + 1
	}

	return fmt.Sprintf("%s %s %s", strings.Repeat("=", equalSigns/2), header, strings.Repeat("=", equalSigns/2))
}

func kvParseVersionsFlags(versions []string) []string {
	versionsOut := make([]string, 0, len(versions))
	for _, v := range versions {
		versionsOut = append(versionsOut, strutil.ParseStringSlice(v, ",")...)
	}

	return versionsOut
}
