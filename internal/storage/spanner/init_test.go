package spanner

import (
	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	v1 "cloud.google.com/go/spanner/apiv1"
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/spanner"
	migrateDB "github.com/golang-migrate/migrate/v4/database/spanner"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/stringtool"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"os"
	"testing"
	"time"
)

var (
	session       *spanner.Client
	dbAdminClient *database.DatabaseAdminClient
	apiClient     *v1.Client
	options       config.Spanner
)

const dir = "deployments/migrations/spanner"

func dbName(opt config.Spanner) string {
	return fmt.Sprintf("projects/%s/instances/%s/databases/%s", opt.ProjectID, opt.Instance, opt.Database)
}

func makeClient(logger *zap.Logger, opt config.Spanner) (*spanner.Client, *database.DatabaseAdminClient, *v1.Client, func(), error) {
	// Despite the docs, this context is also used for auth,
	// so it needs to be long-lived.
	ctx := context.Background()

	dialCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(dialCtx, opt.EndPoint, grpc.WithInsecure())
	if err != nil {
		logger.Sugar().Warnf("Dialing in-memory fake: %v", err)
		return nil, nil, nil, nil, err
	}

	client, err := spanner.NewClient(ctx, dbName(opt), option.WithGRPCConn(conn))
	if err != nil {
		logger.Sugar().Warnf("Connecting to local DB: %v", err)
		return nil, nil, nil, nil, err
	}
	adminClient, err := database.NewDatabaseAdminClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		logger.Sugar().Warnf("Connecting to local DB admin: %v", err)
		return nil, nil, nil, nil, err
	}

	gapicClient, err := v1.NewClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		logger.Sugar().Warnf("Connecting to local DB generated Spanner client: %v", err)
		return nil, nil, nil, nil, err
	}

	return client, adminClient, gapicClient, func() {
			client.Close()
			adminClient.Close()
			gapicClient.Close()
			conn.Close()
		},
		nil
}

func TestMain(m *testing.M) {
	logger := zap.NewExample()

	envOption := config.Spanner{}
	if err := config.LoadFromEnv(&envOption); err != nil {
		logger.Sugar().Warn(err)
		return
	}
	client, adminClient, generatedClient, cleanup, err := makeClient(logger, envOption)
	if err != nil {
		logger.Sugar().Warn(err)
		return
	}
	defer cleanup()

	session = client
	dbAdminClient = adminClient
	apiClient = generatedClient
	options = envOption

	db := migrateDB.NewDB(*dbAdminClient, *session)
	if err := SetupDB(options, db); err != nil && err != migrate.ErrNoChange {
		logger.Sugar().Info(err)
		return
	}

	time.Sleep(1 * time.Second)
	code := m.Run()
	time.Sleep(1 * time.Second)
	if err := TeardownDB(options, db); err != nil {
		logger.Sugar().Info(err)
		return
	}
	time.Sleep(2 * time.Second)
	os.Exit(code)
}

func SetupDB(opt config.Spanner, db *migrateDB.DB) error {
	driver, err := migrateDB.WithInstance(db, &migrateDB.Config{
		DatabaseName:    dbName(opt),
		CleanStatements: true,
	})
	if err != nil {
		return err
	}
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		stringtool.StringJoin("file:///", path, "/../../../", dir),
		fmt.Sprintf("spanner://%s", dbName(opt)), driver)
	if err != nil {
		return err
	}
	return m.Up()
}

func TeardownDB(opt config.Spanner, db *migrateDB.DB) error {
	driver, err := migrateDB.WithInstance(db, &migrateDB.Config{
		DatabaseName:    dbName(opt),
		CleanStatements: true,
	})
	if err != nil {
		return err
	}
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		stringtool.StringJoin("file:", path, "/../../../", dir),
		fmt.Sprintf("spanner://%s", dbName(opt)), driver)
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil {
		return err
	}
	return m.Drop()
}
