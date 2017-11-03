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
)

type MappedMBusReader struct {
	collection    string
	dir           string
	streamDataSet map[string]*StreamData
	streamRecords []*StreamRecord
	xy            Coordinates
	rw            *sync.RWMutex
	running       bool
}

func NewMappedMBusReader(path string, collection string) (*MappedMBusReader, error) {
	m := MappedMBusReader{running: false}
	m.collection = strings.ToUpper(collection)
	m.dir = filepath.Join(path, mbusDir, m.collection)
	m.rw = &sync.RWMutex{}

	ensureDirectory(m.dir)

	coll, err := loadCollections(&m)
	if err != nil {
		return nil, err
	}

	m.streamRecords = coll
	log.Info("total %d stream loaded to collection", len(coll))

	loadStreamSet(&m)

	log.Info("total %d stream data loaded", len(m.streamDataSet))

	m.running = true

	return &m, nil
}

func (m *MappedMBusReader) Close() error {
	m.rw.Lock()
	defer m.rw.Unlock()

	m.running = false

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
		m.startReading()
	} ()

	return nil
}

func (m *MappedMBusReader) startReading() {
	sleepMillis := 1.0

	for true {
		if !m.running	{
			return
		}

		count := m.consumeIncomingData()
		if count == 0 {
			sleepMillis = 1.0
			continue
		}

		sleepMillis = math.Min(sleepMillis * 2, maxConsumingSleepMillis)
		time.Sleep(time.Millisecond * time.Duration(sleepMillis))
	}
}

const (
	maxConsumingSleepMillis = 16.0
)

// TODO
func (m *MappedMBusReader) consumeIncomingData() int {
	m.rw.Lock()
	defer m.rw.Unlock()

	consumeCount := 0
	for _, v := range m.streamRecords {
		data := m.streamDataSet[v.GetProducerName()]
		if data == nil {
			// logging?
			continue
		}
		read, err := data.Read(v.GetReadCoordinates())
		if err != nil {
			// logging?
			log.Error("fail to read : %s", err.Error())
			continue
		}
		if read != nil {
			consumeCount = consumeCount + len(read)
			// TODO : consume...
			for _, v := range read {
				log.Info("%s", string(v))
			}
		}
	}
	return consumeCount
}

func (m *MappedMBusReader) readIncoming() []byte {

	return nil
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
		for _, candi := range coll {
			if strings.Compare(work.GetProducerName(), candi.GetProducerName()) == 0 {
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if found {
		m.reflectCollectionChanges(coll)
	}
}

func (m *MappedMBusReader) reflectCollectionChanges(fresh []*StreamRecord) {
	if fresh == nil {
		return
	}

	log.Info("refresh collections... old[%d], new[%d]", len(m.streamRecords), len(fresh))

	m.rw.Lock()
	m.rw.Unlock()

	removed := make([]*StreamRecord, 0)
	survived := make([]*StreamRecord, 0)

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

	for _, v := range removed {
		name := v.GetProducerName()
		data := m.streamDataSet[name]
		if data != nil {
			data.Close()
		}
		delete(m.streamDataSet, name)
		v.Close()
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
	m, err := lib.NewMmap(collectionFile, size)
	if err != nil {
		return coll, err
	}

	mreader.rw.Lock()
	defer mreader.rw.Unlock()

	ptr := 0
	for true {
		b := make([]byte, streamRecordSize)
		err = m.Read(ptr, b)
		if err != nil {
			m.Close()
			return coll, err
		}
		r := newStreamRecord(ptr, m)
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
