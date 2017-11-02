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
	record     *StreamRecord
	streamList	[]*StreamRecord
	data       *StreamData
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

	for _, v := range m.streamList {
		log.Info("[%s] %v", v.GetProducerName(), v.GetWriteCoordinates())
	}

	m.xy = m.record.GetWriteCoordinates()
	return &m, nil
}


func (m *MappedMBusReader) Close() error {
	if m.record != nil {
		m.record.Close()
		m.record = nil
	}
	if m.data != nil {
		m.data.Close()
		m.data = nil
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
