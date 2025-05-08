package engine

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/engine/apps/core"
	"github.com/lucdrenth/murph_engine/src/engine/apps/renderer"
	"github.com/lucdrenth/murph_engine/src/log"
)

const (
	AppIDCore app.ID = iota
	AppIDRenderer
)

// Engine is the main struct of a Murph application.
type Engine struct {
	apps          map[app.ID]app.SubApp
	logger        log.Logger
	exitChannel   chan struct{}
	isDoneChannel chan bool
}

func Empty() Engine {
	logger := log.NoOp()

	return Engine{
		apps:          map[app.ID]app.SubApp{},
		logger:        &logger,
		exitChannel:   make(chan struct{}),
		isDoneChannel: make(chan bool),
	}
}

func Default() (Engine, error) {
	logger := log.Console()

	coreApp, err := core.New(&logger)
	if err != nil {
		return Engine{}, fmt.Errorf("failed to create core app: %w", err)
	}

	renderApp, err := renderer.New(&logger)
	if err != nil {
		return Engine{}, fmt.Errorf("failed to create renderer app: %w", err)
	}

	engine := Empty()
	engine.SetLogger(&logger)
	engine.AddSubApp(&coreApp, AppIDCore)
	engine.AddSubApp(&renderApp, AppIDRenderer)
	return engine, nil
}

func (e *Engine) AddSubApp(app app.SubApp, id app.ID) {
	if _, exists := e.apps[id]; exists {
		e.logger.Error(fmt.Sprintf("failed to add sub app %s: already exists", reflect.TypeOf(app).String()))
		return
	}

	e.apps[id] = app
}

// SetLogger either sets a new logger. Uses log.NoOp when logger is nil to prevent crashes.
func (e *Engine) SetLogger(logger log.Logger) {
	if logger == nil {
		noOpLogger := log.NoOp()
		e.logger = &noOpLogger
	} else {
		e.logger = logger
	}
}

func (e *Engine) App(id app.ID) app.SubApp {
	return e.apps[id]
}

func (e *Engine) Run() {
	e.logger.Debug("Running Murph application")

	for _, app := range e.apps {
		go app.Run(e.exitChannel, e.isDoneChannel)
	}

	quitSignal := waitForSigterm()
	e.logger.Debug(fmt.Sprintf("Received signal to quit: %s. Starting graceful shutdown", quitSignal))
	close(e.exitChannel)
	e.waitForAppsCleanup()
}

func (e *Engine) waitForAppsCleanup() {
	for range e.apps {
		<-e.isDoneChannel
	}
}

func waitForSigterm() os.Signal {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	return sig
}
