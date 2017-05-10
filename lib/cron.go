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
// @author 1100282
// @date 2017. 5. 11. AM 8:50
//

package lib

import (
	robfig_cron "github.com/robfig/cron"
	"throosea.com/fatima"
	"errors"
	"fmt"
	"sync"
	"throosea.com/log"
	"path/filepath"
	"os"
	"time"
	"io/ioutil"
)

const (
	fileRerun = "cron.rerun"
	configPrefix = "cron."
	configSuffixSpec = ".spec"
	configSuffixDesc = ".desc"
)

var (
	cronCreationLock	sync.RWMutex
	cron 				*robfig_cron.Cron
	cronJobList			[]*CronJob = make([]*CronJob, 0)
	fatimaRuntime		fatima.FatimaRuntime
	oneSecondTick		*time.Ticker
	lastRerunModifiedTime time.Time

	errInvalidConfig = errors.New("invalid fatima config")
)

type CronJob struct {
	name			string
	desc			string
	spec			string
	runnable		func(fatima.FatimaRuntime)
}

func (c CronJob) Run() {
	log.Info("start job [%s]", c.name)
	defer func() {
		if r := recover(); r != nil {
			log.Error("panic to execute : %s", r)
		}
	}()

	startMillis := CurrentTimeMillis()
	c.runnable(fatimaRuntime)
	endMillis := CurrentTimeMillis()

	log.Info("cron job [%s] elapsed %s milli seconds", c.name, endMillis - startMillis)
}


func StartCron() {
	if len(cron.Entries()) == 0 {
		return
	}

	log.Info("total %d cron jobs scheduled", len(cron.Entries()))
	cron.Start()
}

func StopCron() {
	cron.Stop()
	if oneSecondTick != nil {
		oneSecondTick.Stop()
		oneSecondTick = nil
		log.Info("cron jobs stopped")
	}
}

func Rerun(jobName string)	{
	log.Info("try to rerun job [%s]", jobName)
	for _, job := range cronJobList {
		if job.name == jobName {
			job.Run()
			return
		}
	}
}

func RegistCronJob(runtime fatima.FatimaRuntime, jobName string, runnable func(fatima.FatimaRuntime)) error {
	if runtime.GetConfig() == nil {
		return errInvalidConfig
	}

	ensureSingleCronInstance(runtime)

	job, err := newCronJob(runtime.GetConfig(), jobName, runnable)
	if err != nil {
		return err
	}

	err = cron.AddJob(job.spec, job)
	if err != nil {
		return err
	}

	log.Info("job[%s] scheduled : %s", jobName, job.spec)
	cronJobList = append(cronJobList, job)

	return nil
}

func ensureSingleCronInstance(runtime fatima.FatimaRuntime) {
	cronCreationLock.Lock()
	if cron == nil {
		cron = robfig_cron.New()
		fatimaRuntime = runtime
		clearRerunFile()
		startRerunFileScanner()
	}
	cronCreationLock.Unlock()
}

func newCronJob(config fatima.Config, name string, runnable func(fatima.FatimaRuntime)) (*CronJob, error) {
	specKey := fmt.Sprintf("%s%s%s", configPrefix, name, configSuffixSpec)
	spec, ok := config.GetValue(specKey)
	if !ok {
		return nil, errors.New("insufficient config key " + specKey)
	}

	descKey := fmt.Sprintf("%s%s%s", configPrefix, name, configSuffixDesc)
	desc, ok := config.GetValue(descKey)
	if ok {
		desc = name
	}

	job := &CronJob{}
	job.name = name
	job.desc = desc
	job.spec = spec
	job.runnable = runnable
	return job, nil
}

func clearRerunFile() {
	file := filepath.Join(fatimaRuntime.GetEnv().GetFolderGuide().GetDataFolder(), fileRerun)
	os.Remove(file)
}

func startRerunFileScanner() {
	oneSecondTick = time.NewTicker(time.Second * 1)
	go func() {
		for range oneSecondTick.C {
			scanRerunFile()
		}
	}()
}

func scanRerunFile() {
	file := filepath.Join(fatimaRuntime.GetEnv().GetFolderGuide().GetDataFolder(), fileRerun)
	stat, err := os.Stat(file)
	if err != nil {
		return
	}

	if lastRerunModifiedTime == stat.ModTime() {
		return
	}

	lastRerunModifiedTime = stat.ModTime()
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	Rerun(string(data))
}

