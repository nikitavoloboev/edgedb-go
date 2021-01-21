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

package edgedb

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectPool(t *testing.T) {
	ctx := context.Background()
	p, err := Connect(ctx, opts)
	require.Nil(t, err)

	var result string
	err = p.QueryOne(ctx, "SELECT 'hello';", &result)
	assert.Nil(t, err)
	assert.Equal(t, "hello", result)

	err = p.Close()
	assert.Nil(t, err)
}

func TestConnectPoolZeroMinAndMaxConns(t *testing.T) {
	o := opts
	o.MinConns = 0
	o.MaxConns = 0

	ctx := context.Background()
	p, err := Connect(ctx, o)
	require.Nil(t, err)

	require.Equal(t, defaultMinConns, p.(*pool).minConns)
	require.Equal(t, defaultMaxConns, p.(*pool).maxConns)

	var result string
	err = p.QueryOne(ctx, "SELECT 'hello';", &result)
	assert.Nil(t, err)
	assert.Equal(t, "hello", result)

	err = p.Close()
	assert.Nil(t, err)
}

func TestClosePoolConcurently(t *testing.T) {
	ctx := context.Background()
	p, err := Connect(ctx, opts)
	require.Nil(t, err)

	errs := make(chan error)
	go func() { errs <- p.Close() }()
	go func() { errs <- p.Close() }()

	assert.Nil(t, <-errs)
	var closedErr InterfaceError
	assert.True(t, errors.As(<-errs, &closedErr))
}

func TestConnectPoolMinConnLteMaxConn(t *testing.T) {
	ctx := context.Background()
	o := Options{MinConns: 5, MaxConns: 1}
	_, err := Connect(ctx, o)
	assert.EqualError(
		t,
		err,
		"edgedb.ConfigurationError: "+
			"MaxConns (1) may not be less than MinConns (5)",
	)

	var expected ConfigurationError
	assert.True(t, errors.As(err, &expected))
}

func TestAcquireFromClosedPool(t *testing.T) {
	p := &pool{
		isClosed:       true,
		freeConns:      make(chan *baseConn),
		potentialConns: make(chan struct{}),
	}

	conn, err := p.Acquire(context.Background())
	var closedErr InterfaceError
	require.True(t, errors.As(err, &closedErr))
	assert.Nil(t, conn)
}

func TestAcquireFreeConnFromPool(t *testing.T) {
	conn := &baseConn{}
	p := &pool{freeConns: make(chan *baseConn, 1)}
	p.freeConns <- conn

	result, err := p.Acquire(context.Background())
	assert.Nil(t, err)

	pConn, ok := result.(*poolConn)
	require.True(t, ok, "unexpected return type: %T", result)
	assert.Equal(t, conn, pConn.baseConn)
}

func BenchmarkPoolAcquireRelease(b *testing.B) {
	p := &pool{
		maxConns:       2,
		minConns:       2,
		freeConns:      make(chan *baseConn, 2),
		potentialConns: make(chan struct{}, 2),
	}

	for i := 0; i < p.maxConns; i++ {
		p.freeConns <- &baseConn{}
	}

	var conn *baseConn
	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn, _ = p.acquire(ctx)
		_ = p.release(conn, nil)
	}
}

func TestAcquirePotentialConnFromPool(t *testing.T) {
	p, err := Connect(context.Background(), opts)
	require.Nil(t, err)
	defer func() {
		assert.Nil(t, p.Close())
	}()

	// free connection
	a, err := p.Acquire(context.Background())
	require.Nil(t, err)
	require.NotNil(t, a)
	defer func() { assert.Nil(t, a.Release()) }()

	// potential connection
	b, err := p.Acquire(context.Background())
	require.Nil(t, err)
	require.NotNil(t, b)
	defer func() { assert.Nil(t, b.Release()) }()
}

func TestPoolAcquireExpiredContext(t *testing.T) {
	p := &pool{
		freeConns:      make(chan *baseConn, 1),
		potentialConns: make(chan struct{}, 1),
	}
	p.freeConns <- &baseConn{}
	p.potentialConns <- struct{}{}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	cancel()

	conn, err := p.Acquire(ctx)
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
	assert.Nil(t, conn)
}

func TestPoolAcquireThenContextExpires(t *testing.T) {
	p := &pool{}

	deadline := time.Now().Add(10 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	conn, err := p.Acquire(ctx)
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
	assert.Nil(t, conn)
	cancel()
}

func TestClosePool(t *testing.T) {
	p := &pool{
		maxConns:       0,
		minConns:       0,
		freeConns:      make(chan *baseConn),
		potentialConns: make(chan struct{}),
	}

	err := p.Close()
	assert.Nil(t, err)

	err = p.Close()
	var closedErr InterfaceError
	assert.True(t, errors.As(err, &closedErr))
}

func TestPoolRetry(t *testing.T) {
	ctx := context.Background()

	p, err := Connect(ctx, opts)
	require.Nil(t, err, "unexpected error: %v", err)
	defer p.Close() // nolint:errcheck

	var result int64
	err = p.Retry(ctx, func(ctx context.Context, tx Tx) error {
		return tx.QueryOne(ctx, "SELECT 33*21", &result)
	})

	require.Nil(t, err, "unexpected error: %v", err)
	require.Equal(t, int64(693), result, "Pool.Retry() failed")
}