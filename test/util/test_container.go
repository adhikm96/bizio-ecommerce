package test_util

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/net/context"
	"log"
	"log/slog"
	"strings"
	"time"
)

var postgresContainer *postgres.PostgresContainer = nil

func Terminate(ctx context.Context) {
	if err := postgresContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}

func SetUpTestContainers(ctx context.Context) {

	if postgresContainer != nil {
		return
	}

	slog.Info("setting up test-containers")

	dbName := "testdb"
	dbUser := "postgres"
	dbPassword := "postgres"

	postgresContainer2, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		//postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
		//postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	postgresContainer = postgresContainer2

	if err != nil {
		panic(err.Error())
	}

	config.Config["DB_NAME"] = dbName
	config.Config["DB_USERNAME"] = dbUser
	config.Config["DB_PASSWORD"] = dbPassword

	port, err := postgresContainer.MappedPort(ctx, "5432")

	if err != nil {
		panic(err.Error())
	}

	config.Config["DB_PORT"] = strings.Split(string(port), "/")[0]
	host, err := postgresContainer.Host(ctx)

	if err != nil {
		panic(err.Error())
	}

	config.Config["DB_HOST"] = host

	if err != nil {
		slog.Error("failed to start container: " + err.Error())
		panic(err.Error())
	}
}
