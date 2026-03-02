FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
ENV PATH="/go/bin:${PATH}"
RUN protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative api/proto/*.proto
RUN go mod tidy
RUN go build -o app ./cmd/app

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations
EXPOSE 9093
CMD ["./app"]

