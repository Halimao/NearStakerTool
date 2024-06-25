package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type ValidatorFullData struct {
	AccountID    string `json:"accountId"`
	CurrentEpoch struct {
		Progress struct {
			Blocks struct {
				Produced int `json:"produced"`
				Total    int `json:"total"`
			} `json:"blocks"`
			Chunks struct {
				Produced int `json:"produced"`
				Total    int `json:"total"`
			} `json:"chunks"`
		} `json:"progress"`
		Stake string `json:"stake"`
	} `json:"currentEpoch"`
	PublicKey string `json:"publicKey"`
	NextEpoch struct {
		Stake string `json:"stake"`
	} `json:"nextEpoch"`
	AfterNextEpoch struct {
		Stake string `json:"stake"`
	} `json:"afterNextEpoch,omitempty"`
	ContractStake interface{} `json:"contractStake"`
	Description   interface{} `json:"description"`
	Index         int         `json:"index"`
	PoolInfo      struct {
		DelegatorsCount int `json:"delegatorsCount"`
		Fee             struct {
			Numerator   int `json:"numerator"`
			Denominator int `json:"denominator"`
		} `json:"fee"`
	} `json:"poolInfo"`
	CumulativeStake struct {
		AccumulatedPercent string `json:"accumulatedPercent"`
		CumulativePercent  string `json:"cumulativePercent"`
		OwnPercentage      string `json:"ownPercentage"`
	} `json:"cumulativeStake"`
	Percent     string `json:"percent"`
	ShowWarning bool   `json:"showWarning"`
	StakeChange struct {
		Symbol string `json:"symbol"`
		Value  string `json:"value"`
	} `json:"stakeChange"`
	StakingStatus string `json:"stakingStatus"`
	Warning       string `json:"warning"`
}

type ValidatorsRsp struct {
	CurrentValidators int     `json:"currentValidators"`
	ElapsedTimeData   float64 `json:"elapsedTimeData"`
	EpochProgressData float64 `json:"epochProgressData"`
	EpochStatsCheck   string  `json:"epochStatsCheck"`
	LatestBlock       struct {
		Height    int   `json:"height"`
		Timestamp int64 `json:"timestamp"`
	} `json:"latestBlock"`
	Total             int                 `json:"total"`
	TotalSeconds      int                 `json:"totalSeconds"`
	TotalStake        string              `json:"totalStake"`
	ValidatorFullData []ValidatorFullData `json:"validatorFullData"`
}

func getValidators(ctx context.Context, page int) ([]ValidatorFullData, error) {
	url := fmt.Sprintf("https://api3.nearblocks.io/v1/validators?page=%d", page)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("origin", "https://nearblocks.io")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://nearblocks.io/")
	req.Header.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Linux\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.With("err", err, "page", page).Error("getValidators error")
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.With("err", err, "page", page).Error("getValidators read body error")
		return nil, err
	}
	var rsp ValidatorsRsp
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		slog.With("err", err, "page", page).Error("getValidators Unmarshal error")
		return nil, err
	}
	return rsp.ValidatorFullData, nil
}
