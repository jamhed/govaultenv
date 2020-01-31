FROM golang:1.13-alpine AS build
WORKDIR /home
COPY src src
RUN cd src && CGO_ENABLED=0 go build -o ../govaultenv

FROM alpine:3.10
COPY --from=build /home/govaultenv govaultenv
CMD ["./govaultenv"]
