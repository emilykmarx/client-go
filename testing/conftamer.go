package testing

import (
	"strings"

	"github.com/emilykmarx/conftamer/pkg/apimessages"
	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/rpc2"
)

// XXX this is specific to k8s - plan is for module dev to import the lib containing their msgInfo here,
// and pass name of msgInfo func
func GetMessageInfo(client *rpc2.RPCClient, bp *api.Breakpoint, goroutine int64) (*apimessages.APICallID, error) {
	loadcfg := api.LoadConfig{FollowPointers: true, MaxVariableRecurse: 10, MaxStringLen: 10, MaxArrayValues: 1, MaxStructFields: -1}
	scope := api.EvalScope{GoroutineID: goroutine}
	verb, err := client.EvalVariable(scope, "action.ActionImpl.Verb", loadcfg)
	if err != nil {
		return nil, err
	}
	resource, err := client.EvalVariable(scope, "action.ActionImpl.Resource.Resource", loadcfg)
	if err != nil {
		return nil, err
	}

	api_call_id := apimessages.APICallID{
		API:  "k8s.io",
		Verb: strings.ToUpper(verb.Value),
		// Unsure yet if resource group and version should be part of the resource used in API call ID, or just the resource type (e.g. pods)
		Resource: resource.Value,
		// TODO will this ever be called for responses?
		APIMessageType: apimessages.Request,
	}

	// XXX get contents
	return &api_call_id, nil
}
