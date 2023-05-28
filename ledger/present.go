package ledger

import (
	"context"
	"fmt"

	"encore.app/workflow"
	"encore.dev/rlog"
	"go.temporal.io/sdk/client"
)

type PresentResponse struct {
	Result string
}

type PresentParams struct {
	AccountID int
	Amount    int
}

//encore:api public method=POST path=/present
func (s *Service) Present(ctx context.Context, p *PresentParams) (*PresentResponse, error) {
	options := client.StartWorkflowOptions{
		TaskQueue: ledgerTaskQueue,
	}

	we, err := s.client.ExecuteWorkflow(ctx, options, workflow.Present, p.AccountID, p.Amount)
	if err != nil {
		return nil, fmt.Errorf("execute present workflow, err: %v", err)
	}
	rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

	var result string
	err = we.Get(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("get workflow status, err: %v", err)
	}
	return &PresentResponse{Result: result}, nil
}
