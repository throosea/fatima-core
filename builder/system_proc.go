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
	"os"
	"os/user"
	"strconv"
	"strings"
	"throosea.com/fatima"
)

type FatimaSystemProc struct {
	pid         int
	uid         int
	gid         string
	username    string
	homeDir     string
	programName string
}

func (this *FatimaSystemProc) GetPid() int {
	return this.pid
}

func (this *FatimaSystemProc) GetUid() int {
	return this.uid
}

func (this *FatimaSystemProc) GetProgramName() string {
	return this.programName
}

func (this *FatimaSystemProc) GetUsername() string {
	return this.username
}

func (this *FatimaSystemProc) GetHomeDir() string {
	return this.homeDir
}

func (this *FatimaSystemProc) GetGid() string {
	return this.gid
}

func newSystemProc() fatima.SystemProc {
	proc := new(FatimaSystemProc)
	proc.pid = os.Getpid()
	systemUser, _ := user.Current()
	uid, _ := strconv.ParseInt(systemUser.Uid, 10, 32)
	proc.uid = int(uid)
	proc.username = systemUser.Username
	proc.homeDir = systemUser.HomeDir
	proc.gid = systemUser.Gid

	args0 := os.Args[0]
	lastIndex := strings.LastIndex(os.Args[0], "/")
	if lastIndex >= 0 {
		proc.programName = args0[lastIndex+1:]
	} else {
		proc.programName = os.Args[0]
	}

	firstIndex := strings.Index(proc.programName, " ")
	if firstIndex > 0 {
		proc.programName = proc.programName[:firstIndex]
	}

	return proc
}
