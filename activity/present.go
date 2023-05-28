package activity

import (
	"context"
	"fmt"

	"encore.app/storage"
)

type Presenter struct {
	Storage storage.Storage
}

func (p *Presenter) Present(ctx context.Context, accountID, amount int) (string, error) {
	res, err := p.Storage.Present(accountID, amount)
	if err != nil {
		return "", fmt.Errorf("authorize transfer in db, err: %v", err)
	}

	return res, nil
}
