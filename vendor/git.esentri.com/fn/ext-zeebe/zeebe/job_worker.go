package zeebe

import (
    "github.com/zeebe-io/zeebe/clients/go/entities"
    "github.com/zeebe-io/zeebe/clients/go/worker"
    "github.com/zeebe-io/zeebe/clients/go/zbc"
    "log"
)

const brokerAddr = "host.docker.internal:26500" // TODO read from a config / context

type JobWorker struct {
	instance worker.JobWorker
}

// Starts a great hard coded worker with a hard coded job type and a hard coded broker address
func (j JobWorker) Work() {
    client, err := zbc.NewZBClient(brokerAddr)
    if err != nil {
        panic(err)
    }

	jobType := "payment-service"
	log.Println("Creating a Zeebe job worker for type", jobType)
    j.instance = client.NewJobWorker().JobType(jobType).Handler(handleJob).Open()
    defer j.instance.Close()
    j.instance.AwaitClose()
}

func (j JobWorker) StopWorking() {
	if j.instance != nil {
		log.Println("Stopping worker...")
		j.instance.Close()
	} else {
		log.Println("Nothing to stop")
	}
}

func handleJob(client worker.JobClient, job entities.Job) {

    jobKey := job.GetKey()

    headers, err := job.GetCustomHeadersAsMap()
    if err != nil {
        // failed to handle job as we require the custom job headers
        failJob(client, job)
        return
    }

    payload, err := job.GetPayloadAsMap()
    if err != nil {
        // failed to handle job as we require the payload
        failJob(client, job)
        return
    }

	payload["totalPrice"] = 46.50;
	
	// Somewhere around here, the Fn Project function must be be called before completing the job
	// We probably need some input/output mapping for the gRPC call

    request, err := client.NewCompleteJobCommand().JobKey(jobKey).PayloadFromMap(payload)
    if err != nil {
        // failed to set the updated payload
        failJob(client, job)
        return
    }

    log.Println("Complete job", jobKey, "of type", job.Type)
    log.Println("Processing order:", payload["orderId"])
    log.Println("Collect money using payment method:", headers["method"])

    request.Send()
}

func failJob(client worker.JobClient, job entities.Job) {
    log.Println("Failed to complete job", job.GetKey())
    client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send()
}
