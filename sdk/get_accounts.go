package sdk

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
)

type Page struct {
	FromIndex int `json:"from_index"`
	Limit     int `json:"limit"`
}

type GetAccountsReq struct {
	RequestType string `json:"request_type"`
	AccountID   string `json:"account_id"`
	MethodName  string `json:"method_name"`
	ArgsBase64  string `json:"args_base64"`
	Finality    string `json:"finality"`
}

type GetAccountsResult struct {
	Result      []byte `json:"result"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
}

type AccountStakeInfo struct {
	AccountID       string `json:"account_id"`
	UnstakedBalance string `json:"unstaked_balance"`
	StakedBalance   string `json:"staked_balance"`
	CanWithdraw     bool   `json:"can_withdraw"`
}

func GetAccounts(ctx context.Context, validatorID string, from, limit int) ([]AccountStakeInfo, error) {
	pageReq := Page{
		FromIndex: from,
		Limit:     limit,
	}
	pageByte, _ := json.Marshal(pageReq)
	req := GetAccountsReq{
		RequestType: "call_function",
		AccountID:   validatorID,
		MethodName:  "get_accounts",
		ArgsBase64:  base64.StdEncoding.EncodeToString(pageByte),
		Finality:    "optimistic",
	}
	response, err := defaultClient.Call(ctx, "query", &req)
	if err != nil {
		slog.With("err", err, "req", req).Error("GetAccounts error")
		return nil, err
	}
	if response.Error != nil {
		slog.With("err", response.Error.Error(), "p", req).Error("GetAccounts error")
		return nil, err
	}

	var ret *GetAccountsResult
	err = response.GetObject(&ret)
	if err != nil || ret == nil {
		slog.With("err", err, "req", req, "ret", ret).Error("GetAccounts GetObject error")
		return nil, err
	}
	accStakeInfos := make([]AccountStakeInfo, 0)
	err = json.Unmarshal(ret.Result, &accStakeInfos)
	if err != nil {
		slog.With("err", err, "req", req, "ret", string(ret.Result)).Error("GetAccounts Unmarshal error")
		return nil, err
	}
	return accStakeInfos, nil
}
