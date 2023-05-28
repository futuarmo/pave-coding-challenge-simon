package ledger

import (
	"context"
	"fmt"

	"encore.app/workflow"
	"encore.dev/rlog"
	"go.temporal.io/sdk/client"
)

type AuthorizeResponse struct {
	Result string
}

type AuthorizeParams struct {
	AccountID int
	Amount    int
}

//encore:api public method=POST path=/authorize
func (s *Service) Authorize(ctx context.Context, p *AuthorizeParams) (*AuthorizeResponse, error) {
	options := client.StartWorkflowOptions{
		TaskQueue: ledgerTaskQueue,
	}

	we, err := s.client.ExecuteWorkflow(ctx, options, workflow.Authorize, p.AccountID, p.Amount)
	if err != nil {
		return nil, fmt.Errorf("execute authorize workflow, err: %v", err)
	}
	rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

	err = we.Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("get workflow status, err: %v", err)
	}
	return &AuthorizeResponse{Result: "request processed"}, nil
}
