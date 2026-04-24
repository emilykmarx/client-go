package testing

import (
	"strings"

	"github.com/emilykmarx/conftamer"
)

/* Call conftamer's message log function for this Action. */
func LogAction(action Action) {
	api_call_id := conftamer.APICallID{
		API:  "k8s.io",
		Verb: strings.ToUpper(action.GetVerb()),
		// Unsure yet if resource group and version should be part of the resource used in API call ID, or just the resource type (e.g. pods)
		Resource: action.GetResource().String(),
		// TODO will this ever be called for responses?
		APIMessageType: conftamer.Request,
	}

	// Remove Action fields that are part of API call ID
	exclude := map[string]struct{}{"Verb": {}, "Resource": {}}
	// The API probably does fancier serialization than this, but doing it generically seems fine for now
	msg_contents := conftamer.ParseJSONFields(action, exclude)
	conftamer.LogAPIMessage(api_call_id, msg_contents)
}
