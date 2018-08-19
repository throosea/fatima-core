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
	"path/filepath"
	"throosea.com/fatima"
	"throosea.com/log"
	"strconv"
	"strings"
	"os"
)

type PropertyConfigReader struct {
	predefines    fatima.Predefines
	env           fatima.FatimaEnv
	configuration map[string]string
}

func NewPropertyConfigReader(env fatima.FatimaEnv, predefines fatima.Predefines) *PropertyConfigReader {
	instance := new(PropertyConfigReader)
	instance.env = env
	instance.predefines = predefines
	instance.configuration = make(map[string]string)

	propFilePath := loadApplicationProperty(env)
	if len(propFilePath) == 0 {
		log.Warn("cannot load properties file...")
		return instance
	}

	log.Info("using properties file : %s", filepath.Base(propFilePath))
	props, err := readProperties(propFilePath)
	if err != nil {
		log.Warn("cannot load properties file : %s", err.Error())
	}
	if props != nil {
		for k, v := range props {
			instance.configuration[k] = predefines.ResolvePredefine(v)
		}
	}
	return instance
}

func loadApplicationProperty(env fatima.FatimaEnv) string {
	list := make([]string, 0)
	if env.GetProfile() != "" {
		list = append(list, fmt.Sprintf("%s.%s.properties", env.GetSystemProc().GetProgramName(), env.GetProfile()))
		list = append(list, fmt.Sprintf("application.%s.properties", env.GetProfile()))
	}

	list = append(list, fmt.Sprintf("%s.properties", env.GetSystemProc().GetProgramName()))
	list = append(list, "application.properties")

	for _, v := range list {
		propFilePath := filepath.Join(env.GetFolderGuide().GetAppFolder(), v)
		if checkFileAvailable(propFilePath) {
			return propFilePath
		}
	}

	return ""
}

func loadApplicationProperties(env fatima.FatimaEnv) string {
	var filename = ""
	var propFilePath = ""
	if env.GetProfile() == "" {
		filename = fmt.Sprintf("application.properties")
	} else {
		filename = fmt.Sprintf("application.%s.properties", env.GetProfile())
	}
	propFilePath = filepath.Join(env.GetFolderGuide().GetAppFolder(), filename)
	if !checkFileAvailable(propFilePath) {
		return ""
	}
	return propFilePath
}

func loadAppNameProperties(env fatima.FatimaEnv) string {
	var filename = ""
	var propFilePath = ""
	if env.GetProfile() == "" {
		filename = fmt.Sprintf("%s.properties", env.GetSystemProc().GetProgramName())
	} else {
		filename = fmt.Sprintf("%s.%s.properties", env.GetSystemProc().GetProgramName(), env.GetProfile())
	}
	propFilePath = filepath.Join(env.GetFolderGuide().GetAppFolder(), filename)
	if !checkFileAvailable(propFilePath) {
		return ""
	}
	return propFilePath
}

func loadBasicApplicationProperties(env fatima.FatimaEnv) string {
	filename := fmt.Sprintf("application.properties")
	propFilePath := filepath.Join(env.GetFolderGuide().GetAppFolder(), filename)
	if !checkFileAvailable(propFilePath) {
		return ""
	}
	return propFilePath
}

func loadBasicAppNameProperties(env fatima.FatimaEnv) string {
	filename := fmt.Sprintf("%s.properties", env.GetSystemProc().GetProgramName())
	propFilePath := filepath.Join(env.GetFolderGuide().GetAppFolder(), filename)
	if !checkFileAvailable(propFilePath) {
		return ""
	}
	return propFilePath
}

func (this *PropertyConfigReader) GetValue(key string) (string, bool) {
	v, ok := this.configuration[key]
	return v, ok
}

func (this *PropertyConfigReader) GetString(key string) (string, error) {
	v, ok := this.configuration[key]
	if !ok {
		return "", fmt.Errorf("not found key in config : %s", key)
	}
	return v, nil
}


func (this *PropertyConfigReader) GetInt(key string) (int, error) {
	v, ok := this.configuration[key]
	if !ok {
		return 0, fmt.Errorf("not found key in config : %s", key)
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("not numeric value for key %s : %s", key, err.Error())
	}

	return i, nil
}


func (this *PropertyConfigReader) GetBool(key string) (bool, error) {
	v, ok := this.configuration[key]
	if !ok {
		return false, fmt.Errorf("not found key in config : %s", key)
	}

	switch strings.ToUpper(v) {
	case "TRUE" :
		return true, nil
	}

	return false, nil
}




func (this *PropertyConfigReader) ResolvePredefine(value string) string {
	return this.predefines.ResolvePredefine(value)
}

func (this *PropertyConfigReader) GetDefine(key string) (string, bool) {
	return this.predefines.GetDefine(key)
}


func checkFileAvailable(path string) bool	{
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Warn("file [%s] does not exist", filepath.Base(path))
		return false
	}

	return true
}

