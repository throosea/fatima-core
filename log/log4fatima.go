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
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

const (
	LOG4FATIMA_PROP_BACKUP_DAYS           = "log4fatima.backup.days"
	LOG4FATIMA_PROP_SHOW_METHOD           = "log4fatima.method.show"
	LOG4FATIMA_PROP_SOURCE_PRINTSIZE      = "log4fatima.source.printsize"
	LOG4FATIMA_PROP_FILE_SIZE_LIMIT       = "log4fatima.filesize.limit"
	LOG4FATIMA_DEFAULT_BACKUP_FILE_NUMBER = 30
	LOG4FATIMA_DEFAULT_SOURCE_PRINTSIZE = 30
)

var log4FatimaInitialized = false
var log4FatimaWriting = false

var logMessageChannel = make(chan LogMessage, 128)

func Initialize(folder string, procName string) {
	if log4FatimaInitialized {
		return
	}
	log4FatimaInitialized = true

	loggingPreference.logFolder = folder
	loggingPreference.processName = procName
	loggingPreference.keepingDays = LOG4FATIMA_DEFAULT_BACKUP_FILE_NUMBER
	loggingPreference.showMethod = true
	loggingPreference.sourcePrintSize = LOG4FATIMA_DEFAULT_SOURCE_PRINTSIZE
	loggingPreference.logFilePath = fmt.Sprintf("%s.log", filepath.Join(folder, procName))

	go func() {
		for {
			logging := <-logMessageChannel
			log4FatimaWriting = true
			writeLogMessage(logging)
			if len(logMessageChannel) == 0 {
				log4FatimaWriting = false
			}
		}
	}()
}

func SetSourcePrintSize(newValue int) {
	if newValue < 1 {
		return
	}

	loggingPreference.sourcePrintSize = newValue
}

func SetShowMethod(newValue bool) {
	loggingPreference.showMethod = newValue
}

func SetBackupDays(days int)	{
	if days < 1 {
		return
	}

	var old = loggingPreference.keepingDays
	loggingPreference.keepingDays = days
	if old != days {
		Info("logging backup days changed to %d", loggingPreference.keepingDays)
		go func() {
			backupDaysChanged()
		}()
	}
}

func SetFileSizeLimit(mb int)	{
	if mb < 1 {
		return
	}

	logfileSizeLimit = mb
	Info("logging file size limit to %d MB", logfileSizeLimit)
}

func WaitLoggingShutdown() {
	for {
		if len(logMessageChannel) == 0 && !log4FatimaWriting {
			return
		}
		time.Sleep(time.Millisecond * 1)
	}
}

/**
* Log an error
 */
func Error(v ...interface{}) {
	if _level >= LOG_ERROR && len(v) > 0 {
		print(LOG_ERROR, v...)
	}
}

/**
* Log a warning
 */
func Warn(v ...interface{}) {
	if _level >= LOG_WARN && len(v) > 0 {
		print(LOG_WARN, v...)
	}
}

/**
* Log info
 */
func Info(v ...interface{}) {
	if _level >= LOG_INFO && len(v) > 0 {
		print(LOG_INFO, v...)
	}
}

/**
* Log debugigng
 */
func Debug(v ...interface{}) {
	if _level >= LOG_DEBUG && len(v) > 0 {
		print(LOG_DEBUG, v...)
	}
}

/**
* Log a trace
 */
func Trace(v ...interface{}) {
	if _level >= LOG_TRACE && len(v) > 0 {
		print(LOG_TRACE, v...)
	}
}

/**
* Central function for printing
 */
func print(level int, v ...interface{}) {
	pc, file, line, _ := runtime.Caller(2)

	var message LogMessage

	if _, ok := v[len(v)-1].(error); ok {
		errMessage := newErrorTraceLogMessage(pc, file, line)
		for i := 3; i < 10; i++ {
			pc, file, line, exist := runtime.Caller(i)
			if !exist {
				break
			}
			point := TracePoint{pc: pc, file: file, line: line}
			errMessage.append(point)
		}
		message = errMessage
	} else {
		message = newGeneralLogMessage(pc, file, line)
	}

	message.setLevel(level)
	message.setArgs(v...)
	//message.publish()
	logMessageChannel <- message
}
