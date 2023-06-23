package main

import (
	"github.com/spf13/cobra"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/app"
)

var rootCmd = &cobra.Command{
	Use:              "ecommerce-microservices",
	Short:            "ecommerce-microservices based on vertical slice architecture",
	Long:             `This is a command runner or cli for api architecture in golang.`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		app.NewApp().Run()
	},
}

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Catalogs Write-Service Api
// @version 1.0
// @description Catalogs Write-Service Api.
func main() {
	_ = rootCmd.Execute()
}
