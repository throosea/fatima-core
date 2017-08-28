//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with reader work for additional information
// regarding copyright ownership.  The ASF licenses reader file
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
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"throosea.com/fatima"
	"throosea.com/log"
	"bytes"
)

type ProcessItem struct {
	Gid      	int    `yaml:"gid"`
	Name     	string `yaml:"name"`
	Loglevel 	string `yaml:"loglevel"`
	Hb       	bool   `yaml:"hb,omitempty"`
	Path     	string `yaml:"path,omitempty"`
	Grep     	string `yaml:"grep,omitempty"`
	Startmode	int		`yaml:"startmode,omitempty"`
}

func (p ProcessItem) GetGid() int {
	return p.Gid
}

func (p ProcessItem) GetName() string {
	return p.Name
}

func (p ProcessItem) GetHeartbeat() bool {
	return p.Hb
}

func (p ProcessItem) GetPath() string {
	return p.Path
}

func (p ProcessItem) GetGrep() string {
	return p.Grep
}

func (p ProcessItem) GetStartMode() fatima.ProcessStartMode {
	switch p.Startmode  {
	case 1 :		return fatima.StartModeAlone
	case 2 :		return fatima.StartModeByHA
	case 3 :		return fatima.StartModeByPS
	default : 		return fatima.StartModeByJuno
	}
}

func (p ProcessItem) GetLogLevel() log.LogLevel {
	return buildLogLevel(p.Loglevel)
}

type GroupItem struct {
	Id   int    `yaml:"id"`
	Name string `yaml:"name"`
}

type YamlFatimaPackageConfig struct {
	env       fatima.FatimaEnv
	predefines fatima.Predefines
	Groups    []GroupItem   `yaml:"group,flow"`
	Processes []ProcessItem `yaml:"process"`
}

func NewYamlFatimaPackageConfig(env fatima.FatimaEnv) *YamlFatimaPackageConfig {
	instance := new(YamlFatimaPackageConfig)
	instance.env = env
	instance.Reload()
	return instance
}

func (y *YamlFatimaPackageConfig) Save() {
	d, err := yaml.Marshal(y)
	if err != nil {
		log.Warn("fail to create yaml data : %s", err.Error())
		return
	}

	var comment = "---\n" +
		"# this is fatima-package.yaml sample\n" +
		"# group (define column)\n" +
		"# process list (define column)\n" +
		"#  gid, name, path, qclear, qkey, hb\n" +
		"# non-fatima process\n" +
		"# startmode : 0(always started by juno), 1(not started by juno), 2(by HA), 3(by PS)\n"

	var buff bytes.Buffer
	buff.WriteString(comment)
	buff.Write(d)
	err = ioutil.WriteFile(y.env.GetFolderGuide().GetPackageProcFile(), buff.Bytes(), 0644)
	if err != nil {
		log.Warn("fail to save yaml configuration file : %s", err.Error())
		return
	}
}

func (y *YamlFatimaPackageConfig) Reload() {
	data, err := ioutil.ReadFile(y.env.GetFolderGuide().GetPackageProcFile())
	check(err)

	err = yaml.Unmarshal(data, &y)
	check(err)

	if len(y.Groups) == 0 || len(y.Processes) == 0 {
		panic(fmt.Errorf("invalid fatima yaml configuration : %s", y.env.GetFolderGuide().GetPackageProcFile()))
	}
}

func (y *YamlFatimaPackageConfig) GetProcByName(name string) fatima.FatimaPkgProc {
	for _, each := range y.Processes {
		if each.Name == name {
			return each
		}
	}

	return nil
}

func (y *YamlFatimaPackageConfig) GetProcByGroup(name string) []fatima.FatimaPkgProc {
	procList := make([]fatima.FatimaPkgProc, 0)
	gid := y.GetGroupId(name)
	if gid < 0 {
		return procList
	}

	for _, each := range y.Processes {
		if each.Gid == gid {
			procList = append(procList, each)
		}
	}

	return procList
}

func (y *YamlFatimaPackageConfig) GetAllProc(exceptOpmGroup bool) []fatima.FatimaPkgProc {
	procList := make([]fatima.FatimaPkgProc, 0)
	if !exceptOpmGroup {
		for _, v := range y.Processes {
			procList = append(procList, v)
		}
		return procList
	}

	gid := y.GetGroupId("OPM")
	for _, each := range y.Processes {
		if each.Gid != gid {
			procList = append(procList, each)
		}
	}

	return procList
}

func (y *YamlFatimaPackageConfig) GetGroupId(groupName string) int {
	comp := strings.ToLower(groupName)
	for _, each := range y.Groups {
		if comp == strings.ToLower(each.Name) {
			return each.Id
		}
	}

	return -1
}

func (y *YamlFatimaPackageConfig) IsValidGroupId(groupId int) bool {
	for _, each := range y.Groups {
		if each.Id == groupId {
			return true
		}
	}

	return false
}

func buildLogLevel(s string) log.LogLevel {
	switch strings.ToLower(s) {
	case "trace":
		return log.LOG_TRACE
	case "debug":
		return log.LOG_DEBUG
	case "info":
		return log.LOG_INFO
	case "warn":
		return log.LOG_WARN
	case "error":
		return log.LOG_ERROR
	case "none":
		return log.LOG_NONE
	}
	return log.LOG_TRACE
}



type DummyFatimaPackageConfig struct {
	env       fatima.FatimaEnv
	predefines fatima.Predefines
	Groups    []GroupItem   `yaml:"group,flow"`
	Processes []ProcessItem `yaml:"process"`
}

func NewDummyFatimaPackageConfig(env fatima.FatimaEnv) *DummyFatimaPackageConfig {
	instance := new(DummyFatimaPackageConfig)
	instance.env = env
	instance.Reload()
	return instance
}

func (y *DummyFatimaPackageConfig) Reload() {
}

func (y *DummyFatimaPackageConfig) GetProcByName(name string) fatima.FatimaPkgProc {
	item := ProcessItem{}
	item.Name = name
	item.Startmode = fatima.StartModeAlone
	item.Path = "/"
	item.Gid = 0
	item.Hb = false
	item.Loglevel = "debug"
	return item
}

