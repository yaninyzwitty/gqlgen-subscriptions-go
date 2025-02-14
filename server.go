package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/yaninyzwitty/gqlgen-subscriptions-go/database"
	"github.com/yaninyzwitty/gqlgen-subscriptions-go/graph"
	"github.com/yaninyzwitty/gqlgen-subscriptions-go/helpers"
	"github.com/yaninyzwitty/gqlgen-subscriptions-go/kafka"
	"github.com/yaninyzwitty/gqlgen-subscriptions-go/pkg"
	"github.com/yaninyzwitty/gqlgen-subscriptions-go/sonyflake"
)

var (
	cfg pkg.Config
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	file, err := os.Open("my_config.yaml")
	if err != nil {
		slog.Error("failed to open config.yaml", "error", err)
		os.Exit(1)
	}
	defer file.Close()

	if err := cfg.LoadConfig(file); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// initialize sonyflake
	err = sonyflake.InitSonyFlake()
	if err != nil {
		slog.Error("failed to initialize sonyflake", "error", err)
		os.Exit(1)
	}

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	// if err := godotenv.Load(); err != nil {
	// 	slog.Error("Failed to load .env file", "error", err)
	// 	os.Exit(1)
	// }

	password := helpers.GetEnvOrDefault("KAFKA_PASSWORD", "")
	dbPassword := helpers.GetEnvOrDefault("SCYLLA_PASSWORD", "")
	token := helpers.GetEnvOrDefault("ASTRA_TOKEN", "")

	slog.Info("Environment variables loaded successfully", "KAFKA_PASSWORD", password != "", "SCYLLA_PASSWORD", dbPassword != "", "ASTRA_TOKEN", token != "")

	kafkaConfig := &kafka.KafkaInputConfig{
		Username:         cfg.Kafka.Username,
		BootstrapServers: cfg.Kafka.BootstrapServers,
		SecurityProtocol: cfg.Kafka.SecurityProtocol,
		SASLMechanism:    cfg.Kafka.SASLMechanism,
		RegistryUrl:      cfg.Kafka.RegistryUrl,
		Password:         password,
	}

	astraCfg := &database.AstraConfig{
		Username: cfg.AstraDB.Username,
		Path:     cfg.AstraDB.Path,
		Token:    token,
	}

	db := database.NewAstraDB()

	session, err := db.Connect(ctx, astraCfg, time.Duration(cfg.AstraDB.Timeout)*time.Second)
	if err != nil {
		slog.Error("error connecting to database", "error", err)
		os.Exit(1)

	}
	defer session.Close()

	// databaseConfig := database.DBConfig{
	// 	Username:        cfg.Database.Username,
	// 	Hosts:           cfg.Database.Hosts,
	// 	LocalDataCenter: cfg.Database.LocalDataCenter,
	// 	Password:        dbPassword,
	// }

	// dbInstance := database.NewScyllaDB()

	// // create a session
	// session, err := dbInstance.Connect(ctx, &databaseConfig)
	// if err != nil {
	// 	slog.Error("error connecting to database", "error", err)
	// 	os.Exit(1)
	// }
	// defer session.Close()

	kafkaInstance := kafka.NewKafka(kafkaConfig)

	// make a kafka writer

	writer, err := kafkaInstance.CreateKafkaWriter(ctx, cfg.Kafka.Topic)

	if err != nil {
		slog.Error("error connecting to kafka", "error", err)
		os.Exit(1)
	}

	defer writer.Close()
	reader, err := kafkaConfig.CreateKafkaReader(cfg.Kafka.Topic, cfg.Kafka.GroupId)
	if err != nil {
		slog.Error("error connecting to kafka", "error", err)
		os.Exit(1)
	}

	defer reader.Close()

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Session: session,
		Reader:  reader,
	}}))
	srv.AddTransport(transport.Websocket{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: mux,
	}

	stopCH := make(chan os.Signal, 1)
	signal.Notify(stopCH, os.Interrupt, syscall.SIGTERM)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	go func() {
		slog.Info("SERVER starting", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := helpers.ProcessMessages(context.Background(), session, writer); err != nil {
					slog.Error("failed to process messages", "error", err)
					os.Exit(1)
				}

			case <-stopCH:
				return
			}
		}
	}()
	<-stopCH
	slog.Info("shutting down the server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	} else {
		slog.Info("server stopped gracefully")
	}
}

// 227088922700738561, 227088978803748865, 227089006368714753, 227089040426463233, 227089068528300033
