package zeebe

import (
	"context"
	"fmt"

	"github.com/fnproject/fn/api/models"
)

type FnListener struct {
}

func (a *FnListener) BeforeFnCreate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! BeforeFnCreate")
	return nil
}

func (a *FnListener) AfterFnCreate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! AfterFnCreate")
	return nil
}

func (a *FnListener) BeforeFnUpdate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! BeforeFnUpdate")
	return nil
}

func (a *FnListener) AfterFnUpdate(ctx context.Context, fn *models.Fn) error {
	fmt.Println("ZEEBE! AfterFnUpdate")
	return nil
}

func (a *FnListener) BeforeFnDelete(ctx context.Context, fnID string) error {
	fmt.Println("ZEEBE! BeforeFnDelete")
	return nil
}

func (a *FnListener) AfterFnDelete(ctx context.Context, fnID string) error {
	fmt.Println("ZEEBE! BeforeFnCreate")
	return nil
}
