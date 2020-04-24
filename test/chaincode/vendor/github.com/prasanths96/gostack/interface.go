// Copyright (c) 2020 Prasanth Sundaravelu

package gostack

//
type Interface interface {
	Push(v interface{})
	Pop() (interface{}, error)
	Read() (interface{}, error)
	Size() int
	PopBool() (bool, error)
	PopByte() (byte, error)
	PopComplex64() (complex64, error)
	PopComplex128() (complex128, error)
	PopError() (error, error)
	//Floats
	PopFloat32() (float32, error)
	PopFloat64() (float64, error)
	// Ints
	PopInt() (int, error)
	PopInt8() (int8, error)
	PopInt16() (int16, error)
	PopInt32() (int32, error)
	PopInt64() (int64, error)
	PopRune() (rune, error)
	PopUint() (uint, error)
	PopUint8() (uint8, error)
	PopUint16() (uint16, error)
	PopUint32() (uint32, error)
	PopUint64() (uint64, error)
	PopUintptr() (uintptr, error)
	PopString() (string, error)
	// Read will read the last elem without popping
	ReadBool() (bool, error)
	ReadByte() (byte, error)
	ReadComplex64() (complex64, error)
	ReadComplex128() (complex128, error)
	ReadError() (error, error)
	// Floats
	ReadFloat32() (float32, error)
	ReadFloat64() (float64, error)
	// Ints
	ReadInt() (int, error)
	ReadInt8() (int8, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadRune() (rune, error)
	ReadUint() (uint, error)
	ReadUint8() (uint8, error)
	ReadUint16() (uint16, error)
	ReadUint32() (uint32, error)
	ReadUint64() (uint64, error)
	ReadUintptr() (uintptr, error)
	ReadString() (string, error)
}
