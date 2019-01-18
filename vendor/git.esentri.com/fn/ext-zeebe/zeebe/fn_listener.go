package zeebe

import (
	"context"
	"fmt"

	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/api/models"
)

// Function listener for the Zeebe extension implementing the next.FnListener interface
// Listens to the function create, update and delete events and delegates them to the Zeebe adapter
type FnListener struct {
	zeebeAdapter *ZeebeAdapter
	server *server.Server
}

func (fnListener *FnListener) BeforeFnCreate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! BeforeFnCreate")
	return nil
}

func (fnListener *FnListener) AfterFnCreate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! AfterFnCreate")
	fnListener.zeebeAdapter.RegisterFunctionAsWorker(ctx, fn, fnListener.server)
	return nil
}

func (fnListener *FnListener) BeforeFnUpdate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! BeforeFnUpdate")
	return nil
}

func (fnListener *FnListener) AfterFnUpdate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! AfterFnUpdate")
	fnListener.zeebeAdapter.UnregisterFunctionAsWorker(ctx, fn.ID)
	fnListener.zeebeAdapter.RegisterFunctionAsWorker(ctx, fn, fnListener.server)
	return nil
}

func (fnListener *FnListener) BeforeFnDelete(ctx context.Context, fnID string) error {
	fmt.Println("ZEEBE! BeforeFnDelete")
	fnListener.zeebeAdapter.UnregisterFunctionAsWorker(ctx, fnID)
	return nil
}

func (fnListener *FnListener) AfterFnDelete(ctx context.Context, fnID string) error {
	fmt.Println("ZEEBE! BeforeFnCreate")
	return nil
}
