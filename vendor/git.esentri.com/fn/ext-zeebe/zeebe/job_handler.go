package zeebe

import (
	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
    "net/http"
    "log" // TODO log as fn logs
)

// Closure over the function ID and the needed context
// This is needed as worker.JobHandler does not have access to the context such as the function id
func JobHandler(fnID string, loadBalancerHost string) worker.JobHandler {

    return func(client worker.JobClient, job entities.Job) {
        
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
    
        log.Println("Invoking function", fnID)
        invocationUrl := loadBalancerHost + "/invoke/" + fnID
        resp, err := http.Post(invocationUrl, "application/json", nil) // TODO get the host and the port from somewhere
        if err != nil {
            // failed to post
            log.Printf("Failed to send the post request for job %v / error: %v\n", jobKey, err)
            failJob(client, job)
            return
        }
        // TODO check the response code
        log.Printf("Function invocation successful. Response: %v\n", resp)

        // Remove the dummy payload. Use the response, Luke.
        payload["totalPrice"] = 46.50

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
}

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())
	client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send()
}
