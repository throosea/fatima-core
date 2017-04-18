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
	"throosea.com/fatima"
	"throosea.com/fatima/lib"
	"throosea.com/fatima/log"
	"errors"
	"fmt"
)

var compPreInit []fatima.FatimaComponent
var compGeneral []fatima.FatimaComponent
var compReader []fatima.FatimaComponent
var compWriter []fatima.FatimaComponent

func init() {
	compPreInit = make([]fatima.FatimaComponent, 0)
	compGeneral = make([]fatima.FatimaComponent, 0)
	compReader = make([]fatima.FatimaComponent, 0)
	compWriter = make([]fatima.FatimaComponent, 0)
}

func registComponent(comp fatima.FatimaComponent) {
	switch comp.GetType() {
	case fatima.COMP_PRE_INIT:
		compPreInit = append(compPreInit, comp)
	case fatima.COMP_GENERAL:
		compGeneral = append(compGeneral, comp)
	case fatima.COMP_READER:
		compReader = append(compReader, comp)
	case fatima.COMP_WRITER:
		compWriter = append(compWriter, comp)
	}
}

func initializeComponent() (res bool) {
	res = false

	defer func() {
		if r := recover(); r != nil {
			log.Warn("**PANIC** while initializing", errors.New(fmt.Sprintf("%s", r)))
			return
		}
		res = true
	}()

	if !callInitial(compPreInit) {
		return
	}
	if !callInitial(compWriter) {
		return
	}
	if !callInitial(compReader) {
		return
	}
	if !callInitial(compGeneral) {
		return
	}

	res = true
	return
}

func callInitial(list []fatima.FatimaComponent) bool {
	for _, v := range list {
		if !v.Initialize() {
			return false
		}
	}
	return true
}

func bootupNotify() {
	all := make([]fatima.FatimaComponent, 0)
	all = append(all, compPreInit...)
	all = append(all, compGeneral...)
	all = append(all, compReader...)
	all = append(all, compWriter...)

	size := len(all)
	if size > 0 {
		cyBarrier := lib.NewCyclicBarrier(size, func() { log.Info("모든 FatimaComponent에게 Boot 메시지를 전달하였습니다") })
		for _, v := range all {
			t := v
			cyBarrier.Dispatch(func() { t.Bootup() })
		}
		cyBarrier.Wait()
	} else {
		log.Info("모든 FatimaComponent에게 Boot 메시지를 전달하였습니다")
	}
}

func shutdownComponent(program string) {
	log.Info("FatimaComponent들을 shutdown 합니다")
	all := make([]fatima.FatimaComponent, 0)
	all = append(all, compPreInit...)
	all = append(all, compGeneral...)
	all = append(all, compReader...)
	all = append(all, compWriter...)

	defer func() {
		if r := recover(); r != nil {
			log.Warn("**PANIC** while shutdown", errors.New(fmt.Sprintf("%s", r)))
			log.WaitLoggingShutdown()
			return
		}
	}()


	size := len(all)
	if size > 0 {
		cyBarrier := lib.NewCyclicBarrier(size, func() {
			log.Warn("%s 프로그램을 종료합니다", program)
		})
		for _, v := range all {
			t := v
			cyBarrier.Dispatch(func() { t.Shutdown() })
		}
		cyBarrier.Wait()
	} else {
		log.Warn("%s 프로그램을 종료합니다", program)
	}

	log.WaitLoggingShutdown()
}
