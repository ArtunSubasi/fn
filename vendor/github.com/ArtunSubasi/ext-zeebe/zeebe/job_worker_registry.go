package zeebe

import (
	"github.com/zeebe-io/zeebe/clients/go/worker"
	"github.com/zeebe-io/zeebe/clients/go/zbc"
	"github.com/sirupsen/logrus"
	"time"
)

type JobWorkerRegistry struct {
	jobWorkers map[string]worker.JobWorker // FnID -> JobWorker
	loadBalancerAddr string
	zeebeGatewayAddr string
	usePlaintextConnection bool
}

func NewJobWorkerRegistry(loadBalancerAddr string, zeebeGatewayAddr string, usePlaintextConnection bool) JobWorkerRegistry {
	jobWorkerRegistry := JobWorkerRegistry{}
	jobWorkerRegistry.jobWorkers = make(map[string]worker.JobWorker)
	jobWorkerRegistry.loadBalancerAddr = loadBalancerAddr
	jobWorkerRegistry.zeebeGatewayAddr = zeebeGatewayAddr
	jobWorkerRegistry.usePlaintextConnection = usePlaintextConnection
	return jobWorkerRegistry
}

func (jobWorkerRegistry *JobWorkerRegistry) RegisterFunctionAsWorker(fnZeebe *FnTriggerWithZeebeJobType) {
	client, err := zbc.NewZBClientWithConfig(&zbc.ZBClientConfig{
		GatewayAddress: jobWorkerRegistry.zeebeGatewayAddr,
		UsePlaintextConnection: jobWorkerRegistry.usePlaintextConnection})
		
	if err != nil {
		panic(err)
	}
	logrus.Infof("Creating a Zeebe job worker of type %v for function ID %v", fnZeebe.jobType, fnZeebe.fnID)
	jobHandler := JobHandler(fnZeebe, jobWorkerRegistry.loadBalancerAddr)
	// TODO Add more configuration possibilities for the worker (Poll Interval, Timeout, ...)
	jobWorkerRegistry.jobWorkers[fnZeebe.fnID] = client.NewJobWorker().JobType(fnZeebe.jobType).Handler(jobHandler).PollInterval(1 * time.Second).Open()

	// If the zeebe gateway is not available, the Zeebe client keeps logging errors after every poll.
	// There should be a way of controlling the logs. As an alternative, there may be a different interval for checking connections.
	// E.g. check for new jobs every 500 ms, but if the connection is down, check every 10 seconds.
	// TODO Contact Camunda Team with a feature request (Circuit breaker).
}

func (jobWorkerRegistry *JobWorkerRegistry) UnregisterFunctionAsWorker(fnID string) {
	jobWorker, exists := jobWorkerRegistry.jobWorkers[fnID]
	if exists {
		logrus.Infoln("Stopping zeebe job worker for function ID: ", fnID)
		jobWorker.Close()
		delete(jobWorkerRegistry.jobWorkers, fnID)
	} else {
		logrus.Infoln("No zeebe job worker for function ID: ", fnID)
	}
}
