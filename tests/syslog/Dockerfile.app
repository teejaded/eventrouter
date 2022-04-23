FROM registry.redhat.io/ubi8/go-toolset:1.16.12 AS  builder
WORKDIR  /go/src/github.com/openshift/eventrouter
USER 0
COPY . .
