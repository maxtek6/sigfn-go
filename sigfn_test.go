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
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func signalRaise(signum syscall.Signal) {
	syscall.Kill(os.Getpid(), signum)
}

func collectSignal(sigChan chan os.Signal) (os.Signal, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()
	select {
	case sig := <-sigChan:
		return sig, true
	case <-ctx.Done():
		return nil, false
	}
}
func TestHandleIgnore(t *testing.T) {
	sigChan := make(chan os.Signal)
	Handle(syscall.SIGINT, func(sig os.Signal) {
		sigChan <- sig
	})
	signalRaise(syscall.SIGINT)
	sig, ok := collectSignal(sigChan)
	assert.Equal(t, sig.String(), syscall.SIGINT.String())
	assert.True(t, ok)

	Ignore(syscall.SIGINT)
	signalRaise(syscall.SIGINT)
	sig, ok = collectSignal(sigChan)
	assert.Nil(t, sig)
	assert.False(t, ok)
}

func TestHandleReset(t *testing.T) {
	sigChan := make(chan os.Signal)
	Handle(syscall.SIGCHLD, func(sig os.Signal) {
		sigChan <- sig
	})
	signalRaise(syscall.SIGCHLD)
	sig, ok := collectSignal(sigChan)
	assert.Equal(t, sig.String(), syscall.SIGCHLD.String())
	assert.True(t, ok)

	Reset(syscall.SIGCHLD)
	signalRaise(syscall.SIGCHLD)
	sig, ok = collectSignal(sigChan)
	assert.Nil(t, sig)
	assert.False(t, ok)
}
