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
func (l *zapLogger) WithFields(fields map[string]interface{}) loggerc.Logger {
	var args []interface{}
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &zapLogger{
		sugar: l.sugar.With(args...),
	}
}
