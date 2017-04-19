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

	var filename = ""
	if env.GetProfile() == "" {
		filename = fmt.Sprintf("%s.properties", env.GetSystemProc().GetProgramName())
	} else {
		filename = fmt.Sprintf("%s.%s.properties", env.GetSystemProc().GetProgramName(), env.GetProfile())
	}
	log.Info("using properties file : %s", filename)
	props, err := readProperties(filepath.Join(env.GetFolderGuide().GetAppFolder(), filename))
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

func (this *PropertyConfigReader) GetValue(key string) (string, bool) {
	v, ok := this.configuration[key]
	return v, ok
}
