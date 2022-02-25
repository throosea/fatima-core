/*
 * //
 * // Licensed to the Apache Software Foundation (ASF) under one
 * // or more contributor license agreements.  See the NOTICE file
 * // distributed with p work for additional information
 * // regarding copyright ownership.  The ASF licenses p file
 * // to you under the Apache License, Version 2.0 (the
 * // "License"); you may not use p file except in compliance
 * // with the License.  You may obtain a copy of the License at
 * //
 * //   http://www.apache.org/licenses/LICENSE-2.0
 * //
 * // Unless required by applicable law or agreed to in writing,
 * // software distributed under the License is distributed on an
 * // "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * // KIND, either express or implied.  See the License for the
 * // specific language governing permissions and limitations
 * // under the License.
 * //
 * // @project fatima
 * // @author DeockJin Chung (jin.freestyle@gmail.com)
 * // @date 22. 1. 5. 오후 7:05
 * //
 */

package builder

import (
	"context"
	"google.golang.org/grpc"
	"throosea.com/fatima"
	proto "throosea.com/fatima/builder/fatima.message.v1"
	"throosea.com/fatima/monitor"
	"throosea.com/log"
	"time"
)

const (
	propPredefineSaturnPort = "var.saturn.port"
	valueDefaultAddress     = ":4389"
	maxQueueSize            = 4096
	dropQueueSize           = 1024 // drop if queue fulls at least half size
)

type GrpcSystemNotifyHandler struct {
	fatimaRuntime fatima.FatimaRuntime
	saturnAddress string
	conn          *grpc.ClientConn
	queue         chan []byte
}

func NewGrpcSystemNotifyHandler(fatimaRuntime fatima.FatimaRuntime) (monitor.SystemNotifyHandler, error) {
	handler := GrpcSystemNotifyHandler{fatimaRuntime: fatimaRuntime}
	handler.queue = make(chan []byte, maxQueueSize)

	handler.saturnAddress = buildServiceAddress(NewPropertyPredefineReader(fatimaRuntime.GetEnv()))

	go handler.consumeQueue()

	return &handler, nil
}

func (s *GrpcSystemNotifyHandler) consumeQueue() {
	for notifyItem := range s.queue {
		if len(notifyItem) < 3 {
			continue
		}

		req := proto.SendFatimaMessageRequest{}
		req.JsonString = string(notifyItem)

		for true {
			if s.conn == nil {
				s.connectSaturn()
			}

			if s.conn == nil {
				log.Warn("sleep for connecting to saturn...")
				time.Sleep(time.Second * 5)
				continue
			}

			ok := s.sendToSaturn(req)
			if !ok {
				break
			}

			// success
			break
		}
	}
}

func (s *GrpcSystemNotifyHandler) sendToSaturn(req proto.SendFatimaMessageRequest) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := proto.NewFatimaMessageServiceClient(s.conn).SendFatimaMessage(ctx, &req)
	if err != nil {
		log.Warn("SendFatimaMessage grpc exception : %s", err.Error())

		// maybe grpc relative errors...
		s.conn.Close()
		s.conn = nil
		return false
	}

	if errRes, ok := res.Response.(*proto.SendFatimaMessageResponse_Error); ok {
		log.Warn("SendFatimaMessage error : [%s] %s", errRes.Error.Code, errRes.Error.Desc)
	}

	return true
}

func (s *GrpcSystemNotifyHandler) connectSaturn() {
	log.Warn("connecting to saturn %s", s.saturnAddress)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	gConn, err := grpc.DialContext(
		ctx,
		s.saturnAddress,
		grpc.WithBlock(),
		grpc.WithInsecure())

	if err != nil {
		log.Warn("fail to connect saturn : %s", err.Error())
		return
	}

	s.conn = gConn
}

func buildServiceAddress(predefinedReader *PropertyPredefineReader) string {
	address, ok := predefinedReader.GetDefine(propPredefineSaturnPort)
	if !ok {
		return valueDefaultAddress
	}

	return address
}

var (
	messageDropFlag = false
)

func (s *GrpcSystemNotifyHandler) enqueueForSending(bytes []byte) {
	if len(s.queue) >= dropQueueSize {
		if !messageDropFlag {
			messageDropFlag = true
			log.Warn("notify handler drop message....")
		}
		return // DROP...
	}

	messageDropFlag = false
	s.queue <- bytes
}

func (s *GrpcSystemNotifyHandler) SendAlarm(level monitor.AlarmLevel, message string) {
	s.enqueueForSending(buildAlarmMessage(s.fatimaRuntime, level, message, ""))
}

func (s *GrpcSystemNotifyHandler) SendAlarmWithCategory(level monitor.AlarmLevel, message string, category string) {
	s.enqueueForSending(buildAlarmMessage(s.fatimaRuntime, level, message, category))
}

func (s *GrpcSystemNotifyHandler) SendEvent(message string, v ...interface{}) {
	s.enqueueForSending(buildEventMessage(s.fatimaRuntime, message, v...))
}

func (s *GrpcSystemNotifyHandler) SendActivity(json interface{}) {
	s.enqueueForSending(buildActivityMessage(s.fatimaRuntime, json))
}
