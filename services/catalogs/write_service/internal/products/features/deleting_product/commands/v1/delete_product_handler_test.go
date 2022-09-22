package v1

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Delete_Product_Command_Handler(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*DeleteProduct, *mediatr.Unit](NewDeleteProductHandler(fixture.Log, fixture.Cfg, fixture.ProductRepository, fixture.Producer))
	if err != nil {
		return
	}

	fixture.Run()
	defer fixture.Cleanup()

	id, err := uuid.FromString("cb5bc70e-932e-4d35-8479-c6a939529876")
	if err != nil {
		return
	}
	command := NewDeleteProduct(id)
	result, err := mediatr.Send[*DeleteProduct, *mediatr.Unit](context.Background(), command)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}
