package zeebe

import (
	"context"
	"fmt"

	"github.com/fnproject/fn/api/models"
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
)

func init() {
	server.RegisterExtension(&Zeebe{})
}

type Zeebe struct {
}

func (e *Zeebe) Name() string {
	return "git.esentri.com/fn/ext-zeebe/zeebe" // "github.com/treeder/fn-ext-example2/logspam2"
}

func (e *Zeebe) Setup(s fnext.ExtServer) error {
	fmt.Println("Zeebe integration setup!")
	s.AddCallListener(&Zeebe{})
	return nil
}
func (l *Zeebe) BeforeCall(ctx context.Context, call *models.Call) error {
	fmt.Println("Zeebe integration before call!")
	return nil
}

func (l *Zeebe) AfterCall(ctx context.Context, call *models.Call) error {
	fmt.Println("Zeebe integration after call!")
	return nil
}