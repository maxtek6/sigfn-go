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
	"fmt"
	"syscall"

	_ "modernc.org/libc"
)

/*
	void signalHook(int);
*/
import "C"

var (
	handlerMap = map[syscall.Signal]func(syscall.Signal){}
)

//export signalHook
func signalHook(signum C.int) {
	signal := syscall.Signal(signum)
	handler := handlerMap[signal]
	handler(signal)
}

// Handle a signal using the handler function
func Handle(signum syscall.Signal, handler func(syscall.Signal)) error {
	handlerMap[signum] = handler
	err := signal(signum, hookType)
	if err != nil {
		return fmt.Errorf("sigfn.Handle: %v", err)
	}
	return nil
}

// Ignore a signal
func Ignore(signum syscall.Signal) error {
	delete(handlerMap, signum)
	err := signal(signum, ignoreType)
	if err != nil {
		return fmt.Errorf("sigfn.Ignore: %v", err)
	}
	return nil
}

// Reset a signal to its default behavior
func Reset(signum syscall.Signal) error {
	delete(handlerMap, signum)
	err := signal(signum, resetType)
	if err != nil {
		return fmt.Errorf("sigfn.Reset: %v", err)
	}
	return nil
}
