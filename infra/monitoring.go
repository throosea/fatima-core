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
	"throosea.com/fatima/builder"
	"throosea.com/log"
	"throosea.com/fatima/monitor"
)

type SystemAwareManagement struct {
	runtimeProcess *builder.FatimaRuntimeProcess
	monitor        monitor.SystemStatusMonitor
	awareHA        []monitor.FatimaSystemHAAware
	awarePS        []monitor.FatimaSystemPSAware
}

func newSystemAwareManagement(runtimeProcess *builder.FatimaRuntimeProcess, mon monitor.SystemStatusMonitor) *SystemAwareManagement {
	instance := new(SystemAwareManagement)

	instance.runtimeProcess = runtimeProcess
	instance.awareHA = make([]monitor.FatimaSystemHAAware, 0)
	instance.awarePS = make([]monitor.FatimaSystemPSAware, 0)
	instance.monitor = mon
	currentStatus := runtimeProcess.GetSystemStatus().(*builder.FatimaPackageSystemStatus)

	ps, _ := mon.GetPSStatus()
	currentStatus.SetPSStatus(ps)
	ha, _ := mon.GetHAStatus()
	currentStatus.SetHAStatus(ha)

	return instance
}

func (this *SystemAwareManagement) RegistSystemHAAware(aware monitor.FatimaSystemHAAware) {
	this.awareHA = append(this.awareHA, aware)
}

func (this *SystemAwareManagement) RegistSystemPSAware(aware monitor.FatimaSystemPSAware) {
	this.awarePS = append(this.awarePS, aware)
}

func (this *SystemAwareManagement) SystemHAStatusChanged(newHAStatus monitor.HAStatus) {
	log.Warn("new HA Status detected : %s", newHAStatus)
	for _, aware := range this.awareHA {
		aware.SystemHAStatusChanged(newHAStatus)
	}
}

func (this *SystemAwareManagement) SystemPSStatusChanged(newPSStatus monitor.PSStatus) {
	log.Warn("new PS Status detected : %s", newPSStatus)
	for _, aware := range this.awarePS {
		aware.SystemPSStatusChanged(newPSStatus)
	}
}

func (this *SystemAwareManagement) Process() {
	currentStatus := this.runtimeProcess.GetSystemStatus().(*builder.FatimaPackageSystemStatus)

	if ps, ok := this.monitor.GetPSStatus(); ok {
		oldps := currentStatus.GetPSStatus()
		if oldps != ps {
			currentStatus.SetPSStatus(ps)
			go func() {
				this.SystemPSStatusChanged(ps)
			}()
		}
	}
	if ha, ok := this.monitor.GetHAStatus(); ok {
		oldha := currentStatus.GetHAStatus()
		if oldha != ha {
			currentStatus.SetHAStatus(ha)
			go func() {
				this.SystemHAStatusChanged(ha)
			}()
		}
	}

	if logLevel, ok := this.monitor.GetLogLevel(); ok {
		if this.runtimeProcess.GetLogLevel() != logLevel {
			log.SetLevel(logLevel)
			log.Warn("fatima proc log level : %s", logLevel)
			this.runtimeProcess.SetLogLevel(logLevel)
		}
	}
}
