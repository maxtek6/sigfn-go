// Copyright (c) 2024 Maxtek Consulting
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sigfn

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleIgnore(t *testing.T) {
	currentSignum := syscall.Signal(0)

	handler := func(signum syscall.Signal) {
		currentSignum = signum
	}
	err := Handle(syscall.SIGINT, handler)
	assert.NoError(t, err)

	signalRaise(syscall.SIGINT)
	assert.Equal(t, currentSignum, syscall.SIGINT)

	currentSignum = syscall.Signal(0)

	err = Ignore(syscall.SIGINT)
	assert.NoError(t, err)

	signalRaise(syscall.SIGINT)
	assert.Zero(t, currentSignum)
}

func TestHandleReset(t *testing.T) {
	currentSignum := syscall.Signal(0)

	handler := func(signum syscall.Signal) {
		currentSignum = signum
	}
	err := Handle(syscall.SIGCHLD, handler)
	assert.NoError(t, err)

	signalRaise(syscall.SIGCHLD)
	assert.Equal(t, currentSignum, syscall.SIGCHLD)

	currentSignum = syscall.Signal(0)

	err = Reset(syscall.SIGCHLD)
	assert.NoError(t, err)

	signalRaise(syscall.SIGCHLD)
	assert.Zero(t, currentSignum, syscall.SIGCHLD)
}

func TestErrors(t *testing.T) {
	badSignum := syscall.Signal(-1)
	cb := func(signum syscall.Signal) {}
	err := Handle(badSignum, cb)
	assert.Error(t, err)
	err = Ignore(badSignum)
	assert.Error(t, err)
	err = Reset(badSignum)
	assert.Error(t, err)
}
