package httperr

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

const defaultStatusCode = http.StatusInternalServerError

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

type (
	Map      func(err error) (int, bool)
	Mappings []Map
	Option   func(h *Handler)
)

func (h Mappings) find(err error) (int, bool) {
	for _, v := range h {
		if status, ok := v(err); ok {
			return status, true
		}
	}
	return 0, false
}

var DefaultHandler = &Handler{}

type Handler struct {
	httpErrMappings Mappings
	handle          func(err error, c echo.Context)
}

func NewHandler(opts ...Option) *Handler {
	var errHandler Handler
	for _, v := range opts {
		v(&errHandler)
	}
	errHandler.setDefaultProblemDetailsHandle()
	return &errHandler
}

func (h *Handler) WithMap(statusCode int, errs ...error) Option {
	return func(h *Handler) {
		h.httpErrMappings = append(h.httpErrMappings, func(err error) (int, bool) {
			for _, v := range errs {
				if errors.Is(err, v) {
					return statusCode, true
				}
			}
			return statusCode, false
		})
	}
}

func (h *Handler) WithMapFunc(m Map) Option {
	return func(h *Handler) {
		h.httpErrMappings = append(h.httpErrMappings, m)
	}
}

func (h *Handler) Handle() func(err error, c echo.Context) {
	if h.handle != nil {
		return h.handle
	}
	h.setDefaultProblemDetailsHandle()
	return h.handle
}

func Handle() func(err error, c echo.Context) {
	return DefaultHandler.Handle()
}

func (h *Handler) setDefaultProblemDetailsHandle() {
	problemDetailsHandle := func(err error, c echo.Context) {
		code := defaultStatusCode

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		mCode, ok := h.httpErrMappings.find(err)
		if !ok {
			handleErr(c, prepareProblemDetails(err, "internal-server-error", code, c))
			return
		}
		handleErr(c, prepareProblemDetails(err, "application-error", mCode, c))
	}
	h.handle = problemDetailsHandle
}

func handleErr(c echo.Context, pDetails ProblemDetails) {
	if c.Response().Committed {
		return
	}
	if c.Request().Method == http.MethodHead {
		if err := c.NoContent(pDetails.Status); err != nil {
			c.Logger().Error(err)
		}
		return
	}
	if err := c.JSON(pDetails.Status, pDetails); err != nil {
		c.Logger().Error(err)
	}
}

func prepareProblemDetails(err error,
	typ string,
	code int,
	c echo.Context) ProblemDetails {
	return ProblemDetails{
		Type:     typ,
		Title:    err.Error(),
		Status:   code,
		Detail:   err.Error(),
		Instance: c.Request().RequestURI,
	}
}
