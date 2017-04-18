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
	"os"
	"path/filepath"
	"time"
	"io/ioutil"
	"regexp"
	"strings"
)

var logfileSizeLimit = 0	// unit : MB

func writeLogMessage(log LogMessage) {
	log.publish()
	ensureLogFile()
	ensureTodayLog(log.getTime())
	writeString(log.getMessage())
}

func ensureTodayLog(t *time.Time) {
	// logFolder/processName.log 파일의 시각을 확인한다
	if loggingPreference.currentLogFileTime.Year() != t.Year() ||
		loggingPreference.currentLogFileTime.Month() != t.Month() ||
		loggingPreference.currentLogFileTime.Day() != t.Day() {
		moveToBackupLog()
	}
}

// 로그 파일이 존재하는지 확인한다
func ensureLogFile() {
	if loggingPreference.logFileLoaded {
		return
	}

	var err error
	var stat os.FileInfo

	stat, err = os.Stat(loggingPreference.logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			loggingPreference.logFilePtr, err = os.Create(loggingPreference.logFilePath)
			if err != nil {
				fmt.Printf("%s fail to create : %s", loggingPreference.logFilePath, err)
				loggingPreference.logFilePtr = nil
				return
			}
			loggingPreference.currentLogFileTime = time.Now()
		} else if stat.IsDir() {
			fmt.Printf("%s path exist as directory. fail to logging", loggingPreference.logFilePath)
			loggingPreference.logFilePtr = nil
		}
	} else {
		loggingPreference.logFilePtr, err = os.OpenFile(loggingPreference.logFilePath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Printf("fail to open : %s", err)
			loggingPreference.logFilePtr = nil
		}
		loggingPreference.currentLogFileTime = stat.ModTime()
	}

	loggingPreference.logFileLoaded = true
}

// 오래된(지난 날짜) 로그 파일을 이동시키고 신규 로그 파일을 생성한다
func moveToBackupLog() {
	var err error
	var stat os.FileInfo

	stat, err = os.Stat(loggingPreference.logFilePath)
	if err != nil {
		fmt.Printf("fail to stat log file : %s\n", err)
		loggingPreference.logFilePtr = nil
		return
	}

	// close current log file ptr
	if loggingPreference.logFilePtr != nil {
		loggingPreference.logFilePtr.Close()
		loggingPreference.logFilePtr = nil
	}

	// move current file to backup
	backupFilePath := fmt.Sprintf("%s%c%s.%s.log",
		loggingPreference.logFolder,
		filepath.Separator,
		loggingPreference.processName, stat.ModTime().Format("2006-01-02"))
	os.Rename(loggingPreference.logFilePath, backupFilePath)

	// open for new log file
	loggingPreference.logFilePtr, err = os.OpenFile(loggingPreference.logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("fail to open for new log file : %s\n", err)
	}

	stat, _ = loggingPreference.logFilePtr.Stat()
	loggingPreference.currentLogFileTime = stat.ModTime()

	go func() {
		backupDaysChanged()
	}()
}

// 로그 내용을 파일에 기록한다
func writeString(s string) (n int, err error) {
	if loggingPreference.logFilePtr == nil {
		return 0, nil
	}
	return loggingPreference.logFilePtr.WriteString(s)
}

func backupDaysChanged() {
	if loggingPreference.keepingDays < 1 {
		return
	}

	// find files in log path
	files, err := ioutil.ReadDir(loggingPreference.logFolder)
	if err != nil {
		return
	}

	express := fmt.Sprintf("%s\\.[0-9]+-[0-9]+-[0-9]+\\.log", loggingPreference.processName)
	var validLogFileId = regexp.MustCompile(express)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), "log") {
			continue
		}

		if !validLogFileId.MatchString(file.Name()) {
			continue
		}

		dotIndex := strings.Index(file.Name(), ".")
		if dotIndex < 1 {
			continue
		}

		lastDotIndex := strings.LastIndex(file.Name(), ".")
		if lastDotIndex <= dotIndex {
			continue
		}

		createdDateExpression := file.Name()[dotIndex+1:lastDotIndex]
		TIME_YYYYMMDD := "2006-01-02"
		createdDate, err := time.Parse(TIME_YYYYMMDD, createdDateExpression)
		if err != nil {
			continue
		}

		diff := time.Duration(24 * loggingPreference.keepingDays) * time.Hour
		deadline := time.Now().Add(-diff)
		if createdDate.Before(deadline) {
			os.Remove(filepath.Join(loggingPreference.logFolder, file.Name()))
		}
	}
}