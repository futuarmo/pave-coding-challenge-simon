package storage

import (
	"encore.app/types"
)

type Storage interface {
	Authorize(accountID, amount int) (string, error)
	CancelAuthorize(accountID, amount int) (string, error)
	Present(accountID, amount int) (string, error)
	Balance(accountID int) (types.AccountBalance, error)
	Close()
}
