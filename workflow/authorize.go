package workflow

import (
	"fmt"
	"time"

	"encore.app/activity"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

const authorizeTimeout = 100 * time.Second

var a *activity.Authorizator

func Authorize(ctx workflow.Context, accountID, amount int) error {
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID:        fmt.Sprintf("from %d %d", accountID, amount),
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}

	childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)

	future := workflow.ExecuteChildWorkflow(childCtx, WaitForPresent, accountID, amount)

	_ = future.GetChildWorkflowExecution().Get(childCtx, nil)

	return nil
}

func WaitForPresent(ctx workflow.Context, accountID, amount int) (string, error) {
	approveChan := workflow.GetSignalChannel(ctx, fmt.Sprintf("approve %d %d", accountID, amount))

	var result string
	var err error

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(approveChan, func(c workflow.ReceiveChannel, _ bool) {
		var approved struct{}
		c.Receive(ctx, &approved)
	})
	selector.AddFuture(workflow.NewTimer(ctx, authorizeTimeout), func(f workflow.Future) {
		result, err = cancelAuthorize(ctx, accountID, amount)
	})

	return result, err
}

func cancelAuthorize(ctx workflow.Context, accountID, amount int) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, defaultActivityOptions)

	var cancelResult string
	err := workflow.ExecuteActivity(ctx, a.CancelAuthorize, accountID, amount).Get(ctx, &cancelResult)
	if err != nil {
		return "", fmt.Errorf("execute cancel authorize activity, err: %w", err)
	}

	return cancelResult, nil
}

// this is implementation with child workflows, is not used
/* func Authorize(ctx workflow.Context, accountID, amount int) (string, error) {
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID:        fmt.Sprintf("from %d %d", accountID, amount),
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}

	childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)

	future := workflow.ExecuteChildWorkflow(childCtx, TimeoutHandledAuthorize, accountID, amount)
	_ = future.GetChildWorkflowExecution().Get(childCtx, nil)
	var result string

	return result, nil
}

func TimeoutHandledAuthorize(ctx workflow.Context, accountID, amount int) (string, error) {
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID:               fmt.Sprintf("from %d %d", accountID, amount),
		ParentClosePolicy:        enums.PARENT_CLOSE_POLICY_TERMINATE,
		WorkflowExecutionTimeout: 5 * time.Second,
	}
	childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)

	future := workflow.ExecuteChildWorkflow(childCtx, InnerAuthorize, accountID, amount)

	var result string
	err := future.GetChildWorkflowExecution().Get(childCtx, &result)
	if err != nil {
		if temporal.IsTimeoutError(err) {
			cancelResult, err := cancelAuthorize(ctx, accountID, amount)
			if err != nil {
				return "", fmt.Errorf("execute cancel authorize activity, err: %w", err)
			}

			return cancelResult, nil
		}

		return "", fmt.Errorf("execute workflow with timeout, err: %w", err)
	}

	return result, nil
}


func InnerAuthorize(ctx workflow.Context, accountID, amount int) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, defaultActivityOptions)

	var result string
	err := workflow.ExecuteActivity(ctx, a.Authorize, accountID, amount).Get(ctx, &result)
	if err != nil {
		return "", fmt.Errorf("execute authorize activity, err: %w", err)
	}

	var approved struct{}
	approveChan := workflow.GetSignalChannel(ctx, fmt.Sprintf("approve %d %d", accountID, amount))
	approveChan.Receive(ctx, &approved)

	return result, nil
} */
