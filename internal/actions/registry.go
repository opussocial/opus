package actions

import (
	"sync"
	// "context"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
)

type ActionRegistry struct {
	actionRegistry  map[string]ActionFunc
	payloadRegistry map[string]func() payloads.Payload
	regMu           sync.RWMutex
}

func NewActionRegistry() *ActionRegistry {
	return &ActionRegistry{
		actionRegistry:  make(map[string]ActionFunc),
		payloadRegistry: make(map[string]func() payloads.Payload),
	}
}

func (r *ActionRegistry) Metrics() (int, int) {
	return len(r.actionRegistry), len(r.payloadRegistry)
}

// RegisterAction binds a string name to an action function
func (r *ActionRegistry) RegisterAction(name string, fn ActionFunc) {
	r.regMu.Lock()
	defer r.regMu.Unlock()
	r.actionRegistry[name] = fn
}

// Register binds a string name to a payload constructor
func (r *ActionRegistry) RegisterPayload(name string, ctor func() payloads.Payload) {
	r.regMu.Lock()
	defer r.regMu.Unlock()
	r.payloadRegistry[name] = ctor
}

// ResolveAction retrieves an action by name
func (r *ActionRegistry) ResolveAction(name string) (ActionFunc, bool) {
	r.regMu.RLock()
	defer r.regMu.RUnlock()
	fn, ok := r.actionRegistry[name]
	return fn, ok
}

// ResolvePayload retrieves a payload constructor by name
func (r *ActionRegistry) ResolvePayload(name string) (func() payloads.Payload, bool) {
	r.regMu.RLock()
	defer r.regMu.RUnlock()
	ctor, ok := r.payloadRegistry[name]
	return ctor, ok
}
