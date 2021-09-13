FROM golang:alpine as builder
RUN apk add protoc git build-base
COPY src /go/src
WORKDIR /go/src
RUN go env -w GO111MODULE=off && \
    go env -w GOBIN=/go/bin && \
    go env -w CGO_ENABLED=0 && \
    go get google.golang.org/grpc \
    google.golang.org/protobuf/proto \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
    # github.com/go-redis/redis && \
    export PATH="$PATH:$(go env GOPATH)/bin"
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/labstatmsg.proto
RUN go get . && \
    go test -v && \
    go build -o /usr/local/bin/labstats

FROM scratch
COPY --from=builder /usr/local/bin/labstats /labstats
ENTRYPOINT ["/labstats"]
CMD ["s"]
