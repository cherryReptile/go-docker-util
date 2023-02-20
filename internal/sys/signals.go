package sys

import (
	"fmt"
	"os"
	"syscall"
)

func HandleSignal(sigs chan os.Signal) {
	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGINT:
			fmt.Println("\nGood bye")
			os.Exit(0)
		default:
			fmt.Println("Ignoring:", sig)
		}
	}
}
