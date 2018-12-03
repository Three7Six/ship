// Copyright 2018 xgfone <xgfone@126.com>
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

package ship

import (
	"time"
)

// NewLoggerMiddleware returns a new logger middleware that will log the request.
//
// By default getTime is time.Now().
func NewLoggerMiddleware(now ...func() time.Time) Middleware {
	_now := time.Now
	if len(now) > 0 && now[0] != nil {
		_now = now[0]
	}

	return MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(ctx Context) (err error) {
			start := _now()
			err = next.Handle(ctx)
			end := _now()

			logger := ctx.Logger()
			if logger != nil {
				req := ctx.Request()
				logger.Info("method=%s, url=%s, starttime=%d, cost=%s, err=%v",
					req.Method, req.URL.RequestURI(), start.Unix(),
					end.Sub(start).String(), err)
			}
			return
		})
	})
}

// NewPanicMiddleware returns a middleware to wrap the panic.
//
// If missing handle, it will use the default, which logs the panic.
func NewPanicMiddleware(handle ...func(Context, interface{})) Middleware {
	handlePanic := HandlePanic
	if len(handle) > 0 && handle[0] != nil {
		handlePanic = handle[0]
	}

	return MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(ctx Context) (err error) {
			defer func() {
				if err := recover(); err != nil {
					handlePanic(ctx, err)
				}
			}()
			err = next.Handle(ctx)
			return
		})
	})
}