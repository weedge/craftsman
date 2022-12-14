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
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"
)

type Level int

// log level
const (
	Debug Level = iota
	Info
	Warn
	Err
	Sqltrace Level = Level(klog.LevelWarn)
)

var levelName = []string{
	Debug:    "[DEBUG]",
	Info:     "[INFO]",
	Warn:     "[WARN]",
	Err:      "[ERROR]",
	Sqltrace: "[SQLTRACE]",
}

func (lv Level) ToString() string {
	if lv >= Debug && lv <= Sqltrace {
		return levelName[lv]
	}
	return fmt.Sprintf("[?%d] ", lv)
}
