package testing

import (
	"fmt"
	"strings"

	"github.com/emilykmarx/conftamer"
)

/* Call conftamer's message log function for this Action. */
func LogAction(action Action) {
	api_call_id := conftamer.APICallID{
		API:      "k8s.io",
		Verb:     strings.ToUpper(action.GetVerb()),
		Resource: action.GetResource().Resource,
		// TODO will this ever be called for responses?
		APIMessageType: conftamer.Request,
	}

	// Contents that apply to all message types
	// TODO should subresource be contents or type?
	msg_contents := []conftamer.MsgField{
		{Key: "namespace", Value: action.GetNamespace()},
		{Key: "subresource", Value: action.GetSubresource()},
	}

	switch concrete_action := action.(type) {
	case ListActionImpl:
		// LEFT OFF append contents for each type (for now - there may be a more generic way of doing this)
		// ListRestrictions is redundant with ListOptions
		fmt.Printf("%+v\n", concrete_action)
	case GetActionImpl:
		fmt.Printf("%+v\n", concrete_action)
	case CreateActionImpl:
		fmt.Printf("%+v\n", concrete_action)
	case UpdateActionImpl:
		fmt.Printf("%+v\n", concrete_action)
	case DeleteActionImpl:
		fmt.Printf("%+v\n", concrete_action)
	case PatchActionImpl:
		fmt.Printf("%+v\n", concrete_action)
	case DeleteCollectionActionImpl:
		fmt.Printf("%+v\n", concrete_action)
	case WatchActionImpl:
		fmt.Printf("%+v\n", concrete_action)
	case ProxyGetActionImpl:
		fmt.Printf("%+v\n", concrete_action)
	}

	conftamer.LogAPIMessage(api_call_id, msg_contents)
}
