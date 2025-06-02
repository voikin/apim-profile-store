package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
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
	"github.com/voikin/apim-profile-store/pkg/logger"
	profilestorepb "github.com/voikin/apim-proto/gen/go/apim_profile_store/v1"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, trManager, ctxGetter := setupPostgres(ctx, cfg.Postgres)
	neo4jDriver := setupNeo4j(cfg.Neo4J)

	postgresRepo := postgres.New(pool, trManager, ctxGetter)
	neo4jRepo := neo4j_repo.New(neo4jDriver, trManager)
	usecase := usecase_pkg.New(postgresRepo, neo4jRepo, neo4jRepo)
	controller := controller_pkg.New(usecase)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	grpcServer := runGRPCServer(cfg.Server.GRPC, controller)
	httpServer := runHTTPServer(ctx, cfg.Server.HTTP, grpcAddr)

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

func setupPostgres(ctx context.Context, cfg *config.Postgres) (*pgxpool.Pool, *manager.Manager, *trmpgx.CtxGetter) {
	pgCfg, err := pgxpool.ParseConfig(cfg.DSN)
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
	return pool, trManager, ctxGetter
}

func setupNeo4j(cfg *config.Neo4J) neo4j.DriverWithContext {
	driver, err := neo4j.NewDriverWithContext(
		cfg.URI,
		neo4j.BasicAuth(cfg.Username, cfg.Password, ""),
	)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("neo4j.NewDriverWithContext failed")
	}
	return driver
}

func runGRPCServer(cfg *config.GRPC, controller profilestorepb.ProfileStoreServiceServer) *grpc.Server {
	grpcAddr := fmt.Sprintf(":%d", cfg.Port)

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
		grpc.ConnectionTimeout(cfg.MaxConnectionAge()),
	)

	profilestorepb.RegisterProfileStoreServiceServer(grpcServer, controller)

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Logger.Fatal().Err(err).Str("addr", grpcAddr).Msg("failed to listen")
	}

	go func() {
		logger.Logger.Info().Int("port", cfg.Port).Msg("gRPC server listening")
		err = grpcServer.Serve(listener)
		if err != nil {
			logger.Logger.Fatal().Err(err).Msg("gRPC server error")
		}
	}()

	return grpcServer
}

func runHTTPServer(ctx context.Context, cfg *config.HTTP, grpcAddr string) *http.Server {
	httpAddr := fmt.Sprintf(":%d", cfg.Port)

	gwMux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := profilestorepb.RegisterProfileStoreServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, dialOpts); err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to register gRPC-Gateway")
	}

	swaggerMux := http.NewServeMux()
	swaggerMux.HandleFunc("/swagger/swagger.json", func(w http.ResponseWriter, _ *http.Request) {
		swaggerURL := getSwaggerURL()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, swaggerURL, nil)
		if err != nil {
			http.Error(w, "failed to fetch swagger", http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "failed to fetch swagger", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.Copy(w, resp.Body)
	})
	swaggerMux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("pkg/swagger"))))

	mainMux := http.NewServeMux()
	mainMux.Handle("/", gwMux)
	mainMux.Handle("/swagger/", swaggerMux)

	httpServer := &http.Server{
		Addr:              httpAddr,
		Handler:           mainMux,
		ReadTimeout:       cfg.ReadTimeout(),
		WriteTimeout:      cfg.WriteTimeout(),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout(),
	}

	httpServer.Handler = CORSMiddleware(httpServer.Handler)

	go func() {
		logger.Logger.Info().Int("port", cfg.Port).Msg("HTTP server listening")
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logger.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	return httpServer
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
