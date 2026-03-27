package orderbook

import (
	"context"
	"fmt"
	"math/big"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (api *api) GetOrderByMev(ctx context.Context, params GetOrderParams) (*OrderResponse, error) {
	u := fmt.Sprintf("/orderbook/v4.1/%d/order/%s", api.chainId, params.OrderHash)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
	}

	var getOrderByHashResponse *OrderResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &getOrderByHashResponse)
	if err != nil {
		return nil, err
	}

	return getOrderByHashResponse, nil
}

type OrdersByMev struct {
	Meta struct {
		HasMore    bool   `json:"hasMore"`
		NextCursor string `json:"nextCursor"`
		Count      int    `json:"count"`
	} `json:"meta"`
	Items []OrderResponse `json:"items"`
}

func (api *api) GetAllOrdersByMev(ctx context.Context, params GetAllOrdersParams) (*OrdersByMev, error) {
	u := fmt.Sprintf("/orderbook/v4.1/%d/all", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
	}

	var allOrdersResponse OrdersByMev
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &allOrdersResponse)
	if err != nil {
		return nil, err
	}

	return &allOrdersResponse, nil
}

func NormalizeResponse(resp OrderResponse) (*NormalizedLimitOrderData, error) {
	saltBigInt, ok := new(big.Int).SetString(resp.Data.Salt, 10)
	if !ok {
		saltBigInt, ok = new(big.Int).SetString(resp.Data.Salt[2:], 16)
		if !ok {
			return nil, fmt.Errorf("invalid salt: %s", resp.Data.Salt)
		}
	}
	makingAmountBigInt, ok := new(big.Int).SetString(resp.Data.MakingAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid making amount: %s", resp.Data.MakingAmount)
	}
	takingAmountBigInt, ok := new(big.Int).SetString(resp.Data.TakingAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid taking amount: %s", resp.Data.TakingAmount)
	}
	makerAssetBigInt := AddressStringToBigInt(resp.Data.MakerAsset)
	takerAssetBigInt := AddressStringToBigInt(resp.Data.TakerAsset)
	makerBigInt := AddressStringToBigInt(resp.Data.Maker)
	receiverBigInt := AddressStringToBigInt(resp.Data.Receiver)
	makerTraits, err := hexutil.DecodeBig(resp.Data.MakerTraits)
	if err != nil {
		return nil, fmt.Errorf("invalid maker traits: %w", err)
	}
	return &NormalizedLimitOrderData{
		Salt:         saltBigInt,
		MakerAsset:   makerAssetBigInt,
		TakerAsset:   takerAssetBigInt,
		Maker:        makerBigInt,
		Receiver:     receiverBigInt,
		MakingAmount: makingAmountBigInt,
		TakingAmount: takingAmountBigInt,
		MakerTraits:  makerTraits,
	}, nil
}
