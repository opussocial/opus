package elements

import (
	"net/http"
	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

type ElementRequestParams struct {
	Story string
	Definition string
	Element string
}

func (h *ElementRequestParams) FromRequest(r *http.Request) {
	h.Story = r.PathValue("story")
	h.Definition = r.PathValue("definition")
	h.Element = r.PathValue("element")
}

func Register(registry *actions.ActionRegistry) {
	registry.RegisterPayload("element", func() payloads.Payload { return &Element{} })
	registry.RegisterPayload("interaction", func() payloads.Payload { return &Interaction{} })
	registry.RegisterPayload("graph", func() payloads.Payload { return &Graph{} })

	registry.RegisterAction("element:create", CreateElementAction)
	registry.RegisterAction("element:delete", DeleteElementAction)

	registry.RegisterAction("interaction:create", CreateInteractionAction)
	registry.RegisterAction("interaction:delete", DeleteInteractionAction)

	registry.RegisterAction("graph:create", CreateGraphAction)
	registry.RegisterAction("graph:delete", DeleteGraphAction)
}
