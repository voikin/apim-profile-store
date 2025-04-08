package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog"
	"github.com/voikin/apim-profile-store/internal/config"
	controller_pkg "github.com/voikin/apim-profile-store/internal/controller"
	neo4j_repo "github.com/voikin/apim-profile-store/internal/repository/neo4j"
	"github.com/voikin/apim-profile-store/internal/repository/postgres"
	usecase_pkg "github.com/voikin/apim-profile-store/internal/usecase"
	profilestorepb "github.com/voikin/apim-profile-store/pkg/api/v1"
	"github.com/voikin/apim-profile-store/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	shutdownTimeout = 10 * time.Second //nolint:gochecknoglobals // global by design
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to load config")
	}

	logger.InitGlobalLogger(cfg.Logger)

	grpcAddr := fmt.Sprintf(":%d", cfg.Server.GRPC.Port)
	httpAddr := fmt.Sprintf(":%d", cfg.Server.HTTP.Port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgCfg, err := pgxpool.ParseConfig(cfg.Postgres.DSN)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("pgxpool.ParseConfig failed")
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("pgxpool.NewWithConfig failed")
	}

	if err = pool.Ping(ctx); err != nil {
		logger.Logger.Fatal().Err(err).Msg("postgres.Ping failed")
	}

	trManager := manager.Must(trmpgx.NewDefaultFactory(pool))
	ctxGetter := trmpgx.DefaultCtxGetter

	neo4jDriver, err := neo4j.NewDriverWithContext(cfg.Neo4J.URI, neo4j.BasicAuth(cfg.Neo4J.Username, cfg.Neo4J.Password, ""))
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("neo4j.NewDriverWithContext failed")
	}

	postgresRepo := postgres.New(pool, trManager, ctxGetter)
	neo4jRepo := neo4j_repo.New(neo4jDriver, trManager)
	usecase := usecase_pkg.New(postgresRepo, neo4jRepo, neo4jRepo)
	controller := controller_pkg.New(usecase)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	loggableEvents := []logging.LoggableEvent{
		logging.StartCall,
		logging.FinishCall,
	}

	if logger.Logger.GetLevel() == zerolog.DebugLevel {
		loggableEvents = append(loggableEvents, logging.PayloadReceived, logging.PayloadSent)
	}

	loggerOpts := []logging.Option{
		logging.WithLogOnEvents(loggableEvents...),
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(logger.InterceptorLogger(logger.Logger), loggerOpts...),
		),
		grpc.ConnectionTimeout(
			cfg.Server.GRPC.MaxConnectionAge(),
		),
	)
	profilestorepb.RegisterProfileStoreServiceServer(grpcServer, controller)

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Logger.Fatal().Err(err).Str("addr", grpcAddr).Msg("failed to listen")
	}

	go func() {
		logger.Logger.Info().Str("addr", grpcAddr).Msg("gRPC server listening")
		if err = grpcServer.Serve(grpcListener); err != nil {
			logger.Logger.Fatal().Err(err).Msg("gRPC server error")
		}
	}()

	gwMux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err = profilestorepb.RegisterProfileStoreServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, dialOpts); err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to register gRPC-Gateway")
	}

	swaggerMux := http.NewServeMux()
	swaggerMux.HandleFunc("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pkg/api/v1/api.swagger.json")
	})
	swaggerMux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("pkg/swagger"))))

	mainMux := http.NewServeMux()
	mainMux.Handle("/", gwMux)
	mainMux.Handle("/swagger/", swaggerMux)

	httpServer := &http.Server{
		Addr:              httpAddr,
		Handler:           mainMux,
		ReadTimeout:       cfg.Server.HTTP.ReadTimeout(),
		WriteTimeout:      cfg.Server.HTTP.WriteTimeout(),
		ReadHeaderTimeout: cfg.Server.HTTP.ReadHeaderTimeout(),
	}

	go func() {
		logger.Logger.Info().Str("addr", httpAddr).Msg("HTTP server listening")
		if err = httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Logger.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	<-shutdownCh
	logger.Logger.Info().Msg("shutdown signal received")

	go func() {
		logger.Logger.Info().Msg("stopping gRPC server...")
		grpcServer.GracefulStop()
	}()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err = httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Error().Err(err).Msg("HTTP shutdown error")
	} else {
		logger.Logger.Info().Msg("HTTP server shut down cleanly")
	}

	logger.Logger.Info().Msg("Server exited gracefully")
}
