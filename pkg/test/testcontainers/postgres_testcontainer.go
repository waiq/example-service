package testcontainer

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName     string = "test"
	dbUser     string = "root"
	dbPassword string = "secret"
)

type (
	PostgreSQLTestContainerVersion struct {
		Image, WaitFor string
	}

	PostgreSQLTestContainer struct {
		*postgres.PostgresContainer
		ConnectionString string
	}
)

// This is the version we're currently using on GCP
var PostgreSQL_16 = PostgreSQLTestContainerVersion{
	Image:   "postgres:16.1",
	WaitFor: "port: 5432 PostgreSQL Community Server (GPL)",
}

func CreatePostgreSQLContainer(
	version PostgreSQLTestContainerVersion,
) (PostgreSQLTestContainer, func(context.Context)) {
	ctx := context.Background()

	fmt.Println("starting PostgreSQL test container")

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(version.Image),
		postgres.WithDatabase(dbName),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		panic(err)
	}

	connectionString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	return PostgreSQLTestContainer{
			PostgresContainer: pgContainer,
			ConnectionString:  connectionString,
		}, func(ctx context.Context) {
			pgContainer.Terminate(ctx)
		}
}
