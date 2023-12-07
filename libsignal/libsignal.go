package libsignal

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/helloferdie/golib/liblogger"
)

// WaitCTRLC -
func WaitCTRLC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var ch chan os.Signal
	ch = make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch

		s := "Program stopped"
		fmt.Println(s)
		liblogger.Log(nil, false).Info(s)

		endWaiter.Done()
	}()
	endWaiter.Wait()
}
