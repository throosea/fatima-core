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
	"fmt"
	"throosea.com/fatima/lib"
	"throosea.com/log"
	"time"
	"math"
)

const (
	streamDataFileSize = 4 * 1024 * 1024 // 4m
	streamRecordSize   = 100
)

var streamRecordHeader uint32 = 0xd0adbeef
var streamDataFrameHeader uint32 = 0xaaadbe0f
var streamDataEOF uint32 = 0x2211be0f

type Coordinates struct {
	sequence       uint32
	positionOfFile uint32
}

func (c Coordinates) String() string {
	return fmt.Sprintf("seq=[%d], pos=[%d]", c.sequence, c.positionOfFile)
}

type StreamData struct {
	dir        string
	collection string
	name       string
	readThreshold int
	mmap       *lib.Mmap
}

func (r *StreamData) Close() error {
	if r.mmap != nil {
		return r.mmap.Close()
	}
	return nil
}

func (m *StreamData) Read(coord Coordinates) ([][]byte, Coordinates, error) {
	baseIndex := int(coord.positionOfFile)
	magic, err := m.mmap.ReadUint32(baseIndex)
	if err != nil {
		return nil, coord, err
	}

	if magic == streamDataEOF {
		log.Info("rolling next data file. current = %s", coord)
		nextSeq := getNextSequence(coord.sequence)
		coord = Coordinates{nextSeq, 0}
		old := m.mmap
		nextStreamData, err := prepareStreamDataFile(m.dir, m.collection, m.name, coord)
		if err != nil {
			// temporary
			time.Sleep(time.Second)
			return nil, coord, fmt.Errorf("fail to load stream data : %s", err.Error())
		}
		m.mmap = nextStreamData.mmap
		old.Close()
		return nil, coord, nil
	} else if magic != streamDataFrameHeader {
		// return empty
		return nil, coord, nil
	}

	bulk := make([][]byte, 0)

	data, baseIndex, err := m.readData(baseIndex)
	if err != nil {
		// temporary
		time.Sleep(time.Second)
		return nil, coord, fmt.Errorf("fail to read stream data record : %s", err.Error())
	}

	bulk = append(bulk, data)

	for true {
		magic, err := m.mmap.ReadUint32(baseIndex)
		if err != nil {
			log.Warn("fail to read stream header [%s, baseIndex=%d] : %s", coord, baseIndex, err.Error())
			break
		}

		if magic != streamDataFrameHeader {
			m.readThreshold = thresholdReadingMin
			break
		}
		data, baseIndex, err := m.readData(baseIndex)
		if err != nil {
			log.Warn("fail to read stream data record [%s, baseIndex=%d] : %s", coord, baseIndex, err.Error())
			break
		}
		bulk = append(bulk, data)
		if len(bulk) >= m.readThreshold {
			a := float64(m.readThreshold * 2)
			b := float64(thresholdReadingMax)
			m.readThreshold = int(math.Min(a, b))
			break
		}
	}

	return bulk, Coordinates{coord.sequence, uint32(baseIndex)}, nil
}

const (
	thresholdReadingMin = 16
	thresholdReadingMax = 8192
)

func (m *StreamData) readData(baseIndex int) ([]byte, int, error) {
	// case of streamDataFrameHeader
	// write bytes to data(StreamData)
	// 4 : header
	// 4 : length
	// n : data
	dlen, err := m.mmap.ReadUint32(baseIndex + 4)
	if err != nil {
		return nil, baseIndex, err
	}

	data := make([]byte, dlen)
	err = m.mmap.Read(baseIndex+8, data)
	if err != nil {
		return nil, baseIndex, err
	}

	baseIndex = baseIndex + 8 + int(dlen)
	return data, baseIndex, nil
}

//
// StreamRecord structure
// magic			uint32		// 0
// epochTime		int			// 4
// lastWrittenTime	int			// 12
// writeSeq			uint32		// 20
// writePos			uint32		// 24
// readSeq			uint32		// 28
// readPos			uint32		// 32
// producerNameLen	uint32		// 36
// producerName		string		// 40
//
type StreamRecord struct {
	baseline int
	mmap     *lib.Mmap
}

func (r *StreamRecord) String() string {
	return fmt.Sprintf("baseline=%d", r.baseline)
}

func (r *StreamRecord) MarkAsNew(producerName string) {
	log.Debug("creating new producer : %s", producerName)
	pdu := make([]byte, streamRecordSize)
	r.mmap.GetByteOrder().PutUint32(pdu, streamRecordHeader)
	epochTime := int(time.Now().UnixNano() / 1000000)
	r.mmap.GetByteOrder().PutUint64(pdu[4:], uint64(epochTime))
	r.mmap.GetByteOrder().PutUint32(pdu[36:], uint32(len(producerName)))
	copy(pdu[40:], []byte(producerName))
	r.mmap.Write(r.baseline, pdu)
}

func (r *StreamRecord) Close() error {
	if r.mmap != nil {
		return r.mmap.Close()
	}
	return nil
}

func (r *StreamRecord) IsValid() bool {
	v, _ := r.mmap.ReadUint32(r.baseline)
	if v == streamRecordHeader {
		return true
	}
	return false
}

func (r *StreamRecord) markUnused() {
	r.mmap.WriteUint32(r.baseline, 0)
	r.mmap.WriteUint64(r.baseline, 0)
	r.mmap.Flush()
}

func (r *StreamRecord) MarkWriteCoordinates(xy Coordinates) {
	r.mmap.WriteUint32(r.baseline+20, xy.sequence)
	r.mmap.WriteUint32(r.baseline+24, xy.positionOfFile)
	r.mmap.WriteUint64(r.baseline+12, uint64(lib.CurrentTimeMillis()))
}

func (r *StreamRecord) MarkReadCoordinates(xy Coordinates) {
	r.mmap.WriteUint32(r.baseline+28, xy.sequence)
	r.mmap.WriteUint32(r.baseline+32, xy.positionOfFile)
}

func (r *StreamRecord) GetLastWriteTime() int {
	v, _ := r.mmap.ReadUint64(r.baseline + 12)
	return int(v)
}

func (r *StreamRecord) GetWriteCoordinates() Coordinates {
	c := Coordinates{}
	v, _ := r.mmap.ReadUint32(r.baseline + 20)
	c.sequence = v
	v, _ = r.mmap.ReadUint32(r.baseline + 24)
	c.positionOfFile = v
	return c
}

func (r *StreamRecord) GetReadCoordinates() Coordinates {
	c := Coordinates{}
	v, _ := r.mmap.ReadUint32(r.baseline + 28)
	c.sequence = v
	v, _ = r.mmap.ReadUint32(r.baseline + 32)
	c.positionOfFile = v
	return c
}

func (r *StreamRecord) GetProducerName() string {
	size, _ := r.mmap.ReadUint32(r.baseline + 36)
	if size < 1 {
		return ""
	}
	buff := make([]byte, size)
	r.mmap.Read(r.baseline+40, buff)
	return string(buff)
}

func newStreamRecord(baseline int, mmap *lib.Mmap) *StreamRecord {
	r := StreamRecord{}
	r.baseline = baseline
	r.mmap = mmap
	return &r
}
