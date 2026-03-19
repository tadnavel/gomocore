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

package sctx

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tadnavel/gomocore/components/loggerc"
	"github.com/tadnavel/gomocore/components/loggerc/zaplogger"
)

const (
	DevEnv = "dev"
	StgEnv = "stg"
	PrdEnv = "prd"
)

type Component interface {
	ID() string
	InitFlags()
	Activate(ServiceContext) error
	Stop() error
}

type ServiceContext interface {
	Load() error
	Stop() error

	Logger(prefix string) loggerc.Logger
	LogLevel() string

	EnvName() string
	GetName() string
	Get(id string) (interface{}, bool)
	MustGet(id string) interface{}
}

type serviceCtx struct {
	name       string
	env        string
	components []Component
	store      map[string]Component
	logger     loggerc.Logger
}

var defaultLogger, _ = zaplogger.NewZapLogger()

type Option func(*serviceCtx)

func WithName(name string) Option {
	return func(s *serviceCtx) {
		s.name = name
	}
}

func WithComponent(c Component) Option {
	return func(s *serviceCtx) {
		if _, ok := s.store[c.ID()]; !ok {
			s.components = append(s.components, c)
			s.store[c.ID()] = c
		}
	}
}

func NewServiceContext(opts ...Option) ServiceContext {
	s := &serviceCtx{
		store: make(map[string]Component),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.initFlags()
	s.parseFlags()

	return s
}

func (s *serviceCtx) initFlags() {
	flag.StringVar(&s.env, "service-env", DevEnv, "Env for service: dev | stg | prd")
	defaultLogger.InitFlags()
	for _, c := range s.components {
		c.InitFlags()
	}
}

func (s *serviceCtx) Load() error {
	if defaultLogger.GetLevel() == "" {
		switch s.env {
		case DevEnv:
			_ = defaultLogger.SetLevel("debug")
		case StgEnv:
			_ = defaultLogger.SetLevel("info")
		case PrdEnv:
			_ = defaultLogger.SetLevel("warn")
		default:
			_ = defaultLogger.SetLevel("info")
		}
	}

	if err := defaultLogger.Activate(); err != nil {
		return err
	}

	s.logger = s.Logger("service-context")
	s.logger.Infof(
		"service starting env=%s log_level=%s",
		s.env,
		defaultLogger.GetLevel(),
	)

	for _, c := range s.components {
		s.logger.Infof("activating component: %s", c.ID())
		if err := c.Activate(s); err != nil {
			return fmt.Errorf("activate %s: %v", c.ID(), err)
		}
	}
	return nil
}

func (s *serviceCtx) Stop() error {
	s.logger.Info("stopping service context")
	for _, c := range s.components {
		if err := c.Stop(); err != nil {
			return fmt.Errorf("stop %s: %v", c.ID(), err)
		}
	}
	_ = defaultLogger.Stop()
	return nil
}

func (s *serviceCtx) Logger(prefix string) loggerc.Logger {
	return defaultLogger.GetLogger(prefix)
}

func (s *serviceCtx) Get(id string) (interface{}, bool) {
	c, ok := s.store[id]
	return c, ok
}

func (s *serviceCtx) MustGet(id string) interface{} {
	c, ok := s.Get(id)
	if !ok {
		panic(fmt.Sprintf("cannot get component %s", id))
	}
	return c
}

func (s *serviceCtx) LogLevel() string {
	return defaultLogger.GetLevel()
}

func (s *serviceCtx) EnvName() string { return s.env }
func (s *serviceCtx) GetName() string { return s.name }

// automatically parse flags into .env file format
func (s *serviceCtx) parseFlags() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			s.logger.Errorf("load env file %s: %v", envFile, err)
		}
	}

	flag.VisitAll(func(f *flag.Flag) {
		envKey := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
		if v := os.Getenv(envKey); v != "" {
			if err := flag.Set(f.Name, v); err != nil {
				s.logger.Errorf("parseFlags: fail to parse %s, error: %v", envKey, err)
			}
		}
	})

	flag.Parse()
}
