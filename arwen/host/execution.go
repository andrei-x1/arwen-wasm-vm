package host

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/contexts"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (host *vmHost) doRunSmartContractCreate(input *vmcommon.ContractCreateInput) *vmcommon.VMOutput {
	host.InitState()
	defer host.Clean()

	_, blockchain, _, output, runtime, storage := host.GetContexts()

	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}

	runtime.SetVMInput(&input.VMInput)
	runtime.SetSCAddress(address)

	output.AddTxValueToAccount(address, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      address,
		CodeDeployerAddress:  input.CallerAddr,
	}

	vmOutput, err := host.performCodeDeployment(codeDeployInput)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}
	return vmOutput
}

func (host *vmHost) performCodeDeployment(input arwen.CodeDeployInput) (*vmcommon.VMOutput, error) {
	log.Trace("performCodeDeployment", "address", input.ContractAddress, "len(code)", len(input.ContractCode), "metadata", input.ContractCodeMetadata)

	_, _, metering, output, runtime, _ := host.GetContexts()

	err := metering.DeductInitialGasForDirectDeployment(input)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return nil, err
	}

	runtime.MustVerifyNextContractCode()

	vmInput := runtime.GetVMInput()
	err = runtime.StartWasmerInstance(input.ContractCode, vmInput.GasProvided, true)
	if err != nil {
		log.Debug("performCodeDeployment/StartWasmerInstance", "err", err)
		return nil, arwen.ErrContractInvalid
	}

	err = host.callInitFunction()
	if err != nil {
		return nil, err
	}

	output.DeployCode(input)
	vmOutput := output.GetVMOutput()
	runtime.CleanWasmerInstance()
	return vmOutput, nil
}

// doRunSmartContractUpgrade upgrades a contract directly
func (host *vmHost) doRunSmartContractUpgrade(input *vmcommon.ContractCallInput) *vmcommon.VMOutput {
	host.InitState()
	defer host.Clean()

	_, _, _, output, runtime, storage := host.GetContexts()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	code, codeMetadata, err := runtime.ExtractCodeUpgradeFromArgs()
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrInvalidUpgradeArguments)
	}

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
		ContractAddress:      input.RecipientAddr,
		CodeDeployerAddress:  input.CallerAddr,
	}

	vmOutput, err := host.performCodeDeployment(codeDeployInput)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}
	return vmOutput
}

func (host *vmHost) doRunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput) {
	host.InitState()
	defer host.Clean()

	_, _, metering, output, runtime, storage := host.GetContexts()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	contract, err := runtime.GetSCCode()
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrContractNotFound)
	}

	metering.UnlockGasIfAsyncCallback()
	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrNotEnoughGas)
	}

	vmInput := runtime.GetVMInput()
	err = runtime.StartWasmerInstance(contract, vmInput.GasProvided, false)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrContractInvalid)
	}

	err = host.callSCMethod()
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}

	vmOutput = output.GetVMOutput()
	runtime.CleanWasmerInstance()
	return
}

func (host *vmHost) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	log.Trace("ExecuteOnDestContext", "function", input.Function)

	bigInt, _, _, output, runtime, storage := host.GetContexts()

	bigInt.PushState()
	bigInt.InitState()

	output.PushState()
	output.CensorVMOutput()
	output.ResetGas()

	runtime.PushState()
	runtime.InitStateFromContractCallInput(input)

	storage.PushState()
	storage.SetAddress(runtime.GetSCAddress())

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err := output.TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue)
	if err != nil {
		return host.finishExecuteOnDestContext(0, err)
	}

	gasUsed, err := host.execute(input)
	if err != nil {
		return host.finishExecuteOnDestContext(gasUsed, err)
	}

	err = host.executeCurrentAsyncContext()
	return host.finishExecuteOnDestContext(gasUsed, err)
}

func (host *vmHost) finishExecuteOnDestContext(gasUsed uint64, executeErr error) (*vmcommon.VMOutput, error) {
	bigInt, _, _, output, runtime, storage := host.GetContexts()

	// Extract the VMOutput produced by the execution in isolation, before
	// restoring the contexts. This needs to be done before popping any state
	// stacks.
	gasUsedBySC, err := computeGasUsedByCurrentSC(gasUsed, output, executeErr)
	if err != nil {
		// Execution failed: restore contexts as if the execution didn't happen,
		// but first create a vmOutput to capture the error.
		vmOutput := output.CreateVMOutputInCaseOfError(err)

		bigInt.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()
		storage.PopSetActiveState()

		return vmOutput, executeErr
	}

	vmOutput := output.GetVMOutput()

	// Restore the previous context states, except Output, which will be merged
	// into the initial state (VMOutput), but only if it the child execution
	// returned vmcommon.Ok.
	bigInt.PopSetActiveState()
	runtime.PopSetActiveState()
	storage.PopSetActiveState()

	if vmOutput.ReturnCode == vmcommon.Ok {
		output.PopMergeActiveState()
		scAddress := string(runtime.GetSCAddress())
		accumulateGasUsedByContract(vmOutput, scAddress, gasUsedBySC)
	} else {
		output.PopSetActiveState()
	}

	return vmOutput, executeErr
}

func computeGasUsedByCurrentSC(
	gasUsed uint64,
	output arwen.OutputContext,
	executeErr error,
) (uint64, error) {
	if executeErr != nil {
		return 0, executeErr
	}

	vmOutput := output.GetVMOutput()
	if vmOutput.ReturnCode != vmcommon.Ok || gasUsed == 0 {
		return 0, nil
	}

	for _, outAcc := range vmOutput.OutputAccounts {
		accumulatedGasLimit := uint64(0)
		for _, outTransfer := range outAcc.OutputTransfers {
			accumulatedGasLimit += outTransfer.GasLimit
		}

		if gasUsed < outAcc.GasUsed+accumulatedGasLimit {
			return 0, arwen.ErrGasUsageError
		}

		gasUsed -= outAcc.GasUsed
		gasUsed -= accumulatedGasLimit
	}

	return gasUsed, nil
}

func (host *vmHost) ExecuteOnSameContext(input *vmcommon.ContractCallInput) (*arwen.AsyncContext, error) {
	log.Trace("ExecuteOnSameContext", "function", input.Function)

	if host.IsBuiltinFunctionName(input.Function) {
		return nil, arwen.ErrBuiltinCallOnSameContextDisallowed
	}

	var err error

	bigInt, _, _, output, runtime, _ := host.GetContexts()

	// Back up the states of the contexts (except Storage, which isn't affected
	// by ExecuteOnSameContext())
	bigInt.PushState()
	output.PushState()
	runtime.PushState()

	output.ResetGas()
	runtime.InitStateFromContractCallInput(input)

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err = output.TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue)
	if err != nil {
		host.finishExecuteOnSameContext(0, err)
		return nil, err
	}

	gasUsed, err := host.execute(input)
	if err != nil {
		host.finishExecuteOnSameContext(gasUsed, err)
		return nil, err
	}

	host.finishExecuteOnSameContext(gasUsed, err)
	return runtime.GetAsyncContext(), nil
}

func (host *vmHost) finishExecuteOnSameContext(gasUsed uint64, executeErr error) {
	bigInt, _, _, output, runtime, _ := host.GetContexts()

	gasUsedBySC, err := computeGasUsedByCurrentSC(gasUsed, output, executeErr)
	if output.ReturnCode() != vmcommon.Ok || err != nil {
		// Execution failed: restore contexts as if the execution didn't happen.
		bigInt.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()

		return
	}

	scAddress := string(runtime.GetSCAddress())
	// Execution successful: discard the backups made at the beginning and
	// resume from the new state.
	bigInt.PopDiscard()
	output.PopDiscard()
	runtime.PopSetActiveState()

	vmOutput := output.GetVMOutput()
	accumulateGasUsedByContract(vmOutput, scAddress, gasUsedBySC)
}

func accumulateGasUsedByContract(vmOutput *vmcommon.VMOutput, scAddress string, gasUsed uint64) {
	if _, ok := vmOutput.OutputAccounts[scAddress]; !ok {
		vmOutput.OutputAccounts[scAddress] = contexts.NewVMOutputAccount([]byte(scAddress))
	}
	vmOutput.OutputAccounts[scAddress].GasUsed += gasUsed
}

func (host *vmHost) isInitFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
}

func (host *vmHost) isBuiltinFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return host.IsBuiltinFunctionName(functionName)
}

func (host *vmHost) IsBuiltinFunctionName(functionName string) bool {
	_, ok := host.protocolBuiltinFunctions[functionName]
	return ok
}

// CreateNewContract creates a new contract indirectly (from another Smart Contract)
func (host *vmHost) CreateNewContract(input *vmcommon.ContractCreateInput) (newContractAddress []byte, err error) {
	newContractAddress = nil
	err = nil

	defer func() {
		if err != nil {
			newContractAddress = nil
		}
	}()

	_, blockchain, metering, output, runtime, _ := host.GetContexts()

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      nil,
		CodeDeployerAddress:  input.CallerAddr,
	}
	err = metering.DeductInitialGasForIndirectDeployment(codeDeployInput)
	if err != nil {
		return
	}

	if runtime.ReadOnly() {
		err = arwen.ErrInvalidCallOnReadOnlyMode
		return
	}

	newContractAddress, err = blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return
	}

	if blockchain.AccountExists(newContractAddress) {
		err = arwen.ErrDeploymentOverExistingAccount
		return
	}

	codeDeployInput.ContractAddress = newContractAddress
	output.DeployCode(codeDeployInput)

	defer func() {
		if err != nil {
			output.DeleteOutputAccount(newContractAddress)
		}
	}()

	runtime.MustVerifyNextContractCode()

	initCallInput := &vmcommon.ContractCallInput{
		RecipientAddr:     newContractAddress,
		Function:          arwen.InitFunctionName,
		AllowInitFunction: true,
		VMInput:           input.VMInput,
	}
	vmOutput, err := host.ExecuteOnDestContext(initCallInput)
	if err != nil {
		return
	}

	metering.RestoreGas(vmOutput.GasRemaining)
	blockchain.IncreaseNonce(input.CallerAddr)

	return
}

func (host *vmHost) checkUpgradePermission(vmInput *vmcommon.ContractCallInput) error {
	contract, err := host.blockChainHook.GetUserAccount(vmInput.RecipientAddr)
	if err != nil {
		return err
	}
	if check.IfNilReflect(contract) {
		return arwen.ErrNilContract
	}

	codeMetadata := vmcommon.CodeMetadataFromBytes(contract.GetCodeMetadata())
	isUpgradeable := codeMetadata.Upgradeable
	callerAddress := vmInput.CallerAddr
	ownerAddress := contract.GetOwnerAddress()
	isCallerOwner := bytes.Equal(callerAddress, ownerAddress)

	if isUpgradeable && isCallerOwner {
		return nil
	}

	return arwen.ErrUpgradeNotAllowed
}

// executeUpgrade upgrades a contract indirectly (from another contract)
func (host *vmHost) executeUpgrade(input *vmcommon.ContractCallInput) (uint64, error) {
	_, _, metering, output, runtime, _ := host.GetContexts()

	initialGasProvided := input.GasProvided
	err := host.checkUpgradePermission(input)
	if err != nil {
		return 0, err
	}

	code, codeMetadata, err := runtime.ExtractCodeUpgradeFromArgs()
	if err != nil {
		return 0, arwen.ErrInvalidUpgradeArguments
	}

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
		ContractAddress:      input.RecipientAddr,
		CodeDeployerAddress:  input.CallerAddr,
	}

	err = metering.DeductInitialGasForDirectDeployment(codeDeployInput)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return 0, err
	}

	runtime.PushInstance()
	runtime.MustVerifyNextContractCode()

	vmInput := runtime.GetVMInput()
	err = runtime.StartWasmerInstance(codeDeployInput.ContractCode, vmInput.GasProvided, true)
	if err != nil {
		log.Debug("performCodeDeployment/StartWasmerInstance", "err", err)
		runtime.PopInstance()
		return 0, arwen.ErrContractInvalid
	}

	err = host.callInitFunction()
	if err != nil {
		runtime.PopInstance()
		return 0, err
	}

	output.DeployCode(codeDeployInput)
	if output.ReturnCode() != vmcommon.Ok {
		runtime.PopInstance()
		return 0, arwen.ErrReturnCodeNotOk
	}

	gasToRestoreToCaller := metering.GasLeft()

	runtime.PopInstance()
	metering.RestoreGas(gasToRestoreToCaller)

	return initialGasProvided - gasToRestoreToCaller, nil
}

func (host *vmHost) executeSmartContractCall(
	input *vmcommon.ContractCallInput,
	metering arwen.MeteringContext,
	runtime arwen.RuntimeContext,
	output arwen.OutputContext,
	withInitialGasDeduct bool,
) (uint64, error) {
	if host.isInitFunctionBeingCalled() && !input.AllowInitFunction {
		return 0, arwen.ErrInitFuncCalledInRun
	}

	// Use all gas initially, on the Wasmer instance of the caller
	// (runtime.PushInstance() is called later). In case of successful execution,
	// the unused gas will be restored.
	metering.UnlockGasIfAsyncCallback()
	initialGasProvided := input.GasProvided
	metering.UseGas(initialGasProvided)

	isUpgrade := input.Function == arwen.UpgradeFunctionName
	if isUpgrade {
		return host.executeUpgrade(input)
	}

	contract, err := runtime.GetSCCode()
	if err != nil {
		return 0, err
	}

	if withInitialGasDeduct {
		err = metering.DeductInitialGasForExecution(contract)
		if err != nil {
			return 0, err
		}
	}

	runtime.PushInstance()

	gasForExecution := runtime.GetVMInput().GasProvided
	err = runtime.StartWasmerInstance(contract, gasForExecution, false)
	if err != nil {
		runtime.PopInstance()
		return 0, err
	}

	err = host.callSCMethodIndirect()
	if err != nil {
		runtime.PopInstance()
		return 0, err
	}

	if output.ReturnCode() != vmcommon.Ok {
		runtime.PopInstance()
		return 0, arwen.ErrReturnCodeNotOk
	}

	gasToRestoreToCaller := metering.GasLeft()

	runtime.PopInstance()
	metering.RestoreGas(gasToRestoreToCaller)

	return initialGasProvided - gasToRestoreToCaller, nil
}

func (host *vmHost) execute(input *vmcommon.ContractCallInput) (uint64, error) {
	_, _, metering, output, runtime, storage := host.GetContexts()

	if host.isBuiltinFunctionBeingCalled() {
		err := metering.UseGasForAsyncStep()
		if err != nil {
			return 0, err
		}

		newVMInput, err := host.callBuiltinFunction(input)
		if err != nil {
			return 0, err
		}

		if newVMInput != nil {
			runtime.InitStateFromContractCallInput(newVMInput)
			storage.SetAddress(runtime.GetSCAddress())
			return host.executeSmartContractCall(newVMInput, metering, runtime, output, false)
		}

		return 0, nil
	}

	return host.executeSmartContractCall(input, metering, runtime, output, true)
}

func (host *vmHost) callSCMethodIndirect() error {
	function, err := host.Runtime().GetFunctionToCall()
	if err != nil {
		if errors.Is(err, arwen.ErrNilCallbackFunction) {
			return nil
		}
		return err
	}

	_, err = function()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}

	return err
}

func (host *vmHost) callBuiltinFunction(input *vmcommon.ContractCallInput) (*vmcommon.ContractCallInput, error) {
	_, _, metering, output, _, _ := host.GetContexts()

	vmOutput, err := host.blockChainHook.ProcessBuiltInFunction(input)
	if err != nil {
		metering.UseGas(input.GasProvided)
		return nil, err
	}

	gasConsumed := input.GasProvided - vmOutput.GasRemaining
	if vmOutput.GasRemaining < input.GasProvided {
		metering.UseGas(gasConsumed)
	}

	newVMInput, err := host.isSCExecutionAfterBuiltInFunc(input, vmOutput)
	if err != nil {
		return nil, err
	}

	output.AddToActiveState(vmOutput)

	return newVMInput, nil
}

func (host *vmHost) EthereumCallData() []byte {
	if host.ethInput == nil {
		host.ethInput = host.createETHCallInput()
	}
	return host.ethInput
}

func (host *vmHost) callInitFunction() error {
	runtime := host.Runtime()
	init := runtime.GetInitFunction()
	if init == nil {
		return nil
	}

	_, err := init()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}

	return err
}

func (host *vmHost) callSCMethod() error {
	runtime := host.Runtime()

	err := host.verifyAllowedFunctionCall()
	if err != nil {
		return err
	}

	callType := runtime.GetVMInput().CallType
	function, err := host.getFunctionByCallType(callType)
	if err != nil {
		if callType == vmcommon.AsynchronousCallBack && errors.Is(err, arwen.ErrNilCallbackFunction) {
			return host.postprocessCrossShardCallback()
		}
		return err
	}

	_, err = function()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}
	if err != nil {
		return err
	}

	err = host.executeCurrentAsyncContext()
	if err != nil {
		return err
	}

	switch callType {
	case vmcommon.DirectCall:
		break
	case vmcommon.AsynchronousCall:
		err = host.sendAsyncCallbackToCaller()
	case vmcommon.AsynchronousCallBack:
		err = host.postprocessCrossShardCallback()
	default:
		err = arwen.ErrUnknownCallType
	}

	return err
}

// TODO move into RuntimeContext
func (host *vmHost) getFunctionByCallType(callType vmcommon.CallType) (wasmer.ExportedFunctionCallback, error) {
	runtime := host.Runtime()

	if callType != vmcommon.AsynchronousCallBack {
		return runtime.GetFunctionToCall()
	}

	asyncContext, err := host.loadCurrentAsyncContext()
	if err != nil {
		return nil, err
	}

	vmInput := runtime.GetVMInput()

	customCallback := false
	for _, asyncCallGroup := range asyncContext.AsyncCallGroups {
		for _, asyncCall := range asyncCallGroup.AsyncCalls {
			if bytes.Equal(vmInput.CallerAddr, asyncCall.Destination) {
				customCallback = true
				// TODO why asyncCall.SuccessCallback? why not check for error as well,
				// and set asyncCall.ErrorCallback?
				runtime.SetCustomCallFunction(asyncCall.SuccessCallback)
				break
			}
		}

		if customCallback {
			break
		}
	}

	return runtime.GetFunctionToCall()
}

func (host *vmHost) verifyAllowedFunctionCall() error {
	runtime := host.Runtime()
	functionName := runtime.Function()

	isInit := functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
	if isInit {
		return arwen.ErrInitFuncCalledInRun
	}

	isCallBack := functionName == arwen.CallbackFunctionName
	isInAsyncCallBack := runtime.GetVMInput().CallType == vmcommon.AsynchronousCallBack
	if isCallBack && !isInAsyncCallBack {
		return arwen.ErrCallBackFuncCalledInRun
	}

	return nil
}

// The first four bytes is the method selector. The rest of the input data are method arguments in chunks of 32 bytes.
// The method selector is the kecccak256 hash of the method signature.
func (host *vmHost) createETHCallInput() []byte {
	newInput := make([]byte, 0)

	function := host.Runtime().Function()
	if len(function) > 0 {
		hashOfFunction, err := host.cryptoHook.Keccak256([]byte(function))
		if err != nil {
			return nil
		}

		newInput = append(newInput, hashOfFunction[0:4]...)
	}

	for _, arg := range host.Runtime().Arguments() {
		paddedArg := make([]byte, arwen.ArgumentLenEth)
		copy(paddedArg[arwen.ArgumentLenEth-len(arg):], arg)
		newInput = append(newInput, paddedArg...)
	}

	return newInput
}

func (host *vmHost) isSCExecutionAfterBuiltInFunc(
	vmInput *vmcommon.ContractCallInput,
	vmOutput *vmcommon.VMOutput,
) (*vmcommon.ContractCallInput, error) {
	if vmOutput.ReturnCode != vmcommon.Ok {
		return nil, nil
	}

	if !core.IsSmartContractAddress(vmInput.RecipientAddr) {
		return nil, nil
	}

	outAcc, ok := vmOutput.OutputAccounts[string(vmInput.RecipientAddr)]
	if !ok {
		return nil, nil
	}
	if len(outAcc.OutputTransfers) != 1 {
		return nil, nil
	}

	callType := vmInput.CallType
	txData := prependCallbackToTxDataIfAsyncCall(outAcc.OutputTransfers[0].Data, callType)

	function, arguments, err := host.CallArgsParser().ParseData(txData)
	if err != nil {
		return nil, err
	}

	// TODO CurrentTxHash might equal PrevTxHash, because Arwen does not generate
	// SCRs between async steps; analyze and fix this (e.g. create dummy SCRs
	// just for hashing, or hash the VMInput itself)
	newVMInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     vmInput.CallerAddr,
			Arguments:      arguments,
			CallValue:      big.NewInt(0),
			CallType:       callType,
			GasPrice:       vmInput.GasPrice,
			GasProvided:    vmOutput.GasRemaining,
			OriginalTxHash: vmInput.OriginalTxHash,
			CurrentTxHash:  vmInput.CurrentTxHash,
			PrevTxHash:     vmInput.PrevTxHash,
		},
		RecipientAddr:     vmInput.RecipientAddr,
		Function:          function,
		AllowInitFunction: false,
	}

	fillWithESDTValue(vmInput, newVMInput)

	return newVMInput, nil
}

func fillWithESDTValue(fullVMInput *vmcommon.ContractCallInput, newVMInput *vmcommon.ContractCallInput) {
	if fullVMInput.Function != core.BuiltInFunctionESDTTransfer {
		return
	}

	newVMInput.ESDTTokenName = fullVMInput.Arguments[0]
	newVMInput.ESDTValue = big.NewInt(0).SetBytes(fullVMInput.Arguments[1])
}

func prependCallbackToTxDataIfAsyncCall(txData []byte, callType vmcommon.CallType) string {
	if callType == vmcommon.AsynchronousCallBack {
		return string(append([]byte(arwen.CallbackFunctionName), txData...))
	}

	return string(txData)
}
