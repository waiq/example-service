BINARY_NAME = bookshop
API_FILE = books.yml

books_api:
	oapi-codegen --config=./api/books/v1/server.cfg.yml ./api/books/v1/$(API_FILE) 

api: books_api

.PHONY: books_api
