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
	"fmt"
	"net/http"
	"throosea.com/fatima"
	"throosea.com/fatima/builder"
	"throosea.com/fatima/lib"
	"throosea.com/fatima/monitor"
	"throosea.com/log"
	"time"
)

type ProcessCoreWorker interface {
	Process()
}

var oneSecondTickWorkers []ProcessCoreWorker
var fiveSecondTickWorkers []ProcessCoreWorker

func init() {
	oneSecondTickWorkers = make([]ProcessCoreWorker, 0)
	fiveSecondTickWorkers = make([]ProcessCoreWorker, 0)
}

type DefaultProcessInteractor struct {
	runtimeProcess *builder.FatimaRuntimeProcess
	awareManager   *SystemAwareManagement
	monitor        monitor.SystemStatusMonitor
	measurement    *SystemMeasureManagement
	readers        []fatima.FatimaIOReader
}

func NewProcessInteractor(runtimeProcess *builder.FatimaRuntimeProcess) *DefaultProcessInteractor {
	instance := new(DefaultProcessInteractor)
	instance.runtimeProcess = runtimeProcess
	// monitor : Active/Standby, Primary/Secondary
	instance.monitor = newCentralFilebaseManagement(runtimeProcess.GetEnv())
	instance.awareManager = newSystemAwareManagement(runtimeProcess, instance.monitor)
	// special type of FatimaComponent. usually we need create 'Reader' type first
	instance.readers = make([]fatima.FatimaIOReader, 0)
	instance.measurement = newSystemMeasureManagement(runtimeProcess)

	oneSecondTickWorkers = append(oneSecondTickWorkers, instance.awareManager)
	fiveSecondTickWorkers = append(fiveSecondTickWorkers, instance.measurement)

	startTickers()
	return instance
}

func (this *DefaultProcessInteractor) Regist(component fatima.FatimaComponent) {
	registComponent(component)

	if comp, ok := component.(monitor.FatimaSystemHAAware); ok {
		this.RegistSystemHAAware(comp)
	}

	if comp, ok := component.(monitor.FatimaSystemPSAware); ok {
		this.RegistSystemPSAware(comp)
	}

	if comp, ok := component.(fatima.FatimaIOReader); ok {
		this.readers = append(this.readers, comp)
	}
}

func (this *DefaultProcessInteractor) RegistSystemHAAware(aware monitor.FatimaSystemHAAware) {
	this.awareManager.RegistSystemHAAware(aware)
}

func (this *DefaultProcessInteractor) RegistSystemPSAware(aware monitor.FatimaSystemPSAware) {
	this.awareManager.RegistSystemPSAware(aware)
}

func (this *DefaultProcessInteractor) Initialize() bool {
	return initializeComponent()
}

func (this *DefaultProcessInteractor) Goaway() {
	goawayComponent()
}

func (this *DefaultProcessInteractor) startListening() {
	for _, v := range this.readers {
		t := v
		go func() {
			t.StartListening()
		}()
	}
}

func (this *DefaultProcessInteractor) Run() {
	this.startListening()
	lib.StartCron()
	bootupNotify()
	this.pprofService()
	if this.runtimeProcess.GetBuilder().GetProcessType() == fatima.PROCESS_TYPE_GENERAL {
		message := fmt.Sprintf("%s process started", this.runtimeProcess.GetEnv().GetSystemProc().GetProgramName())
		this.runtimeProcess.GetSystemNotifyHandler().SendAlarm(monitor.AlarmLevelMinor, message)
	}
}

func (this *DefaultProcessInteractor) Stop() {

}

func (this *DefaultProcessInteractor) Shutdown() {
	if this.runtimeProcess.GetBuilder().GetProcessType() == fatima.PROCESS_TYPE_GENERAL {
		message := fmt.Sprintf("%s process shutdowned", this.runtimeProcess.GetEnv().GetSystemProc().GetProgramName())
		this.runtimeProcess.GetSystemNotifyHandler().SendAlarm(monitor.AlamLevelMajor, message)
	}
	lib.StopCron()
	shutdownComponent(this.runtimeProcess.GetEnv().GetSystemProc().GetProgramName())
}

func (this *DefaultProcessInteractor) RegistMeasureUnit(unit monitor.SystemMeasurable) {
	this.measurement.registUnit(unit)
}

func (this *DefaultProcessInteractor) pprofService() {
	addr, ok := this.runtimeProcess.GetConfig().GetValue(builder.GOFATIMA_PROP_PPROF_ADDRESS)
	if ok {
		go func() {
			err := http.ListenAndServe(addr, http.DefaultServeMux)
			if err != nil {
				log.Warn("fail to start pprof service : {}", err.Error())
			}
		}()
		log.Info("starting pprof service. address=%s", addr)
	}
}

func startTickers() {
	oneSecondTick := time.NewTicker(time.Second * 1)
	go func() {
		for range oneSecondTick.C {
			iterateWorkers(oneSecondTickWorkers)
		}
	}()
	fiveSecondTick := time.NewTicker(time.Second * 5)
	go func() {
		for range fiveSecondTick.C {
			iterateWorkers(fiveSecondTickWorkers)
		}
	}()
}

func iterateWorkers(workers []ProcessCoreWorker) {
	if workers == nil || len(workers) < 1 {
		return
	}

	for _, v := range workers {
		v.Process()
	}
}
