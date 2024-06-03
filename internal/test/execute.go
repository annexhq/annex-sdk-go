package test

import (
	"errors"
	"fmt"

	testv1 "github.com/annexhq/annex-proto/gen/go/type/test/v1"
	"github.com/annexhq/annex/test"
	"go.temporal.io/sdk/workflow"

	"github.com/annexhq/annex-sdk-go/internal/temporal"
)

type TestExecutor func(ctx workflow.Context, payload *testv1.Payload) error

func ExecuteTest(ctx workflow.Context, wf func(ctx workflow.Context)) error {
	weInfo := workflow.GetInfo(ctx)
	testExecID, err := test.ParseTestWorkflowID(weInfo.WorkflowExecution.ID)
	if err != nil {
		return err
	}

	ctx = temporal.WorkflowContextWithTestLogConfig(ctx, temporal.TestLogConfig{
		TestExecID: testExecID,
	})

	return execWithRecover(func() {
		wf(ctx)
	})
}

func execWithRecover(wrapper func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%+v", r))
		}
	}()
	wrapper()
	return err
}
