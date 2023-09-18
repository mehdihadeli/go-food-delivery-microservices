package main

import (
	"os"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/app"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "catalogs-read-microservices",
	Short:            "catalogs-read-microservices based on vertical slice architecture",
	Long:             `This is a command runner or cli for api architecture in golang.`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		app.NewApp().Run()
	},
}

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Catalogs Read-Service Api
// @version 1.0
// @description Catalogs Read-Service Api.
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
