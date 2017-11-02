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
)

type MappedMBusReader struct {
	collection 	string
	dir        	string
	streamSet	map[string]*StreamData
	streamList	[]*StreamRecord
	xy			Coordinates
	mutex 		*sync.Mutex
}

func NewMappedMBusReader(path string, collection string) (*MappedMBusReader, error) {
	m := MappedMBusReader{}
	m.collection = strings.ToUpper(collection)
	m.dir = filepath.Join(path, mbusDir, m.collection)
	m.mutex = &sync.Mutex{}

	ensureDirectory(m.dir)

	coll, err := loadCollections(&m)
	if err != nil {
		return nil, err
	}

	m.streamList = coll

	loadStreamSet(&m)

	log.Info("total %d stream data loaded", len(m.streamSet))

	return &m, nil
}


func (m *MappedMBusReader) Close() error {
	for _, v := range m.streamList {
		v.Close()
	}

	for _, v := range m.streamSet {
		v.Close()
	}

	return nil
}

func (m *MappedMBusReader) Activate() error {
	// TODO

	return nil
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
	mreader.streamSet = make(map[string]*StreamData)

	for _, v := range mreader.streamList {
		data, err := loadStreamData(mreader.dir, mreader.collection, v)
		if err != nil {
			log.Error("fail to load stream[%s] data : %s", v.GetProducerName(), err.Error())
			continue
		}
		mreader.streamSet[v.GetProducerName()] = data
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