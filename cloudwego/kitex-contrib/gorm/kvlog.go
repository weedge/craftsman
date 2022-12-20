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
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
)

// DefaultLogger is default logger.
var DefaultLogger = NewStdLogger(log.Writer())

type IkvLogger interface {
	// CtxKVLog  kvs must be kv pairs k,v , k,v ...
	CtxKVLog(ctx context.Context, level int, format string, kvs ...interface{})
}

func Map2KvPairs(mapData map[string]interface{}) (kvs []interface{}) {
	for key, val := range mapData {
		kvs = append(kvs, key, val)
	}

	return
}

var _ IkvLogger = (*stdLogger)(nil)

type stdLogger struct {
	log  *log.Logger
	pool *sync.Pool
}

// NewStdLogger new a logger with writer.
func NewStdLogger(w io.Writer) IkvLogger {
	return &stdLogger{
		log: log.New(w, "", 0),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Log print the kv pairs log.
func (l *stdLogger) CtxKVLog(ctx context.Context, level int, format string, kvs ...interface{}) {
	if len(kvs) == 0 {
		return
	}

	if (len(kvs) & 1) == 1 {
		kvs = append(kvs, "KEYVALS UNPAIRED")
	}
	buf := l.pool.Get().(*bytes.Buffer)
	buf.WriteString(Level(level).toString())
	for i := 0; i < len(kvs); i += 2 {
		_, _ = fmt.Fprintf(buf, " %s=%v", kvs[i], kvs[i+1])
	}
	_ = l.log.Output(4, buf.String()) //nolint:gomnd

	buf.Reset()
	l.pool.Put(buf)
}
