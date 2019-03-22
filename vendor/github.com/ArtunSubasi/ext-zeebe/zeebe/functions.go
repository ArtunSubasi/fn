package zeebe

import (
	"encoding/json"
	"github.com/fnproject/fn/api/models"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type FnTriggerWithZeebeJobType struct {
	fnID    		string
	appName			string
	triggerID 		string
	triggerName		string
	triggerSource	string
	jobType 		string
}

func GetZeebeJobType(fn *models.Fn) (string, bool) {
	zeebeJobType, ok := fn.Config["zeebe_job_type"]
	return zeebeJobType, ok
}

// Gets all functions which are deployed and have a configured Zeebe job type
func GetFunctionsWithZeebeJobType(apiServerHost string) []*FnTriggerWithZeebeJobType {
	functionsWithZeebeJobType := make([]*FnTriggerWithZeebeJobType, 0)
	appList := getApps(apiServerHost)
	for _, app := range appList.Items {
		logrus.Debugf("App-ID %v / App-Name: %v\n", app.ID, app.Name)
		fnList := getFunctions(apiServerHost, app.ID)
		// TODO too many loops, too many lines -> refactor, extract the functions
		for _, fn := range fnList.Items {
			logrus.Debugf("Fn-ID %v / Fn-Name: %v\n", fn.ID, fn.Name)
			jobType, hasJobType := GetZeebeJobType(fn)
			if hasJobType {
				trigger, hasTrigger := GetTrigger(apiServerHost, fn)
				if hasTrigger {
					functionsWithZeebeJobType = append(functionsWithZeebeJobType, 
						&FnTriggerWithZeebeJobType{fn.ID, app.Name, trigger.ID, trigger.Name, trigger.Source, jobType})
				} else {
					logrus.Infof("The function %v defines a Zeebe job type but does not have a trigger. Function ID: %v\n", fn.Name, fn.ID)
				}
			} else {
				logrus.Infoln("No Zeebe job type configuration found. Function ID: ", fn.ID)
			}
		}
	}

	for _, fn := range functionsWithZeebeJobType {
		logrus.Infof("Registered triggers: %v\n", fn)
	}

	return functionsWithZeebeJobType
}

func getApps(apiServerHost string) *models.AppList {
	appListUrl := apiServerHost + "/v2/apps"
	logrus.Debugln("Getting apps using the url: ", appListUrl)
	resp, err := http.Get(appListUrl)

	// TODO Better error handling
	if err != nil || resp.StatusCode != http.StatusOK {
		logrus.Errorln("Failed to get the list of apps")
		return &models.AppList{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln("Failed to get the list of apps / can't read the body")
		return &models.AppList{}
	}
	resp.Body.Close()

	appList := models.AppList{}
	err = json.Unmarshal(body, &appList)
	if err != nil {
		logrus.Errorln("Cannot unmarshall body into json")
		return &models.AppList{}
	}

	return &appList
}

func GetApp(apiServerHost string, appID string) *models.App {
	appUrl := apiServerHost + "/v2/apps/" + appID
	logrus.Debugln("Getting app  using the url: ", appUrl)
	resp, err := http.Get(appUrl)

	// TODO Better error handling
	if err != nil || resp.StatusCode != http.StatusOK {
		logrus.Errorln("Failed to get the app")
		return &models.App{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln("Failed to get the app / can't read the body")
		return &models.App{}
	}
	resp.Body.Close()

	app := models.App{}
	err = json.Unmarshal(body, &app)
	if err != nil {
		logrus.Errorln("Cannot unmarshall body into json")
		return &models.App{}
	}

	return &app
}

func getFunctions(apiServerHost string, appID string) *models.FnList {
	fnListUrl := apiServerHost + "/v2/fns?app_id=" + appID
	logrus.Debugln("Getting fns using the url: ", fnListUrl)
	resp, err := http.Get(fnListUrl)

	// TODO Better error handling
	if err != nil || resp.StatusCode != http.StatusOK {
		logrus.Errorln("Failed to get the list of functions")
		return &models.FnList{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln("Failed to get the list of functions / can't read the body")
		return &models.FnList{}
	}
	resp.Body.Close()

	fnList := models.FnList{}
	err = json.Unmarshal(body, &fnList)
	if err != nil {
		logrus.Errorln("Cannot unmarshall body into json")
		return &models.FnList{}
	}

	return &fnList
}

func GetTriggers(apiServerHost string, fn *models.Fn) *models.TriggerList {
	triggerSearchUrl := apiServerHost + "/v2/triggers?app_id=" + fn.AppID + "&fn_id=" + fn.ID
	logrus.Debugln("Searching for triggers using the url: ", triggerSearchUrl)
	resp, err := http.Get(triggerSearchUrl)

	// TODO Better error handling
	if err != nil || resp.StatusCode != http.StatusOK {
		logrus.Errorln("Failed to get the list of triggers")
		return &models.TriggerList{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln("Failed to get the list of functions / can't read the body")
		return &models.TriggerList{}
	}
	resp.Body.Close()

	triggerList := models.TriggerList{}
	err = json.Unmarshal(body, &triggerList)
	if err != nil {
		logrus.Errorln("Cannot unmarshall body into json")
		return &models.TriggerList{}
	}

	return &triggerList
}

func GetTrigger(apiServerHost string, fn *models.Fn) (*models.Trigger, bool) {
	triggers := GetTriggers(apiServerHost, fn)
	if len(triggers.Items) == 1 {
		return triggers.Items[0], true
	} else {
		return nil, false
	}
}