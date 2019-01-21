package zeebe

import (
	"github.com/zeebe-io/zeebe/clients/go/worker"
	"github.com/zeebe-io/zeebe/clients/go/zbc"
	"log" // TODO log as fn logs
)

const brokerAddr = "host.docker.internal:26500"  // TODO read from a config / context (app config?)
const loadBalancerAddr = "http://localhost:8080" // TODO read from a config / context (app config?)

type JobWorkerRegistry struct {
	// TODO add a map of workers for different job types
	instance worker.JobWorker
}

// Starts a great hard coded worker with a hard coded job type and a hard coded broker address
func (jobWorkerRegistry *JobWorkerRegistry) RegisterFunctionAsWorker(fnID string, zeebeJobType string) {
	client, err := zbc.NewZBClient(brokerAddr)
	if err != nil {
		panic(err)
	}
	log.Println("Creating a Zeebe job worker for type", zeebeJobType)
	jobHandler := JobHandler(fnID, loadBalancerAddr)
	jobWorkerRegistry.instance = client.NewJobWorker().JobType(zeebeJobType).Handler(jobHandler).Open()
}

func (zeebeAdapter *JobWorkerRegistry) UnregisterFunctionAsWorker(fnID string) {
	if zeebeAdapter.instance != nil {
		log.Println("Stopping worker for function ID: ", fnID)
		zeebeAdapter.instance.Close()
		zeebeAdapter.instance = nil
	} else {
		log.Println("Nothing to stop")
	}
}
