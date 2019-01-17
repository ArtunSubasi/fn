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
	return "git.esentri.com/fn/ext-zeebe/zeebe"
}

func (e *Zeebe) Setup(s fnext.ExtServer) error {
	fmt.Println("Zeebe integration setup!")
	s.AddCallListener(&Zeebe{})
	

	server := s.(*server.Server) // TODO this type assertion is hacky. ExtServer should implement the AddFnListener interface.

	// TODO now we can add a function listener using the server variable, as shown with the AddAppListener:
	server.AddAppListener(&Zeebe{}) // TODO those should be separate listener in their own files

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

// BeforeAppCreate called right before creating App in the database
func (l *Zeebe) BeforeAppCreate(ctx context.Context, app *models.App) error {
	return nil
}

// AfterAppCreate called after creating App in the database
func (l *Zeebe) AfterAppCreate(ctx context.Context, app *models.App) error {
	return nil
}

// BeforeAppUpdate called right before updating App in the database
func (l *Zeebe) BeforeAppUpdate(ctx context.Context, app *models.App) error {
	fmt.Println("BEFORE APP UPDATE!")
	return nil
}

// AfterAppUpdate called after updating App in the database
func (l *Zeebe) AfterAppUpdate(ctx context.Context, app *models.App) error {
	fmt.Println("AFTER APP UPDATE!")
	return nil
}

// BeforeAppDelete called right before deleting App in the database
func (l *Zeebe) BeforeAppDelete(ctx context.Context, app *models.App) error {
	return nil
}

// AfterAppDelete called after deleting App in the database
func (l *Zeebe) AfterAppDelete(ctx context.Context, app *models.App) error {
	return nil
}

// BeforeAppGet called right before getting an app
func (l *Zeebe) BeforeAppGet(ctx context.Context, appID string) error {
	return nil
}

// AfterAppGet called after getting app from database
func (l *Zeebe) AfterAppGet(ctx context.Context, app *models.App) error {
	return nil
}

// BeforeAppsList called right before getting a list of all user's apps. Modify the filter to adjust what gets returned.
func (l *Zeebe) BeforeAppsList(ctx context.Context, filter *models.AppFilter) error {
	return nil
}

// AfterAppsList called after deleting getting a list of user's apps. apps is the result after applying AppFilter.
func (l *Zeebe) AfterAppsList(ctx context.Context, apps []*models.App) error {
	return nil
}
