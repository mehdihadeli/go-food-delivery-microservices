package customHadnlers

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/problemDetails"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	errorUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/error_utils"

	"github.com/labstack/echo/v4"
)

func ProblemHandlerFunc(err error, c echo.Context, logger logger.Logger) {
	prb := problemDetails.ParseError(err)

	if prb != nil {
		if !c.Response().Committed {
			if _, err := problemDetails.WriteTo(prb, c.Response()); err != nil {
				logger.Error(err)
			}
		}
	} else {
		if !c.Response().Committed {
			prb := problemDetails.NewInternalServerProblemDetail(err.Error(), errorUtils.ErrorsWithStack(err))
			if _, err := problemDetails.WriteTo(prb, c.Response()); err != nil {
				logger.Error(err)
			}
		}
	}
}
