//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with p work for additional information
// regarding copyright ownership.  The ASF licenses p file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use p file except in compliance
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

package builder

import (
	"throosea.com/fatima"
	"throosea.com/log"
	"throosea.com/fatima/monitor"
	"path/filepath"
	"throosea.com/fatima/lib/mbus"
	"encoding/json"
	"throosea.com/fatima/lib"
	"time"
	"fmt"
)

const (
	APPLICATION_CODE = 0x1
	LOGIC_MEASURE = 10
	LOGIC_NOTIFY = 20
)

type DefaultSystemNotifyHandler struct {
	fatimaRuntime	fatima.FatimaRuntime
	mbus			*mbus.MappedMBus
}

func NewDefaultSystemNotifyHandler(fatimaRuntime fatima.FatimaRuntime) (monitor.SystemNotifyHandler, error) {
	handler := DefaultSystemNotifyHandler{fatimaRuntime:fatimaRuntime}

	dataDir := filepath.Join(fatimaRuntime.GetEnv().GetFolderGuide().GetFatimaHome(), FatimaFolderData)
	dest := "saturn"	// fatima monitoring process name is saturn
	proc := fatimaRuntime.GetEnv().GetSystemProc().GetProgramName()

	mbus, err := mbus.NewMappedMBus(dataDir, dest, proc)
	if err != nil {
		return nil, err
	}
	handler.mbus = mbus
	return &handler, nil
}

func (s *DefaultSystemNotifyHandler) SendAlarm(message string) {
	s.mbus.Write(buildAlarmMessage(s.fatimaRuntime, message))
}

func (s *DefaultSystemNotifyHandler) SendEvent(message string, v ...interface{}) {
	s.mbus.Write(buildEventMessage(s.fatimaRuntime, message, v...))
}

func (s *DefaultSystemNotifyHandler) SendActivity(json interface{}) {
	s.mbus.Write(buildActivityMessage(s.fatimaRuntime, json))
}

func buildAlarmMessage(fatimaRuntime fatima.FatimaRuntime, message string) []byte {
	m := make(map[string]interface{})
	header := make(map[string]interface{})
	body := make(map[string]interface{})


	header["application_code"] = APPLICATION_CODE
	header["logic"] = LOGIC_NOTIFY

	body["package_host"] = fatimaRuntime.GetPackaging().GetHost()
	body["package_name"] = fatimaRuntime.GetPackaging().GetName()
	body["package_group"] = fatimaRuntime.GetPackaging().GetGroup()
	body["package_profile"] = fatimaRuntime.GetEnv().GetProfile()
	body["package_process"] = fatimaRuntime.GetEnv().GetSystemProc().GetProgramName()
	body["event_time"] = lib.CurrentTimeMillis()

	alarm := make(map[string]interface{})
	alarm["type"] = "ALARM"
	alarm["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	alarm["alarm_level"] = "WARN"
	alarm["from"] = "go-fatima"
	alarm["initiator"] = "go-fatima"
	alarm["message"] = message

	body["message"] = alarm

	m["header"] = header
	m["body"] = body

	b, err := json.Marshal(m)
	if err != nil {
		log.Warn("fail to make alarm json : %s", err.Error())
		return nil
	}

	return b
}


func buildEventMessage(fatimaRuntime fatima.FatimaRuntime, message string, v ...interface{}) []byte {
	m := make(map[string]interface{})
	header := make(map[string]interface{})
	body := make(map[string]interface{})


	header["application_code"] = APPLICATION_CODE
	header["logic"] = LOGIC_NOTIFY

	body["package_host"] = fatimaRuntime.GetPackaging().GetHost()
	body["package_name"] = fatimaRuntime.GetPackaging().GetName()
	body["package_group"] = fatimaRuntime.GetPackaging().GetGroup()
	body["package_profile"] = fatimaRuntime.GetEnv().GetProfile()
	body["package_process"] = fatimaRuntime.GetEnv().GetSystemProc().GetProgramName()
	body["event_time"] = lib.CurrentTimeMillis()

	alarm := make(map[string]interface{})
	alarm["type"] = "EVENT"
	alarm["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	alarm["from"] = "go-fatima"
	alarm["initiator"] = "go-fatima"
	alarm["message"] = message

	if len(v) > 0 {
		args := make([]string, 0)
		for _, a := range v {
			if e, ok := a.(fmt.Stringer); ok {
				args = append(args, e.String())
			} else if e, ok := a.(string); ok {
				args = append(args, e)
			} else if e, ok := a.(int); ok {
				args = append(args, fmt.Sprintf("%d", e))
			} else if e, ok := a.(float32); ok {
				args = append(args, fmt.Sprintf("%f", e))
			} else if e, ok := a.(float64); ok {
				args = append(args, fmt.Sprintf("%f", e))
			} else if e, ok := a.(int32); ok {
				args = append(args, fmt.Sprintf("%d", e))
			} else if e, ok := a.(uint32); ok {
				args = append(args, fmt.Sprintf("%d", e))
			} else if e, ok := a.(uint64); ok {
				args = append(args, fmt.Sprintf("%d", e))
			} else if e, ok := a.(bool); ok {
				args = append(args, fmt.Sprintf("%b", e))
			} else {
				args = append(args, ".")
			}
		}
		alarm["params"] = args
	}

	body["message"] = alarm

	m["header"] = header
	m["body"] = body

	b, err := json.Marshal(m)
	if err != nil {
		log.Warn("fail to make event json : %s", err.Error())
		return nil
	}

	return b
}


func buildActivityMessage(fatimaRuntime fatima.FatimaRuntime, v interface{}) []byte {
	m := make(map[string]interface{})
	header := make(map[string]interface{})
	body := make(map[string]interface{})


	header["application_code"] = APPLICATION_CODE
	header["logic"] = LOGIC_MEASURE

	body["package_host"] = fatimaRuntime.GetPackaging().GetHost()
	body["package_name"] = fatimaRuntime.GetPackaging().GetName()
	body["package_group"] = fatimaRuntime.GetPackaging().GetGroup()
	body["package_profile"] = fatimaRuntime.GetEnv().GetProfile()
	body["package_process"] = fatimaRuntime.GetEnv().GetSystemProc().GetProgramName()
	body["event_time"] = lib.CurrentTimeMillis()
	body["message"] = v

	m["header"] = header
	m["body"] = body

	b, err := json.Marshal(m)
	if err != nil {
		log.Warn("fail to make alarm json : %s", err.Error())
		return nil
	}

	return b
}

