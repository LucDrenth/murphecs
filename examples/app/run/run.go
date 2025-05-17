// Demonstrate how to run an app
package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lucdrenth/murphecs/src/app"
)

func RunApp(subApp *app.SubApp) {
	// run app
	exitChannel := make(chan struct{})
	isDoneChannel := make(chan bool)
	go subApp.Run(exitChannel, isDoneChannel)

	// wait for sigterm
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	<-cancelChan
	close(exitChannel)

	// wait for app to finish
	<-isDoneChannel
}
