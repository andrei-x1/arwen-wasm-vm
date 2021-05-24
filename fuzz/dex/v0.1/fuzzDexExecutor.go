package dex

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	am "github.com/ElrondNetwork/arwen-wasm-vm/arwenmandos"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/fileresolver"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	mjparse "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/parse"
	mjwrite "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/write"
	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

type fuzzDexExecutorInitArgs struct {
	wegldTokenId            string
	mexTokenId              string
	busdTokenId				string
	wemeLpTokenId			string
	webuLpTokenId			string
	wemeFarmTokenId			string
	webuFarmTokenId			string
	mexFarmTokenId			string
	numUsers                int
	numEvents               int
	removeLiquidityProb     float32
	addLiquidityProb        float32
	swapProb                float32
	queryPairsProb          float32
	enterFarmProb           float32
	exitFarmProb            float32
	increaseEpochProb       float32
	removeLiquidityMaxValue int
	addLiquidityMaxValue    int
	swapMaxValue            int
	enterFarmMaxValue       int
	exitFarmMaxValue        int
	blockEpochIncrease      int
}

type SwapPair struct {
	firstToken 	string
	secondToken string
	lpToken		string
	address		string
}

type Farm struct {
	farmingToken string
	farmToken	 string
	rewardToken  string
	address 	 string
}

type FarmerInfo struct {
	user    string
	value   int64
	farm 	Farm
}

type fuzzDexExecutor struct {
	arwenTestExecutor *am.ArwenTestExecutor
	world             *worldhook.MockWorld
	vm                vmi.VMExecutionHandler
	mandosParser      mjparse.Parser
	txIndex           int

	wegldTokenId            string
	mexTokenId              string
	busdTokenId				string
	wemeLpTokenId			string
	webuLpTokenId			string
	wemeFarmTokenId			string
	webuFarmTokenId			string
	mexFarmTokenId			string
	ownerAddress            string
	wemeFarmAddress			string
	webuFarmAddress			string
	mexFarmAddress			string
	wemeSwapAddress			string
	webuSwapAddress			string
	numUsers                int
	numTokens               int
	numEvents               int
	removeLiquidityProb     float32
	addLiquidityProb        float32
	swapProb                float32
	queryPairsProb          float32
	enterFarmProb           float32
	exitFarmProb            float32
	unbondProb              float32
	increaseEpochProb       float32
	removeLiquidityMaxValue int
	addLiquidityMaxValue    int
	swapMaxValue            int
	enterFarmMaxValue       int
	exitFarmMaxValue        int
	unbondMaxValue          int
	blockEpochIncrease      int
	tokensCheckFrequency    int
	currentFarmTokenNonce   map[string]int
	farmers                 map[int]FarmerInfo
	generatedScenario       *mj.Scenario
}

type eventsStatistics struct {
	swapFixedInputHits   int
	swapFixedInputMisses int

	swapFixedOutputHits   int
	swapFixedOutputMisses int

	addLiquidityHits        int
	addLiquidityMisses      int
	addLiquidityPriceChecks int

	removeLiquidityHits        int
	removeLiquidityMisses      int
	removeLiquidityPriceChecks int

	queryPairsHits   int
	queryPairsMisses int

	enterFarmHits   int
	enterFarmMisses int

	exitFarmHits        int
	exitFarmMisses      int
	exitFarmWithRewards int
}

func newFuzzDexExecutor(fileResolver fr.FileResolver) (*fuzzDexExecutor, error) {
	arwenTestExecutor, err := am.NewArwenTestExecutor("../../arwenmandos")
	if err != nil {
		return nil, err
	}

	parser := mjparse.NewParser(fileResolver)

	return &fuzzDexExecutor{
		arwenTestExecutor: arwenTestExecutor,
		world:             arwenTestExecutor.World,
		vm:                arwenTestExecutor.GetVM(),
		mandosParser:      parser,
		txIndex:           0,
		generatedScenario: &mj.Scenario{
			Name: "fuzz generated",
		},
	}, nil
}

func (pfe *fuzzDexExecutor) saveGeneratedScenario() {
	serialized := mjwrite.ScenarioToJSONString(pfe.generatedScenario)

	err := ioutil.WriteFile("fuzz_gen.scen.json", []byte(serialized), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (pfe *fuzzDexExecutor) executeStep(stepSnippet string) error {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return err
	}

	pfe.addStep(step)
	return pfe.arwenTestExecutor.ExecuteStep(step)
}

func (pfe *fuzzDexExecutor) addStep(step mj.Step) {
	pfe.generatedScenario.Steps = append(pfe.generatedScenario.Steps, step)
}

func (pfe *fuzzDexExecutor) executeTxStep(stepSnippet string) (*vmi.VMOutput, error) {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return nil, err
	}

	txStep, isTx := step.(*mj.TxStep)
	if !isTx {
		return nil, errors.New("tx step expected")
	}

	pfe.addStep(step)

	return pfe.arwenTestExecutor.ExecuteTxStep(txStep)
}

func (pfe *fuzzDexExecutor) log(info string, args ...interface{}) {
	fmt.Printf(info+"\n", args...)
}

func (pfe *fuzzDexExecutor) userAddress(userIndex int) string {
	return fmt.Sprintf("address:user%06d", userIndex)
}

func (pfe *fuzzDexExecutor) fullOfEsdtWalletString() string {
	esdtString := ""

	esdtString += fmt.Sprintf(`
						"str:%s": "1,000,000,000,000,000,000,000,000,000,000",`, pfe.wegldTokenId)
	esdtString += fmt.Sprintf(`
						"str:%s": "1,000,000,000,000,000,000,000,000,000,000",`, pfe.mexTokenId)
	esdtString += fmt.Sprintf(`
						"str:%s": "1,000,000,000,000,000,000,000,000,000,000",`, pfe.busdTokenId)
	esdtString += fmt.Sprintf(`
						"str:%s": "1,000,000,000,000,000,000,000,000,000,000",`, pfe.wemeLpTokenId)
	esdtString += fmt.Sprintf(`
						"str:%s": "1,000,000,000,000,000,000,000,000,000,000"`, pfe.webuLpTokenId)

	return esdtString
}

func (pfe *fuzzDexExecutor) querySingleResult(from, to, funcName, args string) ([][]byte, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%s",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "%s",
			"arguments": [
				%s
			],
			"gasLimit": "10,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [ "*" ],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		funcName,
		from,
		to,
		funcName,
		args,
	))
	if err != nil {
		return [][]byte{}, err
	}

	return output.ReturnData, nil
}

func (pfe *fuzzDexExecutor) querySingleResultStringAddr(from string, to string, funcName string, args string) ([][]byte, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%s",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "%s",
			"arguments": [
				%s
			],
			"gasLimit": "10,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [ "*" ],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		funcName,
		from,
		to,
		funcName,
		args,
	))
	if err != nil {
		return [][]byte{}, err
	}

	return output.ReturnData, nil
}

func (pfe *fuzzDexExecutor) increaseBlockEpoch(epochDelta int) error {
	currentBlockEpoch := uint32(0)
	if pfe.world.CurrentBlockInfo != nil {
		currentBlockEpoch = pfe.world.CurrentBlockInfo.BlockEpoch
	}

	err := pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"comment": "%d - increase block epoch",
		"currentBlockInfo": {
			"blockEpoch": "%d"
		}
	}`,
		pfe.nextTxIndex(),
		currentBlockEpoch+uint32(epochDelta),
	))
	if err != nil {
		return err
	}

	pfe.log("block epoch: %d ---> %d", currentBlockEpoch, currentBlockEpoch+uint32(epochDelta))
	return nil
}

func (pfe *fuzzDexExecutor) nextTxIndex() int {
	pfe.txIndex++
	return pfe.txIndex
}

func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

// This function allows equality with a += 1
func equalMatrix(left [][]byte, right [][]byte) bool {
	if len(left) != len(right) {
		return false
	}

	for i := 0; i < len(left); i++ {
		if !bytes.Equal(left[i], right[i]) {
			if i == len(left)-1 {
				leftIncreased := make([]byte, len(left[i]))
				copy(leftIncreased, left[i])
				if len(leftIncreased) > 0 {
					leftIncreased[len(leftIncreased)-1] += 1
				}

				rightIncreased := make([]byte, len(right[i]))
				copy(rightIncreased, right[i])
				if len(rightIncreased) > 0 {
					rightIncreased[len(rightIncreased)-1] += 1
				}

				return bytes.Equal(leftIncreased, right[i]) || bytes.Equal(left[i], rightIncreased)
			}
		}
	}

	return true
}