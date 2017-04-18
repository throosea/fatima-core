//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
// @project fatima
// @author DeockJin Chung (jin.freestyle@gmail.com)
// @date 2017. 3. 6. PM 7:42
//

package log

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"time"
)

func newGeneralLogMessage(pc uintptr, file string, line int) *GeneralLogMessage {
	message := GeneralLogMessage{}
	message.t = time.Now()
	message.pc = pc
	message.file = file
	message.line = line
	if loggingPreference.showMethod {
		message.funcName = findFunctionName(pc)
	}
	return &message
}

type GeneralLogMessage struct {
	t         time.Time
	pc        uintptr
	level     string
	file      string
	funcName  string
	line      int
	message   []interface{}
	published string
}

func (this *GeneralLogMessage) getMessage() string {
	return this.published
}

func (this *GeneralLogMessage) publish() {
	this.published = this.buildMessage(func() string {
		if format, ok := this.message[0].(string); ok {
			return fmt.Sprintf(format, this.message[1:]...)
		} else {
			return fmt.Sprintf("%s", this.message[0])
		}
	})
}

func (this *GeneralLogMessage) buildMessage(f func() string) string {
	var location = this.file
	var found = strings.LastIndex(this.file, "/src/")
	if found > 0 {
		location = string(this.file[found+5:])
	}

	var buffer bytes.Buffer
	var tokens = strings.Split(location, "/")
	var length = len(tokens)
	for i, s := range tokens {
		if i < length-1 {
			buffer.WriteByte(s[0])
			buffer.WriteByte('.')
		} else {
			buffer.WriteString(s[:len(s)-3])
		}
	}

	if loggingPreference.showMethod {
		return fmt.Sprintf("%s %s [%s] %s\n",
			this.t.Format("2006-01-02 15:04:05.000"),
			this.level,
			this.buildSourceDescription(buffer.String()),
			f())
	} else {
		return fmt.Sprintf("%s %s [%s] %s\n",
			this.t.Format("2006-01-02 15:04:05.000"),
			this.level,
			this.buildSourceDescription(buffer.String()),
			f())
	}
}

func (this *GeneralLogMessage) buildSourceDescription(source string) string {
	var message string
	if loggingPreference.showMethod {
		message = fmt.Sprintf("%s.%s():%d", source,  this.funcName, this.line)
	} else {
		message = fmt.Sprintf("%s:%d", source,  this.line)
	}

	startIndex := len(message) - loggingPreference.sourcePrintSize
	if startIndex >= 0 {
		return message[startIndex:]
	}

	var buffer bytes.Buffer
	for i:=startIndex; i<0; i++ {
		buffer.WriteByte(' ')
	}
	buffer.WriteString(message)
	return buffer.String()
}

func (this *GeneralLogMessage) getTime() *time.Time {
	return &this.t
}

func (this *GeneralLogMessage) setLevel(level int) {
	switch level {
	case LOG_ERROR:
		this.level = "ERROR"
	case LOG_WARN:
		this.level = "WARN "
	case LOG_INFO:
		this.level = "INFO "
	case LOG_DEBUG:
		this.level = "DEBUG"
	case LOG_TRACE:
		this.level = "TRACE"
	}
}

func (this *GeneralLogMessage) setArgs(args ...interface{}) {
	this.message = args
}

func findFunctionName(pc uintptr) string {
	var funcName = runtime.FuncForPC(pc).Name()
	var found = strings.LastIndexByte(funcName, '.')
	if found < 0 {
		return funcName
	}
	return funcName[found+1:]
}
