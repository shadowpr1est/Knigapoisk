PROTO_DIRS := services/auth-service/api/proto \
              services/book-service/api/proto \
              services/file-service/api/proto \
              services/reading-service/api/proto \
              services/review-service/api/proto

proto-gen:
	@for dir in $(PROTO_DIRS); do \
		for file in $$dir/*.proto; do \
			if [ -f $$file ]; then \
				echo "Generating $$file"; \
				protoc --go_out=$$(dirname $$file) --go-grpc_out=$$(dirname $$file) $$file; \
			fi; \
		done; \
	done

build:
	docker compose build

run:
	docker compose up -d

down:
	docker compose down

