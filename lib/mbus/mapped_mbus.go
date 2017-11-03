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
// @author DeockJin Chung (jin.freestyle@gmail.com)
// @date 2017. 3. 25. PM 8:56
//

package mbus

import (
	"path/filepath"
	"strings"
	"sync"
	"fmt"
	"throosea.com/log"
	"throosea.com/fatima/lib"
	"os"
)

const (
	mbusDir = "mbus"
	maxStreamSize = 100
)

type MappedMBus struct {
	stream     	string
	collection 	string
	dir        	string
	record     *StreamRecord
	data       *StreamData
	xy			Coordinates
	mutex 		*sync.Mutex
}

func (m *MappedMBus) Write(bytes []byte) error {
	if bytes == nil || m.record == nil || m.data == nil || len(bytes) < 1 {
		return nil
	}

	blen := len(bytes)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if int(m.xy.positionOfFile) + 12 + blen > streamDataFileSize {
		log.Debug("rolling next data file. current = %s", m.xy)
		oldData := m.data
		lastPos := m.xy.positionOfFile
		nextSeq := getNextSequence(m.xy.sequence)
		newCord := Coordinates{nextSeq, 0}
		streamData, err := prepareStreamDataFile(m.dir, m.collection, m.stream, newCord)
		if err != nil {
			return err
		}
		m.data = streamData
		m.xy = newCord
		m.markEOFToPreviousFile(oldData, lastPos)
		oldData.Close()
	}

	// write bytes to data(StreamData)
	// 4 : header
	// 4 : length
	// n : data
	buff := make([]byte, 8+blen)
	m.record.mmap.GetByteOrder().PutUint32(buff, streamDataFrameHeader)
	m.record.mmap.GetByteOrder().PutUint32(buff[4:], uint32(blen))
	copy(buff[8:], bytes)
	e := m.data.mmap.Write(int(m.xy.positionOfFile), buff)
	if e != nil {
		return e
	}

	// write(mark) coordinate to record
	size := int(m.xy.positionOfFile) + 8 + blen
	m.xy.positionOfFile = uint32(size)
	m.record.MarkWriteCoordinates(m.xy)
	return nil
}

func getNextSequence(currentSeq uint32) uint32 {
	if currentSeq < 2147483647 {
		return currentSeq + 1
	}
	return 0
}

func prepareStreamDataFile(dir string, collection string, stream string, newCord Coordinates) (*StreamData, error) {
	file := fmt.Sprintf("%s.%s.%d", collection, stream, newCord.sequence)
	path := filepath.Join(dir, file)

	log.Debug("prepare data file path : %s", path)

	// streamDataFileSize
	mm, err := lib.NewMmap(path, streamDataFileSize)
	if err != nil {
		return nil, fmt.Errorf("fail to open new mmap : %s", err.Error())
	}

	data := new(StreamData)
	data.mmap = mm
	data.dir = dir
	data.collection = collection
	data.readThreshold = thresholdReadingMin
	data.name = stream
	return data, nil
}

func (m *MappedMBus) markEOFToPreviousFile(data *StreamData, pos uint32) error {
	buff := make([]byte, 4)
	m.record.mmap.GetByteOrder().PutUint32(buff, streamDataEOF)
	e := data.mmap.Write(int(pos), buff)
	if e != nil {
		return e
	}
	return nil
}

func (m *MappedMBus) Close() error {
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

func NewMappedMBus(path string, collection string, stream string) (*MappedMBus, error) {
	m := MappedMBus{}
	m.collection = strings.ToUpper(collection)
	m.dir = filepath.Join(path, mbusDir, m.collection)
	m.stream = stream
	m.mutex = &sync.Mutex{}

	ensureDirectory(m.dir)

	err := loadCollectionFile(&m)
	if err != nil {
		return nil, err
	}

	log.Debug("loaded record : %s", m.record)

	err = loadDataFile(&m)
	if err != nil {
		m.Close()
		return nil, err
	}

	m.xy = m.record.GetWriteCoordinates()
	return &m, nil
}


func loadCollectionFile(mbus *MappedMBus) error {
	name := fmt.Sprintf("%s.COLLECTION", mbus.collection)
	collectionFile := filepath.Join(mbus.dir, strings.ToUpper(name))
	size := maxStreamSize * streamRecordSize
	m, err := lib.NewMmap(collectionFile, size)
	if err != nil {
		return err
	}

	var newRecord *StreamRecord
	ptr := 0
	for true {
		b := make([]byte, streamRecordSize)
		err = m.Read(ptr, b)
		if err != nil {
			m.Close()
			return err
		}
		r := newStreamRecord(ptr, m)
		if r.IsValid() {
			if r.GetProducerName() == mbus.stream {
				mbus.record = r
				return nil
			}
		} else {
			if newRecord == nil {
				newRecord = r
			}
		}

		ptr += streamRecordSize
		if ptr >= size {
			break
		}
	}

	if newRecord != nil {
		newRecord.mmap = m
		mbus.record = newRecord
		newRecord.MarkAsNew(mbus.stream)
		return nil
	}

	// not found
	m.Close()
	return fmt.Errorf("collection is full but not found stream [%s] in collection", mbus.stream)
}


func loadDataFile(mbus *MappedMBus) error {
	xy := mbus.record.GetWriteCoordinates()
	file := fmt.Sprintf("%s.%s.%d", mbus.collection, mbus.stream, xy.sequence)
	path := filepath.Join(mbus.dir, file)

	m, err := lib.NewMmap(path, streamDataFileSize)
	if err != nil {
		return err
	}

	mbus.data = new(StreamData)
	mbus.data.mmap = m

	return nil
}


func ensureDirectory(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, 0755)
		}
	}

	return nil
}