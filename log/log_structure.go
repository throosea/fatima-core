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
	"errors"
	"os"
	"strconv"
	"time"
	"strings"
)

const (
	LOG_NONE  = 0x0  // 0000 0000
	LOG_ERROR = 0x7  // 0000 0111
	LOG_WARN  = 0xF  // 0000 1111
	LOG_INFO  = 0x1F // 0001 1111
	LOG_DEBUG = 0x2F // 0010 1111
	LOG_TRACE = 0xFF // 1111 1111
)

var loggingPreference = new(log4FatimaPreference)

type log4FatimaPreference struct {
	showMethod         bool
	logFolder          string
	keepingDays        int
	sourcePrintSize		int
	processName        string
	logFileLoaded      bool
	logFilePath        string
	currentLogFileTime time.Time
	logFilePtr         *os.File
}

type LogMessage interface {
	getTime() *time.Time
	getMessage() string
	setLevel(level int)
	setArgs(args ...interface{})
	publish()
}

type LogLevel uint8

var _level LogLevel

func (this LogLevel) String() string {
	switch this {
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_INFO:
		return "INFO"
	case LOG_TRACE:
		return "TRACE"
	case LOG_WARN:
		return "WARN"
	case LOG_ERROR:
		return "ERROR"
	}
	return "LOG_NONE"
}

func ToLogLevel(value string) (LogLevel, error) {
	if len(value) < 3 || (value[1] != 'x' && value[1] != 'X') {
		return LOG_NONE, errors.New("invalid value format")
	}
	parsed, err := strconv.ParseInt(value[2:], 16, 64)
	if err != nil {
		return LOG_NONE, err
	}

	switch parsed {
	case LOG_ERROR:
		return LOG_ERROR, nil
	case LOG_WARN:
		return LOG_WARN, nil
	case LOG_INFO:
		return LOG_INFO, nil
	case LOG_DEBUG:
		return LOG_DEBUG, nil
	case LOG_TRACE:
		return LOG_TRACE, nil
	}
	return LOG_NONE, nil
}

func ToLogLevelString(value string) string {
	if len(value) < 0 {
		return "0x0"
	}

	switch strings.ToLower(value) {
	case "info" :
		return "0x1F"
	case "debug" :
		return "0x2F"
	case "warn" :
		return "0xF"
	case "error" :
		return "0x7"
	case "trace" :
		return "0xFF"
	}

	return "0x0"
}

/**
* Set our logging level.
* @param {int} The level to set to.
 */
func SetLevel(level LogLevel) {
	_level = level
}

func GetLevel() LogLevel {
	return (_level)
}
