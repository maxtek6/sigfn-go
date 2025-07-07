package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/maxtek6/sigfn-go"
)

func main() {
	done := make(chan struct{})
	handleSignal := func(sig os.Signal) {
		fmt.Printf("Received signal: %s\n", sig)
		done <- struct{}{}
	}
	sigfn.Handle(syscall.SIGINT, handleSignal)
	fmt.Println("Paused. Press Ctrl+C to continue...")
	<-done
}
