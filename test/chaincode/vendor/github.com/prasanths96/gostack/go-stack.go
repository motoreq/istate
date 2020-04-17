// Copyright (c) 2020 Prasanth Sundaravelu

/*
	Synchronized Stack accepting dynamic type
*/

//

package gostack

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type stack struct {
	mux sync.Mutex // you don't have to do this if you don't want thread safety
	s   []interface{}
}

//
func NewStack() (sInterface Interface) {
	sInterface = &stack{sync.Mutex{}, nil}
	return
}

func (s *stack) Push(v interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.s = append(s.s, v)
}

func (s *stack) Pop() (interface{}, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *stack) Read() (interface{}, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	return res, nil
}

func (s *stack) Size() int {
	s.mux.Lock()
	defer s.mux.Unlock()

	return len(s.s)
}

func (s *stack) PopBool() (val bool, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(bool); !ok {
		return false, fmt.Errorf("Cannot pop %v value as bool", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopByte() (val byte, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(byte); !ok {
		return byte('\000'), fmt.Errorf("Cannot pop %v value as byte", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopComplex64() (val complex64, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(complex64); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as complex64", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopComplex128() (val complex128, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(complex128); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as complex128", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopError() (val error, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(error); !ok {
		return nil, fmt.Errorf("Cannot pop %v value as error", reflect.TypeOf(result))
	}
	return
}

// Floats
func (s *stack) PopFloat32() (val float32, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(float32); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as float32", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopFloat64() (val float64, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(float64); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as float64", reflect.TypeOf(result))
	}
	return
}

// Ints
func (s *stack) PopInt() (val int, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(int); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as int", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopInt8() (val int8, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(int8); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as int8", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopInt16() (val int16, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(int16); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as int16", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopInt32() (val int32, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(int32); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as int32", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopInt64() (val int64, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(int64); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as int64", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopRune() (val rune, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(rune); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as rune", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopUint() (val uint, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(uint); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as uint", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopUint8() (val uint8, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(uint8); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as uint8", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopUint16() (val uint16, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(uint16); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as uint16", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopUint32() (val uint32, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(uint32); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as uint32", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopUint64() (val uint64, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(uint64); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as uint64", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopUintptr() (val uintptr, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(uintptr); !ok {
		return 0, fmt.Errorf("Cannot pop %v value as uintptr", reflect.TypeOf(result))
	}
	return
}

func (s *stack) PopString() (val string, err error) {

	var ok bool
	result, err := s.Pop()
	if err != nil {
		return
	}
	if val, ok = result.(string); !ok {
		return "", fmt.Errorf("Cannot pop %v value as string", reflect.TypeOf(result))
	}
	return
}

// ReadBool will return the top element of stack as bool
// without popping it
func (s *stack) ReadBool() (val bool, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(bool); !ok {
		return false, fmt.Errorf("Cannot Read %v value as bool", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadByte() (val byte, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(byte); !ok {
		return byte('\000'), fmt.Errorf("Cannot Read %v value as byte", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadComplex64() (val complex64, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(complex64); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as complex64", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadComplex128() (val complex128, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(complex128); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as complex128", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadError() (val error, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(error); !ok {
		return nil, fmt.Errorf("Cannot Read %v value as error", reflect.TypeOf(result))
	}
	return
}

// Floats
func (s *stack) ReadFloat32() (val float32, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(float32); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as float32", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadFloat64() (val float64, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(float64); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as float64", reflect.TypeOf(result))
	}
	return
}

// Ints

func (s *stack) ReadInt() (val int, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(int); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as int", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadInt8() (val int8, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(int8); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as int8", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadInt16() (val int16, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(int16); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as int16", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadInt32() (val int32, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(int32); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as int32", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadInt64() (val int64, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(int64); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as int64", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadRune() (val rune, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(rune); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as rune", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadUint() (val uint, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(uint); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as uint", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadUint8() (val uint8, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(uint8); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as uint8", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadUint16() (val uint16, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(uint16); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as uint16", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadUint32() (val uint32, err error) {

	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(uint32); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as uint32", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadUint64() (val uint64, err error) {
	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(uint64); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as uint64", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadUintptr() (val uintptr, err error) {
	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(uintptr); !ok {
		return 0, fmt.Errorf("Cannot Read %v value as uintptr", reflect.TypeOf(result))
	}
	return
}

func (s *stack) ReadString() (val string, err error) {
	var ok bool
	result, err := s.Read()
	if err != nil {
		return
	}
	if val, ok = result.(string); !ok {
		return "", fmt.Errorf("Cannot Read %v value as string", reflect.TypeOf(result))
	}
	return
}
