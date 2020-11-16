package arwen

import "github.com/ElrondNetwork/elrond-go/core/vmcommon"

// AsyncCall holds the information about an individual async call
type AsyncCall struct {
	Status          AsyncCallStatus
	Destination     []byte
	Data            []byte
	GasLimit        uint64
	GasLocked       uint64
	ValueBytes      []byte
	SuccessCallback string
	ErrorCallback   string
	ProvidedGas     uint64
}

// GetDestination returns the destination of an async call
func (ac *AsyncCall) GetDestination() []byte {
	return ac.Destination
}

// GetData returns the transaction data of the async call
func (ac *AsyncCall) GetData() []byte {
	return ac.Data
}

// GetGasLimit returns the gas limit of the current async call
func (ac *AsyncCall) GetGasLimit() uint64 {
	return ac.GasLimit
}

// GetGasLocked returns the gas locked for the async callback
func (ac *AsyncCall) GetGasLocked() uint64 {
	return ac.GasLocked
}

// GetValueBytes returns the byte representation of the value of the async call
func (ac *AsyncCall) GetValueBytes() []byte {
	return ac.ValueBytes
}

// IsInterfaceNil returns true if there is no value under the interface
func (ac *AsyncCall) IsInterfaceNil() bool {
	return ac == nil
}

// UpdateStatus sets the status of the async call depending on the provided ReturnCode
func (ac *AsyncCall) UpdateStatus(returnCode vmcommon.ReturnCode) {
	ac.Status = AsyncCallResolved
	if returnCode != vmcommon.Ok {
		ac.Status = AsyncCallRejected
	}
}

// GetCallbackName returns the name of the callback to execute, depending on
// the status of the async call
func (ac *AsyncCall) GetCallbackName() string {
	if ac.Status == AsyncCallResolved {
		return ac.SuccessCallback
	}

	return ac.ErrorCallback
}