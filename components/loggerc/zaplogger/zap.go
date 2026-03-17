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
	"github.com/tadnavel/gomocore/components/loggerc"
	"go.uber.org/zap"
)

type zapLogger struct {
	sugar *zap.SugaredLogger
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}
func (l *zapLogger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}
func (l *zapLogger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}
func (l *zapLogger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}
func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

func (l *zapLogger) With(key string, value interface{}) loggerc.Logger {
	return &zapLogger{
		sugar: l.sugar.With(key, value),
	}
}
