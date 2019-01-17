package zeebe

import (
	"context"
	"fmt"

	"github.com/fnproject/fn/api/models"
)

// Function listener for the Zeebe extension implementing the next.FnListener interface
// Listens to the function create, update and delete events and delegates them to the Zeebe adapter
type FnListener struct {
	zeebeAdapter *ZeebeAdapter
}

func NewFnListener() FnListener {
	f := FnListener{}
	f.zeebeAdapter = &ZeebeAdapter{}
	return f
}

func (a *FnListener) BeforeFnCreate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! BeforeFnCreate")
	return nil
}

func (a *FnListener) AfterFnCreate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! AfterFnCreate")
	a.zeebeAdapter.RegisterFunctionAsWorker()
	return nil
}

func (a *FnListener) BeforeFnUpdate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! BeforeFnUpdate")
	return nil
}

func (a *FnListener) AfterFnUpdate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! AfterFnUpdate")
	a.zeebeAdapter.UnregisterFunctionAsWorker()
	a.zeebeAdapter.RegisterFunctionAsWorker()
	return nil
}

func (a *FnListener) BeforeFnDelete(ctx context.Context, fnID string) error {
	fmt.Println("ZEEBE! BeforeFnDelete")
	a.zeebeAdapter.UnregisterFunctionAsWorker()
	return nil
}

func (a *FnListener) AfterFnDelete(ctx context.Context, fnID string) error {
	fmt.Println("ZEEBE! BeforeFnCreate")
	return nil
}
