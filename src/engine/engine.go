package engine

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/engine/apps/core"
	"github.com/lucdrenth/murph_engine/src/engine/apps/renderer"
	"github.com/lucdrenth/murph_engine/src/log"
)

// Engine is the main struct of a Murph application.
type Engine struct {
	subApps       []app.SubApp
	logger        log.Logger
	exitChannel   chan struct{}
	isDoneChannel chan bool
}

func Empty() Engine {
	logger := log.NoOp()

	return Engine{
		subApps:       []app.SubApp{},
		logger:        &logger,
		exitChannel:   make(chan struct{}),
		isDoneChannel: make(chan bool),
	}
}

func Default() Engine {
	logger := log.Console()

	coreApp := core.New(&logger)
	rendererApp := renderer.New(&logger)

	return Engine{
		subApps: []app.SubApp{
			&coreApp,
			&rendererApp,
		},
		logger:        &logger,
		exitChannel:   make(chan struct{}),
		isDoneChannel: make(chan bool),
	}
}

func (e *Engine) AddSubApp(app app.SubApp) {
	e.subApps = append(e.subApps, app)
}

func (e *Engine) SetLogger(logger log.Logger) {
	if logger != nil {
		e.logger = logger
	}
}

func (e *Engine) Run() {
	e.logger.Debug("Running Murph application")

	for _, app := range e.subApps {
		go app.Run(e.exitChannel, e.isDoneChannel)
	}

	quitSignal := waitForSigterm()
	e.logger.Debug(fmt.Sprintf("Received signal to quit: %s. Starting graceful shutdown", quitSignal))
	close(e.exitChannel)
	e.waitForAppsCleanup()
}

func (e *Engine) waitForAppsCleanup() {
	for range e.subApps {
		<-e.isDoneChannel
	}
}

func waitForSigterm() os.Signal {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	return sig
}
