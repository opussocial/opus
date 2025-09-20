package actions

import (
	"fmt"
	"sync"

	"gitlab.com/pedrokoblitz/opus-go/quality"
)

type Executable interface {
	Execute() error
}

type Action struct {
	name    string
	payload Payload
	container Container
	execute ActionFunc
}

type ActionFunc func(p Payload, container Container) error

func NewAction(
	name string,
	payload Payload,
	container Container,
	execute ActionFunc,
) *Action {
	return &Action{
		name:    name,
		payload: payload,
		container: container,
		execute: execute,
	}
}

func (a *Action) Execute() error {
	quality.LogStart(a.name, a.payload)

	err := a.execute(a.payload, a.container)
	if err != nil {
		quality.LogError(a.name, err)
		return err
	}

	quality.LogSuccess(a.name, a.payload)
	return nil
}

type ActionRunner struct {
	Actions []Executable
}

func NewActionRunner(actions []Executable) *ActionRunner {
	return &ActionRunner{
		Actions: actions,
	}
}

func (r *ActionRunner) Execute() error {
	errors := make(chan error, len(r.Actions))
	var wg sync.WaitGroup

	for _, action := range r.Actions {
		wg.Add(1)
		go func(action Executable) {
			defer wg.Done()
			if err := action.Execute(); err != nil {
				errors <- fmt.Errorf("action failed: %w", err)
			}
		}(action)
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	for err := range errors {
		return err
	}
	return nil
}
