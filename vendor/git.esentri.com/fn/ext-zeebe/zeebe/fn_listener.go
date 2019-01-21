package zeebe

import (
	"context"
	"github.com/fnproject/fn/api/models"
	"log"
)

// Function listener for the Zeebe extension implementing the next.FnListener interface
// Listens to the function create, update and delete events and delegates them to the Zeebe adapter
type FnListener struct {
	jobWorkerRegistry *JobWorkerRegistry
}

func (fnListener *FnListener) BeforeFnCreate(ctx context.Context, fn *models.Fn) error {
	return nil
}

func (fnListener *FnListener) AfterFnCreate(ctx context.Context, fn *models.Fn) error {
	fnListener.registerFunctionAsWorkerIfConfigured(fn)
	return nil
}

func (fnListener *FnListener) BeforeFnUpdate(ctx context.Context, fn *models.Fn) error {
	return nil
}

func (fnListener *FnListener) AfterFnUpdate(ctx context.Context, fn *models.Fn) error {
	fnListener.jobWorkerRegistry.UnregisterFunctionAsWorker(fn.ID)
	fnListener.registerFunctionAsWorkerIfConfigured(fn)
	return nil
}

func (fnListener *FnListener) BeforeFnDelete(ctx context.Context, fnID string) error {
	fnListener.jobWorkerRegistry.UnregisterFunctionAsWorker(fnID)
	return nil
}

func (fnListener *FnListener) AfterFnDelete(ctx context.Context, fnID string) error {
	return nil
}

func (fnListener *FnListener) registerFunctionAsWorkerIfConfigured(fn *models.Fn) {
	jobType, ok := fn.Config["zeebe_job_type"]
	if ok {
		fnListener.jobWorkerRegistry.RegisterFunctionAsWorker(fn.ID, jobType)
	} else {
		log.Println("No Zeebe job type configuration found. Function ID: ", fn.ID)
	}
}
