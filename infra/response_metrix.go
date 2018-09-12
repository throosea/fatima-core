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
// @date 2018. 9. 12. PM 7:30
//

package infra

import (
	"bytes"
	"fmt"
	"throosea.com/fatima"
	"throosea.com/fatima/monitor"
)

type ResponseMarker interface {
	Mark(score int)
}
type ResponseMetrix interface {
	monitor.SystemMeasurable
	ResponseMarker
}

const (
	defaultMax	= 10000
	defaultPart = 10
)

func NewResponseMetrix(fatimaRuntime fatima.FatimaRuntime, name string)	ResponseMarker {
	return NewCustomResponseMetrix(fatimaRuntime, name, defaultMax, defaultPart)
}

func NewCustomResponseMetrix(fatimaRuntime fatima.FatimaRuntime, name string, max, part int) ResponseMarker	{
	m := &basicResponseMetrix{}
	m.name = name
	if max < 1000 {
		max = 1000
	}
	m.max = max

	if part < 1 {
		part = 1
	}
	m.part = part
	m.build()
	fatimaRuntime.RegistMeasureUnit(m)
	return m
}

type basicResponseMetrix struct {
	name 	string
	max 	int
	part 	int
	bunch 	int
	list 	[]int
	labels 	[]string
}

func (b *basicResponseMetrix) build()	{
	b.list = make([]int, b.part)
	b.labels = make([]string, b.part)
	b.bunch = b.max / b.part
	for i:=0; i<b.part; i++	{
		b.labels[i] = fmt.Sprintf("%5d", (i * b.bunch) + b.bunch)
	}
}

func (b *basicResponseMetrix) reset()	{
	for i:=0; i<b.part; i++	{
		b.list[i] = 0
	}
}

func (b *basicResponseMetrix) Mark(score int)	{
	if score <= 0 {
		b.list[0]++
		return
	}

	if score >= b.max {
		b.list[len(b.list)-1]++
		return
	}

	pos := score / b.bunch
	b.list[pos]++
}

func (b *basicResponseMetrix) GetKeyName()	string {
	return b.name
}

func (b *basicResponseMetrix) GetMeasure()	string {
	defer b.reset()

	buff := bytes.Buffer{}
	for i:=0; i<b.part; i++ {
		buff.WriteString(b.labels[i])
		buff.WriteString("|")
	}
	buff.WriteString("\n")
	for i:=0; i<b.part; i++ {
		buff.WriteString(fmt.Sprintf("%5d", b.list[i]))
		buff.WriteString("|")
	}

	return buff.String()
}

