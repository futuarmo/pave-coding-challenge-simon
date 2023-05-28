package storage

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"

	"encore.app/types"
	tb "github.com/tigerbeetledb/tigerbeetle-go"
	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

const bank = "3"

var errNoPengingTransfer = errors.New("no pending transfer")

type TBStorage struct {
	client tb.Client
}

func NewTBStorage(cli tb.Client) (*TBStorage, error) {
	tbStor := &TBStorage{
		client: cli,
	}

	err := tbStor.insertSampleAccounts()
	if err != nil {
		return nil, fmt.Errorf("create sample accounts, err: %v", err)
	}

	return tbStor, nil
}

func (tb *TBStorage) Authorize(accountID, amount int) (string, error) {
	ID := genID(fmt.Sprintf("from %d %d", accountID, amount))

	transfer := []tb_types.Transfer{
		{
			ID:              ID,
			DebitAccountID:  uint128(bank),
			CreditAccountID: uint128(fmt.Sprintf("%d", accountID)),
			Amount:          uint64(amount),
			Ledger:          uint32(1),
			Code:            1,
			Flags: tb_types.TransferFlags{
				Pending: true,
			}.ToUint16(),
		},
	}

	res, err := tb.client.CreateTransfers(transfer)
	if err != nil {
		return "", fmt.Errorf("make transfer, err: %v", err)
	}

	return resultFrom(res[0])
}

func (tb *TBStorage) CancelAuthorize(accountID, amount int) (string, error) {
	pendingID := genID(fmt.Sprintf("from %d %d", accountID, amount))

	res, err := tb.client.LookupTransfers([]tb_types.Uint128{pendingID})
	if err != nil {
		return "", fmt.Errorf("lookup pending transfer, err: %v", err)
	}

	if len(res) == 0 {
		return "", errNoPengingTransfer
	}

	transfer := []tb_types.Transfer{
		{
			ID:        genID(fmt.Sprintf("cancel %d %d", accountID, amount)),
			PendingID: pendingID,
			Ledger:    uint32(1),
			Code:      1,
			Flags: tb_types.TransferFlags{
				VoidPendingTransfer: true,
			}.ToUint16(),
		},
	}

	result, err := tb.client.CreateTransfers(transfer)
	if err != nil {
		return "", fmt.Errorf("make transfer, err: %v", err)
	}

	return resultFrom(result[0])
}

func (tb *TBStorage) Present(accountID, amount int) (string, error) {
	pendingID := genID(fmt.Sprintf("from %d %d", accountID, amount))

	res, err := tb.client.LookupTransfers([]tb_types.Uint128{pendingID})
	if err != nil {
		return "", fmt.Errorf("lookup pending transfer, err: %v", err)
	}

	if isPended := len(res) == 1; isPended {
		transfer := []tb_types.Transfer{
			{
				ID:        genID(fmt.Sprintf("approve %d %d", accountID, amount)),
				PendingID: pendingID,
				Ledger:    uint32(1),
				Code:      1,
				Flags: tb_types.TransferFlags{
					PostPendingTransfer: true,
				}.ToUint16(),
			},
		}

		res, err := tb.client.CreateTransfers(transfer)
		if err != nil {
			return "", fmt.Errorf("make transfer, err: %v", err)
		}

		resultStr, err := resultFrom(res[0])
		if resultStr == "success" {
			resultStr = "accepted pending"
		}

		return resultStr, err
	}

	// Presentments can come in without authorizations,
	// but should be declined if there’s not available fund

	account, err := tb.getAccount(amount)
	if err != nil {
		return "", fmt.Errorf("get account info, err: %v", err)
	}

	if account.CreditsPosted < uint64(amount) {
		return "", fmt.Errorf("not enough money, need %d, has %d", amount, account.CreditsPosted)
	}

	transfer := []tb_types.Transfer{
		{
			ID:              genID(fmt.Sprintf("from %d %d", accountID, amount)),
			DebitAccountID:  uint128(bank),
			CreditAccountID: uint128(fmt.Sprintf("%d", accountID)),
			Ledger:          uint32(1),
			Code:            1,
			Amount:          uint64(amount),
		},
	}

	result, err := tb.client.CreateTransfers(transfer)
	if err != nil {
		return "", fmt.Errorf("make transfer, err: %v", err)
	}

	return resultFrom(result[0])
}

func (tb *TBStorage) Balance(accountID int) (types.AccountBalance, error) {
	acc, err := tb.getAccount(accountID)
	if err != nil {
		return types.AccountBalance{}, fmt.Errorf("get account info, err: %v", err)
	}

	return types.AccountBalance{
		CreditsPosted:  int(acc.CreditsPosted),
		CreditsPending: int(acc.CreditsPending),
		CreditsTotal:   int(acc.CreditsPosted + acc.CreditsPending),
		DebitsPosted:   int(acc.DebitsPosted),
		DebitsPending:  int(acc.DebitsPending),
		DebitsTotal:    int(acc.DebitsPosted + acc.DebitsPending),
	}, nil
}

func (tb *TBStorage) Close() {
	tb.client.Close()
}

func (tb *TBStorage) insertSampleAccounts() error {
	_, err := tb.client.CreateAccounts([]tb_types.Account{
		{
			ID:             uint128("1"),
			UserData:       tb_types.Uint128{},
			Reserved:       [48]uint8{},
			Ledger:         1,
			Code:           1,
			Flags:          0,
			DebitsPending:  0,
			DebitsPosted:   100,
			CreditsPending: 0,
			CreditsPosted:  100,
			Timestamp:      0,
		},
		{
			ID:             uint128("2"),
			UserData:       tb_types.Uint128{},
			Reserved:       [48]uint8{},
			Ledger:         1,
			Code:           1,
			Flags:          0,
			DebitsPending:  0,
			DebitsPosted:   100,
			CreditsPending: 0,
			CreditsPosted:  100,
			Timestamp:      0,
		},
		{
			ID:             uint128(bank),
			UserData:       tb_types.Uint128{},
			Reserved:       [48]uint8{},
			Ledger:         1,
			Code:           1,
			Flags:          0,
			DebitsPending:  0,
			DebitsPosted:   100,
			CreditsPending: 0,
			CreditsPosted:  100,
			Timestamp:      0,
		},
	})

	return err
}

func (tb *TBStorage) getAccount(accountID int) (tb_types.Account, error) {
	acc, err := tb.client.LookupAccounts([]tb_types.Uint128{uint128(strconv.Itoa(accountID))})
	if err != nil {
		return tb_types.Account{}, fmt.Errorf("get account balance, err: %v", err)
	}

	if count := len(acc); count != 1 {
		return tb_types.Account{}, fmt.Errorf("found %d accounts instead 1", count)
	}

	return acc[0], nil
}

func uint128(value string) tb_types.Uint128 {
	x, err := tb_types.HexStringToUint128(value)
	if err != nil {
		panic(err)
	}
	return x
}

func genID(idStr string) tb_types.Uint128 {
	return uint128(fmt.Sprintf("%x", md5.Sum([]byte(idStr))))
}

func resultFrom(transferResult tb_types.TransferEventResult) (string, error) {
	if transferResult.Result == tb_types.TransferOK {
		return "success", nil
	} else {
		return "", fmt.Errorf("%s", transferResult.Result.String())
	}
}
