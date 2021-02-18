package api

import "errors"

// ResultStatus result of a command
type ResultStatus string

var (
	// ResultStatusSuccess command completed without errors
	ResultStatusSuccess ResultStatus = "Success"
	// ResultStatusFailure command completed with errors
	ResultStatusFailure ResultStatus = "Failed"
)

// ResultMessage is amessage that contains a result from the command. The struture is a key value, where value is a string
type ResultMessage map[string]string

// NewResultFromKeyPair build a result message from key value pairs
func NewResultFromKeyPair(kp ...string) (*ForgeOpsResult, error) {
	result := ResultMessage{}
	// make sure we are given key pairs
	if len(kp)%2 != 0 {
		return &ForgeOpsResult{}, errors.New("only keypairs are accepted")
	}
	// loop over key pairs by 2
	for i := 0; i < len(kp); i = i + 2 {
		result[kp[i]] = kp[i+1]
	}
	return &ForgeOpsResult{
		Results: []ResultMessage{
			result,
		},
		Version: "v1alpha1",
	}, nil
}

// ForgeOpsResult structure of a result message
type ForgeOpsResult struct {
	// Results are outputs from the given
	Results []ResultMessage
	// Status completion status of command
	Status ResultStatus
	// Version of this message structure
	Version string
}

// Success set the status to succeded
func (r *ForgeOpsResult) Success() {
	r.Status = ResultStatusSuccess
}

// Failed set the status to failed
func (r *ForgeOpsResult) Failed() {
	r.Status = ResultStatusFailure
}
