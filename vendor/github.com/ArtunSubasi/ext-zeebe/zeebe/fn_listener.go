package zeebe

import (
	"context"
	"github.com/fnproject/fn/api/models"
	"github.com/sirupsen/logrus"
	"time"
)

// Function listener for the Zeebe extension implementing the next.FnListener interface
// Listens to the function create, update and delete events and delegates them to the Zeebe adapter
type FnListener struct {
	jobWorkerRegistry 	*JobWorkerRegistry
	apiServerAddr 		string
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
	jobType, ok := GetZeebeJobType(fn)
	if ok {
		go fnListener.waitAndRegisterTriggers(fn, jobType)
	} else {
		logrus.Infoln("No Zeebe job type configuration found. Function ID: ", fn.ID)
	}
}

func (fnListener *FnListener) waitAndRegisterTriggers(fn *models.Fn, jobType string) {
	// Sleeping and waiting for the trigger to be registered before searching for it 
	// since the Extension Setup does not have any callback such as OnTriggerCreated
	// TODO Get in touch with the Fn Project: Create a feature request, maybe a pull request
	time.Sleep(2 * time.Second)
	trigger, hasTrigger := GetTrigger(fnListener.apiServerAddr, fn)
	if hasTrigger {
		app := GetApp(fnListener.apiServerAddr, fn.AppID)
		fnZeebe := &FnTriggerWithZeebeJobType{fn.ID, app.Name, trigger.ID, trigger.Name, trigger.Source, jobType}
		fnListener.jobWorkerRegistry.RegisterFunctionAsWorker(fnZeebe)
	} else {
		logrus.Infof("The function %v defines a Zeebe job type but does not have a trigger. Function ID: %v", fn.Name, fn.ID)
	}
}
