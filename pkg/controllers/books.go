package controllers

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/waiq/example-service/api/books/v1"
	"github.com/waiq/example-service/pkg/models"
	"github.com/waiq/example-service/pkg/service"
	"github.com/waiq/example-service/pkg/util"
)

type BooksController struct {
	Lock sync.Mutex

	service *service.BooksService
}

func NewBooksController(service *service.BooksService) *BooksController {
	return &BooksController{
		service: service,
	}
}

// List all books
// (GET /books)
func (b *BooksController) GetBooks(
	ctx context.Context,
	request books.GetBooksRequestObject,
) (books.GetBooksResponseObject, error) {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	result, err := b.service.GetBooks(ctx)
	if err != nil {
		return nil, err
	}

	var response books.GetBooks200JSONResponse
	for _, r := range result {
		response = append(response, books.Books{
			Id:     util.Ptr(r.UUID.String()),
			Title:  &r.Title,
			Author: &r.Author,
		})
	}

	return response, nil
}

// Add a new book
// (POST /books)
func (b *BooksController) PostBooks(
	ctx context.Context,
	request books.PostBooksRequestObject,
) (books.PostBooksResponseObject, error) {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	book := models.Book{
		UUID:   uuid.New(),
		Title:  *request.Body.Title,
		Author: *request.Body.Author,
	}

	err := b.service.AddBook(ctx, &book)
	if err != nil {
		return nil, err
	}

	return books.PostBooks201JSONResponse{
		Id:     util.Ptr(book.UUID.String()),
		Title:  &book.Title,
		Author: &book.Author,
	}, nil
}

// Get details of a specific book
// (GET /books/{bookId})
func (b *BooksController) GetBooksBookId(
	ctx context.Context,
	request books.GetBooksBookIdRequestObject,
) (books.GetBooksBookIdResponseObject, error) {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	// id, err := uuid.
	id, err := uuid.Parse(request.BookId)
	if err != nil {
		return books.GetBooksBookId400Response{}, err
	}

	r, err := b.service.FindBookByUUID(ctx, id)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return books.GetBooksBookId404Response{}, nil
	}

	return books.GetBooksBookId200JSONResponse{
		Id:     util.Ptr(r.UUID.String()),
		Title:  &r.Title,
		Author: &r.Author,
	}, nil
}

func (b *BooksController) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
	})

	handler := books.NewStrictHandlerWithOptions(
		&BooksController{
			service: b.service,
		},
		[]books.StrictMiddlewareFunc{},
		books.StrictHTTPServerOptions{})

	return books.HandlerFromMux(handler, r)
}
