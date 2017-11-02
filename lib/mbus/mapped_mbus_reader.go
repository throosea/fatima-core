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
	"sync"
	"strings"
	"path/filepath"
	"throosea.com/log"
	"fmt"
	"throosea.com/fatima/lib"
	"time"
)

type MappedMBusReader struct {
	collection    string
	dir           string
	streamDataSet map[string]*StreamData
	streamRecords []*StreamRecord
	xy            Coordinates
	mutex         *sync.Mutex
	running		 bool
}

func NewMappedMBusReader(path string, collection string) (*MappedMBusReader, error) {
	m := MappedMBusReader{running:false}
	m.collection = strings.ToUpper(collection)
	m.dir = filepath.Join(path, mbusDir, m.collection)
	m.mutex = &sync.Mutex{}

	ensureDirectory(m.dir)

	coll, err := loadCollections(&m)
	if err != nil {
		return nil, err
	}

	m.streamRecords = coll

	loadStreamSet(&m)

	log.Info("total %d stream data loaded", len(m.streamDataSet))

	m.running = true

	// startCollectionCleaning
	hourTick := time.NewTicker(time.Hour * 1)
	go func() {
		m.collectionCleaning()
		for range hourTick.C {
			m.collectionCleaning()
		}
	}()

	return &m, nil
}


func (m *MappedMBusReader) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

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
	// TODO

	return nil
}


func (m *MappedMBusReader) collectionCleaning()	{
	if !m.running {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	initial := len(m.streamRecords)
	log.Info("start collection cleaning. initial=%d", initial)
	test := 0
	for true {
		removed := searchRetiredRecord(m.streamRecords)
		if removed >= 0 {
			m.streamRecords = append(m.streamRecords[:removed], m.streamRecords[removed+1:]...)
			log.Info("after removed : %d", len(m.streamRecords))
		}
		test++
		if test > 9 {
			break
		}
	}

	diff := initial - len(m.streamRecords)
	if diff > 0 {
		log.Info("maintain %d stream records", len(m.streamRecords))
	}

}

func searchRetiredRecord(list []*StreamRecord) int	{
	oldMillis := int(time.Now().AddDate(0, 0, -7).UnixNano() / 1000000)
	for i, v := range list	{
		if v.GetLastWriteTime() < oldMillis {
			log.Info("mbus %s marks unused", v.GetProducerName())
			v.markUnused()
			return i
		}
	}

	log.Info("not found...")
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

	log.Info("total %d stream loaded to collection", len(coll))
	return coll, nil
}

func loadStreamSet(mreader *MappedMBusReader)	{
	mreader.streamDataSet = make(map[string]*StreamData)

	for _, v := range mreader.streamRecords {
		data, err := loadStreamData(mreader.dir, mreader.collection, v)
		if err != nil {
			log.Error("fail to load stream[%s] data : %s", v.GetProducerName(), err.Error())
			continue
		}
		mreader.streamDataSet[v.GetProducerName()] = data
	}
}

func loadStreamData(dir string, collection string, stream *StreamRecord) (*StreamData, error) {
	xy := stream.GetReadCoordinates()
	file := fmt.Sprintf("%s.%s.%d", collection, stream.GetProducerName(), xy.sequence)
	path := filepath.Join(dir, file)

	m, err := lib.NewMmap(path, streamDataFileSize)
	if err != nil {
		return nil, err
	}

	data := new(StreamData)
	data.mmap = m

	return data, nil
}
