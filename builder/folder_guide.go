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
	"os"
	"path/filepath"
	"throosea.com/fatima"
	"throosea.com/fatima/lib"
)

const (
	FatimaFolderApp      = "app"
	FatimaFolderBin      = "bin"
	FatimaFolderConf     = "conf"
	FatimaFolderData     = "data"
	FatimaFolderJavalib  = "javalib"
	FatimaFolderLib      = "lib"
	FatimaFolderLog      = "log"
	FatimaFolderPackage  = "package"
	FatimaFolderStat     = "stat"
	FatimaFolderProc     = "proc"
	FatimaFileProcConfig = "fatima-package.yaml"
)

type FatimaFolderGuide struct {
	fatimaHomePath string
	app            string
	bin            string
	conf           string
	data           string
	javalib        string
	lib            string
	log            string
	pack           string
	stat           string
	proc           string
}

func (this *FatimaFolderGuide) GetFatimaHome() string {
	return this.fatimaHomePath
}

func (this *FatimaFolderGuide) GetPackageProcFile() string {
	return fmt.Sprintf("%s%c%s", this.conf, os.PathSeparator, FatimaFileProcConfig)
}

func (this *FatimaFolderGuide) GetAppProcFolder() string {
	return this.proc
}

func (this *FatimaFolderGuide) GetLogFolder() string {
	return this.log
}

func (this *FatimaFolderGuide) GetConfFolder() string {
	return this.conf
}

func (this *FatimaFolderGuide) GetDataFolder() string {
	return this.data
}

func (this *FatimaFolderGuide) GetAppFolder() string {
	return this.app
}

func (this *FatimaFolderGuide) CreateTmpFolder() string {
	seed := lib.RandomAlphanumeric(16)
	tmp := filepath.Join(this.data, ".tmp", seed)
	checkDirectory(tmp, true)
	return tmp
}

func (this *FatimaFolderGuide) CreateTmpFilePath() string {
	seed := lib.RandomAlphanumeric(16)
	tmpDir := filepath.Join(this.data, ".tmp", seed)
	checkDirectory(tmpDir, true)
	seed = lib.RandomAlphanumeric(16)
	return filepath.Join(tmpDir, seed)
}

func (this *FatimaFolderGuide) resolveFolder(programName string) {
	this.app = filepath.Join(this.fatimaHomePath, FatimaFolderApp, programName)
	checkDirectory(this.app, true)
	this.bin = filepath.Join(this.fatimaHomePath, FatimaFolderBin, programName)
	checkDirectory(this.bin, false)
	this.conf = filepath.Join(this.fatimaHomePath, FatimaFolderConf)
	checkDirectory(this.conf, false)
	this.data = filepath.Join(this.fatimaHomePath, FatimaFolderData, programName)
	checkDirectory(this.data, true)
	this.javalib = filepath.Join(this.fatimaHomePath, FatimaFolderJavalib)
	checkDirectory(this.javalib, false)
	this.lib = filepath.Join(this.fatimaHomePath, FatimaFolderLib)
	checkDirectory(this.lib, false)
	this.log = filepath.Join(this.fatimaHomePath, FatimaFolderLog, programName)
	checkDirectory(this.log, true)
	this.pack = filepath.Join(this.fatimaHomePath, FatimaFolderPackage)
	checkDirectory(this.pack, false)
	this.stat = filepath.Join(this.fatimaHomePath, FatimaFolderStat, programName)
	checkDirectory(this.stat, true)
	this.proc = filepath.Join(this.app, FatimaFolderProc)
	checkDirectory(this.proc, true)

	os.RemoveAll(filepath.Join(this.data, ".tmp"))
}

func checkDirectory(path string, forceCreate bool) {
	if err := ensureDirectory(path, forceCreate); err != nil {
		panic(err.Error())
	}
}

func newFolderGuide(proc fatima.SystemProc) fatima.FolderGuide {
	folderGuide := new(FatimaFolderGuide)
	folderGuide.fatimaHomePath = os.Getenv(fatima.ENV_FATIMA_HOME)
	if folderGuide.fatimaHomePath == "" {
		panic("Not found FATIMA_HOME")
	}

	folderGuide.resolveFolder(proc.GetProgramName())

	os.Chdir(folderGuide.proc)

	return folderGuide
}
