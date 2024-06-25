package main

import (
	"context"
	"strconv"
	"testing"

	"github.com/Halimao/NearStakerTool/sdk"
)

func TestGetAccounts(t *testing.T) {
	allValidators := make([]ValidatorFullData, 0)
	for page := 1; page < 100; page++ {
		t.Log("start query validators", "page", page)
		validators, err := getValidators(context.Background(), page)
		if err != nil {
			t.Error(err)
		}
		size := len(validators)
		if size == 0 {
			break
		}
		cumulativePercent, err := strconv.ParseFloat(validators[size-1].CumulativeStake.CumulativePercent, 64)
		if err != nil {
			t.Error(err)
		}
		allValidators = append(allValidators, validators...)
		if cumulativePercent > 96 {
			break
		}
	}
	validator := allValidators[0]
	stakeInfos, err := sdk.GetAccounts(context.Background(), validator.AccountID, 0, 1)
	if err != nil {
		t.Error(err)
	}
	t.Logf("acc: %s, amt: %.8f", stakeInfos[0].AccountID, FormatHumanityAmount(stakeInfos[0].StakedBalance))
	stakeInfos, err = sdk.GetAccounts(context.Background(), validator.AccountID, 1, 2)
	if err != nil {
		t.Error(err)
	}
	t.Logf("acc: %s, amt: %.8f", stakeInfos[0].AccountID, FormatHumanityAmount(stakeInfos[0].StakedBalance))
}
