package zerodowntime

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

func HandleApp(app *fx.App) error {
	startCtx, cancel := context.WithTimeout(context.Background(), app.StartTimeout())
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		return err
	}

	sigs := app.Done()

	fmt.Println(fmt.Sprintf("received signal: %v", <-sigs))
	fmt.Println(fmt.Sprintf("exiting in %s", app.StartTimeout().String()))

	stopCtx, cancel := context.WithTimeout(context.Background(), app.StopTimeout())
	defer cancel()

	fmt.Println("stopping app")
	return app.Stop(stopCtx)
}
