package main

import (
	"context"

	"github.com/jae2274/auth-service/auth_service/app"
)

func main() {

	app.Run(context.Background())
}
