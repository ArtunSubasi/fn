package zeebe

import (
	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
    "net/http"
    "log" // TODO log as fn logs
    "io/ioutil"
    "bytes"
    "encoding/json"
)

// Closure over the function ID and the needed context
// This is needed as worker.JobHandler does not have access to the context such as the function id
func JobHandler(fnID string, loadBalancerHost string) worker.JobHandler {

    return func(client worker.JobClient, job entities.Job) {
        
        jobKey := job.GetKey()
    
        log.Println("Invoking function", fnID)
        invocationUrl := loadBalancerHost + "/invoke/" + fnID
        log.Println("InvocationUrl:", invocationUrl)
        resp, err := http.Post(invocationUrl, "application/json", bytes.NewBuffer([]byte(job.Payload)))
        if err != nil {
            // failed to post
            log.Printf("Failed to send the post request for job %v / error: %v\n", jobKey, err)
            failJob(client, job)
            return
        }

        log.Println("Function invocation successful.")

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Println("Failed to read the body")
            return
        }

        var responseMap map[string]interface{}
        err = json.Unmarshal(body, &responseMap)
        if err != nil {
            log.Println("Failed to unmarshall the response. Response will be ignored.")
            responseMap = make(map[string]interface{})
        }

        request, err := client.NewCompleteJobCommand().JobKey(jobKey).PayloadFromMap(responseMap) 
        if err != nil {
            // failed to set the updated payload
            failJob(client, job)
            return
        }
    
        log.Println("Completed job", jobKey, "of type", job.Type)
    
        request.Send()
    }
}

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())
	client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send()
}
