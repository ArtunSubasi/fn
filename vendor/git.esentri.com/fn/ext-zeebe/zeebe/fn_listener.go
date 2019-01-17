package zeebe

import (
	"context"
	"fmt"

	"github.com/fnproject/fn/api/models"
)

type FnListener struct {
	jobWorker JobWorker
}

func (a *FnListener) BeforeFnCreate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! BeforeFnCreate")
	fmt.Printf("Function: %v\n", fn)
	return nil
}

func (a *FnListener) AfterFnCreate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! AfterFnCreate")
	fmt.Println("Config:")

    for key, val := range fn.Config {
		s := fmt.Sprintf("%s=\"%s\"", key, val)
        fmt.Println(s)
	}
		
	fmt.Println("Annotations:")
	for key, val := range fn.Annotations {
		s := fmt.Sprintf("%s=\"%s\"", key, val)
        fmt.Println(s)
	}

	// TODO this must be started using goroutines and probably synchronized and stopped using channels
	a.jobWorker.Work()

	return nil
}

func (a *FnListener) BeforeFnUpdate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! BeforeFnUpdate")
	return nil
}

func (a *FnListener) AfterFnUpdate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! AfterFnUpdate")
	a.jobWorker.StopWorking()
	a.jobWorker.Work()
	return nil
}

func (a *FnListener) BeforeFnDelete(ctx context.Context, fnID string) error {
	fmt.Println("ZEEBE! BeforeFnDelete")
	a.jobWorker.StopWorking()
	return nil
}

func (a *FnListener) AfterFnDelete(ctx context.Context, fnID string) error {
	fmt.Println("ZEEBE! BeforeFnCreate")
	return nil
}
