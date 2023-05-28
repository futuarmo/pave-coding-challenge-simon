package activity

import (
	"context"
	"fmt"

	"encore.app/storage"
)

type Authorizator struct {
	Storage storage.Storage
}

func (a *Authorizator) Authorize(ctx context.Context, accountID, amount int) (string, error) {
	res, err := a.Storage.Authorize(accountID, amount)
	if err != nil {
		return "", fmt.Errorf("authorize transfer, err: %v", err)
	}

	return res, nil
}

func (a *Authorizator) CancelAuthorize(ctx context.Context, accountID, amount int) (string, error) {
	res, err := a.Storage.CancelAuthorize(accountID, amount)
	if err != nil {
		return "", fmt.Errorf("cancel authorized transfer in db, err: %v", err)
	}

	return res, nil
}
