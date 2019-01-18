package zeebe

import (
    "fmt"
    "context"
    "github.com/fnproject/fn/api/models"
    "github.com/fnproject/fn/api/server"
	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
    "github.com/zeebe-io/zeebe/clients/go/zbc"
    "net/http"
    "log" // TODO log as fn logs
)

const brokerAddr = "host.docker.internal:26500" // TODO read from a config / context

type ZeebeAdapter struct {
	instance worker.JobWorker
}

// Starts a great hard coded worker with a hard coded job type and a hard coded broker address
func (zeebeAdapter *ZeebeAdapter) RegisterFunctionAsWorker(ctx context.Context, fn *models.Fn, server *server.Server) {
	client, err := zbc.NewZBClient(brokerAddr)
	if err != nil {
		panic(err)
    }

	jobType := "payment-service"
    log.Println("Creating a Zeebe job worker for type", jobType)
    jobHandler := contextAwareJobHandler(ctx, fn, server)
	zeebeAdapter.instance = client.NewJobWorker().JobType(jobType).Handler(jobHandler).Open()
}

func (zeebeAdapter *ZeebeAdapter) UnregisterFunctionAsWorker(ctx context.Context, fnID string) {
	if zeebeAdapter.instance != nil {
		log.Println("Stopping worker")
		zeebeAdapter.instance.Close()
		zeebeAdapter.instance = nil
	} else {
		log.Println("Nothing to stop")
	}
}

// TODO extract to a handler file

// Closure over the context and the function
func contextAwareJobHandler(ctx context.Context, fn *models.Fn, server *server.Server) worker.JobHandler {

    fmt.Printf("contextAwareJobHandler running\n")

    return func(client worker.JobClient, job entities.Job) {
        
        jobKey := job.GetKey()
        fmt.Println("Handling job", jobKey)
        fmt.Printf("ctx: %v / fn: %v\n", ctx, fn)
    
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
    
        payload["totalPrice"] = 46.50

        // We probably need some input/output mapping for the gRPC call
        fmt.Println("Invoking function", fn.ID)
        resp, err := http.Post("http://localhost:8080/invoke/" + fn.ID, "application/json", nil) // TODO get the host and the port from somewhere
        if err != nil {
            // failed to post
            fmt.Printf("Failed to send the post request for job %v / error: %v\n", jobKey, err)
            failJob(client, job)
            return
        }
        fmt.Printf("Function invocation successful. Response: %v\n", resp)

        request, err := client.NewCompleteJobCommand().JobKey(jobKey).PayloadFromMap(payload) // TODO use the response, Luke
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
}

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())
	client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send()
}
