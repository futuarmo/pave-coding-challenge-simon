package ledger

import (
	"context"
	"fmt"

	"encore.app/types"
	"encore.app/workflow"
	"encore.dev/rlog"
	"go.temporal.io/sdk/client"
)

type BalanceParams struct {
	AccountID int
}

//encore:api public method=POST path=/balance
func (s *Service) Balance(ctx context.Context, p *BalanceParams) (*types.AccountBalance, error) {
	options := client.StartWorkflowOptions{
		TaskQueue: ledgerTaskQueue,
	}
	we, err := s.client.ExecuteWorkflow(ctx, options, workflow.Balance, p.AccountID)
	if err != nil {
		return nil, fmt.Errorf("execute balance workflow, err: %v", err)
	}
	rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

	var result types.AccountBalance
	err = we.Get(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("get workflow status, err: %v", err)
	}
	return &result, nil
}
