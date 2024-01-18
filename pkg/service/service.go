package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/waiq/example-service/pkg/repository/models"
)

type Repository interface {
	StoreBook(context.Context, *models.Book) error
	FindBookByUUID(context.Context, uuid.UUID) (*models.Book, error)
	GetBooks(context.Context) ([]models.Book, error)
}

type BooksService struct {
	repo Repository
}

func NewBookService(ctx context.Context, repo Repository) *BooksService {
	return &BooksService{
		repo: repo,
	}
}

func (s *BooksService) AddBook(ctx context.Context, book *models.Book) error {
	return s.repo.StoreBook(ctx, book)
}

func (s *BooksService) GetBooks(ctx context.Context) ([]models.Book, error) {
	return s.repo.GetBooks(ctx)
}

func (s *BooksService) FindBookByUUID(ctx context.Context, id uuid.UUID) (*models.Book, error) {
	return s.repo.FindBookByUUID(ctx, id)
}
