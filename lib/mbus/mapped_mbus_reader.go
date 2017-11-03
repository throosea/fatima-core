//
// Copyright (c) 2017 SK TECHX.
// All right reserved.
//
// This software is the confidential and proprietary information of SK TECHX.
// You shall not disclose such Confidential Information and
// shall use it only in accordance with the terms of the license agreement
// you entered into with SK TECHX.
//
//
// @project fatima
// @author 1100282
// @date 2017. 11. 2. PM 3:33
//

package mbus

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"throosea.com/fatima/lib"
	"throosea.com/log"
	"time"
	"bytes"
	"math"
	"os"
)

type MappedMBusReader struct {
	collection    string
	dir           string
	streamDataSet map[string]*StreamData
	streamRecords []*StreamRecord
	xy            Coordinates
	rw            *sync.RWMutex
	running       bool
	consume			func([]byte)
	recordChan		chan IncomingRecord
}

type IncomingRecord struct {
	last  	bool
	data	[]byte
}

func NewMappedMBusReader(path string, collection string, consume func([]byte)) (*MappedMBusReader, error) {
	m := MappedMBusReader{running: false}
	m.collection = strings.ToUpper(collection)
	m.dir = filepath.Join(path, mbusDir, m.collection)
	m.rw = &sync.RWMutex{}
	m.consume = consume

	ensureDirectory(m.dir)

	coll, err := loadCollections(&m)
	if err != nil {
		return nil, err
	}

	m.streamRecords = coll
	log.Info("total %d stream loaded to collection", len(coll))

	for _, v := range coll {
		log.Trace("[%s] write=[%s], read=[%s]", v.GetProducerName(), v.GetWriteCoordinates(), v.GetReadCoordinates())
	}

	loadStreamSet(&m)

	log.Info("total %d stream data loaded", len(m.streamDataSet))

	if consume != nil {
		m.recordChan = make(chan IncomingRecord, math.MaxUint16)
	}

	m.running = true

	return &m, nil
}

func (m *MappedMBusReader) Close() error {
	m.running = false

	m.rw.Lock()
	defer m.rw.Unlock()

	if m.consume != nil {
		m.recordChan <- IncomingRecord{last:true}
	}

	for _, v := range m.streamRecords {
		v.Close()
	}

	for _, v := range m.streamDataSet {
		v.Close()
	}

	return nil
}

func (m *MappedMBusReader) Activate() error {
	if !m.running {
		return fmt.Errorf("mbusreader stopped. give up activate")
	}

	m.startCollectionCleaning()
	m.startCheckingCollectionModified()

	go func()	{
		m.startRecordConsume()
	} ()

	go func()	{
		m.startReading()
	} ()

	log.Info("mbus reader activated")
	return nil
}

func (m *MappedMBusReader) startReading() {
	sleepMillis := 1.0

	logCnt := 0
	for true {
		if !m.running	{
			return
		}

		count := m.readIncomingData()
		logCnt++
		//if logCnt < 100 {
		//	log.Info("count : %d", count)
		//}
		if count == 0 {
			sleepMillis = math.Min(sleepMillis * 2, maxConsumingSleepMillis)
			time.Sleep(time.Millisecond * time.Duration(sleepMillis))
		} else {
			sleepMillis = 1.0
		}
	}
}

const (
	maxConsumingSleepMillis = 16.0
)

func (m *MappedMBusReader) readIncomingData() int {
	m.rw.Lock()
	defer m.rw.Unlock()

	consumeCount := 0
	for _, v := range m.streamRecords {
		if !m.running {
			break
		}

		data := m.streamDataSet[v.GetProducerName()]
		if data == nil {
			// logging?
			continue
		}

		readCoord := v.GetReadCoordinates()
		read, newCoord, err := data.Read(readCoord)
		if err != nil {
			// logging?
			log.Error("fail to read : %s", err.Error())
			continue
		}

		if read != nil || (readCoord.sequence != newCoord.sequence) {
			v.MarkReadCoordinates(newCoord)
		}

		if read != nil {
			consumeCount = consumeCount + len(read)
			if m.consume != nil {
				for _, v := range read {
					m.recordChan <- IncomingRecord{last:false, data:v}
				}
			}
		}
	}
	return consumeCount
}

func (m *MappedMBusReader) startRecordConsume()  {
	if m.consume == nil {
		return
	}

	for {
		r := <- m.recordChan
		if r.last {
			return
		}

		m.consume(r.data)
	}
}

func (m *MappedMBusReader) startCollectionCleaning() {
	hourTick := time.NewTicker(time.Hour * 1)
	go func() {
		m.collectionCleaning()
		for range hourTick.C {
			m.collectionCleaning()
		}
	}()
}

func (m *MappedMBusReader) startCheckingCollectionModified() {
	secondTick := time.NewTicker(time.Second * 1)
	go func() {
		for range secondTick.C {
			m.checkCollectionModified()
		}
	}()
}

func (m *MappedMBusReader) collectionCleaning() {
	if !m.running {
		return
	}

	m.rw.Lock()
	defer m.rw.Unlock()

	initial := len(m.streamRecords)
	log.Debug("start collection cleaning. initial=%d", initial)

	for true {
		removed := searchRetiredRecord(m.streamRecords)
		if removed < 0 {
			break
		}

		m.streamRecords = append(m.streamRecords[:removed], m.streamRecords[removed+1:]...)
	}

	diff := initial - len(m.streamRecords)
	if diff > 0 {
		log.Info("maintain %d stream records", len(m.streamRecords))
	}

}

func (m *MappedMBusReader) checkCollectionModified() {
	if !m.running {
		return
	}

	coll, err := loadCollections(m)
	if err != nil {
		return
	}

	if len(m.streamRecords) == 0 && len(coll) == 0 {
		return
	}

	if len(m.streamRecords) != len(coll) {
		m.reflectCollectionChanges(coll)
		return
	}

	m.rw.RLock()
	defer m.rw.RUnlock()

	found := false
	for _, work := range m.streamRecords {
		found = false
		for _, candi := range coll {
			if strings.Compare(work.GetProducerName(), candi.GetProducerName()) == 0 {
				found = true
				break
			}
		}
		if !found {
			break
		}
	}

	if !found {
		m.reflectCollectionChanges(coll)
	} else {
		if len(coll) > 0 {
			coll[0].Close()
		}
	}
}

func (m *MappedMBusReader) reflectCollectionChanges(fresh []*StreamRecord) {
	if fresh == nil {
		return
	}

	log.Info("refresh collections... old[%d], new[%d]", len(m.streamRecords), len(fresh))

	m.rw.Lock()
	defer m.rw.Unlock()

	removed := make([]*StreamRecord, 0)
	survived := make([]*StreamRecord, 0)

	//oldMaster := m.streamRecords
	// remove unused stream
	for _, v := range m.streamRecords {
		found := false
		for _, t := range fresh {
			if strings.Compare(v.GetProducerName(), t.GetProducerName()) == 0 {
				found = true
				break
			}
		}
		if found {
			survived = append(survived, v)
		} else {
			removed = append(removed, v)
		}
	}

	log.Info("collection. survived=%d, removed=%d", len(survived), len(removed))
	// find new stream
	for _, v := range fresh {
		found := false
		for _, t := range m.streamRecords {
			if strings.Compare(v.GetProducerName(), t.GetProducerName()) == 0 {
				found = true
				break
			}
		}

		if !found {
			data, err := prepareStreamDataFile(m.dir, m.collection, v.GetProducerName(), v.GetReadCoordinates())
			if err != nil {
				log.Error("fail to load stream[%s] data : %s", v.GetProducerName(), err.Error())
			} else {
				m.streamDataSet[v.GetProducerName()] = data
				survived = append(survived, v)
			}
		}
	}

	m.streamRecords = survived
	//if oldMaster != nil && len(oldMaster) > 0 {
	//	oldMaster[0].Close()
	//}

	for _, v := range removed {
		name := v.GetProducerName()
		log.Info("try to delete %s", name)
		data := m.streamDataSet[name]
		if data != nil {
			data.Close()
			file := fmt.Sprintf("%s.%s.%d", data.collection, name, v.GetReadCoordinates().sequence)
			log.Info("removing file %s", file)
			err := os.Remove(filepath.Join(m.dir, file))
			if err != nil {
				log.Warn("fail to remove stream data file [%s] : %s", file, err.Error())
			}
		}

		delete(m.streamDataSet, name)
		v.markUnused()
	}

	var buff bytes.Buffer
	for _, v := range survived {
		buff.WriteString(v.GetProducerName())
		buff.WriteByte(' ')
	}

	log.Warn("after collection refined : %s", buff.String())
}

func searchRetiredRecord(list []*StreamRecord) int {
	oldMillis := int(time.Now().AddDate(0, 0, -7).UnixNano() / 1000000)
	for i, v := range list {
		if v.GetLastWriteTime() < oldMillis {
			log.Info("mbus %s marks unused", v.GetProducerName())
			v.markUnused()
			return i
		}
	}

	return -1
}

func loadCollections(mreader *MappedMBusReader) ([]*StreamRecord, error) {
	coll := make([]*StreamRecord, 0)

	name := fmt.Sprintf("%s.COLLECTION", mreader.collection)
	collectionFile := filepath.Join(mreader.dir, strings.ToUpper(name))
	size := maxStreamSize * streamRecordSize
	recordMaster, err := lib.NewMmap(collectionFile, size)
	if err != nil {
		return nil, err
	}

	mreader.rw.Lock()
	defer mreader.rw.Unlock()

	ptr := 0
	for true {
		b := make([]byte, streamRecordSize)
		err = recordMaster.Read(ptr, b)
		if err != nil {
			recordMaster.Close()
			return nil, err
		}
		r := newStreamRecord(ptr, recordMaster)
		if r.IsValid() {
			coll = append(coll, r)
		}

		ptr += streamRecordSize
		if ptr >= size {
			break
		}
	}

	return coll, nil
}

func loadStreamSet(mreader *MappedMBusReader) {
	mreader.streamDataSet = make(map[string]*StreamData)

	for _, v := range mreader.streamRecords {
		data, err := prepareStreamDataFile(mreader.dir, mreader.collection, v.GetProducerName(), v.GetReadCoordinates())
		if err != nil {
			log.Error("fail to load stream[%s] data : %s", v.GetProducerName(), err.Error())
			continue
		}
		mreader.streamDataSet[v.GetProducerName()] = data
	}
}
