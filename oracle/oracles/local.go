// Copyright 2021 TiKV Authors
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

// NOTE: The code in this file is based on code from the
// TiDB project, licensed under the Apache License v 2.0
//
// https://github.com/pingcap/tidb/tree/cc5e161ac06827589c4966674597c137cc9e809c/store/tikv/oracle/oracles/local.go
//

// Copyright 2016 PingCAP, Inc.
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

package oracles

import (
	"context"
	"sync"
	"time"

	"github.com/tikv/client-go/v2/oracle"
)

var _ oracle.Oracle = &localOracle{}

type localOracle struct {
	sync.Mutex
	lastTimeStampTS uint64
	n               uint64
	hook            *struct {
		currentTime time.Time
	}
}

// NewLocalOracle creates an Oracle that uses local time as data source.
func NewLocalOracle() oracle.Oracle {
	return &localOracle{}
}

func (l *localOracle) IsExpired(lockTS, TTL uint64, _ *oracle.Option) bool {
	now := time.Now()
	if l.hook != nil {
		now = l.hook.currentTime
	}
	expire := oracle.GetTimeFromTS(lockTS).Add(time.Duration(TTL) * time.Millisecond)
	return !now.Before(expire)
}

func (l *localOracle) GetTimestamp(ctx context.Context, _ *oracle.Option) (uint64, error) {
	l.Lock()
	defer l.Unlock()
	now := time.Now()
	if l.hook != nil {
		now = l.hook.currentTime
	}
	ts := oracle.GoTimeToTS(now)
	if l.lastTimeStampTS == ts {
		l.n++
		return ts + l.n, nil
	}
	l.lastTimeStampTS = ts
	l.n = 0
	return ts, nil
}

func (l *localOracle) GetTimestampAsync(ctx context.Context, _ *oracle.Option) oracle.Future {
	return &future{
		ctx: ctx,
		l:   l,
	}
}

func (l *localOracle) GetLowResolutionTimestamp(ctx context.Context, opt *oracle.Option) (uint64, error) {
	return l.GetTimestamp(ctx, opt)
}

func (l *localOracle) GetLowResolutionTimestampAsync(ctx context.Context, opt *oracle.Option) oracle.Future {
	return l.GetTimestampAsync(ctx, opt)
}

// GetStaleTimestamp return physical
func (l *localOracle) GetStaleTimestamp(ctx context.Context, txnScope string, prevSecond uint64) (ts uint64, err error) {
	return oracle.GoTimeToTS(time.Now().Add(-time.Second * time.Duration(prevSecond))), nil
}

type future struct {
	ctx context.Context
	l   *localOracle
}

func (f *future) Wait() (uint64, error) {
	return f.l.GetTimestamp(f.ctx, &oracle.Option{})
}

// UntilExpired implement oracle.Oracle interface.
func (l *localOracle) UntilExpired(lockTimeStamp, TTL uint64, opt *oracle.Option) int64 {
	now := time.Now()
	if l.hook != nil {
		now = l.hook.currentTime
	}
	return oracle.ExtractPhysical(lockTimeStamp) + int64(TTL) - oracle.GetPhysical(now)
}

func (l *localOracle) Close() {
}
