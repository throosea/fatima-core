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
)

type MappedMBusReader struct {
	stream     	string
	collection 	string
	dir        	string
	record     *StreamRecord
	data       *StreamData
	xy			Coordinates
	mutex 		*sync.Mutex
}

func NewMappedMBusReader(path string, collection string, stream string) (*MappedMBusReader, error) {
	m := MappedMBusReader{}
	m.collection = strings.ToUpper(collection)
	m.dir = filepath.Join(path, mbusDir, m.collection)
	m.stream = stream
	m.mutex = &sync.Mutex{}

	ensureDirectory(m.dir)

	//err := loadCollectionFile(&m)
	//if err != nil {
	//	return nil, err
	//}

	log.Debug("loaded record : %s", m.record)

	//err = loadDataFile(&m)
	//if err != nil {
	//	m.Close()
	//	return nil, err
	//}

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
