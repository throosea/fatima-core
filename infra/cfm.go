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

package infra

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"throosea.com/fatima"
	"throosea.com/fatima/log"
	"throosea.com/fatima/monitor"
)

type logLevelItem struct {
	process string
	level   string
}

type CentralFilebaseManagement struct {
	env fatima.FatimaEnv
}

func newCentralFilebaseManagement(env fatima.FatimaEnv) *CentralFilebaseManagement {
	instance := new(CentralFilebaseManagement)
	instance.env = env
	return instance
}

func (this *CentralFilebaseManagement) GetPSStatus() (monitor.PSStatus, bool) {
	filePath := filepath.Join(
		this.env.GetFolderGuide().GetFatimaHome(),
		"package",
		"cfm",
		"ha",
		"system.ps")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			ioutil.WriteFile(filePath, []byte(strconv.Itoa(monitor.PS_STATUS_SECONDARY)), 0644)
			return monitor.PS_STATUS_SECONDARY, true
		}
		return monitor.PS_STATUS_UNKNOWN, true
	}
	value, err1 := strconv.Atoi(strings.Trim(string(data), "\r\n"))
	if err1 != nil {
		return monitor.PS_STATUS_UNKNOWN, false
	}
	return monitor.ToPSStatus(value), true
}

func (this *CentralFilebaseManagement) GetHAStatus() (monitor.HAStatus, bool) {
	filePath := filepath.Join(
		this.env.GetFolderGuide().GetFatimaHome(),
		"package",
		"cfm",
		"ha",
		"system.ha")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			ioutil.WriteFile(filePath, []byte(strconv.Itoa(monitor.HA_STATUS_STANDBY)), 0644)
			return monitor.HA_STATUS_STANDBY, true
		}
		return monitor.HA_STATUS_UNKNOWN, true
	}
	value, err1 := strconv.Atoi(strings.Trim(string(data), "\r\n"))
	if err1 != nil {
		return monitor.HA_STATUS_UNKNOWN, false
	}
	return monitor.ToHAStatus(value), true
}

func (this *CentralFilebaseManagement) GetLogLevel() (log.LogLevel, bool) {
	filePath := filepath.Join(
		this.env.GetFolderGuide().GetFatimaHome(),
		"package",
		"cfm",
		"loglevels")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("read file fail", err)
		return log.LOG_NONE, false
	}

	//	var items []logLevelItem
	var items map[string]string
	err = json.Unmarshal(data, &items)
	if err != nil {
		log.Warn("fail to unmarshal loglevel to json : %s", err.Error())
		return log.LOG_NONE, false
	}

	value, ok := items[this.env.GetSystemProc().GetProgramName()]
	if !ok {
		// not found
		return log.LOG_NONE, false
	}

	loglevel, err1 := log.ToLogLevel(value)
	if err1 != nil {
		return log.LOG_NONE, false
	}

	return loglevel, true
}
