package handlers_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	testclient "github.com/waiq/example-service/api/books/v1/test"
	"github.com/waiq/example-service/pkg/handlers"
	"github.com/waiq/example-service/pkg/models"
	"github.com/waiq/example-service/pkg/repository"
	"github.com/waiq/example-service/pkg/service"
	testcontainer "github.com/waiq/example-service/pkg/test/testcontainers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"net/http"
	"net/http/httptest"
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

func TestIntegrationFindBook(t *testing.T) {

	ctx := context.Background()

	//Test setup
	require := require.New(t)
	tx := database.Begin()
	defer tx.Rollback()

	repo, err := repository.New(context.Background(), tx)
	require.NoError(err)

	bookService := service.NewBookService(context.Background(), repo)
	handler := handlers.NewBooksHandler(bookService)

	ts := httptest.NewServer(handler.Handler())
	defer ts.Close()

	client, err := testclient.NewClientWithResponses(ts.URL)
	require.Nil(err)

	m := &models.Book{
		UUID:   uuid.New(),
		Title:  "Mega Man Man",
		Author: "Bill the Bill",
	}

	err = bookService.AddBook(context.Background(), m)
	require.NoError(err)

	// Test begin
	resp, err := client.GetBooksBookIdWithResponse(ctx, m.UUID.String())
	require.Nil(err)
	require.NotNil(resp)
	require.Equal(http.StatusOK, resp.StatusCode())

	require.Equal(m.UUID.String(), *resp.JSON200.Id)
	require.Equal(m.Title, *resp.JSON200.Title)
	require.Equal(m.Author, *resp.JSON200.Author)

	resp, err = client.GetBooksBookIdWithResponse(ctx, uuid.NewString())
	require.Nil(err)
	require.Equal(http.StatusNotFound, resp.StatusCode())
}

func TestIntegrationAddBook(t *testing.T) {

	//Test setup
	ctx := context.Background()

	require := require.New(t)
	tx := database.Begin()
	defer tx.Rollback()

	repo, err := repository.New(context.Background(), tx)
	require.NoError(err)

	bookService := service.NewBookService(context.Background(), repo)
	handler := handlers.NewBooksHandler(bookService)

	ts := httptest.NewServer(handler.Handler())
	defer ts.Close()

	client, err := testclient.NewClientWithResponses(ts.URL)
	require.Nil(err)

	var Title = "Mega Man Man"
	var Author = "Bill the Bill"

	resp, err := client.PostBooksWithResponse(ctx, testclient.Book{
		Title:  &Title,
		Author: &Author,
	})

	require.Nil(err)
	require.NotNil(resp.JSON201.Id)
	require.NoError(uuid.Validate(*resp.JSON201.Id))

	require.Equal(Title, *resp.JSON201.Title)
	require.Equal(Author, *resp.JSON201.Author)

}

func TestIntegrationGetBooks(t *testing.T) {

	//Test setup
	ctx := context.Background()

	require := require.New(t)
	tx := database.Begin()
	defer tx.Rollback()

	repo, err := repository.New(context.Background(), tx)
	require.NoError(err)

	bookService := service.NewBookService(context.Background(), repo)
	handler := handlers.NewBooksHandler(bookService)

	data := []struct {
		UUID   uuid.UUID
		Title  string
		Author string
	}{
		{uuid.New(), "Mega She She", "Bull the Bull"},
		{uuid.New(), "Bending noses", "Marty mac Smart"},
		{uuid.New(), "Bending noses second edition", "Marty mac Smart"},
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

	ts := httptest.NewServer(handler.Handler())
	defer ts.Close()
	// Test begin

	client, err := testclient.NewClientWithResponses(ts.URL)
	require.Nil(err)

	resp, err := client.GetBooksWithResponse(ctx)
	require.Nil(err)
	require.NotNil(resp.JSON200)
	require.Len(*resp.JSON200, len(data))

}
