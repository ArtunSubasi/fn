package zeebe

import (
	"fmt"
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
	fmt.Println("Zeebe integration setup!") // TODO better logging
	server := s.(*server.Server) // TODO this type assertion is hacky. ExtServer should implement the AddFnListener interface.

	// The second parameter 'server' is even more hacky, but we need the server in the function listener in order to invoke functions
	// Otherwise all functions have to be called over http. It this really so bad as it sounds like?
	// TODO change the server parameter to ExtServer since we will invoke the functions over http using the invoke endpoint.
	server.AddFnListener(&FnListener{&ZeebeAdapter{}, server}) 
	return nil
}
