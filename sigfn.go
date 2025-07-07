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
	"os"
	"os/signal"
	"sync"

	_ "modernc.org/libc"
)

var (
	table      *signalTable
	once       sync.Once
	bufferSize = 1
)

type signalTable struct {
	handlerMap map[os.Signal]func(os.Signal)
	sigChan    chan os.Signal
}

func getTable() *signalTable {
	once.Do(func() {
		// Initialize the signal table or any other necessary setup
		table = &signalTable{
			handlerMap: map[os.Signal]func(os.Signal){},
			sigChan:    make(chan os.Signal, bufferSize),
		}
		go table.mainLoop()
	})
	return table
}

func (s *signalTable) mainLoop() {
	for {
		sig := <-s.sigChan
		if handler, ok := s.handlerMap[sig]; ok {
			handler(sig)
		}
	}
}

func (s *signalTable) addHandler(signum os.Signal, handler func(os.Signal)) {
	signal.Notify(s.sigChan, signum)
	s.handlerMap[signum] = handler
}

func (s *signalTable) removeHandler(signum os.Signal) {
	delete(s.handlerMap, signum)
}

// Handle a signal using the handler function
func Handle(signum os.Signal, handler func(os.Signal)) {
	getTable().addHandler(signum, handler)
}

// Ignore a signal
func Ignore(signum os.Signal) {
	getTable().removeHandler(signum)
	signal.Ignore(signum)
}

// Reset a signal to its default behavior
func Reset(signum os.Signal) {
	getTable().removeHandler(signum)
	signal.Reset(signum)
}
