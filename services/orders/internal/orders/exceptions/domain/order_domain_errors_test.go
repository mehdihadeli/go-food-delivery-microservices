package domain

import (
	"fmt"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Order_Shop_Items_Required_Error(t *testing.T) {
	err := NewOrderShopItemsRequiredError("order items required")
	assert.True(t, IsOrderShopItemsRequiredError(err))
	fmt.Println(httpErrors.ErrorsWithStack(err))
}

func Test_Order_Not_Found_Error(t *testing.T) {
	err := NewOrderNotFoundError(1)
	assert.True(t, IsOrderNotFoundError(err))
	fmt.Println(httpErrors.ErrorsWithStack(err))
}

func Test_Invalid_Delivery_Address_Error(t *testing.T) {
	err := NewInvalidDeliveryAddressError("address is not valid")
	assert.True(t, IsInvalidDeliveryAddressError(err))
	fmt.Println(httpErrors.ErrorsWithStack(err))
}

func Test_InvalidEmail_Address_Error(t *testing.T) {
	err := NewInvalidEmailAddressError("email address is not valid")
	assert.True(t, IsInvalidEmailAddressError(err))
	fmt.Println(httpErrors.ErrorsWithStack(err))
}
