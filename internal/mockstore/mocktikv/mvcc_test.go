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
// https://github.com/pingcap/tidb/tree/cc5e161ac06827589c4966674597c137cc9e809c/store/tikv/mockstore/mocktikv/mvcc_test.go
//

// Copyright 2018-present, PingCAP, Inc.
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

package mocktikv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegionContains(t *testing.T) {
	assert := assert.New(t)
	assert.True(regionContains([]byte{}, []byte{}, []byte{}))
	assert.True(regionContains([]byte{}, []byte{}, []byte{1}))
	assert.False(regionContains([]byte{1, 1, 1}, []byte{}, []byte{1, 1, 0}))
	assert.True(regionContains([]byte{1, 1, 1}, []byte{}, []byte{1, 1, 1}))
	assert.True(regionContains([]byte{}, []byte{2, 2, 2}, []byte{2, 2, 1}))
	assert.False(regionContains([]byte{}, []byte{2, 2, 2}, []byte{2, 2, 2}))
	assert.False(regionContains([]byte{1, 1, 1}, []byte{2, 2, 2}, []byte{1, 1, 0}))
	assert.True(regionContains([]byte{1, 1, 1}, []byte{2, 2, 2}, []byte{1, 1, 1}))
	assert.True(regionContains([]byte{1, 1, 1}, []byte{2, 2, 2}, []byte{2, 2, 1}))
	assert.False(regionContains([]byte{1, 1, 1}, []byte{2, 2, 2}, []byte{2, 2, 2}))
}
