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

package zaplogger

import (
	"errors"
	"flag"

	"github.com/tadnavel/gomocore/components/loggerc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type serviceLogger struct {
	level       string
	atomicLevel zap.AtomicLevel
	logger      *zap.Logger
}

func NewZapLogger() (loggerc.ServiceLogger, error) {
	zl, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return &serviceLogger{
		atomicLevel: zap.NewAtomicLevel(),
		logger:      zl,
	}, nil
}

func (a *serviceLogger) InitFlags() {
	flag.StringVar(
		&a.level,
		"log-level",
		"",
		"Log level: debug | info | warn | error",
	)
}

func (a *serviceLogger) Activate() error {
	if a.level == "" {
		return errors.New("log level cannot be empty")
	}

	lv, err := zapcore.ParseLevel(a.level)
	if err != nil {
		return err
	}

	a.atomicLevel.SetLevel(lv)

	var cfg zap.Config
	if lv == zap.DebugLevel {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.Level = a.atomicLevel

	logger, err := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return err
	}

	a.logger = logger
	return nil
}

func (a *serviceLogger) Stop() error {
	if a.logger != nil {
		_ = a.logger.Sync()
	}
	return nil
}

func (a *serviceLogger) GetLogger(prefix string) loggerc.Logger {
	return &zapLogger{
		sugar: a.logger.With(zap.String("prefix", prefix)).Sugar(),
	}
}

func (a *serviceLogger) GetLevel() string {
	return a.level
}

func (a *serviceLogger) SetLevel(level string) error {
	lv, err := zapcore.ParseLevel(level)
	if err != nil {
		return err
	}
	a.level = level
	a.atomicLevel.SetLevel(lv)
	return nil
}
