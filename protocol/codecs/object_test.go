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
	"reflect"
	"runtime/debug"
	"testing"
	"unsafe"

	"github.com/edgedb/edgedb-go/protocol/buff"
	"github.com/edgedb/edgedb-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetObjectType(t *testing.T) {
	type Thing struct {
		Bool  bool       `edgedb:"bool"`
		Small int16      `edgedb:"small"`
		Med   int32      `edgedb:"med"`
		Large int64      `edgedb:"large"`
		Name  string     `edgedb:"name"`
		ID    types.UUID `edgedb:"id"`
	}

	codec := &Object{fields: []*objectField{
		{name: "bool", codec: &Bool{typ: boolType}},
		{name: "small", codec: &Int16{typ: int16Type}},
		{name: "med", codec: &Int32{typ: int32Type}},
		{name: "large", codec: &Int64{typ: int64Type}},
		{name: "name", codec: &Str{typ: strType}},
		{name: "id", codec: &UUID{typ: uuidType}},
	}}

	err := codec.setType(reflect.TypeOf(Thing{}))
	require.Nil(t, err)

	assert.Equal(t, uintptr(0), codec.fields[0].offset)
	assert.Equal(t, uintptr(2), codec.fields[1].offset)
	assert.Equal(t, uintptr(4), codec.fields[2].offset)
	assert.Equal(t, uintptr(8), codec.fields[3].offset)
	assert.Equal(t, uintptr(16), codec.fields[4].offset)
	assert.Equal(t, uintptr(32), codec.fields[5].offset)
}

func TestDecodeObject(t *testing.T) {
	msg := buff.NewMessage([]byte{
		0, 0, 0, 36, // data length
		0, 0, 0, 2, // element count
		// field 0
		0, 0, 0, 0, // reserved
		0, 0, 0, 4, // data length
		102, 111, 117, 114, // utf-8 data
		// field 1
		0, 0, 0, 0, // reserved
		0, 0, 0, 4, // data length
		0, 0, 0, 4, // int32
		// field 2
		0, 0, 0, 0, // reserved
		0xff, 0xff, 0xff, 0xff, // data length (-1)
	})

	type SomeThing struct {
		A string
		B int32
		C int64
	}

	var result SomeThing

	codec := &Object{fields: []*objectField{
		{name: "A", codec: &Str{typ: strType}},
		{name: "B", codec: &Int32{typ: int32Type}},
		{name: "C", codec: &Int64{typ: int64Type}},
	}}
	err := codec.setType(reflect.TypeOf(result))
	require.Nil(t, err)
	codec.Decode(msg, unsafe.Pointer(&result))

	// force garbage collection to be sure that
	// references are durable.
	debug.FreeOSMemory()

	expected := SomeThing{A: "four", B: 4, C: 0}
	assert.Equal(t, expected, result)
}

func BenchmarkDecodeObject(b *testing.B) {
	data := []byte{
		0, 0, 0, 36, // data length
		0, 0, 0, 2, // element count
		// field 0
		0, 0, 0, 0, // reserved
		0, 0, 0, 4, // data length
		102, 111, 117, 114, // utf-8 data
		// field 1
		0, 0, 0, 0, // reserved
		0, 0, 0, 4, // data length
		0, 0, 0, 4, // int32
		// field 2
		0, 0, 0, 0, // reserved
		0xff, 0xff, 0xff, 0xff, // data length (-1)
	}
	msg := buff.NewMessage(data)

	type SomeThing struct {
		A string
		B int32
		C int64
	}

	var result SomeThing
	ptr := unsafe.Pointer(&result)

	codec := &Object{fields: []*objectField{
		{name: "A", codec: &Str{typ: strType}},
		{name: "B", codec: &Int32{typ: int32Type}},
		{name: "C", codec: &Int64{typ: int64Type}},
	}}
	err := codec.setType(reflect.TypeOf(result))
	require.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Bts = data
		codec.Decode(msg, ptr)
	}
}
