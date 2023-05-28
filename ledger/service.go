package ledger

import (
	"context"
	"fmt"
	"os"

	"encore.app/activity"
	"encore.app/storage"
	"encore.app/utils"
	"encore.app/workflow"
	tb "github.com/tigerbeetledb/tigerbeetle-go"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const ledgerTaskQueue = "ledgerTaskQueue"

//encore:service
type Service struct {
	client client.Client
	worker worker.Worker
	db     storage.Storage
}

func newService() (*Service, error) {
	svc := new(Service)

	err := svc.initAndStart()
	if err != nil {
		svc.Shutdown(context.TODO())

		return nil, fmt.Errorf("init and start service, err: %v", err)
	}

	return svc, nil
}

func (s *Service) initAndStart() error {
	tbAddress := os.Getenv("TB_ADDRESS")
	if len(tbAddress) == 0 {
		tbAddress = "3000"
	}

	var err error
	db, err := tb.NewClient(0, []string{tbAddress}, 1)
	if err != nil {
		return fmt.Errorf("create tiger beetle client, err: %v", err)
	}

	s.db, err = storage.NewTBStorage(db)
	if err != nil {
		return fmt.Errorf("create storage, err: %v", err)
	}

	s.client, err = client.Dial(client.Options{})
	if err != nil {
		return fmt.Errorf("create temporal client: %v", err)
	}

	s.worker = worker.New(s.client, ledgerTaskQueue, worker.Options{})

	s.worker.RegisterWorkflow(workflow.Authorize)
	s.worker.RegisterWorkflow(workflow.WaitForPresent)
	//s.worker.RegisterWorkflow(workflow.TimeoutHandledAuthorize)
	//s.worker.RegisterWorkflow(workflow.InnerAuthorize)
	s.worker.RegisterWorkflow(workflow.Present)
	s.worker.RegisterWorkflow(workflow.Balance)
	s.worker.RegisterActivity(&activity.Authorizator{Storage: s.db})
	s.worker.RegisterActivity(&activity.Presenter{Storage: s.db})
	s.worker.RegisterActivity(&activity.Balancer{Storage: s.db})

	err = s.worker.Start()
	if err != nil {
		return fmt.Errorf("start workflow worker: %v", err)
	}

	return nil
}

func initService() (*Service, error) {
	ledgerSvc, err := newService()
	if err != nil {

		return nil, fmt.Errorf("create ledger service, err: %v", err)
	}

	return ledgerSvc, nil
}

func (s *Service) Shutdown(force context.Context) {
	utils.CloseIfNotNil(s.client)
	utils.StopIfNotNil(s.worker)
	utils.CloseIfNotNil(s.db)
}
