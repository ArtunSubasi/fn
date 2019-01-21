package zeebe

import (
	"github.com/fnproject/fn/api/models"
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
)

type FnWithZeebeJobType struct {
	fnID string
	jobType string
}

// Gets all functions which are deployed and have a configured Zeebe job type
func ListFunctions(loadBalancerHost string) []FnWithZeebeJobType {
	appList := GetApps(loadBalancerHost)
	for _, app := range appList.Items {
		log.Printf("App-ID %v / App-Name: %v\n", app.ID, app.Name)
	}
	// TODO now iterate over all apps and fetch/filter he functions (and remove the dummy array)
	return make([]FnWithZeebeJobType, 0)
}

func GetApps(loadBalancerHost string) *models.AppList {
	appListUrl := loadBalancerHost + "/v2/apps"
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
