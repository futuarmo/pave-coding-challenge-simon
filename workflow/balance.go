package workflow

import (
	"fmt"

	"encore.app/activity"
	"encore.app/types"
	"go.temporal.io/sdk/workflow"
)

var b *activity.Balancer

func Balance(ctx workflow.Context, accountID int) (types.AccountBalance, error) {
	ctx = workflow.WithActivityOptions(ctx, defaultActivityOptions)

	var balance types.AccountBalance
	err := workflow.ExecuteActivity(ctx, b.Balance, accountID).Get(ctx, &balance)
	if err != nil {
		return types.AccountBalance{}, fmt.Errorf("execute balance activity, err: %w", err)
	}

	return balance, nil
}
