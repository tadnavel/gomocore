// -----------------------------------------------------------------------------
// Copyright (C) 2026 tadnavel
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
// -----------------------------------------------------------------------------

package ginc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
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
			g.logger.Errorf("gin server forced to shutdown: %w", err)
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

func (g *ginEngine) Run() error {
	addr := fmt.Sprintf(":%d", g.port)
	g.httpServer = &http.Server{
		Addr:           addr,
		Handler:        g.router,
		ReadTimeout:    defaultReadTimeout * time.Second,
		WriteTimeout:   defaultWriteTimeout * time.Second,
		IdleTimeout:    defaultIdleTimeout * time.Second,
		MaxHeaderBytes: defaultMaxHeaderBytes,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", g.port, err)
	}

	go func() {
		g.logger.Infof("gin server is running on port %d...", g.port)
		if err := g.httpServer.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			g.logger.Errorf("serve error: %s\n", err)
		}
	}()

	return nil
}
