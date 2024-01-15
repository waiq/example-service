package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/waiq/example-service/pkg/models"
	"github.com/waiq/example-service/pkg/repository"
	"github.com/waiq/example-service/pkg/service"
	testcontainer "github.com/waiq/example-service/pkg/test/testcontainers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	container testcontainer.PostgreSQLTestContainer
	terminate func(context.Context)

	repo *repository.Repository

	database *gorm.DB
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func setup() {
	ctx := context.Background()

	container, terminate = testcontainer.CreatePostgreSQLContainer(testcontainer.PostgreSQL_16)

	var err error
	database, err = gorm.Open(postgres.Open(container.ConnectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	repo, err = repository.New(context.Background(), database)
	if err != nil {
		panic(err)
	}

	if err := repo.Migration(ctx); err != nil {
		panic(err)
	}
}

func teardown() {
	terminate(context.Background())
}

func TestFindBook(t *testing.T) {
	//Test setup
	require := require.New(t)
	tx := database.Begin()
	defer tx.Rollback()

	repo, err := repository.New(context.Background(), tx)
	require.NoError(err)

	bookService := service.NewBookService(context.Background(), repo)

	id := uuid.New()
	err = bookService.AddBook(context.Background(), &models.Book{
		UUID:   id,
		Title:  "Mega Man Man",
		Author: "Bill the Bill",
	})
	require.NoError(err)

	// Test begin
	book, err := bookService.FindBookByUUID(context.Background(), id)
	require.NoError(err)
	require.Equal(id, book.UUID)

	// return nil if not found
	book, err = bookService.FindBookByUUID(context.Background(), uuid.New())
	require.NoError(err)
	require.Nil(book)
}

func TestAddBook(t *testing.T) {
	require := require.New(t)
	tx := database.Begin()
	defer tx.Rollback()

	repo, err := repository.New(context.Background(), tx)
	require.NoError(err)

	bookService := service.NewBookService(context.Background(), repo)

	book := &models.Book{
		UUID:   uuid.New(),
		Title:  "Mega Man Man",
		Author: "Bill the Bill",
	}
	err = bookService.AddBook(context.Background(), book)
	require.NoError(err)
	require.Greater(int(book.ID), 0)
}

func TestGetBooks(t *testing.T) {
	require := require.New(t)
	tx := database.Begin()
	defer tx.Rollback()

	repo, err := repository.New(context.Background(), tx)
	require.NoError(err)

	bookService := service.NewBookService(context.Background(), repo)

	data := []struct {
		UUID   uuid.UUID
		Title  string
		Author string
	}{
		{uuid.New(), "Mega She She", "Bull the Bull"},
		{uuid.New(), "Bending noses", "Marty mac Smart"},
	}

	for _, v := range data {
		err = bookService.AddBook(context.Background(),
			&models.Book{
				UUID:   v.UUID,
				Title:  v.Title,
				Author: v.Author,
			},
		)
		require.NoError(err)
	}

	books, err := bookService.GetBooks(context.Background())
	require.NoError(err)
	require.Len(books, 2)
}
