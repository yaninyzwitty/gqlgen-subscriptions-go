package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/gocql/gocql"
)

type DBConfig struct {
	Username        string
	Hosts           []string
	LocalDataCenter string
	Password        string
}

type DatabaseMethods interface {
	Connect(ctx context.Context, cfg *DBConfig) (*gocql.Session, error)
}

type ScyllaDB struct{}

func NewScyllaDB() DatabaseMethods {
	return &ScyllaDB{}
}

func (db *ScyllaDB) Connect(ctx context.Context, cfg *DBConfig) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	cluster.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy(cfg.LocalDataCenter)
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		slog.Error("Failed to create ScyllaDB session", "error", err)
		return nil, err
	}

	slog.Info("ScyllaDB session created successfully")
	return session, nil
}
