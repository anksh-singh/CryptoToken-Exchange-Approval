package rpc

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/adapters/evm/application/core"
	"bridge-allowance/internal/common"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/pkg/jsonrpc"
	"bridge-allowance/utils"
	// "bridge-allowance/utils/models"
	// "encoding/json"
	"fmt"
	"github.com/umbracle/ethgo/builtin/erc20"
	// "math"
	"math/big"
	"strconv"
	"strings"
	// "sync"

	// "bridge-allowance/pkg/unmarshal"
	"github.com/onrik/ethrpc"
	"github.com/umbracle/ethgo"
	_ "github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	_ "github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
	_ "github.com/umbracle/ethgo/contract"
	ethgoJsonRPC "github.com/umbracle/ethgo/jsonrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EvmCore interface {
	GetTokenAllowance(request *pb.AllowanceRequest) (*pb.AllowanceResponse, error)
	
}


type evmCore struct {
	rpc         map[string]*ethrpc.EthRPC
	env         *config.Config
	logger      *zap.SugaredLogger
	services    common.Services
	util        *utils.UtilConf
	ethgoRpc    map[string]*ethgoJsonRPC.Client
	httpRequest utils.IHttpRequest
	rpcHandler  *jsonrpc.RPCHandler
}

// NewEVMCore Manager to initialize EVM specific rpc node endpoints
func NewEVMCore(config *config.Config, logger *zap.SugaredLogger, services common.Services) *evmCore {
	rpc := make(map[string]*ethrpc.EthRPC)
	ethgoRpc := make(map[string]*ethgoJsonRPC.Client)
	//Do not continue if no EVM configurations are provided
	logger.Info("Supported EVM chains:", len(config.EVM.Cfg.Wallets))
	if len(config.EVM.Cfg.Wallets) < 1 {
		logger.Fatal("No EVM wallet configurations found")
	}
	//Initialize EVM RPC configurations
	for i, w := range config.EVM.Cfg.Wallets {
		i++
		rpc[w.ChainName] = ethrpc.New(w.RPC)
		var err error
		ethgoRpc[w.ChainName], err = ethgoJsonRPC.NewClient(w.RPC)
		if err != nil {
			logger.Errorf(err.Error())
			logger.Fatalf("Error initializing go RPC client for `%s` chain", w.ChainName)
		}
		logger.Infof("%v. %v EVM chain initialized with configurations %v", i, w.ChainName, w)
	}
	utilsManager := utils.NewUtils(logger, config)
	httpRequest := utils.NewHttpRequest(logger)
	rpcHandler := jsonrpc.NewJsonRPCHandler(config, logger, httpRequest)
	return &evmCore{rpc, config, logger, services, utilsManager, ethgoRpc,
		httpRequest, rpcHandler}
}


// Clean the given argument to remove nil references
func Clean(arg interface{}) interface{} {
	if arg == nil {
		return ""
	} else {
		return arg
	}
}

var TransactionTypes = func() map[string]string {
	return map[string]string{
		"swap":         "swap",
		"transfer":     "send",
		"mint":         "addLiquidity",
		"addLiquidity": "addLiquidity",
		"withdrawal":   "withdraw",
		"withdrawn":    "withdraw",
		"withdraw":     "withdraw",
		"approval":     "approve",
	}
}


// getEthTokenDecimals retrieve token decimals from a contract address
func (evm *evmCore) getEthTokenDecimals(contractAddr string, chain string) (int, error) {
	erc20 := erc20.NewERC20(ethgo.HexToAddress(contractAddr), contract.WithJsonRPC(evm.ethgoRpc[chain].Eth()))
	decimals, err := erc20.Decimals()
	if err != nil {
		evm.logger.Info("Error while fetching token decimals", err.Error())
		decimals = 0
	}
	return int(decimals), err
}

func (evm *evmCore) createContractABI(request *pb.GasLimitRequest) ([]byte, error) {
	transferFunction := []string{core.AbiTransferFunction}
	abiContract, err := abi.NewABIFromList(transferFunction)
	if err != nil {
		evm.logger.Error("Method: transfer not found")
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	addr := ethgo.HexToAddress(request.To)
	contractInstance := contract.NewContract(addr, abiContract, contract.WithJsonRPC(evm.ethgoRpc[request.Chain].Eth()))
	method := contractInstance.GetABI().GetMethod("transfer")
	if method == nil {
		evm.logger.Error("Method: transfer not found")
	}
	data, err := method.Encode(request.Value)
	return data, nil
}

// optContractABI get opL1 fee from contact abi
func (evm *evmCore) optContractABI(chain string) float64 {
	if chain != "optimism" {
		//L1 fee not applicable to chains other than optimism
		return 0.0
	}
	abiInstance, err := abi.NewABI(core.OpGasPriceOracleABI)
	if err != nil {
		evm.logger.Errorf("Error generating ABI interface: %v", err)
		return core.DefaultOPL1Fee
	}
	contractInstance := contract.NewContract(ethgo.HexToAddress("0x420000000000000000000000000000000000000F"), abiInstance, contract.WithJsonRPC(evm.ethgoRpc[chain].Eth()))
	method := contractInstance.GetABI().GetMethod(core.OpL1FeeABIMethod)
	if method == nil {
		evm.logger.Error("Error fetching method: getL1Fee")
		return core.DefaultOPL1Fee
	}
	fee, err := contractInstance.Call("getL1Fee", ethgo.Latest, "0x")
	//ABI contract should return a map with only one entry
	if fee == nil || len(fee) != 1 {
		return core.DefaultOPL1Fee
	}
	var formattedFee = core.DefaultOPL1Fee
	for _, v := range fee {
		value := v.(*big.Int)
		//Convert to Gwei
		etherFee := new(big.Float).Quo(new(big.Float).SetInt(value), new(big.Float).SetFloat64(utils.Ether))
		//Multiply by a delta to handle higher fee
		etherFeeMulDelta := new(big.Float).Mul(etherFee, new(big.Float).SetFloat64(2.0))
		finalFee, _ := etherFeeMulDelta.Float64()
		finalFeeStr := strconv.FormatFloat(finalFee, 'f', 5, 64)
		formattedFee, _ = strconv.ParseFloat(finalFeeStr, 64)
	}
	if formattedFee > 0.0 {
		//Handle precision loss
		return formattedFee
	}
	return core.DefaultOPL1Fee
}

func (evm *evmCore) contractABI(request ContractABIRequest) (string, error) {
	abiContract, err := abi.NewABI(core.TokenABI)
	if err != nil {
		evm.logger.Errorf("Error in generating ABI Interface %v", err)
		return "", err
	}
	addr := ethgo.HexToAddress(request.Contract) //contract address
	c := contract.NewContract(addr, abiContract, contract.WithJsonRPC(evm.ethgoRpc[request.Chain].Eth()))
	switch request.Method {
	case "approve":
		//write call
		n := new(big.Int)
		value, ok := n.SetString(request.Data, 10)
		if !ok {
			return "", err
		}
		method := c.GetABI().GetMethod("approve")
		if method == nil {
			evm.logger.Error("Method: approve not found")
		}
		data, err := method.Encode(map[string]interface{}{
			"_spender": ethgo.HexToAddress(request.To),
			"_value":   value,
		})
		if err != nil {
			evm.logger.Error(err)
		}
		return fmt.Sprintf("0x%x", data), nil
	case "allowance":
		//Read call
		res, err := c.Call("allowance", ethgo.Latest, ethgo.HexToAddress(request.From), ethgo.HexToAddress(request.To))
		if err != nil {
			evm.logger.Error(err)
			return "", status.Errorf(codes.Internal, err.Error(), "Internal Error")
		}
		return res["amount"].(*big.Int).String(), nil
	default:
		return "", status.Errorf(codes.InvalidArgument, "", "Unsupported Method")
	}
}

func (evm *evmCore) isChainNativeToken(tokenAddress string, chain string) bool {
	var isNativeChain bool
	for _, c := range evm.env.EVM.Cfg.Wallets {
		if chain == c.ChainName && tokenAddress == c.NativeTokenInfo.Address {
			isNativeChain = true
			return isNativeChain
		} else {
			isNativeChain = false
		}
	}
	return isNativeChain
}

func (evm *evmCore) GetTokenAllowance(request *pb.AllowanceRequest) (*pb.AllowanceResponse, error) {
	allowanceForNativeToken := 999999999999 //standard value for native token
	var allowance string
	if request.Chain == "xinfin" {
		request.Contract = evm.util.ResolveXDCAddress(request.Contract)
		request.Owner = evm.util.ResolveXDCAddress(request.Owner)
		if request.Contract == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
			request.Contract = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
	}
	isNativeChain := evm.isChainNativeToken(strings.ToLower(request.Contract), request.Chain)
	if isNativeChain {
		allowance = strconv.Itoa(allowanceForNativeToken)
	} else {
		data, err := evm.contractABI(ContractABIRequest{
			From:     request.Owner,
			To:       request.Spender,
			Contract: request.Contract,
			Chain:    request.Chain,
			Method:   "allowance",
		})
		if err != nil {
			return nil, err
		}
		allowance = data
	}
	return &pb.AllowanceResponse{
		Allowance: evm.util.ToDecimal(allowance, 18).String(),
	}, nil
}
