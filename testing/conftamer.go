package testing

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
)

/*
Support for ConfTamer logging of parameters, methods, and messages.
*/

type CType interface {
	// For any fields accessing params via copy or alias,
	// log the key and value of the corresponding param
	// (maybe also for fields set "because of" a param)
	LogCTypeParams()
}

// info on caller of function that called this one
func GetCaller() runtime.Frame {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame
}

func LogCTypesMethodEntry(ctype CType) {
	w := csv.NewWriter(os.Stdout)
	// Log method name
	w.WriteAll([][]string{
		{"ENTER CTYPES METHOD", GetCaller().Func.Name()},
	})

	// Log params
	ctype.LogCTypeParams()
}

func LogCTypesMethodExit() {
	w := csv.NewWriter(os.Stdout)
	w.WriteAll([][]string{
		{"EXIT CTYPES METHOD", GetCaller().Func.Name()},
	})
}

// Log message info: which API call this message corresponds to, and contents
func LogAction(action Action) {
	w := csv.NewWriter(os.Stdout)
	// Message type (API call info)
	api_call_id := fmt.Sprintf("API: k8s.io, TYPE: REQUEST, VERB: %v RESOURCE.SUB: %v.%v NAMESPACE: %v", action.GetVerb(), action.GetResource().String(), action.GetSubresource(), action.GetNamespace())

	w.WriteAll([][]string{
		{"SEND API MESSAGE", api_call_id},
	})

	// Message contents (TODO format as csv: key,value)
	switch concrete_action := action.(type) {
	case ListActionImpl:
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
}
