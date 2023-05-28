package workflow

import (
	"fmt"

	"encore.app/activity"
	"go.temporal.io/sdk/workflow"
)

var p *activity.Presenter

func Present(ctx workflow.Context, accountID, amount int) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, defaultActivityOptions)

	var result string
	err := workflow.ExecuteActivity(ctx, p.Present, accountID, amount).Get(ctx, &result)
	if err != nil {
		return "", fmt.Errorf("execute balance activity, err: %w", err)
	}

	if result == "accepted pending" {
		err := workflow.SignalExternalWorkflow(ctx, fmt.Sprintf("from %d %d", accountID, amount), "", fmt.Sprintf("approve %d %d", accountID, amount), struct{}{}).Get(ctx, nil)
		if err != nil {
			return "", fmt.Errorf("send present executed signal, err: %v", err)
		}
	}

	return "success", nil
}
