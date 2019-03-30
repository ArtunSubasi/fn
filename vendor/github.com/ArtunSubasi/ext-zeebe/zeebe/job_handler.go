package zeebe

import (
	"github.com/zeebe-io/zeebe/clients/go/entities"
    "github.com/zeebe-io/zeebe/clients/go/worker"
    "github.com/sirupsen/logrus"
    "net/http"
    "io/ioutil"
    "bytes"
    "encoding/json"
    "strconv"
)

// Closure over the function ID and the needed context
// This is needed as worker.JobHandler of the Zeebe package does not have access to the context such as the function id
func JobHandler(fnZeebe *FnTriggerWithZeebeJobType, loadBalancerHost string) worker.JobHandler {

    return func(client worker.JobClient, job entities.Job) {
        
        jobKey := job.GetKey()

        // TODO refactor: extract function invocation as a separate function, also extract other functions =)
        invocationUrl := loadBalancerHost + "/t/" + fnZeebe.appName + fnZeebe.triggerSource
        logrus.Infoln("Invoking function", fnZeebe.fnID, "/ InvocationUrl:", invocationUrl)
        logrus.Debugln("Payload:", job.Payload)
        resp, err := http.Post(invocationUrl, "application/json", bytes.NewBuffer([]byte(job.Payload)))
        if err != nil {
            logrus.Errorf("Failed to send the post request for job key %v / error: %v", jobKey, err)
            failJob(client, job, "Failed to invoke Fn function " + fnZeebe.fnID)
            return
        }

        logrus.Infof("Function invocation returned the response: %v. Job key: %v / fnID: %v", resp.Status, jobKey, fnZeebe.fnID)

        if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
            logrus.Infof("Function invocation returned the HTTP status code %v. Job key: %v / fnID: %v", resp.Status, jobKey, fnZeebe.fnID)
        } else {
            logrus.Infof("Function invocation returned the HTTP status code %v. Failing job. Job key: %v / fnID: %v", resp.Status, jobKey, fnZeebe.fnID)
            errorPrefix := "HTTP Response code: " + strconv.Itoa(resp.StatusCode)
            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                failJob(client, job, errorPrefix)
            } else {
                failJob(client, job, errorPrefix + " / Msg: " + string(body))
            }
            return
        }

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            logrus.Errorf("Failed to read the body after invoking function. Failing job. Job key: %v / fnID: %v", jobKey, fnZeebe.fnID)
            failJob(client, job, "Failed to read the HTTP response body")
            return
        }

        var responseJsonObject map[string]interface{}
        err = json.Unmarshal(body, &responseJsonObject)
        if err != nil {
            logrus.Warnf("Failed to unmarshall the response. Zeebe supports only JSON objects on root level. Response will be ignored. Job key: %v / fnID: %v\n", jobKey, fnZeebe.fnID)
            logrus.Debugln("Response:", string(body))
            responseJsonObject = nil
        } else {
            logrus.Debugln("Response:", responseJsonObject)
        }

        request, err := client.NewCompleteJobCommand().JobKey(jobKey).PayloadFromObject(responseJsonObject) 
        if err != nil {
            logrus.Errorf("Failed to set the updated payload. Failing job. Job key: %v / fnID: %v", jobKey, fnZeebe.fnID)
            failJob(client, job, "Failed to set to updated payload")
            return
        }
    
        logrus.Println("Completed job", jobKey, "of type", job.Type)
        request.Send()
    }
}

func failJob(client worker.JobClient, job entities.Job, errorMessage string) {
	logrus.Println("Failed to complete job", job.GetKey())
	client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).ErrorMessage(errorMessage).Send()
}
