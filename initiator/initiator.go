package initiator

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	middleware "github.com/TNAHOM/ATS-system-main/internal/glue/middleware"
	"github.com/gin-gonic/gin"
)

func Initiator() {
	logger := InitialzeLog()
	logger.Info("Logger initialized successfully")

	logger.Info("Loading coinfig")
	LoadEnvVariables()
	logger.Info("Config successfully initialized")

	logger.Info("Loadoing DB")
	db := InitializeDB()
	logger.Info("DB successfully initialized")

	logger.Info("Loading syncDB(migration)")
	SyncDatabase(db)
	logger.Info("SyncDB successfully initialized")

	logger.Info("Loading Encryption")
	platform := InitPlatform(logger)
	logger.Info("Platform Successfully initiated")

	logger.Info("Loading Persistance")
	persistance := InitPersistance(db, logger)
	logger.Info("Persistance successfully initialized")

	logger.Info("Loading Module")
	module := InitModule(logger, persistance, platform)
	logger.Info("Module successfully initialized")

	logger.Info("Loading Handler")
	handler := InitHandler(logger, module)
	logger.Info("Handler successfully initialized")

	logger.Info("initializing http server")
	server := gin.New()

	// Use modular CORS middleware (defined in internal/glue/middleware)
	server.Use(middleware.CORS())
	// server.Use(middleware.GinLogger(*logger))
	// server.Use(middleware.CORS())
	// server.Use(middleware.ErrorHandler())
	// Mount all routes under /api
	ginsrv := server.Group("/api")

	logger.Info("initializing route")
	InitRoute(ginsrv, handler, module, logger)
	logger.Info("done initializing route")

	ginsrv.GET("/debug/pprof/*any", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))

	logger.Info("initialize swagger")
	InitSwagger(ginsrv)

	logger.Info("initializing server")
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		Handler:           server,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT)
		<-sigint
		log.Fatal("HTTP server Shutdown")

	}()
	logger.Info(fmt.Sprintf("http server listening on port : %s", os.Getenv("PORT")))

	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Could not start HTTP server: %s", err))
	}

}
