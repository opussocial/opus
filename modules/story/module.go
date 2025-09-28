package story

import (
	"strings"
	"unicode"
	"gitlab.com/pedrokoblitz/opus-go/actions"
	// "gitlab.com/pedrokoblitz/opus-go/quality"
)

type BaseRequest struct {
	Story string
	Definition string
	Page int
	PageSize int
}

type BaseResponse struct {}

func Register(registry *actions.ActionRegistry) {
	registry.RegisterPayload("default:id", func() actions.Payload { return &actions.DefaultID{} })

	registry.RegisterPayload("story", func() actions.Payload { return &Story{} })
	registry.RegisterPayload("story:results", func() actions.Payload { return &StoryResults{} })
	registry.RegisterPayload("definition", func() actions.Payload { return &Definition{} })
	registry.RegisterPayload("role", func() actions.Payload { return &Role{} })
	registry.RegisterPayload("relationship", func() actions.Payload { return &GraphRelationship{} })

	registry.RegisterAction("story:index", IndexStoryAction)
	registry.RegisterAction("story:show", ShowStoryAction)
	registry.RegisterAction("story:create", CreateStoryAction)
	registry.RegisterAction("story:update", UpdateStoryAction)
	registry.RegisterAction("story:delete", DeleteStoryAction)

	registry.RegisterAction("definition:create", CreateDefinitionAction)
	registry.RegisterAction("definition:update", UpdateDefinitionAction)
	registry.RegisterAction("definition:delete", DeleteDefinitionAction)

	registry.RegisterAction("role:create", CreateRoleAction)
	registry.RegisterAction("role:update", UpdateRoleAction)
	registry.RegisterAction("role:delete", DeleteRoleAction)

	registry.RegisterAction("relationship:create", CreateGraphRelationshipAction)
	registry.RegisterAction("relationship:delete", DeleteGraphRelationshipAction)
}

func generateSlug(s string) string {
	var replacements = map[rune]rune{
		'à': 'a', 'á': 'a', 'À': 'a', 'Á': 'a',
		'è': 'e', 'é': 'e', 'È': 'e', 'É': 'e',
		'ì': 'i', 'í': 'i', 'Ì': 'i', 'Í': 'i',
		'ò': 'o', 'ó': 'o', 'Ò': 'o', 'Ó': 'o',
		'ù': 'u', 'ú': 'u', 'Ù': 'u', 'Ú': 'u',
		'â': 'a', 'ê': 'e', 'ô': 'o', 'Â': 'a', 'Ê': 'e', 'Ô': 'o',
		'ã': 'a', 'õ': 'o', 'Ã': 'a', 'Õ': 'o',
		'ç': 'c', 'Ç': 'c',
	}

	// Convert to lower case
	slug := strings.ToLower(s)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Replace accented characters and remove non-alphanumeric characters
	slug = strings.Map(func(r rune) rune {
		if replacement, ok := replacements[r]; ok {
			return replacement
		}
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-' {
			return -1 // Remove the character
		}
		return r
	}, slug)

	return slug
}
