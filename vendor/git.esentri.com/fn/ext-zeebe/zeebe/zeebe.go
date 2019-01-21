package zeebe

import (
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
	"log" // TODO log as fn logs
	"time"
)

// Extension for Zeebe integration
func init() {
	// TODO only register the extension for the modes "Full" and "API Server"
	server.RegisterExtension(&Zeebe{})
}

type Zeebe struct {
}

func (e *Zeebe) Name() string {
	return "git.esentri.com/fn/ext-zeebe/zeebe"
}

func (e *Zeebe) Setup(s fnext.ExtServer) error {
	log.Println("Zeebe integration setup!")
	server := s.(*server.Server) // TODO this type assertion is hacky. ExtServer should implement the AddFnListener interface.
	server.AddFnListener(&FnListener{&ZeebeAdapter{}})

	// Waiting for the REST endpoints to come up since the Extension Setup does not have any callback such as OnServerStarted
	// TODO Get in touch with the Fn Project: Create a feature request, maybe a pull request
	go waitAndListFunctions()

	return nil
}

func waitAndListFunctions() {
	time.Sleep(5 * time.Second)
	functions := GetFunctionsWithZeebeJobType(loadBalancerAddr)
	log.Println("Functions: ", functions)
}
