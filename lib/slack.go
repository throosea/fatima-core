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
// @author 1100282
// @date 2017. 5. 2. PM 3:03
//

package lib

import (
	"throosea.com/fatima"
	"time"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
	"sync"
	"bytes"
	"throosea.com/log"
	"net/http"
)

const (
	fileWebhookSlack = "webhook.slack"
	attachmentsColor = "#00FF00"
	userName = "FATIMA"
	footerIcon = "https://platform.slack-edge.com/img/default_application_icon.png"
	applicationJsonUtf8Value = "application/json;charset=UTF-8"
)

func NewSlackNotification(fatimaRuntime fatima.FatimaRuntime) *SlackNotification {
	return NewSlackNotificationWithKey(fatimaRuntime, "default")
}

func NewSlackNotificationWithKey(fatimaRuntime fatima.FatimaRuntime, key string) *SlackNotification {
	slack := SlackNotification{}
	slack.fatimaRuntime = fatimaRuntime
	slack.key = key
	slack.mutex = &sync.Mutex{}
	slack.config.Active = false
	return &slack
}

type SlackNotification struct {
	fatimaRuntime	fatima.FatimaRuntime
	key				string
	lastLoadingTime	time.Time
	config			SlackConfig
	mutex			*sync.Mutex
}

type SlackConfig struct {
	Active		bool
	Alarm		bool
	Event		bool
	Url			string
}

func (s *SlackNotification) loading()	{
	s.lastLoadingTime = time.Now()
	webhookConfigFile := filepath.Join(s.fatimaRuntime.GetEnv().GetFolderGuide().GetDataFolder(), fileWebhookSlack)
	dataBytes, err := ioutil.ReadFile(webhookConfigFile)
	if err != nil {
		return
	}

	var data map[string]SlackConfig
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return
	}

	c, ok := data[s.key]
	if !ok {
		return
	}

	s.config.Active = c.Active
	s.config.Alarm = c.Alarm
	s.config.Event = c.Event
	s.config.Url = c.Url

	log.Debug("slack config loaded : %v", s.config)
}

func (s *SlackNotification) isEventWritable() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	deadline := time.Now().Add(-time.Minute)
	if s.lastLoadingTime.Before(deadline) {
		s.loading()
	}
	if !s.config.Active || !s.config.Event || len(s.config.Url) < 6 {
		return false
	}
	return true
}

func (s *SlackNotification) SendEvent(message string) {
	if !s.isEventWritable() {
		return
	}

	m := make(map[string]interface{})
	m["username"] = userName
	list := make([]interface{}, 0)
	list = append(list, buildAttachment(s.fatimaRuntime, message))
	m["attachments"] = list

	b, err := json.Marshal(m)
	if err != nil {
		log.Warn("fail to build json : %s", err.Error())
		return
	}

	resp, err := http.Post("http://example.com/upload", applicationJsonUtf8Value, bytes.NewReader(b))
	if err != nil {
		log.Warn("fail to send slack notification : %s", err.Error())
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Debug("successfully send to slack : %s", message)
	} else {
		log.Info("slack response : %s", resp.Status)
	}

}

func buildAttachment(fatimaRuntime fatima.FatimaRuntime, message string) map[string]interface{} {
	m := make(map[string]interface{})
	m["pretext"] = buildPretext(fatimaRuntime)
	m["color"] = attachmentsColor
	m["text"] = message
	m["footer"] = fatimaRuntime.GetEnv().GetSystemProc().GetProgramName()
	m["footer_icon"] = footerIcon
	m["ts"] = CurrentTimeMillis() / 1000
	return m
}

func buildPretext(fatimaRuntime fatima.FatimaRuntime) string {
	var buff bytes.Buffer
	if len(fatimaRuntime.GetEnv().GetProfile()) > 0 {
		buff.WriteByte('[')
		buff.WriteString(fatimaRuntime.GetEnv().GetProfile())
		buff.WriteByte(']')
		buff.WriteByte(' ')
	}
	buff.WriteString(fatimaRuntime.GetPackaging().GetGroup())
	buff.WriteByte(':')
	buff.WriteString(fatimaRuntime.GetPackaging().GetHost())
	if fatimaRuntime.GetPackaging().GetName() != "default" {
		buff.WriteByte(':')
		buff.WriteString(fatimaRuntime.GetPackaging().GetName())
	}
	return buff.String()
}
