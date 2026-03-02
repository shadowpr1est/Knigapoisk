FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/services/api-gateway
RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
ENV PATH="/go/bin:${PATH}"
WORKDIR /app
RUN protoc -I /app --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative services/auth-service/api/proto/*.proto
RUN protoc -I /app --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative services/book-service/api/proto/*.proto
RUN protoc -I /app --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative services/file-service/api/proto/*.proto
RUN protoc -I /app --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative services/reading-service/api/proto/*.proto
RUN protoc -I /app --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative services/review-service/api/proto/*.proto
WORKDIR /app/services/api-gateway
RUN go mod tidy
RUN go build -o /app/app ./cmd/app

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/services/api-gateway/openapi.yaml ./openapi.yaml
ENV GIN_MODE=release
EXPOSE 8080
CMD ["./app"]

