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

/*
	#include <signal.h>
	#include <stdlib.h>

	void cgoSignalHook(int signum)
	{
		void signalHook(int);
		signalHook(signum);
	}

	int signalHandle(int signum)
	{
		int result;

		result = EXIT_SUCCESS;

		if(signal(signum, cgoSignalHook) == SIG_ERR)
		{
			result = EXIT_FAILURE;
		}
		return result;
	}

	int signalIgnore(int signum)
	{
		int result;

		result = EXIT_SUCCESS;

		if(signal(signum, SIG_IGN) == SIG_ERR)
		{
			result = EXIT_FAILURE;
		}
		return result;
	}

	int signalReset(int signum)
	{
		int result;

		result = EXIT_SUCCESS;

		if(signal(signum, SIG_DFL) == SIG_ERR)
		{
			result = EXIT_FAILURE;
		}
		return result;
	}

	void signalRaise(int signum)
	{
		(void)raise(signum);
	}
*/
import "C"
import (
	"errors"
	"syscall"
)

type handlerType int

const (
	hookType   handlerType = 0
	ignoreType handlerType = 1
	resetType  handlerType = 2
)

func signal(signum syscall.Signal, handler handlerType) error {
	target := C.int(signum)
	status := C.int(0)
	switch handler {
	case hookType:
		status = C.signalHandle(target)
	case ignoreType:
		status = C.signalIgnore(target)
	case resetType:
		status = C.signalReset(target)
	default:
		return errors.New("invalid handler type")
	}
	if status != 0 {
		return errors.New("signal() failed")
	}
	return nil
}

func signalRaise(signum syscall.Signal) {
	C.signalRaise(C.int(signum))
}
