package zeebe

import (
	"github.com/zeebe-io/zeebe/clients/go/worker"
	"github.com/zeebe-io/zeebe/clients/go/zbc"
	"time"
	"log" // TODO log as fn logs
)

const brokerAddr = "host.docker.internal:26500"  // TODO read from a config / context (app config?)
const loadBalancerAddr = "http://localhost:8080" // TODO read from a config / context (app config?)

type JobWorkerRegistry struct {
	// FnID is not unique without the app ID. The AppID should be used as well.
	// But in this case, workers cannot be unregistered because the FnListener.AfterFnDelete only has the FnID as parameter
	// TODO Contact the Fn Project. FnListener.AfterFnDelete should have more context, at least the AppID of the App in which the function lived.
	jobWorkers map[string]worker.JobWorker // FnID -> JobWorker
}

func NewJobWorkerRegistry() JobWorkerRegistry {
	jobWorkerRegistry := JobWorkerRegistry{}
	jobWorkerRegistry.jobWorkers = make(map[string]worker.JobWorker)
	return jobWorkerRegistry
}

// Starts a great hard coded worker with a hard coded job type and a hard coded broker address
func (jobWorkerRegistry *JobWorkerRegistry) RegisterFunctionAsWorker(fnID string, zeebeJobType string) {
	client, err := zbc.NewZBClient(brokerAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("Creating a Zeebe job worker of type %v for function ID %v\n", zeebeJobType, fnID)
	jobHandler := JobHandler(fnID, loadBalancerAddr)
	// TODO Add more configuration possibilities for the worker (Poll Interval, Timeout, ...)
	jobWorkerRegistry.jobWorkers[fnID] = client.NewJobWorker().JobType(zeebeJobType).Handler(jobHandler).PollInterval(1 * time.Second).Open()

	// If the zeebe gateway is not available, the Zeebe client keeps logging errors after every poll.
	// There should be a way of controlling the logs. As an alternative, there may be a different interval for checking connections.
	// E.g. check for new jobs every 500 ms, but if the connection is down, check every 10 seconds.
	// TODO Contact Camunda Team with a feature request.
}

func (jobWorkerRegistry *JobWorkerRegistry) UnregisterFunctionAsWorker(fnID string) {
	jobWorker, exists := jobWorkerRegistry.jobWorkers[fnID]
	if exists {
		log.Println("Stopping zeebe job worker for function ID: ", fnID)
		jobWorker.Close()
		delete(jobWorkerRegistry.jobWorkers, fnID)
	} else {
		log.Println("No zeebe job worker for function ID: ", fnID)
	}
}
