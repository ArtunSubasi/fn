package zeebe

import (
	"log" // TODO log as fn logs
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
)

// Extension for Zeebe integration
func init() {
	// TODO only register the extension for the modes "Full" and "Load Balancer"
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

	functions := ListFunctions(loadBalancerAddr)
	log.Println("Functions: ", functions)

	return nil
}
