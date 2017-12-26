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
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"throosea.com/fatima"
	"throosea.com/log"
	"throosea.com/fatima/monitor"
	"errors"
	"strconv"
	"strings"
	"encoding/json"
)

const (
	proc_status_created = 1 << iota
	proc_status_initializing
	proc_status_ready
	proc_status_running
	proc_status_shutdown
)

const (
	LOG4FATIMA_PROP_BACKUP_DAYS           = "log4fatima.backup.days"
	LOG4FATIMA_PROP_SHOW_METHOD           = "log4fatima.method.show"
	LOG4FATIMA_PROP_SOURCE_PRINTSIZE      = "log4fatima.source.printsize"
	LOG4FATIMA_PROP_FILE_SIZE_LIMIT       = "log4fatima.filesize.limit"
	LOG4FATIMA_DEFAULT_BACKUP_FILE_NUMBER = 30
	LOG4FATIMA_DEFAULT_SOURCE_PRINTSIZE = 30
)

type FatimaProcessStatus uint8

var fatimaProcess *FatimaRuntimeProcess = new(FatimaRuntimeProcess)

func NewFatimaRuntime() *FatimaRuntimeProcess {
	return fatimaProcess
}

type FatimaProcessEnv struct {
	systemProc  fatima.SystemProc
	folderGuide fatima.FolderGuide
	profile     string
}

func (env *FatimaProcessEnv) GetSystemProc() fatima.SystemProc {
	return env.systemProc
}

func (env *FatimaProcessEnv) GetFolderGuide() fatima.FolderGuide {
	return env.folderGuide
}

func (env *FatimaProcessEnv) GetProfile() string {
	return env.profile
}

type FatimaRuntimeBuilder interface {
	GetPkgProcConfig() fatima.FatimaPkgProcConfig
	GetPredefines() fatima.Predefines
	GetConfig() fatima.Config
	GetProcessType() fatima.FatimaProcessType
}

type FatimaPackaging struct {
	name	string
	host	string
	group	string
}

func (p *FatimaPackaging) GetName()	string	{
	return p.name
}

func (p *FatimaPackaging) GetHost()	string	{
	return p.host
}

func (p *FatimaPackaging) GetGroup() string	{
	return p.group
}

type FatimaRuntimeProcess struct {
	env           fatima.FatimaEnv
	platform      fatima.PlatformSupport
	systemStatus  FatimaPackageSystemStatus
	sigs          chan os.Signal
	logLevel      log.LogLevel
	builder       FatimaRuntimeBuilder
	packaging     *FatimaPackaging
	interactor    fatima.ProcessInteractor
	notifyHandler monitor.SystemNotifyHandler
	status        FatimaProcessStatus
}

func (process *FatimaRuntimeProcess) GetEnv() fatima.FatimaEnv {
	return process.env
}

func (process *FatimaRuntimeProcess) GetLogLevel() log.LogLevel {
	return process.logLevel
}

func (process *FatimaRuntimeProcess) SetLogLevel(logLevel log.LogLevel) {
	process.logLevel = logLevel
}

func (process *FatimaRuntimeProcess) SetInteractor(interactor fatima.ProcessInteractor) {
	process.interactor = interactor
}

func (process *FatimaRuntimeProcess) GetConfig() fatima.Config {
	return process.builder.GetConfig()
}

func (process *FatimaRuntimeProcess) GetPackaging() fatima.Packaging {
	if process.packaging == nil {
		pack := FatimaPackaging{name: "default", host: "unknown", group: "basic"}
		v, ok := process.builder.GetPredefines().GetDefine(GLOBAL_DEFINE_PACKAGE_GROUPNAME)
		if ok {
			pack.group = v
		}
		v, ok = process.builder.GetPredefines().GetDefine(GLOBAL_DEFINE_PACKAGE_NAME)
		if ok {
			pack.name = v
		}
		v, ok = process.builder.GetPredefines().GetDefine(GLOBAL_DEFINE_PACKAGE_HOSTNAME)
		if ok {
			pack.host = v
		} else {
			n, err := os.Hostname()
			if err != nil {
				pack.host = "unknown"
			} else {
				pack.host = n
			}
		}
		process.packaging = &pack
	}

	return process.packaging
}

func (process *FatimaRuntimeProcess) GetSystemStatus() monitor.FatimaSystemStatus {
	return &process.systemStatus
}

func (process *FatimaRuntimeProcess) GetSystemNotifyHandler() monitor.SystemNotifyHandler {
	return process.notifyHandler
}

func (process *FatimaRuntimeProcess) GetBuilder() FatimaRuntimeBuilder {
	return process.builder
}

func (process *FatimaRuntimeProcess) IsRunning() bool {
	if process.status == proc_status_running || process.status == proc_status_ready {
		return true
	}

	return false
}

func (process *FatimaRuntimeProcess) Run() {
	if process.status >= proc_status_running {
		log.Warn("aleady process run")
		return
	}

	process.status = proc_status_running

	sigs := make(chan os.Signal, 1)
	go func() {
		sig := <-process.sigs
		process.status = proc_status_shutdown
		sigs <- sig
	}()

	if !process.interactor.Initialize() {
		log.Warn("프로세스 초기화에 실패하였습니다. %s 프로그램을 종료합니다", process.env.GetSystemProc().GetProgramName())
		log.Close()
		return
	}

	process.interactor.Run()

	defer func() {
		if r := recover(); r != nil {
			log.Warn("**PANIC** while running", errors.New(fmt.Sprintf("%s", r)))
			process.status = proc_status_shutdown
			process.interactor.Shutdown()
			log.Close()
			return
		}
	}()

	<-sigs
	process.interactor.Shutdown()
}

func (process *FatimaRuntimeProcess) Stop() {
	p, _ := os.FindProcess(process.env.GetSystemProc().GetPid())
	p.Signal(os.Interrupt)
}

func (process *FatimaRuntimeProcess) Regist(component fatima.FatimaComponent) {
	if process.IsRunning()	{
		process.interactor.Regist(component)
	}
}

func (process *FatimaRuntimeProcess) RegistSystemHAAware(aware monitor.FatimaSystemHAAware) {
	if process.IsRunning()	{
		process.interactor.RegistSystemHAAware(aware)
	}
}

func (process *FatimaRuntimeProcess) RegistSystemPSAware(aware monitor.FatimaSystemPSAware) {
	if process.IsRunning()	{
		process.interactor.RegistSystemPSAware(aware)
	}
}

func (process *FatimaRuntimeProcess) RegistMeasureUnit(unit monitor.SystemMeasurable) {
	if process.IsRunning()	{
		process.interactor.RegistMeasureUnit(unit)
	}
}

func (process *FatimaRuntimeProcess) Initialize(builder FatimaRuntimeBuilder)  {
	if process.status >= proc_status_initializing {
		return
	}

	process.status = proc_status_initializing
	process.builder = builder

	pkgProc := process.getThisPkgProc()

	buildLogging(builder)

	process.logLevel = pkgProc.GetLogLevel()
	if process.logLevel != log.GetLevel() {
		log.SetLevel(process.logLevel)
		log.Info("로그레벨을 변경합니다 : %s", process.logLevel)
	}

	process.parepareProcFolder(pkgProc, builder.GetProcessType())
	process.status = proc_status_ready
}

func buildLogging(builder FatimaRuntimeBuilder) {
	// log4fatima show method preference
	v, ok := builder.GetConfig().GetValue(LOG4FATIMA_PROP_SHOW_METHOD)
	if ok {
		if strings.ToLower(v) == "false" {
			log.SetShowMethod(false)
		} else {
			log.SetShowMethod(true)
		}
	}

	// log4fatima source printsize
	v, ok = builder.GetConfig().GetValue(LOG4FATIMA_PROP_SOURCE_PRINTSIZE)
	if ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Warn("[%s] invalid value format : %s", LOG4FATIMA_PROP_SOURCE_PRINTSIZE, v)
		} else {
			log.SetSourcePrintSize(uint8(i))
		}
	}

	// log4fatima backup days
	v, ok = builder.GetConfig().GetValue(LOG4FATIMA_PROP_BACKUP_DAYS)
	if ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Warn("[%s] invalid value format : %s", LOG4FATIMA_PROP_BACKUP_DAYS, v)
		} else {
			log.SetKeepingFileDays(uint16(i))
		}
	}

	// log4fatima file size limit
	v, ok = builder.GetConfig().GetValue(LOG4FATIMA_PROP_FILE_SIZE_LIMIT)
	if ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Warn("[%s] invalid value format : %s", LOG4FATIMA_PROP_FILE_SIZE_LIMIT, v)
		} else {
			log.SetFileSizeLimitMB(uint16(i))
		}
	}
}

func (process *FatimaRuntimeProcess) getThisPkgProc() fatima.FatimaPkgProc {
	fatimaProc := process.builder.GetPkgProcConfig().GetProcByName(process.env.GetSystemProc().GetProgramName())
	if fatimaProc == nil {
		panic("not found " + process.env.GetSystemProc().GetProgramName() + " proc configuration")
	}

	return fatimaProc
}

func (process *FatimaRuntimeProcess) parepareProcFolder(proc fatima.FatimaPkgProc, processType fatima.FatimaProcessType) {
	procFolder := process.env.GetFolderGuide().GetAppProcFolder()

	// remove old output files
	files, _ := filepath.Glob(fmt.Sprintf("%s%c%s.*.output", procFolder, filepath.Separator, proc.GetName()))
	for _, v := range files {
		os.Remove(v)
	}

	// remove old pid files
	files, _ = filepath.Glob(fmt.Sprintf("%s%c%s.pid", procFolder, filepath.Separator, proc.GetName()))
	for _, v := range files {
		os.Remove(v)
	}

	// create my pid file
	pid := []byte(fmt.Sprintf("%d", process.env.GetSystemProc().GetPid()))
	err := ioutil.WriteFile(filepath.Join(procFolder, process.env.GetSystemProc().GetProgramName()+".pid"), pid, 0644)
	check(err)

	if processType == fatima.PROCESS_TYPE_GENERAL {
		// redirect output to file
		outfile, err := os.Create(
			filepath.Join(
				procFolder,
				fmt.Sprintf("%s.%d.output", process.env.GetSystemProc().GetProgramName(), process.env.GetSystemProc().GetPid())))
		check(err)

		os.Stdout = outfile
		os.Stderr = outfile
	}
}

func init() {
	log.SetLevel(log.LOG_TRACE)

	fatimaProcess.sigs = make(chan os.Signal, 1)
	signal.Notify(fatimaProcess.sigs, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	fatimaProcess.status = proc_status_created
	fatimaProcess.env = newFatimaProcessEnv()
	fatimaProcess.platform = createPlatformSupport()
	err := fatimaProcess.platform.EnsureSingleInstance(fatimaProcess.env.GetSystemProc())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(0)
		//fatimaProcess.status = proc_status_shutdown
	}
	fatimaProcess.notifyHandler, err = NewDefaultSystemNotifyHandler(fatimaProcess)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(0)
	}

	logPref := log.NewPreference(fatimaProcess.env.GetFolderGuide().GetLogFolder())
	logPref.DeliveryMode = log.DELIVERY_MODE_ASYNC
	log.Initialize(logPref)

	log.Warn("%s 프로세스를 시작합니다", fatimaProcess.env.GetSystemProc().GetProgramName())

	displayDeploymentInfo(fatimaProcess.env)

	//if fatimaProcess.status == proc_status_shutdown {
	//	log.Warn("%s 프로세스를 종료합니다", fatimaProcess.env.GetSystemProc().GetProgramName())
	//}
}

func newFatimaProcessEnv() *FatimaProcessEnv {
	processEnv := new(FatimaProcessEnv)
	processEnv.systemProc = newSystemProc()
	processEnv.folderGuide = newFolderGuide(processEnv.systemProc)
	processEnv.profile = os.Getenv(fatima.ENV_FATIMA_PROFILE)
	return processEnv
}

func createPlatformSupport() fatima.PlatformSupport {
	return new(OSPlatform)
	/*
	switch runtime.GOOS {
	case "linux":
		return new(PlatformLinux)
	case "darwin":
		return new(PlatformOSX)
	default:
		// windows, freebsd
		panic("Unsupported fatima arch")
	}
	//	return support
	*/
}

func check(e error) {
	if e != nil {
		panic(fmt.Errorf("fail to build runtime : ", e))
	}
}

const (
	deploymentJsonFile = "deployment.json"
)

func displayDeploymentInfo(env fatima.FatimaEnv) {
	deploymentFile := filepath.Join(env.GetFolderGuide().GetAppFolder(), deploymentJsonFile)
	file, err := ioutil.ReadFile(deploymentFile)
	if err != nil {
		fmt.Printf("readfile err : %s\n", err.Error())
		return
	}

	deployment := Deployment{}
	err = json.Unmarshal(file, &deployment)
	if err != nil {
		fmt.Printf("json unmarshal err : %s\n", err.Error())
		return
	}

	if deployment.HasBuildInfo() {
		log.Info("패키지 빌드 시각 : %s", deployment.Build.BuildTime)
		if deployment.Build.HasGit() {
			log.Info("패키지 빌드 (git) : %s", deployment.Build.Git)
		}
	}
}

type Deployment struct {
	Process		string		`json:"process"`
	ProcessType string		`json:"process_type,omitempty"`
	Build 		DeploymentBuild `json:"build,omitempty"`
}

func (d Deployment) HasBuildInfo()	bool 	{
	if len(d.Build.BuildTime) == 0 {
		return false
	}
	return true
}

type DeploymentBuild struct {
	Git			DeploymentBuildGit `json:"git,omitempty"`
	BuildTime 	string		`json:"time,omitempty"`
}

func (d DeploymentBuild) HasGit()	bool 	{
	if len(d.Git.Branch) == 0 {
		return false
	}
	return true
}

type DeploymentBuildGit struct {
	Branch		string		`json:"branch"`
	Commit		string		`json:"commit"`
}

func (d DeploymentBuildGit) String()	string 	{
	return fmt.Sprintf("Branch=[%s], Commit=[%s]", d.Branch, d.Commit)
}
