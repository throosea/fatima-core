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

package monitor

type ResponseMetrix interface {
	SystemMeasurable
	Mark(score int)
}

const (
	defaultMax	= 10000
	defaultPart = 10
)

func NewResponseMetrix(name string)	ResponseMetrix {
	return NewCustomResponseMetrix(name, defaultMax, defaultPart)
}

func NewCustomResponseMetrix(name string, max, part int) ResponseMetrix	{
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
	return m
}

type basicResponseMetrix struct {
	name 	string
	max 	int
	part 	int
}

func (b *basicResponseMetrix) Mark(score int)	{
	// TODO
}

func (b *basicResponseMetrix) GetKeyName()	string {
	// TODO

	return b.name
}

func (b *basicResponseMetrix) GetMeasure()	string {
	// TODO

	return "basicResponseMetrix sample"
}

