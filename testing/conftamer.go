package testing

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"k8s.io/apimachinery/pkg/util/json"
)

/*
Support for ConfTamer logging of parameters, methods, and messages.
*/

const (
	methodEntryLog = "ENTER CTYPES METHOD"
	methodExitLog  = "EXIT CTYPES METHOD"
	apiMessageLog  = "SEND API MESSAGE"
)

type CTypeParams struct {
	Key   string
	Value string
}

type CType interface {
	// For any fields accessing params via copy or alias,
	// return the key and value of the corresponding params
	// (maybe also for fields set "because of" a param)
	CTypeParams() []CTypeParams
}

// info on caller of function that called this one
func GetCaller() runtime.Frame {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame
}

// Log method name and params
func LogCTypesMethodEntry(ctype CType) {
	w := csv.NewWriter(os.Stdout)

	params, err := json.Marshal(ctype.CTypeParams())
	if err != nil {
		log.Panicf("marshaling %v: %v\n", ctype.CTypeParams(), err.Error())
	}
	w.WriteAll([][]string{
		{methodEntryLog, GetCaller().Func.Name(),
			string(params)},
	})
}

func LogCTypesMethodExit() {
	w := csv.NewWriter(os.Stdout)
	w.WriteAll([][]string{
		{methodExitLog, GetCaller().Func.Name()},
	})
}

// Log message info: which API call this message corresponds to, and contents
func LogAction(action Action) {
	w := csv.NewWriter(os.Stdout)
	// Message type (API call info)
	api_call_id := fmt.Sprintf("API: k8s.io, TYPE: REQUEST, VERB: %v RESOURCE.SUB: %v.%v NAMESPACE: %v", action.GetVerb(), action.GetResource().String(), action.GetSubresource(), action.GetNamespace())

	w.WriteAll([][]string{
		{apiMessageLog, api_call_id},
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

type DataFlow struct {
	paramKey string
	msgField string
}
type testMethod struct {
	test   string
	method string
}

// A CType method and corresponding params
type MethodParams struct {
	method string
	params []CTypeParams
}
type APIMessageInfo struct {
	controlFlow map[string][]testMethod   // param key => tests that found CF from param to msg
	dataFlow    map[DataFlow][]testMethod // {param key, msg field} => tests that found DF from param to msg field
}

// Taint info for each msg gathered across all tests (API call ID => influence)
type AllTaint map[string]APIMessageInfo

// Eventually may want something more graphable
func (m *AllTaint) prettyPrint(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	/*
			// TODO
		for api_call_id, info := range m {
				 format:
				 	API call
						CF
							Param key
								Test/method
									Param value
								<repeat>
						DF
							Msg field
								Param key
									Test/method
										Param value (same as field value)
									<repeat>
	*/
	return nil
}

func (m *AllTaint) add(test string, api_call_id string, params []MethodParams) {
	// Can't edit map value in place => get it (initializing its maps if needed) and put it back
	existing_flow := (*m)[api_call_id]
	if existing_flow.controlFlow == nil {
		existing_flow.controlFlow = make(map[string][]testMethod)
	}
	if existing_flow.dataFlow == nil {
		existing_flow.dataFlow = make(map[DataFlow][]testMethod)
	}

	for _, methodParam := range params {
		for _, param := range methodParam.params {

			/* LEFT OFF learn about them selectors confirm this is right -
			// Treat selectors.role as a prefix to name the label and field by the role they correspond to, not as its own param
			*/
			// CF: Msg is CF-tainted by all params
			existing_flow.controlFlow[param.Key] = append(existing_flow.controlFlow[param.Key],
				testMethod{test, methodParam.method})

			// DF: Msg field is DF-tainted by any params whose content match the field
			// TODO compare params to each message field
		}

	}
	(*m)[api_call_id] = existing_flow
}

func ParseTestOutput(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	r := csv.NewReader(file)
	// Allow variable number of fields
	r.FieldsPerRecord = -1

	// Stack of in-scope methods and their params (last element = most recent)
	cur_ctype_params := []MethodParams{}
	msg_taint := make(AllTaint)
	cur_test := ""

	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(row)
				continue
			}
		}

		// Enter method
		if row[0] == methodEntryLog {
			// TODO nested methods: if same CType, doesn't matter - but if different CTypes, take union of params?
			params := []CTypeParams{}
			err := json.Unmarshal([]byte(row[2]), &params)
			if err != nil {
				log.Panicf("unmarshaling %v: %v\n", row[2], err.Error())
			}
			cur_ctype_params = append(cur_ctype_params, MethodParams{method: row[1], params: params})
		} else if row[0] == methodExitLog {
			cur_ctype_params = cur_ctype_params[:len(cur_ctype_params)-1]
		} else if row[0] == apiMessageLog {
			msg_taint.add(cur_test, row[1], cur_ctype_params)
		} else if strings.HasPrefix(row[0], "=== RUN") {
			fields := strings.Fields(row[0])
			test := fields[len(fields)-1]
			cur_test = test
		} else {
			// Some other test log
		}
	}

	return nil
}
