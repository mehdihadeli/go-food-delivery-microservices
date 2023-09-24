package otelmetrics

import "github.com/labstack/echo/v4/middleware"

type config struct {
	Skipper middleware.Skipper
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}
