# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /hook
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY pod_displayer.gotmpl pod_displayer.gotmpl
COPY sample.gotmpl sample.gotmpl
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
ENV KUBE_NAMESPACE=default
ENV IS_IN_CLUSTER_DEP=true

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o hook main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
#FROM bitnami/kubectl:1.20
#WORKDIR /hook
#COPY --from=builder /workspace/ .
USER 65532:65532

CMD ["./hook"]
