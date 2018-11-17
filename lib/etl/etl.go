//
// Copyright (c) 2018 SK Planet.
// All right reserved.
//
// This software is the confidential and proprietary information of K Planet.
// You shall not disclose such Confidential Information and
// shall use it only in accordance with the terms of the license agreement
// you entered into with SK Planet.
//
//
// @project fatima
// @author 1100282
// @date 2018. 11. 17. PM 5:02
//

package etl

import (
	"container/list"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"throosea.com/fatima/lib"
	"throosea.com/log"
	"time"
)

type SimpleETL interface {
	Process()
}

type SimpleETLBuilder interface {
	SetExtractor(func(ExtractionDeliver) error) SimpleETLBuilder
	SetLogger(Logger) SimpleETLBuilder
	SetTransformQueueSize(int) SimpleETLBuilder
	SetTransformer(int, func(interface{}, Loader)) SimpleETLBuilder
	SetLoader(func(interface{})) SimpleETLBuilder
	Build() (SimpleETL, error)
}

type ExtractionDeliver interface {
	Deliver(interface{})
}

type Loader interface {
	Load(interface{})
}

type Logger interface {
	Printf(string, ...interface{})
}

func NewSimpleETLBuilder() SimpleETLBuilder {
	builder := new(dataFetchBuilder)
	builder.transformQueueSize = math.MaxInt16
	builder.transformWorkerSize = runtime.NumCPU()
	return builder
}

type dataFetchBuilder struct {
	extractor           func(ExtractionDeliver) error
	logger              Logger
	transformFunc       func(interface{}, Loader)
	loader 				func(interface{})
	transformWorkerSize int
	transformQueueSize	int
}

func (f *dataFetchBuilder) SetExtractor(extractor func(ExtractionDeliver) error) SimpleETLBuilder {
	f.extractor = extractor
	return f
}

func (f *dataFetchBuilder) SetLogger(logger Logger) SimpleETLBuilder {
	f.logger = logger
	return f
}

func (f *dataFetchBuilder) SetTransformer(workerSize int, transformFunc func(interface{}, Loader)) SimpleETLBuilder {
	if workerSize > 0 {
		f.transformWorkerSize = workerSize
	}
	f.transformFunc = transformFunc
	return f
}

func (f *dataFetchBuilder) SetTransformQueueSize(size int) SimpleETLBuilder	{
	if size > 0 {
		f.transformQueueSize = size
	}
	return f
}

func (f *dataFetchBuilder) SetLoader(loader func(interface{})) SimpleETLBuilder	{
	f.loader = loader
	return f
}

func (f *dataFetchBuilder) Build() (SimpleETL, error) {
	etl := new(simpleETL)
	etl.extractFunc = f.extractor
	etl.logger = f.logger
	etl.loaderFunc = f.loader
	etl.transformFunc = f.transformFunc
	etl.ingestList = list.New()
	etl.ingestFinish = false
	etl.dataChan = make(chan interface{}, math.MaxInt16)
	etl.executor = lib.NewExecutorBuilder(etl.transform).
		SetQueueSize(f.transformQueueSize).
		SetWorkerSize(f.transformWorkerSize).
		Build()

	return etl, nil
}

type simpleETL struct {
	executor      lib.Executor
	extractFunc   func(ExtractionDeliver) error
	transformFunc func(interface{}, Loader)
	dataChan      chan interface{}
	loaderFunc    func(interface{})
	logger        Logger
	loadWg        sync.WaitGroup
	deliverWg 		sync.WaitGroup
	loadCount     uint32
	ingestList    *list.List
	ingestFinish	bool
}

func (d *simpleETL) transform(val interface{})  {
	d.transformFunc(val, d)
}

func (d *simpleETL) Process()  {
	defer func() {
		d.ingestList = nil
		d.executor.Close()
		close(d.dataChan)

		if r := recover(); r != nil {
			log.Info("process panic : %v", r)
			if d.logger != nil {
				d.logger.Printf("process panic : %v\n", r)
			}
			return
		}
	}()

	go d.startLoading()
	go d.startDeliverToTransform()

	// wait for a second (go func started...)
	time.Sleep(time.Second)

	startMillis := lib.CurrentTimeMillis()
	log.Info("start loading....")
	if d.logger != nil {
		d.logger.Printf("%s", "start loading....\n")
	}

	d.startExtracting()

	d.loadWg.Wait()

	log.Info("ETL finish. total %d/%d done : %s", d.executor.Count(), d.loadCount, lib.ExpressDuration(startMillis))
	if d.logger != nil {
		d.logger.Printf("ETL finish. total %d/%d done : %s\n", d.executor.Count(), d.loadCount, lib.ExpressDuration(startMillis))
	}
}

func (d *simpleETL) startExtracting()  {
	log.Info("start extracting with concurrent extractor")
	if d.logger != nil {
		d.logger.Printf("%s", "start extracting with concurrent extractor\n")
	}

	startMillis := lib.CurrentTimeMillis()
	err := d.extractFunc(d)
	if err != nil {
		log.Warn("fail to extract : %s", err)
		if d.logger != nil {
			d.logger.Printf("fail to extract : %s\n", err)
		}
		return
	}

	d.ingestFinish = true
	d.deliverWg.Wait()
	log.Warn("waiting extract done. %s", lib.ExpressDuration(startMillis))
	if d.logger != nil {
		d.logger.Printf("waiting extract done. %s", lib.ExpressDuration(startMillis))
	}

	d.executor.Wait()

	elapsed := lib.ExpressDuration(startMillis)
	log.Warn("extract and transform finished. %d record. %s", d.executor.Count(), elapsed)
	if d.logger != nil {
		d.logger.Printf("extract and transform finished. %d record. %s\n", d.executor.Count(), elapsed)
	}
}

func (d *simpleETL) Deliver(v interface{})	{
	d.ingestList.PushBack(v)
}

func (d *simpleETL) startLoading()	{
	for elem := range d.dataChan {
		d.loaderFunc(elem)
		d.loadWg.Done()
	}
}

func (d *simpleETL) startDeliverToTransform()	{
	log.Info("startDeliverToTransform...")
	d.deliverWg.Add(1)
	defer func() {
		d.deliverWg.Done()
	}()

	for true {
		if d.ingestList == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		if d.ingestList.Len() == 0 && d.ingestFinish {
			return
		}

		elem := d.ingestList.Front()
		if elem == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		d.executor.Insert(elem.Value)
		d.ingestList.Remove(elem)
	}
}

func (d *simpleETL) Load(v interface{})	{
	d.loadWg.Add(1)
	atomic.AddUint32(&d.loadCount, 1)
	d.dataChan <- v
}

