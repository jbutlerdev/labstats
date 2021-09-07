FROM golang:alpine as builder
RUN apk add protoc git build-base
RUN go env -w GO111MODULE=off && \
    go env -w CGO_ENABLED=0 && \
    go get google.golang.org/grpc \
    google.golang.org/protobuf/proto \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/go-redis/redis && \
    export PATH="$PATH:$(go env GOPATH)/bin"
COPY . .
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative src/protos/labstatmsg.proto && \
    cd src && \
    go test -v && \
    go build -o /usr/local/bin/labstats

FROM scratch
COPY --from=builder /usr/local/bin/labstats /labstats
ENTRYPOINT ["/labstats"]
CMD ["s"]
