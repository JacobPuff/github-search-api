FROM golang:1.17 as builder
WORKDIR /workspace
# install go plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# install buf
ENV BIN="/usr/local/bin"
ENV VERSION="1.0.0-rc6"
ENV BINARY_NAME="buf"
RUN  curl -sSL "https://github.com/bufbuild/buf/releases/download/v${VERSION}/${BINARY_NAME}-$(uname -s)-$(uname -m)"  -o "${BIN}/${BINARY_NAME}" && chmod +x "${BIN}/${BINARY_NAME}"
# copy relevant things
COPY buf.gen.yaml .
COPY buf.work.yaml .
COPY go.mod .
COPY go.sum .
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
COPY appconfig/ ./appconfig
COPY githubsearchapis/ ./githubsearchapis
COPY server/ ./server
# generate compiled protos
RUN buf generate
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o main server/main.go

# Use distroless as minimal base image
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/main .
USER 65532:65532
EXPOSE 9090

ENTRYPOINT ["/main"]
