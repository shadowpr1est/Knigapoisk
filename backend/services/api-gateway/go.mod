module github.com/shadowpr1est/knigapoisk-api-gateway

go 1.24.0

toolchain go1.24.5

require (
	github.com/gin-gonic/gin v1.11.0
	github.com/go-playground/validator/v10 v10.27.0
	github.com/joho/godotenv v1.5.1
	github.com/shadowpr1est/knigapoisk-auth-service v0.0.0
	github.com/shadowpr1est/knigapoisk-book-service v0.0.0
	github.com/shadowpr1est/knigapoisk-file-service v0.0.0
	github.com/shadowpr1est/knigapoisk-reading-service v0.0.0
	github.com/shadowpr1est/knigapoisk-review-service v0.0.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.67.0
	google.golang.org/protobuf v1.36.9
)

replace github.com/shadowpr1est/knigapoisk-auth-service => ../auth-service

replace github.com/shadowpr1est/knigapoisk-book-service => ../book-service

replace github.com/shadowpr1est/knigapoisk-file-service => ../file-service

replace github.com/shadowpr1est/knigapoisk-reading-service => ../reading-service

replace github.com/shadowpr1est/knigapoisk-review-service => ../review-service
