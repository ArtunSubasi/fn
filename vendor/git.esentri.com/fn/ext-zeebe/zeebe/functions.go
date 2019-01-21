package zeebe

import (
	"encoding/json"
	"github.com/fnproject/fn/api/models"
	"io/ioutil"
	"log"
	"net/http"
)

type FnWithZeebeJobType struct {
	fnID    string
	jobType string
}

func GetZeebeJobType(fn *models.Fn) (string, bool) {
	zeebeJobType, ok := fn.Config["zeebe_job_type"]
	return zeebeJobType, ok
}

// Gets all functions which are deployed and have a configured Zeebe job type
func GetFunctionsWithZeebeJobType(apiServerHost string) []*FnWithZeebeJobType {
	functionsWithZeebeJobType := make([]*FnWithZeebeJobType, 0)
	appList := getApps(apiServerHost)
	for _, app := range appList.Items {
		log.Printf("App-ID %v / App-Name: %v\n", app.ID, app.Name)
		fnList := getFunctions(apiServerHost, app.ID)
		for _, fn := range fnList.Items {
			log.Printf("Fn-ID %v / Fn-Name: %v\n", fn.ID, fn.Name)
			jobType, ok := GetZeebeJobType(fn)
			if ok {
				functionsWithZeebeJobType = append(functionsWithZeebeJobType, &FnWithZeebeJobType{fn.ID, jobType})
			} else {
				log.Println("No Zeebe job type configuration found. Function ID: ", fn.ID)
			}
		}
	}

	for _, fn := range functionsWithZeebeJobType {
		log.Printf("Fn-ID %v / Fn-JobType: %v\n", fn.fnID, fn.jobType)
	}

	// TODO now iterate over all apps and fetch/filter he functions (and remove the dummy array)
	return functionsWithZeebeJobType
}

func getApps(apiServerHost string) *models.AppList {
	appListUrl := apiServerHost + "/v2/apps"
	log.Println("Getting apps using the url: ", appListUrl)
	resp, err := http.Get(appListUrl)

	// TODO Better error handling
	if err != nil || resp.StatusCode != http.StatusOK {
		// failed to post
		log.Println("Failed to get the list of apps")
		return &models.AppList{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to get the list of apps / can't read the body")
		return &models.AppList{}
	}
	resp.Body.Close()

	appList := models.AppList{}
	err = json.Unmarshal(body, &appList)
	if err != nil {
		log.Println("Cannot unmarshall body into json")
		return &models.AppList{}
	}

	return &appList
}

func getFunctions(apiServerHost string, appID string) *models.FnList {
	fnListUrl := apiServerHost + "/v2/fns?app_id=" + appID
	log.Println("Getting fns using the url: ", fnListUrl)
	resp, err := http.Get(fnListUrl)

	// TODO Better error handling
	if err != nil || resp.StatusCode != http.StatusOK {
		// failed to post
		log.Println("Failed to get the list of functions")
		return &models.FnList{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to get the list of functions / can't read the body")
		return &models.FnList{}
	}
	resp.Body.Close()

	fnList := models.FnList{}
	err = json.Unmarshal(body, &fnList)
	if err != nil {
		log.Println("Cannot unmarshall body into json")
		return &models.FnList{}
	}

	return &fnList
}
