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

func TestTupleSetType(t *testing.T) {
	codec := &Tuple{fields: []Codec{
		&Int64{typ: int64Type},
		&Int32{typ: int32Type},
	}}
	err := codec.setType(reflect.TypeOf([]interface{}{}))
	require.Nil(t, err)

	assert.Equal(t, 16, codec.step)
}

func TestDecodeTuple(t *testing.T) {
	buf := buff.New([]byte{
		0,
		0, 0, 0, 40,
		0, 0, 0, 32, // data length
		0, 0, 0, 2, // number of elements
		// element 0
		0, 0, 0, 0, // reserved
		0, 0, 0, 8, // data length
		0, 0, 0, 0, 0, 0, 0, 2,
		// element 1
		0, 0, 0, 0, // reserved
		0, 0, 0, 4, // data length
		0, 0, 0, 3,
	})
	buf.Next()

	var result []interface{}

	codec := &Tuple{fields: []Codec{
		&Int64{typ: int64Type},
		&Int32{typ: int32Type},
	}}
	err := codec.setType(reflect.TypeOf(result))
	require.Nil(t, err)
	codec.Decode(buf, unsafe.Pointer(&result))

	// force garbage collection to be sure that
	// references are durable.
	debug.FreeOSMemory()

	expected := []interface{}{int64(2), int32(3)}
	assert.Equal(t, expected, result)
}

func TestEncodeNullTuple(t *testing.T) {
	buf := buff.New(nil)
	buf.BeginMessage(0xff)
	(&Tuple{}).Encode(buf, []interface{}{})
	buf.EndMessage()

	expected := []byte{
		0xff,         // message type
		0, 0, 0, 0xc, // message length
		0, 0, 0, 4, // data length
		0, 0, 0, 0, // number of elements
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func TestEncodeTuple(t *testing.T) {
	buf := buff.New(nil)
	buf.BeginMessage(0xff)

	codec := &Tuple{fields: []Codec{&Int64{}, &Int64{}}}
	codec.Encode(buf, []interface{}{int64(2), int64(3)})
	buf.EndMessage()

	expected := []byte{
		0xff,          // message type
		0, 0, 0, 0x2c, // message length
		0, 0, 0, 36, // data length
		0, 0, 0, 2, // number of elements
		// element 0
		0, 0, 0, 0, // reserved
		0, 0, 0, 8, // data length
		0, 0, 0, 0, 0, 0, 0, 2,
		// element 1
		0, 0, 0, 0, // reserved
		0, 0, 0, 8, // data length
		0, 0, 0, 0, 0, 0, 0, 3,
	}

	assert.Equal(t, expected, *buf.Unwrap())
}

func BenchmarkEncodeTuple(b *testing.B) {
	codec := Tuple{fields: []Codec{&UUID{}}}
	id := types.UUID{1, 2, 3, 4, 5, 6, 7, 8, 8, 7, 6, 5, 4, 3, 2, 1}
	ids := []interface{}{id}
	buf := buff.New(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf.BeginMessage(0)
		codec.Encode(buf, ids)
	}
}