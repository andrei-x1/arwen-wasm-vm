package host

import (
	"math/big"
	"testing"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/testscommon/txDataBuilder"
	"github.com/stretchr/testify/require"
)

func createTestAsyncRecursiveChildContract(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock, testConfig *asyncCallBaseTestConfig) {
	childInstance := imb.CreateAndStoreInstanceMock(t, host, childAddress, testConfig.parentBalance)
	addAsyncRecursiveChildMethodsToInstanceMock(childInstance, testConfig)
}

func addAsyncRecursiveChildMethodsToInstanceMock(instanceMock *mock.InstanceMock, testConfig *asyncCallBaseTestConfig) {
	input := DefaultTestContractCallInput()
	input.GasProvided = testConfig.gasProvidedToChild

	instanceMock.AddMockMethod("recursiveAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.gasUsedByChild)

		recursiveChildCalls := big.NewInt(0).SetBytes(arguments[0])
		recursiveChildCalls.Sub(recursiveChildCalls, big.NewInt(1))
		if recursiveChildCalls.Int64() == 0 {
			return instance
		}

		destination := host.Runtime().GetSCAddress()
		function := string("recursiveAsyncCall")
		value := big.NewInt(testConfig.transferFromParentToChild).Bytes()

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)
		callData.BigInt(recursiveChildCalls)

		err := host.Runtime().ExecuteAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})

	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		host.Metering().UseGas(testConfig.gasUsedByCallback)
		instance := mock.GetMockInstance(host)
		return instance
	})
}
