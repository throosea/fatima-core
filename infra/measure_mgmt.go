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
	"throosea.com/fatima/monitor"
	"time"
)

func newSystemMeasureManagement(runtimeProcess *builder.FatimaRuntimeProcess) *SystemMeasureManagement {
	mgmt := new(SystemMeasureManagement)
	mgmt.runtimeProcess = runtimeProcess
	mgmt.writer = newMeasureFileWriter(runtimeProcess.GetEnv())
	mgmt.units = make([]monitor.SystemMeasurable, 0)
	mgmt.registUnit(newProcessMeasurement())
	return mgmt
}

type SystemMeasureManagement struct {
	runtimeProcess *builder.FatimaRuntimeProcess
	units          []monitor.SystemMeasurable
	writer         MeasurementWriter
}

func (this *SystemMeasureManagement) registUnit(unit monitor.SystemMeasurable) {
	this.units = append(this.units, unit)
}

// 5초에 1번씩 호출된다
// 1분에 1번씩은 saturn으로 보내자
// saturn 입장에서 7일동안 전송이 없는 프로세스는 파일을 제거하기 때문에
// 주기적으로 메시지를 보내주어야 할 필요가 있다
// 이 부분은 추후 개선되어야 한다
var measureTick uint64 = 0
func (this *SystemMeasureManagement) Process() {
	msr := measurement{eventTime: time.Now()}
	msr.items = make([]measureItem, 0)
	for _, v := range this.units {
		msr.items = append(msr.items, measureItem{v.GetKeyName(), v.GetMeasure()})
	}
	this.writer.write(msr)

	measureTick += 1
	if measureTick % 12 == 0 {
		activity := make(map[string]string)
		for _,v := range msr.items {
			activity[v.keyName] = v.value
		}
		this.runtimeProcess.GetSystemNotifyHandler().SendActivity(activity)
	}
}

type measurement struct {
	eventTime time.Time
	items     []measureItem
}

type measureItem struct {
	keyName string
	value   string
}
