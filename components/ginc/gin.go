package ginc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sctx "github.com/tadnavel/gomocore"
	"github.com/tadnavel/gomocore/components/loggerc"
)

const (
	defaultPort            = 8080
	defaultMode            = "debug" // debug | release
	defaultShutdownTimeout = 10      // seconds
	defaultReadTimeout     = 10
	defaultWriteTimeout    = 10
	defaultIdleTimeout     = 60
	defaultMaxHeaderBytes  = 1 << 20
)

type config struct {
	port    int
	ginMode string
}

type ginEngine struct {
	*config
	id         string
	logger     loggerc.Logger
	router     *gin.Engine
	httpServer *http.Server
}

func NewGin(id string) *ginEngine {
	return &ginEngine{
		config: new(config),
		id:     id,
	}
}

func (g *ginEngine) ID() string {
	return g.id
}

func (g *ginEngine) Activate(serviceCtx sctx.ServiceContext) error {
	g.logger = serviceCtx.Logger(g.id)

	g.logger.Info("activating gin engine...")
	env := serviceCtx.EnvName()
	mode := gin.ReleaseMode

	if env == sctx.DevEnv {
		mode = gin.DebugMode
	}

	if g.ginMode != "" {
		switch g.ginMode {
		case gin.DebugMode, gin.ReleaseMode:
			mode = g.ginMode
		default:
			return fmt.Errorf("invalid gin mode: %s (allowed: debug | release", g.ginMode)
		}
	}

	gin.SetMode(mode)
	g.router = gin.New()
	g.logger.Info("gin engine started!")

	return nil
}

func (g *ginEngine) Stop() error {
	if g.httpServer != nil {
		g.logger.Info("shutting down gin server...")
		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout*time.Second)

		defer cancel()

		if err := g.httpServer.Shutdown(ctx); err != nil {
			g.logger.Errorf("gin server forced to shutdown: %v", err)
			return err
		}
		g.logger.Info("gin server is turned off!")
	}
	return nil
}

func (g *ginEngine) InitFlags() {
	flag.IntVar(&g.port, "gin-port", defaultPort, fmt.Sprintf("gin server port. Default %d", defaultPort))
	flag.StringVar(&g.ginMode, "gin-mode", defaultMode, fmt.Sprintf("gin server (debug | release). Default %s", defaultMode))
}

func (g *ginEngine) GetPort() int {
	return g.port
}

func (g *ginEngine) GetRouter() *gin.Engine {
	return g.router
}

func (g *ginEngine) Run() {
	addr := fmt.Sprintf(":%d", g.port)
	g.httpServer = &http.Server{
		Addr:           addr,
		Handler:        g.router,
		ReadTimeout:    defaultReadTimeout * time.Second,
		WriteTimeout:   defaultWriteTimeout * time.Second,
		IdleTimeout:    defaultIdleTimeout * time.Second,
		MaxHeaderBytes: defaultMaxHeaderBytes,
	}

	go func() {
		g.logger.Infof("gin server is running on port %d...", g.port)
		if err := g.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			g.logger.Errorf("listen: %s\n", err)
		}
	}()
}
