package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern void getOwner(void *context, int32_t resultOffset);
// extern void getExternalBalance(void *context, int32_t addressOffset, int32_t resultOffset);
// extern int32_t blockHash(void *context, long long nonce, int32_t resultOffset);
// extern int32_t transferValue(void *context, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t getArgumentLength(void *context, int32_t id);
// extern int32_t getArgument(void *context, int32_t id, int32_t argOffset);
// extern int32_t getFunction(void *context, int32_t functionOffset);
// extern int32_t getNumArguments(void *context);
// extern int32_t storageStore(void *context, int32_t keyOffset, int32_t dataOffset, int32_t dataLength);
// extern int32_t storageGetValueLength(void *context, int32_t keyOffset);
// extern int32_t storageLoad(void *context, int32_t keyOffset, int32_t dataOffset);
// extern void getCaller(void *context, int32_t resultOffset);
// extern int32_t callValue(void *context, int32_t resultOffset);
// extern void writeLog(void *context, int32_t pointer, int32_t length, int32_t topicPtr, int32_t numTopics);
// extern void returnData(void* context, int32_t dataOffset, int32_t length);
// extern void signalError(void* context, int32_t messageOffset, int32_t messageLength);
// extern long long getGasLeft(void *context);
//
// extern int32_t executeOnDestContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t executeOnSameContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t delegateExecution(void *context, long long gas, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t executeReadOnly(void *context, long long gas, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t createContract(void *context, int32_t valueOffset, int32_t codeOffset, int32_t length, int32_t resultOffset, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t asyncCall(void *context, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
//
// extern int32_t getNumReturnData(void *context);
// extern int32_t getReturnDataSize(void *context, int32_t resultId);
// extern int32_t getReturnData(void *context, int32_t resultId, int32_t dataOffset);
//
// extern long long getBlockTimestamp(void *context);
// extern long long getBlockNonce(void *context);
// extern long long getBlockRound(void *context);
// extern long long getBlockEpoch(void *context);
// extern void getBlockRandomSeed(void *context, int32_t resultOffset);
// extern void getStateRootHash(void *context, int32_t resultOffset);
//
// extern long long getPrevBlockTimestamp(void *context);
// extern long long getPrevBlockNonce(void *context);
// extern long long getPrevBlockRound(void *context);
// extern long long getPrevBlockEpoch(void *context);
// extern void getPrevBlockRandomSeed(void *context, int32_t resultOffset);
//
// extern long long int64getArgument(void *context, int32_t id);
// extern int32_t int64storageStore(void *context, int32_t keyOffset, long long value);
// extern long long int64storageLoad(void *context, int32_t keyOffset);
// extern void int64finish(void* context, long long value);
import "C"

import (
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func ElrondEImports() (*wasmer.Imports, error) {
	imports := wasmer.NewImports()
	imports = imports.Namespace("env")

	imports, err := imports.Append("getOwner", getOwner, C.getOwner)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getExternalBalance", getExternalBalance, C.getExternalBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockHash", blockHash, C.blockHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transferValue", transferValue, C.transferValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("asyncCall", asyncCall, C.asyncCall)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgumentLength", getArgumentLength, C.getArgumentLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgument", getArgument, C.getArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getFunction", getFunction, C.getFunction)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumArguments", getNumArguments, C.getNumArguments)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageStore", storageStore, C.storageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageGetValueLength", storageGetValueLength, C.storageGetValueLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoad", storageLoad, C.storageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCaller", getCaller, C.getCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallValue", callValue, C.callValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("writeLog", writeLog, C.writeLog)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("finish", returnData, C.returnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("signalError", signalError, C.signalError)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockTimestamp", getBlockTimestamp, C.getBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockNonce", getBlockNonce, C.getBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRound", getBlockRound, C.getBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockEpoch", getBlockEpoch, C.getBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRandomSeed", getBlockRandomSeed, C.getBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getStateRootHash", getStateRootHash, C.getStateRootHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockTimestamp", getPrevBlockTimestamp, C.getPrevBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockNonce", getPrevBlockNonce, C.getPrevBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRound", getPrevBlockRound, C.getPrevBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockEpoch", getPrevBlockEpoch, C.getPrevBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRandomSeed", getPrevBlockRandomSeed, C.getPrevBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getGasLeft", getGasLeft, C.getGasLeft)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnDestContext", executeOnDestContext, C.executeOnDestContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnSameContext", executeOnSameContext, C.executeOnSameContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("delegateExecution", delegateExecution, C.delegateExecution)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("createContract", createContract, C.createContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeReadOnly", executeReadOnly, C.executeReadOnly)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumReturnData", getNumReturnData, C.getNumReturnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnDataSize", getReturnDataSize, C.getReturnDataSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnData", getReturnData, C.getReturnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64getArgument", int64getArgument, C.int64getArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64storageStore", int64storageStore, C.int64storageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64storageLoad", int64storageLoad, C.int64storageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64finish", int64finish, C.int64finish)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export getGasLeft
func getGasLeft(context unsafe.Pointer) int64 {
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetGasLeft
	metering.UseGas(gasToUse)

	return int64(metering.GasLeft())
}

//export getOwner
func getOwner(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	owner := runtime.GetSCAddress()
	err := runtime.MemStore(resultOffset, owner)
	if withFault(err, context) {
		return
	}

	gasToUse := metering.GasSchedule().ElrondAPICost.GetOwner
	metering.UseGas(gasToUse)
}

//export signalError
func signalError(context unsafe.Pointer, messageOffset int32, messageLength int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if withFault(err, context) {
		return
	}
	runtime.SignalUserError(string(message))

	gasToUse := metering.GasSchedule().ElrondAPICost.SignalError
	metering.UseGas(gasToUse)
}

//export getExternalBalance
func getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if withFault(err, context) {
		return
	}

	balance := blockchain.GetBalance(address)

	err = runtime.MemStore(resultOffset, balance)
	if withFault(err, context) {
		return
	}

	gasToUse := metering.GasSchedule().ElrondAPICost.GetExternalBalance
	metering.UseGas(gasToUse)
}

//export blockHash
func blockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockHash
	metering.UseGas(gasToUse)

	//TODO: change blockchain hook to treat actual nonce - not the offset.
	hash := blockchain.BlockHash(nonce)
	err := runtime.MemStore(resultOffset, hash)
	if withFault(err, context) {
		return 1
	}

	return 0
}

//export transferValue
func transferValue(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	output := arwen.GetOutputContext(context)

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	if withFault(err, context) {
		return 1
	}

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if withFault(err, context) {
		return 1
	}

	data, err := runtime.MemLoad(dataOffset, length)
	if withFault(err, context) {
		return 1
	}

	gasToUse := metering.GasSchedule().ElrondAPICost.TransferValue
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	metering.UseGas(gasToUse)

	invBytes := arwen.InverseBytes(value)

	output.Transfer(dest, send, 0, big.NewInt(0).SetBytes(invBytes), data)

	return 0
}

//export asyncCall
func asyncCall(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	output := arwen.GetOutputContext(context)

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	if withFault(err, context) {
		return 1
	}

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if withFault(err, context) {
		return 1
	}

	data, err := runtime.MemLoad(dataOffset, length)
	if withFault(err, context) {
		return 1
	}

	// TODO define gas cost specific to async calls - asyncCall should cost more
	// than ExecuteOnDestContext, especially cross-shard async calls
	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteOnDestContext
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
	metering.UseGas(gasToUse)

	// TODO verify if gasLimit is large enough for sender call, destination call and callback call
	gasLimit := metering.GasLeft()

	// Set up the async call as if it is not known whether the called SC
	// is in the same shard with the caller or not. This will be later resolved
	// in the handler for BreakpointAsyncCall.
	invValueBytes := arwen.InverseBytes(value)
	output.Transfer(dest, send, gasLimit, big.NewInt(0).SetBytes(invValueBytes), data)

	runtime.SetAsyncCallInfo(dest, invValueBytes, gasLimit, data)

	// Instruct Wasmer to interrupt the execution of the caller SC.
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointAsyncCall)

	return 0
}

//export getArgumentLength
func getArgumentLength(context unsafe.Pointer, id int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		return -1
	}

	return int32(len(args[id]))
}

//export getArgument
func getArgument(context unsafe.Pointer, id int32, argOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		return -1
	}

	err := runtime.MemStore(argOffset, args[id])
	if withFault(err, context) {
		return -1
	}

	return int32(len(args[id]))
}

//export getFunction
func getFunction(context unsafe.Pointer, functionOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetFunction
	metering.UseGas(gasToUse)

	function := runtime.Function()
	err := runtime.MemStore(functionOffset, []byte(function))
	if withFault(err, context) {
		return -1
	}

	return int32(len(function))
}

//export getNumArguments
func getNumArguments(context unsafe.Pointer) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetNumArguments
	metering.UseGas(gasToUse)

	return int32(len(runtime.Arguments()))
}

//export storageStore
func storageStore(context unsafe.Pointer, keyOffset int32, dataOffset int32, dataLength int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, arwen.HashLen)
	if withFault(err, context) {
		return -1
	}

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if withFault(err, context) {
		return -1
	}

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageStore
	metering.UseGas(gasToUse)

	return storage.SetStorage(runtime.GetSCAddress(), key, data)
}

//export storageGetValueLength
func storageGetValueLength(context unsafe.Pointer, keyOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, arwen.HashLen)
	if withFault(err, context) {
		return -1
	}

	data := storage.GetStorage(runtime.GetSCAddress(), key)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageLoad
	metering.UseGas(gasToUse)

	return int32(len(data))
}

//export storageLoad
func storageLoad(context unsafe.Pointer, keyOffset int32, dataOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, arwen.HashLen)
	if withFault(err, context) {
		return -1
	}

	data := storage.GetStorage(runtime.GetSCAddress(), key)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageLoad
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	metering.UseGas(gasToUse)

	err = runtime.MemStore(dataOffset, data)
	if withFault(err, context) {
		return -1
	}

	return int32(len(data))
}

//export getCaller
func getCaller(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	caller := runtime.GetVMInput().CallerAddr

	err := runtime.MemStore(resultOffset, caller)
	if withFault(err, context) {
		return
	}

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCaller
	metering.UseGas(gasToUse)
}

//export callValue
func callValue(context unsafe.Pointer, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	value := runtime.GetVMInput().CallValue.Bytes()
	invBytes := arwen.InverseBytes(value)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	err := runtime.MemStore(resultOffset, invBytes)
	if withFault(err, context) {
		return -1
	}

	return int32(len(value))
}

//export writeLog
func writeLog(context unsafe.Pointer, pointer int32, length int32, topicPtr int32, numTopics int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	log, err := runtime.MemLoad(pointer, length)
	if withFault(err, context) {
		return
	}

	topics, err := arwen.GuardedMakeByteSlice2D(numTopics)
	if withFault(err, context) {
		return
	}

	for i := int32(0); i < numTopics; i++ {
		topics[i], err = runtime.MemLoad(topicPtr+i*arwen.HashLen, arwen.HashLen)
		if withFault(err, context) {
			return
		}
	}

	output.WriteLog(runtime.GetSCAddress(), topics, log)

	gasToUse := metering.GasSchedule().ElrondAPICost.Log
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(numTopics*arwen.HashLen+length)
	metering.UseGas(gasToUse)
}

//export getBlockTimestamp
func getBlockTimestamp(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockTimeStamp
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentTimeStamp())
}

//export getBlockNonce
func getBlockNonce(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockNonce
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentNonce())
}

//export getBlockRound
func getBlockRound(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRound
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentRound())
}

//export getBlockEpoch
func getBlockEpoch(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockEpoch
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentEpoch())
}

//export getBlockRandomSeed
func getBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	metering.UseGas(gasToUse)

	randomSeed := blockchain.CurrentRandomSeed()
	err := runtime.MemStore(pointer, randomSeed)
	withFault(err, context)
}

//export getStateRootHash
func getStateRootHash(context unsafe.Pointer, pointer int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetStateRootHash
	metering.UseGas(gasToUse)

	stateRootHash := blockchain.GetStateRootHash()
	err := runtime.MemStore(pointer, stateRootHash)
	withFault(err, context)
}

//export getPrevBlockTimestamp
func getPrevBlockTimestamp(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockTimeStamp
	metering.UseGas(gasToUse)

	return int64(blockchain.LastTimeStamp())
}

//export getPrevBlockNonce
func getPrevBlockNonce(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockNonce
	metering.UseGas(gasToUse)

	return int64(blockchain.LastNonce())
}

//export getPrevBlockRound
func getPrevBlockRound(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRound
	metering.UseGas(gasToUse)

	return int64(blockchain.LastRound())
}

//export getPrevBlockEpoch
func getPrevBlockEpoch(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockEpoch
	metering.UseGas(gasToUse)

	return int64(blockchain.LastEpoch())
}

//export getPrevBlockRandomSeed
func getPrevBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	metering.UseGas(gasToUse)

	randomSeed := blockchain.LastRandomSeed()
	err := runtime.MemStore(pointer, randomSeed)
	withFault(err, context)
}

//export returnData
func returnData(context unsafe.Pointer, pointer int32, length int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	data, err := runtime.MemLoad(pointer, length)
	if withFault(err, context) {
		return
	}

	output.Finish(data)

	gasToUse := metering.GasSchedule().ElrondAPICost.Finish
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	metering.UseGas(gasToUse)
}

//export int64getArgument
func int64getArgument(context unsafe.Pointer, id int32) int64 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id {
		return -1
	}

	argBigInt := big.NewInt(0).SetBytes(args[id])
	return argBigInt.Int64()
}

//export int64storageStore
func int64storageStore(context unsafe.Pointer, keyOffset int32, value int64) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, arwen.HashLen)
	if withFault(err, context) {
		return -1
	}

	data := big.NewInt(value)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64StorageStore
	metering.UseGas(gasToUse)

	return storage.SetStorage(runtime.GetSCAddress(), key, data.Bytes())
}

//export int64storageLoad
func int64storageLoad(context unsafe.Pointer, keyOffset int32) int64 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, arwen.HashLen)
	if withFault(err, context) {
		return 0
	}

	data := storage.GetStorage(runtime.GetSCAddress(), key)

	bigInt := big.NewInt(0).SetBytes(data)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64StorageLoad
	metering.UseGas(gasToUse)

	return bigInt.Int64()
}

//export int64finish
func int64finish(context unsafe.Pointer, value int64) {
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	var valueBytes []byte
	if value == 0 {
		valueBytes = []byte{0}
	} else {
		valueBytes = big.NewInt(0).SetInt64(value).Bytes()
	}
	output.Finish(valueBytes)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64Finish
	metering.UseGas(gasToUse)
}

//export executeOnSameContext
func executeOnSameContext(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVmContext(context)
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if withFault(err, context) {
		return 1
	}

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if withFault(err, context) {
		return 1
	}

	function, data, actualLen := getArgumentsFromMemory(context, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)

	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteOnSameContext
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	metering.UseGas(gasToUse)

	invBytes := arwen.InverseBytes(value)
	bigIntVal := big.NewInt(0).SetBytes(invBytes)
	output.Transfer(dest, send, 0, bigIntVal, nil)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   data,
			CallValue:   bigIntVal,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: dest,
		Function:      function,
	}

	_, err = host.ExecuteOnDestContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

//export executeOnDestContext
func executeOnDestContext(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVmContext(context)
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if withFault(err, context) {
		return 1
	}

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if withFault(err, context) {
		return 1
	}

	function, data, actualLen := getArgumentsFromMemory(context, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)

	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteOnDestContext
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	metering.UseGas(gasToUse)

	invBytes := arwen.InverseBytes(value)
	output.Transfer(dest, send, 0, big.NewInt(0).SetBytes(invBytes), nil)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   data,
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: dest,
		Function:      function,
	}

	_, err = host.ExecuteOnDestContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

func getArgumentsFromMemory(
	context unsafe.Pointer,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) (string, [][]byte, int32) {
	runtime := arwen.GetRuntimeContext(context)

	argumentsLengthData, err := runtime.MemLoad(argumentsLengthOffset, numArguments*4)
	if withFault(err, context) {
		return "", nil, 0
	}

	currOffset := dataOffset
	data, err := arwen.GuardedMakeByteSlice2D(numArguments)
	if withFault(err, context) {
		return "", nil, 0
	}

	for i := int32(0); i < numArguments; i++ {
		currArgLenData := argumentsLengthData[i*4 : i*4+4]
		actualLen := dataToInt32(currArgLenData)

		data[i], err = runtime.MemLoad(currOffset, actualLen)
		if withFault(err, context) {
			return "", nil, 0
		}

		currOffset += actualLen
	}

	function, err := runtime.MemLoad(functionOffset, functionLength)
	if withFault(err, context) {
		return "", nil, 0
	}

	return string(function), data, currOffset - dataOffset
}

//export delegateExecution
func delegateExecution(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVmContext(context)
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()

	address, err := runtime.MemLoad(addressOffset, arwen.HashLen)
	if withFault(err, context) {
		return 1
	}

	function, data, actualLen := getArgumentsFromMemory(context, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)

	value := runtime.GetVMInput().CallValue
	sender := runtime.GetVMInput().CallerAddr

	gasToUse := metering.GasSchedule().ElrondAPICost.DelegateExecution
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	metering.UseGas(gasToUse)

	output.Transfer(address, sender, 0, value, nil)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   value,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: address,
		Function:      function,
	}

	err = host.ExecuteOnSameContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

func dataToInt32(data []byte) int32 {
	actualLen := int32(0)
	for i := len(data) - 1; i >= 0; i-- {
		actualLen = (actualLen << 8) + int32(data[i])
	}

	return actualLen
}

//export executeReadOnly
func executeReadOnly(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVmContext(context)
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()

	address, err := runtime.MemLoad(addressOffset, arwen.HashLen)
	if withFault(err, context) {
		return 1
	}

	function, data, actualLen := getArgumentsFromMemory(context, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)

	value := runtime.GetVMInput().CallValue
	sender := runtime.GetVMInput().CallerAddr

	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteReadOnly
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	metering.UseGas(gasToUse)

	output.Transfer(address, sender, 0, value, nil)

	runtime.SetReadOnly(true)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   value,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: address,
		Function:      function,
	}

	err = host.ExecuteOnSameContext(contractCallInput)
	runtime.SetReadOnly(false)
	if err != nil {
		return 1
	}

	return 0
}

//export createContract
func createContract(
	context unsafe.Pointer,
	valueOffset int32,
	codeOffset int32,
	length int32,
	resultOffset int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVmContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	sender := runtime.GetSCAddress()
	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if withFault(err, context) {
		return 1
	}

	code, err := runtime.MemLoad(codeOffset, length)
	if withFault(err, context) {
		return 1
	}

	_, data, actualLen := getArgumentsFromMemory(context, 0, 0, numArguments, argumentsLengthOffset, dataOffset)

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	metering.UseGas(gasToUse)
	gasLimit := metering.GasLeft()

	contractCreate := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: gasLimit,
		},
		ContractCode: code,
	}

	newAddress, err := host.CreateNewContract(contractCreate)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, newAddress)
	if withFault(err, context) {
		return 1
	}

	return 0
}

//export getNumReturnData
func getNumReturnData(context unsafe.Pointer) int32 {
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetNumReturnData
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	return int32(len(returnData))
}

//export getReturnDataSize
func getReturnDataSize(context unsafe.Pointer, resultId int32) int32 {
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetReturnDataSize
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	if resultId >= int32(len(returnData)) {
		return 0
	}

	return int32(len(returnData[resultId]))
}

//export getReturnData
func getReturnData(context unsafe.Pointer, resultId int32, dataOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetReturnData
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	if resultId >= int32(len(returnData)) {
		return 0
	}

	err := runtime.MemStore(dataOffset, returnData[resultId])
	if withFault(err, context) {
		return 0
	}

	return int32(len(returnData[resultId]))
}

func withFault(err error, context unsafe.Pointer) bool {
	if err != nil {
		runtime := arwen.GetRuntimeContext(context)
		metering := arwen.GetMeteringContext(context)

		runtime.SignalUserError(err.Error())
		metering.UseGas(metering.GasLeft())

		return true
	}

	return false
}
