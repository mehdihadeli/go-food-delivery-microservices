package main

import (
	"github.com/spf13/cobra"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/app"
)

var rootCmd = &cobra.Command{
	Use:              "orders-microservice",
	Short:            "orders-microservice based on vertical slice architecture",
	Long:             `This is a command runner or cli for api architecture in golang.`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		app.NewApp().Run()
	},
}

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Orders Service Api
// @version 1.0
// @description Orders Service Api
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
