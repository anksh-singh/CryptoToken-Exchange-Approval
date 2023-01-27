package debankAPI

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
)

type IDebankAPI interface {
	GetPosition(request *pb.PositionRequest) (PositionResponse, error)
}

type DebankAPIService struct {
	env          *config.Config
	logger       *zap.SugaredLogger
	httpRequest  utils.IHttpRequest
	helper       *utils.Helpers
	coinGecko    coingecko.ICoinGecko
	util         *utils.UtilConf
	defaultToken []string
}

func NewDebankAPIService(env *config.Config, logger *zap.SugaredLogger, util *utils.UtilConf) *DebankAPIService {
	httRequest := utils.NewHttpRequest(logger)
	coingecko2 := coingecko.NewCoinGecko(env, logger, httRequest)
	helper2 := utils.Helpers{}
	return &DebankAPIService{
		env:         env,
		logger:      logger,
		httpRequest: httRequest,
		util:        util,
		coinGecko:   coingecko2,
		helper:      &helper2,
	}
}

// GetPosition get position info from debank
func (d *DebankAPIService) GetPosition(request *pb.PositionRequest) (PositionResponse, error) {
	var response PositionResponse
	var chainIds []string

	// if chain ids are empty then will set default chain ids
	if len(request.PositionPayload.ChainIds) == 0 {
		request.PositionPayload.ChainIds = supportedChainIds
	}

	// ProtocolIds
	var protocolIds []string
	if len(request.PositionPayload.ProtocolIds) > 0 {
		protocolIds = request.PositionPayload.ProtocolIds
		for _, protocolId := range protocolIds {
			protocolIds = append(protocolIds, strings.TrimSpace(protocolId))
		}
	}

	// Getting Debank ChainIds which we support
	for _, chainId := range request.PositionPayload.ChainIds {
		debankChaindId := DebankChainNetworkId[chainId]
		chainIds = append(chainIds, debankChaindId)
	}
	// Calling AllComplexProtocolList API
	url := d.env.DebankAPI.EndPoint + "user/all_complex_protocol_list?id=" + request.Address + "&chain_ids=" + strings.Join(chainIds, ",")
	body, err := d.httpRequest.GetRequestWithHeaders(url, "AccessKey", d.env.DebankAPI.AccessKey)
	if err != nil {
		d.logger.Error("Getting an error while fetching Debank API")
		return response, err
	}
	// Unmarshal the response
	var debankPositionResponse DebankPositionResponse
	err = json.Unmarshal(body, &debankPositionResponse)
	if err != nil {
		return response, fmt.Errorf("error while unmarshal debank position response: %v", err)
	}
	var chainWiseData = make(map[string]Positions)
	for _, item := range debankPositionResponse {
		if protocolIds != nil && len(protocolIds) > 0 {
			if slices.Contains(protocolIds, item.ID) {
				chainWiseData = d.GetChainWiseData(item, chainWiseData)
			}
		} else {
			chainWiseData = d.GetChainWiseData(item, chainWiseData)
		}
	}
	for key, value := range chainWiseData {
		data := Positions{
			Chain:     d.MapKey(DebankChainNetworkId, key),
			Address:   request.Address,
			Protocols: value.Protocols,
		}
		response.Positions = append(response.Positions, &data)
	}
	if len(response.Positions) == 0 {
		return response, fmt.Errorf("ChainIds / ProtocolIds not supported")
	}
	return response, nil
}

func (d *DebankAPIService) GetChainWiseData(item ResponseItem, chainWiseData map[string]Positions) map[string]Positions {
	portfolio, _ := d.GetPortfolioInfo(item.PortfolioItemList, item)
	overAllActionable, err := d.getOverlAllActionableValueForProtocol(portfolio)
	if err != nil {
		d.logger.Error(err)
		overAllActionable = false
	}
	protocols := d.GetProtocols(item, portfolio, overAllActionable)
	if chainWiseDataExists, ok := chainWiseData[item.Chain]; ok {
		chainWiseData[item.Chain] = Positions{
			Chain:     d.MapKey(DebankChainNetworkId, chainWiseDataExists.Chain),
			Protocols: append(chainWiseData[item.Chain].Protocols, protocols...),
		}
	} else {
		chainWiseData[item.Chain] = Positions{
			Chain:     d.MapKey(DebankChainNetworkId, item.Chain),
			Protocols: append(chainWiseData[item.Chain].Protocols, protocols...),
		}
	}
	return chainWiseData
}

// getOverlAllActionableValueForProtocol get overall actionable value for protocol
func (d *DebankAPIService) getOverlAllActionableValueForProtocol(portfolio *Portfolio) (bool, error) {
	if portfolio.Deposit != nil && portfolio.Deposit[0].IsActionable {
		return true, nil
	} else if portfolio.Locked != nil && portfolio.Locked[0].IsActionable {
		return true, nil
	} else if portfolio.Farm != nil && portfolio.Farm[0].IsActionable {
		return true, nil
	} else if portfolio.Investment != nil && portfolio.Investment[0].IsActionable {
		return true, nil
	} else if portfolio.Lending != nil && portfolio.Lending[0].IsActionable {
		return true, nil
	} else if portfolio.LiquidityPool != nil && portfolio.LiquidityPool[0].IsActionable {
		return true, nil
	} else if portfolio.OptionsSeller != nil && portfolio.OptionsSeller[0].IsActionable {
		return true, nil
	} else if portfolio.Rewards != nil && portfolio.Rewards[0].IsActionable {
		return true, nil
	} else if portfolio.Vesting != nil && portfolio.Vesting[0].IsActionable {
		return true, nil
	} else if portfolio.Yield != nil && portfolio.Yield[0].IsActionable {
		return true, nil
	} else if portfolio.Staked != nil && portfolio.Staked[0].IsActionable {
		return true, nil
	} else if portfolio.LeaveragedFarming != nil && portfolio.LeaveragedFarming[0].IsActionable {
		return true, nil
	} else if portfolio.OptionsBuyer != nil && portfolio.OptionsBuyer[0].IsActionable {
		return true, nil
	} else if portfolio.InsuranceSeller != nil && portfolio.InsuranceSeller[0].IsActionable {
		return true, nil
	} else if portfolio.InsuranceBuyer != nil && portfolio.InsuranceBuyer[0].IsActionable {
		return true, nil
	} else if portfolio.Perpetuals != nil && portfolio.Perpetuals[0].IsActionable {
		return true, nil
	} else if portfolio.Others != nil && portfolio.Others[0].IsActionable {
		return true, nil
	}
	return false, nil
}

func (d *DebankAPIService) GetProtocols(item ResponseItem, portfolio *Portfolio, overAllActionable bool) []*Protocols {
	var protocols []*Protocols
	protocols = append(protocols, &Protocols{
		ProtocolId:   item.ID,
		Name:         item.Name,
		SiteUrl:      item.SiteURL,
		LogoUrl:      item.LogoURL,
		IsActionable: overAllActionable,
		Portfolio:    portfolio,
	})
	return protocols
}

func (d *DebankAPIService) GetPortfolioInfo(portfolioItemList []PortfolioItem, responseItem ResponseItem) (*Portfolio, error) {
	var portfolioInfo Portfolio
	actionable := supportedProtocolActions[responseItem.ID]
	var protocolActions ProtocolActions
	// convert interface to string
	strObj := fmt.Sprintf("%v", actionable)
	err := json.Unmarshal([]byte(strObj), &protocolActions)
	if err != nil {
		d.logger.Error("error while unmarshalling protocol actions: ", err)
		protocolActions = ProtocolActions{}
	}
	for _, item := range portfolioItemList {
		switch item.Name {
		case LENDING:
			lending, netUsdValue := d.LendingInfo(item, protocolActions.Lending)
			portfolioInfo.Lending = append(portfolioInfo.Lending, lending...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case LIQUIDITY_POOL:
			liquidityPool, netUsdValue := d.LiquidityPoolInfo(item, protocolActions.LiquidityPool)
			portfolioInfo.LiquidityPool = append(portfolioInfo.LiquidityPool, liquidityPool...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case FARMING:
			farming, netUsdValue := d.FarmingInfo(item, protocolActions.Farming)
			portfolioInfo.Farm = append(portfolioInfo.Farm, farming...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case DEPOSIT:
			deposit, netUsdValue := d.DepositInfo(item, protocolActions.Deposit)
			portfolioInfo.Deposit = append(portfolioInfo.Deposit, deposit...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case INVESTMENT:
			inverstment, netUsdValue := d.InvestmentInfo(item, protocolActions.Investment)
			portfolioInfo.Investment = append(portfolioInfo.Investment, inverstment...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case OPTIONS_SELLER:
			optionsSeller, netUsdValue := d.OptionSellerInfo(item, protocolActions.OptionsSeller)
			portfolioInfo.OptionsSeller = append(portfolioInfo.OptionsSeller, optionsSeller...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case REWARDS:
			rewards, netUsdValue := d.RewardsInfo(item, protocolActions.Rewards)
			portfolioInfo.Rewards = append(portfolioInfo.Rewards, rewards...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case VESTING:
			vesting, netUsdValue := d.VestingInfo(item, protocolActions.Vesting)
			portfolioInfo.Vesting = append(portfolioInfo.Vesting, vesting...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case YIELD:
			yield, netUsdValue := d.YieldInfo(item, protocolActions.Yield)
			portfolioInfo.Yield = append(portfolioInfo.Yield, yield...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case STAKED:
			staked, netUsdValue := d.StakedInfo(item, protocolActions.Staked)
			portfolioInfo.Staked = append(portfolioInfo.Staked, staked...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case LOCKED:
			locked, netUsdValue := d.LockedInfo(item, protocolActions.Locked)
			portfolioInfo.Locked = append(portfolioInfo.Locked, locked...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case LEAVERAGEDFARMING:
			leaveragedFarming, netUsdValue := d.LeaveragedFarmingInfo(item, protocolActions.LeaveragedFarming)
			portfolioInfo.LeaveragedFarming = append(portfolioInfo.LeaveragedFarming, leaveragedFarming...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case OPTIONSBUYER:
			optionsBuyer, netUsdValue := d.OptionBuyerInfo(item, protocolActions.OptionsBuyer)
			portfolioInfo.OptionsBuyer = append(portfolioInfo.OptionsBuyer, optionsBuyer...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case INSURANCESELLER:
			insuranceSeller, netUsdValue := d.InsuranceSellerInfo(item, protocolActions.InsuranceSeller)
			portfolioInfo.InsuranceSeller = append(portfolioInfo.InsuranceSeller, insuranceSeller...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case INSURANCEBUYER:
			insuranceBuyer, netUsdValue := d.InsuranceBuyerInfo(item, protocolActions.InsuranceBuyer)
			portfolioInfo.InsuranceBuyer = append(portfolioInfo.InsuranceBuyer, insuranceBuyer...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		case PERPETUALS:
			perpetuals, netUsdValue := d.PerpetualsInfo(item, protocolActions.Perpetuals)
			portfolioInfo.Perpetuals = append(portfolioInfo.Perpetuals, perpetuals...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		default:
			others, netUsdValue := d.OthersInfo(item, protocolActions.Others)
			d.logger.Warn("New Action Found from the Debank Positions ==> ", item.Name)
			portfolioInfo.Others = append(portfolioInfo.Others, others...)
			portfolioInfo.NetUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(portfolioInfo.NetUsdValue)+d.helper.ConvertStringToFloat64(netUsdValue), 'f', -1, 64)
		}
	}
	// checking for empty portfolio for each action
	//portfolioInfo = d.GetEmptyPortfolioInfo(portfolioInfo, protocolActions)
	return &portfolioInfo, nil
}

func (d *DebankAPIService) GetNativeTokenAddress(chain string) string {
	chainName := d.MapKey(DebankChainNetworkId, chain)
	info := d.util.GetWalletInfo(chainName)
	return info.NativeTokenInfo.Address
}
func (d *DebankAPIService) MapKey(hashMap map[string]string, value string) string {
	var key string
	for k, v := range hashMap {
		if v == value {
			key = k
		}
	}
	return key
}

// GetEmptyPortfolioInfo get empty portfolio info
func (d *DebankAPIService) GetEmptyPortfolioInfo(portfolioInfo Portfolio, protocolActions ProtocolActions) Portfolio {
	portfolioInfo.Lending = append(portfolioInfo.Lending, d.GetLendingEmptyActionInfo(portfolioInfo.Lending, LENDING, protocolActions.Lending)...)
	portfolioInfo.LiquidityPool = append(portfolioInfo.LiquidityPool, d.GetLiquidityEmptyActionInfo(portfolioInfo.LiquidityPool, LIQUIDITY_POOL, protocolActions.LiquidityPool)...)
	portfolioInfo.Deposit = append(portfolioInfo.Deposit, d.GetDepositEmptyActionInfo(portfolioInfo.Deposit, DEPOSIT, protocolActions.Deposit)...)
	portfolioInfo.Farm = append(portfolioInfo.Farm, d.GetFarmEmptyActionInfo(portfolioInfo.Farm, FARMING, protocolActions.Farming)...)
	portfolioInfo.Investment = append(portfolioInfo.Investment, d.GetInvestmentEmptyActionInfo(portfolioInfo.Investment, INVESTMENT, protocolActions.Investment)...)
	portfolioInfo.OptionsSeller = append(portfolioInfo.OptionsSeller, d.GetOptionsSellerEmptyActionInfo(portfolioInfo.OptionsSeller, OPTIONS_SELLER, protocolActions.OptionsSeller)...)
	portfolioInfo.Rewards = append(portfolioInfo.Rewards, d.GetRewardsEmptyActionInfo(portfolioInfo.Rewards, REWARDS, protocolActions.Rewards)...)
	portfolioInfo.Vesting = append(portfolioInfo.Vesting, d.GetVestingEmptyActionInfo(portfolioInfo.Vesting, VESTING, protocolActions.Vesting)...)
	portfolioInfo.Yield = append(portfolioInfo.Yield, d.GetYieldEmptyActionInfo(portfolioInfo.Yield, YIELD, protocolActions.Yield)...)
	portfolioInfo.Staked = append(portfolioInfo.Staked, d.GetStakedEmptyActionInfo(portfolioInfo.Staked, STAKED, protocolActions.Staked)...)
	portfolioInfo.Locked = append(portfolioInfo.Locked, d.GetLockedEmptyActionInfo(portfolioInfo.Locked, LOCKED, protocolActions.Locked)...)
	portfolioInfo.LeaveragedFarming = append(portfolioInfo.LeaveragedFarming, d.GetLeaveragedFarmingEmptyActionInfo(portfolioInfo.LeaveragedFarming, LEAVERAGEDFARMING, protocolActions.LeaveragedFarming)...)
	portfolioInfo.InsuranceSeller = append(portfolioInfo.InsuranceSeller, d.GetInsuranceSellerEmptyActionInfo(portfolioInfo.InsuranceSeller, INSURANCESELLER, protocolActions.OptionsSeller)...)
	portfolioInfo.OptionsBuyer = append(portfolioInfo.OptionsBuyer, d.GetOptionsBuyerEmptyActionInfo(portfolioInfo.OptionsBuyer, OPTIONSBUYER, protocolActions.OptionsBuyer)...)
	portfolioInfo.InsuranceBuyer = append(portfolioInfo.InsuranceBuyer, d.GetInsuranceBuyerEmptyActionInfo(portfolioInfo.InsuranceBuyer, INSURANCEBUYER, protocolActions.InsuranceBuyer)...)
	portfolioInfo.Perpetuals = append(portfolioInfo.Perpetuals, d.GetPerpetualsEmptyActionInfo(portfolioInfo.Perpetuals, PERPETUALS, protocolActions.Perpetuals)...)
	portfolioInfo.Others = append(portfolioInfo.Others, d.GetOthersEmptyActionInfo(portfolioInfo.Others, OTHERS, protocolActions.Others)...)
	return portfolioInfo
}
