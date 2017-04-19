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
	"time"
	"throosea.com/log"
	"throosea.com/fatima/lib"
)


const (
	streamDataFileSize = 4 * 1024 * 1024 // 4m
	streamRecordSize = 100
)

var streamRecordHeader uint32 = 0xd0adbeef
var streamDataFrameHeader uint32 = 0xaaadbe0f
var streamDataEOF uint32 = 0x2211be0f

type Coordinates struct {
	sequence		uint32
	positionOfFile	uint32
}

func (c Coordinates) String() string {
	return fmt.Sprintf("seq=[%d], pos=[%d]", c.sequence, c.positionOfFile)
}

type StreamData struct {
	mmap	*lib.Mmap
}

func (r *StreamData) Close() error {
	if r.mmap != nil {
		return r.mmap.Close()
	}
	return nil
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
	baseline		int
	mmap			*lib.Mmap
}

func (r *StreamRecord) String() string {
	return fmt.Sprintf("baseline=%d", r.baseline)
}

func (r *StreamRecord) MarkAsNew(producerName string)  {
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

func (r *StreamRecord) WriteCoordinates(xy Coordinates)  {
	r.mmap.WriteUint32(r.baseline + 20, xy.sequence)
	r.mmap.WriteUint32(r.baseline + 24, xy.positionOfFile)
	r.mmap.WriteUint64(r.baseline + 12, uint64(lib.CurrentTimeMillis()))
}

func (r *StreamRecord) GetWriteCoordinates() Coordinates {
	c := Coordinates{}
	v, _ := r.mmap.ReadUint32(r.baseline + 20)
	c.sequence = v
	v, _ = r.mmap.ReadUint32(r.baseline + 24)
	c.positionOfFile = v
	return c
}

func (r *StreamRecord) GetProducerName() string {
	size, _ := r.mmap.ReadUint32(r.baseline + 36)
	if size < 1 {
		return ""
	}
	buff := make([]byte, size)
	r.mmap.Read(r.baseline + 40, buff)
	return string(buff)
}

func newStreamRecord(baseline int, mmap *lib.Mmap) *StreamRecord {
	r := StreamRecord{}
	r.baseline = baseline
	r.mmap = mmap
	return &r
}
