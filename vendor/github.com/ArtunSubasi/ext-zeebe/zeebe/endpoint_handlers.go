package zeebe

import (
	"fmt"
	"net/http"
)

// Handler for the endpoint "zeebe" which lists all registered zeebe job workers as plain text
// TODO return the job workers in a structured form (JSON etc.) so that it can be processed in a UI
type zeebeEndpointHandler struct{
	apiServerAddr string
}

func (h *zeebeEndpointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Zeebe Endpoint Handler called")
	functionsWithZeebeJobType := GetFunctionsWithZeebeJobType(h.apiServerAddr)
	fmt.Fprintf(w, "Registered zeebe job workers:\n")
	if len(functionsWithZeebeJobType) == 0 {
		fmt.Fprintf(w, "N/A\n")
	} else {
		for _, fn := range functionsWithZeebeJobType {
			fmt.Fprintf(w, "Fn Function-ID %q -> Zeebe Job Type: %q\n", fn.fnID, fn.jobType)
		}
	}
	
}
