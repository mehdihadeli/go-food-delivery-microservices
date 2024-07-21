package problemdetail

import (
	problemDetails "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/problemdetails"

	"github.com/labstack/echo/v4/middleware"
)

type config struct {
	Skipper       middleware.Skipper
	ProblemParser problemDetails.ErrorParserFunc
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func WithSkipper(skipper middleware.Skipper) Option {
	return optionFunc(func(cfg *config) {
		cfg.Skipper = skipper
	})
}

func WithErrorParser(errorParser problemDetails.ErrorParserFunc) Option {
	return optionFunc(func(cfg *config) {
		cfg.ProblemParser = errorParser
	})
}
