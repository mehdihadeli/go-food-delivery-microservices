package problemdetail

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/problemDetails"

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
