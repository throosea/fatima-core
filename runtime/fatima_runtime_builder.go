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

package runtime

import (
	"throosea.com/fatima"
	"throosea.com/fatima/builder"
	"throosea.com/fatima/monitor"
)

type DefaultProcessBuilder struct {
	pkgProcConfig fatima.FatimaPkgProcConfig
	predefines    fatima.Predefines
	config        fatima.Config
	monitor       monitor.SystemStatusMonitor
	systemAware   monitor.FatimaSystemAware
	processType	  fatima.FatimaProcessType
}

func (this *DefaultProcessBuilder) GetPkgProcConfig() fatima.FatimaPkgProcConfig {
	return this.pkgProcConfig
}

func (this *DefaultProcessBuilder) GetPredefines() fatima.Predefines {
	return this.predefines
}

func (this *DefaultProcessBuilder) GetConfig() fatima.Config {
	return this.config
}

func (this *DefaultProcessBuilder) GetProcessType() fatima.FatimaProcessType {
	return this.processType
}

func (this *DefaultProcessBuilder) GetSystemStatusMonitor() monitor.SystemStatusMonitor {
	return this.monitor
}

func (this *DefaultProcessBuilder) GetSystemAware() monitor.FatimaSystemAware {
	return this.systemAware
}

func getRuntimeBuilder(env fatima.FatimaEnv, processType fatima.FatimaProcessType) builder.FatimaRuntimeBuilder {
	processBuilder := new(DefaultProcessBuilder)
	processBuilder.processType = processType
	if processType == fatima.PROCESS_TYPE_GENERAL {
		processBuilder.pkgProcConfig = builder.NewYamlFatimaPackageConfig(env)
	} else {
		// USER INTERACTIVE
		processBuilder.pkgProcConfig = builder.NewDummyFatimaPackageConfig(env)
	}
	processBuilder.predefines = builder.NewPropertyPredefineReader(env)
	processBuilder.config = builder.NewPropertyConfigReader(env, processBuilder.predefines)

	return processBuilder
}
