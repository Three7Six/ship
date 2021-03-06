// Copyright 2018 xgfone
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

package utils

import (
	"os"
	"sync/atomic"
)

var (
	exited int32
	funcs  = make([]func(), 0)
)

// OnExit registers some exit function.
func OnExit(f ...func()) {
	funcs = append(funcs, f...)
}

// CallOnExit calls the exit functions by the reversed added order.
//
// This function can be called many times.
func CallOnExit() {
	if atomic.CompareAndSwapInt32(&exited, 0, 1) {
		for _len := len(funcs) - 1; _len >= 0; _len-- {
			if f := funcs[_len]; f != nil {
				f()
			}
		}
	}
}

// Exit exits the process with the code, but calling the exit functions
// before exiting.
func Exit(code int) {
	CallOnExit()
	os.Exit(code)
}
