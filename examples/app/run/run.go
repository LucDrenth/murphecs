// Demonstrate how to run 1 or multiple apps
package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lucdrenth/murphecs/src/app"
)

func RunApps(subApps ...*app.SubApp) {
	// run apps
	exitChannel := make(chan struct{})
	isDoneChannels := []chan bool{}

	for _, subApp := range subApps {
		isDoneChannel := make(chan bool)
		isDoneChannels = append(isDoneChannels, isDoneChannel)
		go subApp.Run(exitChannel, isDoneChannel)
	}

	// wait for sigterm
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	<-cancelChan
	close(exitChannel)

	// wait for apps to finish
	for _, isDoneChannel := range isDoneChannels {
		<-isDoneChannel
	}
}
