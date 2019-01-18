package zeebe

import (
	"log"
    "net/http"
)

type FnWithZeebeJobType struct {
	fnID string
	jobType string
}

// Gets all functions which are deployed and have a configured Zeebe job type
func ListFunctions(loadBalancerHost string) []FnWithZeebeJobType {
	appListUrl := loadBalancerHost + "/v2/apps"
	log.Println("Getting apps using the url: ", appListUrl)
	resp, err := http.Get(appListUrl)
	if err != nil {
		// failed to post
		log.Println("Failed to get the list of apps")
	}
	log.Printf("Apps-Response: %v\n", resp)

	// TODO now iterate over all apps and fetch/filter he functions (and remove the dummy array)
	return make([]FnWithZeebeJobType, 0)
}
