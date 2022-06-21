package v1

//
//import (
//	"github.com/go-playground/validator"
//	"github.com/labstack/echo/v4"
//	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
//	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
//	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
//	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
//	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
//	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
//	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/dto"
//	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/deleting_product"
//	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/getting_product_by_id"
//	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/updating_product"
//	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared"
//	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/middlewares"
//	"github.com/opentracing/opentracing-go"
//	uuid "github.com/satori/go.uuid"
//	"net/http"
//)
//
//type productsController struct {
//	e       *echo.Echo
//	log     logger.Logger
//	mw      middlewares.MiddlewareManager
//	cfg     *config.Config
//	m       *mediatr.Mediator
//	v       *validator.Validate
//	metrics *shared.CatalogsServiceMetrics
//}
//
//func NewProductsController(
//	e *echo.Echo,
//	log logger.Logger,
//	mw middlewares.MiddlewareManager,
//	cfg *config.Config,
//	m *mediatr.Mediator,
//	v *validator.Validate,
//	metrics *shared.CatalogsServiceMetrics,
//) *productsController {
//	return &productsController{e: e, log: log, mw: mw, cfg: cfg, m: m, v: v, metrics: metrics}
//}
//
//// GetAllProducts
//// @Tags Products
//// @Summary Get all product
//// @Description Get all products
//// @Accept json
//// @Produce json
//// @Success 200
//// @Router /products [get]
//func (h *productsController) GetAllProducts() echo.HandlerFunc {
//	return nil
//}
//
////// CreateProduct
////// @Tags Products
////// @Summary Create product
////// @Description Create new product item
////// @Accept json
////// @Produce json
////// @Param CreateProductRequestDto body dto.CreateProductRequestDto true "Product data"
////// @Success 201 {object} dto.CreateProductResponseDto
////// @Router /products [post]
////func (h *productsController) CreateProduct() echo.HandlerFunc {
////	return func(c echo.Context) error {
////		h.metrics.CreateProductHttpRequests.Inc()
////
////		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.CreateProduct")
////		defer span.Finish()
////
////		request := &dto.CreateProductRequestDto{}
////		if err := c.Bind(request); err != nil {
////			h.log.WarnMsg("Bind", err)
////			h.traceErr(span, err)
////			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
////		}
////
////		if err := h.v.StructCtx(ctx, request); err != nil {
////			h.log.Errorf("(validate) err: {%v}", err)
////			tracing.TraceErr(span, err)
////			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
////		}
////
////		command := creating_product.NewCreateProduct(request.Name, request.Description, request.Price)
////		_, err := h.m.Send(ctx, *command)
////
////		if err != nil {
////			h.log.Errorf("(CreateOrder.Handle) id: {%s}, err: {%v}", command.ProductID, err)
////			tracing.TraceErr(span, err)
////			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
////		}
////
////		h.log.Infof("(order created) id: {%s}", command.ProductID)
////		//h.metrics.SuccessHttpRequests.Inc()
////		return c.JSON(http.StatusCreated, dto.CreateProductResponseDto{ProductID: command.ProductID})
////	}
////}
//
//// GetProductByID
//// @Tags Products
//// @Summary Get product
//// @Description Get product by id
//// @Accept json
//// @Produce json
//// @Param id path string true "Product ID"
//// @Success 200
//// @Router /products/{id} [get]
//func (h *productsController) GetProductByID() echo.HandlerFunc {
//	return func(c echo.Context) error {
//		h.metrics.GetProductByIdHttpRequests.Inc()
//
//		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.GetProductByID")
//		defer span.Finish()
//
//		productUUID, err := uuid.FromString(c.Param(constants.ID))
//		if err != nil {
//			h.log.WarnMsg("uuid.FromString", err)
//			h.traceErr(span, err)
//			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
//		}
//
//		query := getting_product_by_id.NewGetProductById(productUUID)
//		response, err := h.m.Send(ctx, query)
//
//		if err != nil {
//			h.log.WarnMsg("GetProductById", err)
//			h.metrics.ErrorHttpRequests.Inc()
//			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
//		}
//
//		h.metrics.SuccessHttpRequests.Inc()
//		return c.JSON(http.StatusOK, response)
//	}
//}
//
//// UpdateProduct
//// @Tags Products
//// @Summary Update product
//// @Description Update existing product
//// @Accept json
//// @Produce json
//// @Param id path string true "Product ID"
//// @Success 200 {object} dto.UpdateProductResponseDto
//// @Router /products/{id} [put]
//func (h *productsController) UpdateProduct() echo.HandlerFunc {
//	return func(c echo.Context) error {
//		h.metrics.UpdateProductHttpRequests.Inc()
//
//		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.UpdateProduct")
//		defer span.Finish()
//
//		productUUID, err := uuid.FromString(c.Param(constants.ID))
//		if err != nil {
//			h.log.WarnMsg("uuid.FromString", err)
//			h.traceErr(span, err)
//			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
//		}
//
//		request := &dto.UpdateProductRequestDto{ProductID: productUUID}
//		if err := c.Bind(request); err != nil {
//			h.log.WarnMsg("Bind", err)
//			h.traceErr(span, err)
//			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
//		}
//
//		if err := h.v.StructCtx(ctx, request); err != nil {
//			h.log.WarnMsg("validate", err)
//			h.traceErr(span, err)
//			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
//		}
//
//		command := updating_product.NewUpdateProduct(productUUID, request.Name, request.Description, request.Price)
//
//		_, err = h.m.Send(ctx, command)
//
//		if err != nil {
//			h.log.WarnMsg("UpdateProduct", err)
//			h.metrics.ErrorHttpRequests.Inc()
//			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
//		}
//
//		h.log.Infof("(order updated) id: {%s}", productUUID.String())
//		h.metrics.SuccessHttpRequests.Inc()
//		return c.JSON(http.StatusOK, dto.UpdateProductResponseDto{ProductID: productUUID, Name: request.Name, Description: request.Description, Price: request.Price})
//	}
//}
//
//// DeleteProduct
//// @Tags Products
//// @Summary Delete product
//// @Description Delete existing product
//// @Accept json
//// @Produce json
//// @Success 200 ""
//// @Param id path string true "Product ID"
//// @Router /products/{id} [delete]
//func (h *productsController) DeleteProduct() echo.HandlerFunc {
//	return func(c echo.Context) error {
//		h.metrics.DeleteProductHttpRequests.Inc()
//
//		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.DeleteProduct")
//		defer span.Finish()
//
//		productUUID, err := uuid.FromString(c.Param(constants.ID))
//		if err != nil {
//			h.log.WarnMsg("uuid.FromString", err)
//			h.traceErr(span, err)
//			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
//		}
//
//		command := deleting_product.NewDeleteProduct(productUUID)
//		_, err = h.m.Send(ctx, command)
//
//		if err != nil {
//			h.log.WarnMsg("DeleteProduct", err)
//			h.metrics.ErrorHttpRequests.Inc()
//			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
//		}
//
//		h.metrics.SuccessHttpRequests.Inc()
//		return c.NoContent(http.StatusOK)
//	}
//}
//
//func (h *productsController) traceErr(span opentracing.Span, err error) {
//	span.SetTag("error", true)
//	span.LogKV("error_code", err.Error())
//	h.metrics.ErrorHttpRequests.Inc()
//}
