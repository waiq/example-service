package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/waiq/example-service/pkg/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(ctx context.Context, db *gorm.DB) (*Repository, error) {
	return &Repository{
		db: db.WithContext(ctx),
	}, nil
}

func (r *Repository) Migration(ctx context.Context) error {
	return r.db.WithContext(ctx).AutoMigrate(
		models.Book{},
	)
}

func (r *Repository) StoreBook(ctx context.Context, book *models.Book) error {
	if err := r.db.WithContext(ctx).Create(&book).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindBookByUUID(ctx context.Context, id uuid.UUID) (*models.Book, error) {
	var results models.Book
	tx := r.db.WithContext(ctx).Where("UUID = ?", id.String()).First(&results)

	if tx.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &results, tx.Error
}

func (r *Repository) GetBooks(ctx context.Context) ([]models.Book, error) {
	var results []models.Book
	tx := r.db.WithContext(ctx).Find(&results)
	return results, tx.Error
}
