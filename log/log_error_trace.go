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

package log

import (
	"bytes"
	//	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func newErrorTraceLogMessage(pc uintptr, file string, line int) *ErrorTraceLogMessage {
	message := ErrorTraceLogMessage{}
	message.t = time.Now()
	message.announce = false
	message.pc = pc
	message.file = file
	message.line = line
	message.tracePoint = make([]TracePoint, 0)
	if loggingPreference.showMethod {
		message.funcName = findFunctionName(pc)
	}
	return &message
}

type TracePoint struct {
	pc   uintptr
	file string
	line int
}

type ErrorTraceLogMessage struct {
	GeneralLogMessage
	announce   bool
	tracePoint []TracePoint
}

func (this *ErrorTraceLogMessage) append(point TracePoint) {
	this.tracePoint = append(this.tracePoint, point)
}

func (this *ErrorTraceLogMessage) publish() {
	var buffer bytes.Buffer

	codeLine := this.buildMessage(func() string {
		size := len(this.message)
		if size == 1 {
			this.announce = true
			return fmt.Sprintf("(%s) :: %s", reflect.TypeOf(this.message[0]).String(), this.message[0])
		} else {
			if format, ok := this.message[0].(string); ok {
				if size == 2 {
					return format
				} else {
					return fmt.Sprintf(format, this.message[1:size-1]...)
				}
			}
			this.announce = true
			return fmt.Sprintf("(%s) :: %s", reflect.TypeOf(this.message[size-1]).String(), this.message[size-1])
		}
	})

	buffer.WriteString(codeLine)
	buffer.WriteString(this.getTrace())
	this.published = buffer.String()
}

func (this *ErrorTraceLogMessage) getTrace() string {
	/*
	   TRACE <<<
	         [inject(), org.springframework.beans.factory.annotation.AutowiredAnnotationBeanPostProcessor$AutowiredFieldElement:569]
	         [inject(), org.springframework.beans.factory.annotation.InjectionMetadata:88]
	         [postProcessPropertyValues(), org.springframework.beans.factory.annotation.AutowiredAnnotationBeanPostProcessor:349]
	         [populateBean(), org.springframework.beans.factory.support.AbstractAutowireCapableBeanFactory:1214]
	         [doCreateBean(), org.springframework.beans.factory.support.AbstractAutowireCapableBeanFactory:543]
	         [createBean(), org.springframework.beans.factory.support.AbstractAutowireCapableBeanFactory:482]
	         [getObject(), org.springframework.beans.factory.support.AbstractBeanFactory$1:306]
	         [getSingleton(), org.springframework.beans.factory.support.DefaultSingletonBeanRegistry:230]
	         [doGetBean(), org.springframework.beans.factory.support.AbstractBeanFactory:302]
	         [getBean(), org.springframework.beans.factory.support.AbstractBeanFactory:197]
	         [preInstantiateSingletons(), org.springframework.beans.factory.support.DefaultListableBeanFactory:776]
	         ......
	         [main(), org.fatima.core.process.launch.DefaultApplicationLauncher:36]
	*/
	var buffer bytes.Buffer

	if this.announce {
		buffer.WriteString("\tTRACE <<<\n")
	} else {
		err := this.message[len(this.message)-1]
		buffer.WriteString(fmt.Sprintf("\t(%s) :: %s\n\tTRACE <<<\n", reflect.TypeOf(err).String(), err))
	}
	for _, v := range this.tracePoint {
		buffer.WriteString(fmt.Sprintf("\t[%s(), %s:%d]\n", findFunctionName(v.pc), buildSourcePath(v.file), v.line))
	}
	return buffer.String()

}

func buildSourcePath(file string) string {
	var location = file
	var found = strings.LastIndex(file, "/src/")
	if found > 0 {
		location = string(file[found+5:])
	}

	return strings.Replace(location, "/", ".", -1)
}
