package activity

import (
	"context"
	"fmt"

	"encore.app/storage"
	"encore.app/types"
)

type Balancer struct {
	Storage storage.Storage
}

func (b *Balancer) Balance(ctx context.Context, accountID int) (types.AccountBalance, error) {
	res, err := b.Storage.Balance(accountID)
	if err != nil {
		return types.AccountBalance{}, fmt.Errorf("get account balance from storage, err: %v", err)
	}

	return res, nil
}
