BINARY_NAME = bookshop
API_FILE = books.yml


.PHONY: test_client
test_client:
	oapi-codegen --config=./api/books/v1/test/test_client.cfg.yml ./api/books/v1/$(API_FILE) 

.PHONY: books_api
books_api:
	oapi-codegen --config=./api/books/v1/server.cfg.yml ./api/books/v1/$(API_FILE) 

.PHONY: books_models
books_models:
	oapi-codegen --config=./api/books/v1/models.cfg.yml ./api/books/v1/$(API_FILE) 

.PHONY: build
build: test_client books_api books_models

.DEFAULT_GOAL := build


