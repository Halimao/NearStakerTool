package main

import (
	"context"
	"database/sql"
	"log/slog"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/Halimao/NearStakerTool/sdk"
	_ "modernc.org/sqlite"
)

var (
	sqliteDB *sql.DB
)

type StakeInfo struct {
	AccID       string  `json:"acc_id"`
	ValidatorID string  `json:"validator_id"`
	Amount      float64 `json:"amount"`
}

type SimpleValidatorInfo struct {
	ValidatorID string  `json:"validator_id"`
	TotalStaker int     `json:"total_staker"`
	TotalStaked float64 `json:"total_staked"`
}

func init() {
	var err error
	sqliteDB, err = sql.Open("sqlite", "stake.db")
	if err != nil {
		panic(err)
	}
	_, err = sqliteDB.Exec(`CREATE TABLE IF NOT EXISTS "stake" (
		"acc_id" VARCHAR(200) NOT NULL,
		"validator_id" VARCHAR(200) NOT NULL,
		"amount" FLOAT DEFAULT 0,
		PRIMARY KEY(acc_id,validator_id))`)
	if err != nil {
		panic(err)
	}
}

func main() {
	allValidators := make([]ValidatorFullData, 0)
	for page := 1; page < 100; page++ {
		slog.Info("start query validators", "page", page)
		validators, err := getValidators(context.Background(), page)
		if err != nil {
			panic(err)
		}
		size := len(validators)
		if size == 0 {
			break
		}
		cumulativePercent, err := strconv.ParseFloat(validators[size-1].CumulativeStake.CumulativePercent, 64)
		if err != nil {
			slog.With("err", err, "val", validators[size-1].CumulativeStake.CumulativePercent).Error("ParseFloat error")
			panic(err)
		}
		allValidators = append(allValidators, validators...)
		if cumulativePercent > 96 {
			break
		}
		slog.Info("end query validators", "page", page)
		time.Sleep(1 * time.Second)
	}
	allValidatorIDs := make([]string, len(allValidators))
	allValidatorSimpleInfos := make([]SimpleValidatorInfo, len(allValidators))
	for i, validator := range allValidators {
		allValidatorIDs[i] = validator.AccountID
		allValidatorSimpleInfos[i] = SimpleValidatorInfo{
			ValidatorID: validator.AccountID,
			TotalStaker: validator.PoolInfo.DelegatorsCount,
			TotalStaked: FormatHumanityAmount(validator.CurrentEpoch.Stake),
		}
	}
	slog.With("ids", allValidatorIDs, "num", len(allValidatorIDs)).Info("all need check validators")
	slog.With("summary", allValidatorSimpleInfos).Info("all need check validators")
	for idx, validator := range allValidators {
		validatorID := validator.AccountID
		slog.With("validator", validatorID, "idx", idx).Info("start query validator stake info")
		totalNum := 0
		from, limit := 0, 100
		for {
			slog.With("validator", validatorID, "idx", idx, "from", from, "limit", limit).Info("query validator stake accounts")
			accStakeInfos, err := sdk.GetAccounts(context.Background(), validatorID, from, limit)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			for _, stakeInfo := range accStakeInfos {
				totalNum = totalNum + 1
				_, err = sqliteDB.Exec("REPLACE INTO `stake`(acc_id, validator_id, amount) values(?, ?, ?)", stakeInfo.AccountID, validatorID, FormatHumanityAmount(stakeInfo.StakedBalance))
				if err != nil {
					slog.With("err", err, "acc", stakeInfo.AccountID, "validator", validatorID, "amt", stakeInfo.StakedBalance).Error("insert into db error")
					continue
				}
			}
			from = from + limit
			if len(accStakeInfos) < limit {
				break
			}
			time.Sleep(1 * time.Second)
		}
		slog.With("validator", validatorID, "idx", idx, "totalStakeAccNum", totalNum).Info("end query validator stake info")
	}
}

func FormatHumanityAmount(amt string) float64 {
	amtF, _ := new(big.Float).SetString(amt)
	hamanAmtF := new(big.Float).Quo(amtF, big.NewFloat(math.Pow10(24)))
	hamanAmt, _ := hamanAmtF.Float64()
	return hamanAmt
}
