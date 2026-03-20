// -----------------------------------------------------------------------------
// Copyright (C) 2026 tadnavel
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// -----------------------------------------------------------------------------

package zaplogger

import (
	"github.com/tadnavel/gomocore/components/loggerc"
	"go.uber.org/zap"
)

type zapLogger struct {
	sugar *zap.SugaredLogger
}

func (l *zapLogger) Debug(args ...any) {
	l.sugar.Debug(args...)
}
func (l *zapLogger) Info(args ...any) {
	l.sugar.Info(args...)
}
func (l *zapLogger) Warn(args ...any) {
	l.sugar.Warn(args...)
}
func (l *zapLogger) Error(args ...any) {
	l.sugar.Error(args...)
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.sugar.Debugf(format, args...)
}
func (l *zapLogger) Infof(format string, args ...any) {
	l.sugar.Infof(format, args...)
}
func (l *zapLogger) Warnf(format string, args ...any) {
	l.sugar.Warnf(format, args...)
}
func (l *zapLogger) Errorf(format string, args ...any) {
	l.sugar.Errorf(format, args...)
}

func (l *zapLogger) With(key string, value any) loggerc.Logger {
	return &zapLogger{
		sugar: l.sugar.With(key, value),
	}
}
func (l *zapLogger) WithFields(fields map[string]any) loggerc.Logger {
	var args []any
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &zapLogger{
		sugar: l.sugar.With(args...),
	}
}
