package zeebe

import (
	"github.com/zeebe-io/zeebe/clients/go/entities"
    "github.com/zeebe-io/zeebe/clients/go/worker"
    "github.com/sirupsen/logrus"
    "net/http"
    "io/ioutil"
    "bytes"
    "encoding/json"
)

// Closure over the function ID and the needed context
// This is needed as worker.JobHandler of the Zeebe package does not have access to the context such as the function id
func JobHandler(fnID string, loadBalancerHost string) worker.JobHandler {

    return func(client worker.JobClient, job entities.Job) {
        
        jobKey := job.GetKey()
    
        // TODO refactor: extract function invocation as a separate function
        logrus.Infoln("Invoking function", fnID)
        invocationUrl := loadBalancerHost + "/invoke/" + fnID
        logrus.Debugln("InvocationUrl:", invocationUrl)
        logrus.Debugln("Payload:", job.Payload)
        resp, err := http.Post(invocationUrl, "application/json", bytes.NewBuffer([]byte(job.Payload)))
        if err != nil {
            logrus.Errorf("Failed to send the post request for job key %v / error: %v\n", jobKey, err)
            failJob(client, job)
            return
        }

        if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
            logrus.Infof("Function invocation returned the HTTP status code %v. Job key: %v / fnID: %v\n", resp.Status, jobKey, fnID)
        } else {
            logrus.Infof("Function invocation returned the HTTP status code %v. Failing job. Job key: %v / fnID: %v\n", resp.Status, jobKey,fnID)
            failJob(client, job)
            return
        }

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            logrus.Errorf("Failed to read the body after invoking function. Failing job. Job key: %v / fnID: %v\n", jobKey, fnID)
            failJob(client, job)
            return
        }

        var responseJsonObject map[string]interface{}
        err = json.Unmarshal(body, &responseJsonObject)
        if err != nil {
            logrus.Warnf("Failed to unmarshall the response. Zeebe supports only JSON objects on root level. Response will be ignored. Job key: %v / fnID: %v\n", jobKey, fnID)
            logrus.Debugln("Response:", string(body))
            responseJsonObject = nil
        } else {
            logrus.Debugln("Response:", responseJsonObject)
        }

        request, err := client.NewCompleteJobCommand().JobKey(jobKey).PayloadFromObject(responseJsonObject) 
        if err != nil {
            logrus.Errorf("Failed to set the updated payload. Failing job. Job key: %v / fnID: %v\n", jobKey, fnID)
            failJob(client, job)
            return
        }
    
        logrus.Println("Completed job", jobKey, "of type", job.Type)
        request.Send()
    }
}

func failJob(client worker.JobClient, job entities.Job) {
	logrus.Println("Failed to complete job", job.GetKey())
	client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send()
}
