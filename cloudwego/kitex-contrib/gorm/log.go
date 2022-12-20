// Copyright 2022 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type gormLogger struct {
	*config
}

func NewGormLogger(opts ...Option) *gormLogger {
	conf := defaultConfig()
	// apply options
	for _, opt := range opts {
		opt.apply(conf)
	}

	return &gormLogger{config: conf}
}

// LogMode  don't to set gorm log level mode, just use biz log level
func (m *gormLogger) LogMode(logger.LogLevel) logger.Interface {
	return m
}

// Info gorm log info for debug
func (m *gormLogger) Info(ctx context.Context, f string, v ...interface{}) {
	f = fmt.Sprintf(f, append([]interface{}{utils.FileWithLineNum()}, v...)...)
	klog.CtxDebugf(ctx, f, v)
}

func (m *gormLogger) Warn(ctx context.Context, f string, v ...interface{}) {
	f = fmt.Sprintf(f, append([]interface{}{utils.FileWithLineNum()}, v...)...)
	klog.CtxWarnf(ctx, f, v)
}

func (m *gormLogger) Error(ctx context.Context, f string, v ...interface{}) {
	f = fmt.Sprintf(f, append([]interface{}{utils.FileWithLineNum()}, v...)...)
	klog.CtxErrorf(ctx, f, v)
}

func (m *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	traceMeta := map[string]interface{}{}
	end := time.Now()
	latency := end.Sub(begin)
	cost := float64(latency.Nanoseconds()/1e4) / 1e2
	traceMeta["cost"] = cost

	msg := "success"
	if err != nil {
		traceMeta["sqlErrMsg"] = err.Error()
	}

	traceMeta["sql"], traceMeta["rowsAffected"] = fc()
	traceMeta["fileLine"] = utils.FileWithLineNum()

	if m.slowThreshold != 0 && m.slowThreshold.Milliseconds() < latency.Milliseconds() {
		m.traceLogLevel = klog.LevelWarn
		msg = "slow sql, please check sql wheather using index"
	}

	m.kvLogger.CtxKVLog(ctx, m.traceLogLevel, fmt.Sprint(levelName[Sqltrace], msg), Map2KvPairs(traceMeta))
}
