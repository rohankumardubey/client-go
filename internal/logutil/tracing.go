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
// https://github.com/pingcap/tidb/tree/cc5e161ac06827589c4966674597c137cc9e809c/store/tikv/logutil/tracing.go
//

// Copyright 2021 PingCAP, Inc.
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

package logutil

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// TraceEventKey presents the TraceEventKey in span log.
var TraceEventKey = "event"

// Event records event in current tracing span.
func Event(ctx context.Context, event string) {
	if span := opentracing.SpanFromContext(ctx); span != nil && span.Tracer() != nil {
		span.LogFields(log.String(TraceEventKey, event))
	}
}

// Eventf records event in current tracing span with format support.
func Eventf(ctx context.Context, format string, args ...interface{}) {
	if span := opentracing.SpanFromContext(ctx); span != nil && span.Tracer() != nil {
		span.LogFields(log.String(TraceEventKey, fmt.Sprintf(format, args...)))
	}
}

// SetTag sets tag kv-pair in current tracing span
func SetTag(ctx context.Context, key string, value interface{}) {
	if span := opentracing.SpanFromContext(ctx); span != nil && span.Tracer() != nil {
		span.SetTag(key, value)
	}
}
