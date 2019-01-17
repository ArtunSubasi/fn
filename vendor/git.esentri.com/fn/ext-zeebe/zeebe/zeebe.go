package zeebe

import (
	"fmt"
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
)

// Extension for Zeebe integration
func init() {
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

	fnListener := NewFnListener()
	server.AddFnListener(&fnListener)
	return nil
}
