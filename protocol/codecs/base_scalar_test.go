// This source file is part of the EdgeDB open source project.
//
// Copyright 2020-present EdgeDB Inc. and the EdgeDB authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codecs

import (
	"testing"
	"time"
	"unsafe"

	"github.com/edgedb/edgedb-go/protocol/buff"
	"github.com/edgedb/edgedb-go/types"
	"github.com/stretchr/testify/assert"
)

func TestDecodeUUID(t *testing.T) {
	msg := buff.NewMessage([]byte{
		0, 0, 0, 16, // data length
		0, 1, 2, 3, 3, 2, 1, 0, 8, 7, 6, 5, 5, 6, 7, 8,
	})

	codec := &UUID{}

	var result types.UUID
	codec.Decode(msg, unsafe.Pointer(&result))

	expected := types.UUID{0, 1, 2, 3, 3, 2, 1, 0, 8, 7, 6, 5, 5, 6, 7, 8}
	assert.Equal(t, expected, result)
}

func BenchmarkDecodeUUID(b *testing.B) {
	data := []byte{
		0, 0, 0, 16, // data length
		0, 1, 2, 3, 3, 2, 1, 0, 8, 7, 6, 5, 5, 6, 7, 8,
	}
	msg := buff.NewMessage(data)

	var result types.UUID
	ptr := unsafe.Pointer(&result)
	codec := &UUID{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Bts = data
		codec.Decode(msg, ptr)
	}
}

func TestEncodeUUID(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&UUID{}).Encode(buf, types.UUID{
		0, 1, 2, 3, 3, 2, 1, 0,
		8, 7, 6, 5, 5, 6, 7, 8,
	})

	expected := []byte{
		0, 0, 0, 16, // data length
		0, 1, 2, 3, 3, 2, 1, 0, 8, 7, 6, 5, 5, 6, 7, 8,
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func BenchmarkEncodeUUID(b *testing.B) {
	codec := &UUID{}
	id := types.UUID{1, 2, 3, 4, 5, 6, 7, 8, 8, 7, 6, 5, 4, 3, 2, 1}
	buf := buff.NewWriter(make([]byte, 2000))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		codec.Encode(buf, id)
	}
}

func TestDecodeString(t *testing.T) {
	data := []byte{
		0, 0, 0, 5, // data length
		104, 101, 108, 108, 111,
	}
	msg := buff.NewMessage(data)

	var result string
	(&Str{}).Decode(msg, unsafe.Pointer(&result))

	assert.Equal(t, "hello", result)

	// make sure that the string value is not tied to the buffer.
	data[5] = 0
	assert.Equal(t, "hello", result)
}

func BenchmarkDecodeString(b *testing.B) {
	data := []byte{
		0, 0, 0, 5, // data length
		104, 101, 108, 108, 111,
	}
	msg := buff.NewMessage(data)

	var result string
	ptr := unsafe.Pointer(&result)
	codec := &Str{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Bts = data
		codec.Decode(msg, ptr)
	}
}

func TestEncodeString(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&Str{}).Encode(buf, "hello")

	expected := []byte{
		0, 0, 0, 5, // data length
		104, 101, 108, 108, 111,
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func TestDecodeBytes(t *testing.T) {
	data := []byte{
		0, 0, 0, 5, // data length
		104, 101, 108, 108, 111,
	}
	msg := buff.NewMessage(data)

	codec := Bytes{}

	var result []byte
	codec.Decode(msg, unsafe.Pointer(&result))

	expected := []byte{104, 101, 108, 108, 111}

	assert.Equal(t, expected, result)

	// assert that memory is not shared with the buffer
	data[5] = 0
	assert.Equal(t, expected, result)
}

func BenchmarkDecodeBytes(b *testing.B) {
	data := []byte{
		0, 0, 0, 5, // data length
		104, 101, 108, 108, 111,
	}
	msg := buff.NewMessage(data)

	var result []byte
	ptr := unsafe.Pointer(&result)
	codec := &Bytes{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Bts = data
		codec.Decode(msg, ptr)
	}
}

func TestEncodeBytes(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&Bytes{}).Encode(buf, []byte{104, 101, 108, 108, 111})

	expected := []byte{
		0, 0, 0, 5, // data length
		104, 101, 108, 108, 111,
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func TestDecodeInt16(t *testing.T) {
	msg := buff.NewMessage([]byte{
		0, 0, 0, 2, // data length
		0, 7, // int16
	})

	var result int16
	codec := Int16{}
	codec.Decode(msg, unsafe.Pointer(&result))

	assert.Equal(t, int16(7), result)
}

func BenchmarkDecodeInt16(b *testing.B) {
	data := []byte{
		0, 0, 0, 2, // data length
		1, 2, // int16
	}
	msg := buff.NewMessage(data)

	var result int16
	ptr := unsafe.Pointer(&result)
	codec := &Int16{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Bts = data
		codec.Decode(msg, ptr)
	}
}

func TestEncodeInt16(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&Int16{}).Encode(buf, int16(7))

	expected := []byte{
		0, 0, 0, 2, // data length
		0, 7, // int16
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func TestDecodeInt32(t *testing.T) {
	msg := buff.NewMessage([]byte{
		0, 0, 0, 4, // data length
		0, 0, 0, 7, // int32
	})

	var result int32
	(&Int32{}).Decode(msg, unsafe.Pointer(&result))

	assert.Equal(t, int32(7), result)
}

func BenchmarkDecodeInt32(b *testing.B) {
	data := []byte{
		0, 0, 0, 4, // data length
		1, 2, 3, 4, // int32
	}
	msg := buff.NewMessage(data)

	var result int32
	ptr := unsafe.Pointer(&result)
	codec := &Int32{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Bts = data
		codec.Decode(msg, ptr)
	}
}

func TestEncodeInt32(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&Int32{}).Encode(buf, int32(7))

	expected := []byte{
		0, 0, 0, 4, // data length
		0, 0, 0, 7, // int32
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func TestDecodeInt64(t *testing.T) {
	msg := buff.NewMessage([]byte{
		0, 0, 0, 8, // data length
		1, 2, 3, 4, 5, 6, 7, 8, // int64
	})

	var result int64
	(&Int64{}).Decode(msg, unsafe.Pointer(&result))

	assert.Equal(t, int64(72623859790382856), result)
}

func BenchmarkDecodeInt64(b *testing.B) {
	data := []byte{
		0, 0, 0, 8, // data length
		1, 2, 3, 4, 5, 6, 7, 8, // int64
	}
	msg := buff.NewMessage(data)

	var result int64
	ptr := unsafe.Pointer(&result)
	codec := &Int64{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Bts = data
		codec.Decode(msg, ptr)
	}
}

func TestEncodeInt64(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&Int64{}).Encode(buf, int64(27))

	expected := []byte{
		0, 0, 0, 8, // data length
		0, 0, 0, 0, 0, 0, 0, 27, // int64
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func TestDecodeFloat32(t *testing.T) {
	msg := buff.NewMessage([]byte{
		0, 0, 0, 4, // data length
		0xc2, 0, 0, 0,
	})

	var result float32
	codec := &Float32{}
	codec.Decode(msg, unsafe.Pointer(&result))

	assert.Equal(t, float32(-32), result)
}

func TestEncodeFloat32(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&Float32{}).Encode(buf, float32(-32))

	expected := []byte{
		0, 0, 0, 4, // data length
		0xc2, 0, 0, 0,
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func TestDecodeFloat64(t *testing.T) {
	msg := buff.NewMessage([]byte{
		0, 0, 0, 8, // data length
		0xc0, 0x50, 0, 0, 0, 0, 0, 0,
	})

	var result float64
	codec := &Float64{}
	codec.Decode(msg, unsafe.Pointer(&result))

	assert.Equal(t, float64(-64), result)
}

func TestEncodeFloat64(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&Float64{}).Encode(buf, float64(-64))

	expected := []byte{
		0, 0, 0, 8, // data length
		0xc0, 0x50, 0, 0, 0, 0, 0, 0,
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func TestDecodeBool(t *testing.T) {
	msg := buff.NewMessage([]byte{
		0, 0, 0, 1, // data length
		1,
	})

	var result bool
	codec := &Bool{}
	codec.Decode(msg, unsafe.Pointer(&result))

	assert.Equal(t, true, result)
}

func TestEncodeBool(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&Bool{}).Encode(buf, true)

	expected := []byte{
		0, 0, 0, 1, // data length
		1,
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func TestDecodeDateTime(t *testing.T) {
	msg := buff.NewMessage([]byte{
		0, 0, 0, 8, // data length
		0xff, 0xfc, 0xa2, 0xfe, 0xc4, 0xc8, 0x20, 0x0,
	})

	var result time.Time
	codec := &DateTime{}
	codec.Decode(msg, unsafe.Pointer(&result))

	expected := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, result)
}

func TestEncodeDateTime(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&DateTime{}).Encode(buf, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))

	expected := []byte{
		0, 0, 0, 8, // data length
		0xff, 0xfc, 0xa2, 0xfe, 0xc4, 0xc8, 0x20, 0x0,
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func TestDecodeDuration(t *testing.T) {
	msg := buff.NewMessage([]byte{
		0, 0, 0, 0x10, // data length
		0, 0, 0, 0, 0, 0xf, 0x42, 0x40,
		0, 0, 0, 0, // reserved
		0, 0, 0, 0, // reserved
	})

	var result time.Duration
	codec := &Duration{}
	codec.Decode(msg, unsafe.Pointer(&result))

	assert.Equal(t, time.Duration(1_000_000_000), result)
}

func TestEncodeDuration(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&Duration{}).Encode(buf, time.Duration(1_000_000_000))

	expected := []byte{
		0, 0, 0, 0x10, // data length
		0, 0, 0, 0, 0, 0xf, 0x42, 0x40,
		0, 0, 0, 0, // reserved
		0, 0, 0, 0, // reserved
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func TestDecodeJSON(t *testing.T) {
	// todo
	t.SkipNow()

	msg := buff.NewMessage([]byte{
		0, 0, 0, 0x12, // data length
		1, // json format
		0x7b, 0x22, 0x68, 0x65,
		0x6c, 0x6c, 0x6f, 0x22,
		0x3a, 0x22, 0x77, 0x6f,
		0x72, 0x6c, 0x64, 0x22,
		0x7d,
	})

	var result interface{}
	(&JSON{}).Decode(msg, unsafe.Pointer(&result))
	expected := map[string]interface{}{"hello": "world"}

	assert.Equal(t, expected, result)
}

func TestEncodeJSON(t *testing.T) {
	buf := buff.NewWriter(nil)
	(&JSON{}).Encode(buf, map[string]string{"hello": "world"})

	expected := []byte{
		0, 0, 0, 0x12, // data length
		1, // json format
		0x7b, 0x22, 0x68, 0x65,
		0x6c, 0x6c, 0x6f, 0x22,
		0x3a, 0x22, 0x77, 0x6f,
		0x72, 0x6c, 0x64, 0x22,
		0x7d,
	}

	assert.Equal(t, expected, *buf.Unwrap())
}
