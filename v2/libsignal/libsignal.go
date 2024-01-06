package libsignal

import (
	"os"
	"os/signal"
	"sync"

	"github.com/helloferdie/golib/v2/liblogger"
)

// WaitCTRLC to stop console program with ctrl + C
func WaitCTRLC(doLog bool) {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var ch chan os.Signal
	ch = make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch

		if doLog {
			liblogger.Infow("Program stopped")
		}
		endWaiter.Done()
	}()
	endWaiter.Wait()
}
