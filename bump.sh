#!/bin/bash -ae
TAG=$(git tag | sort | tail -1)
MINOR=$(echo $TAG | sed -E 's|^.*\.([0-9]+)$|\1|')
MAJOR=$(echo $TAG | sed -E 's|^v(.*)\.[0-9]+$|\1|')
NEWVER="$MAJOR.$((MINOR+1))"
git tag v$NEWVER
git push origin v$NEWVER
