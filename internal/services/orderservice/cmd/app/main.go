package main

import (
	"os"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/shared/app"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
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
	pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Orders", pterm.FgLightGreen.ToStyle()),
		putils.LettersFromStringWithStyle(" Service", pterm.FgLightMagenta.ToStyle())).
		Render()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
