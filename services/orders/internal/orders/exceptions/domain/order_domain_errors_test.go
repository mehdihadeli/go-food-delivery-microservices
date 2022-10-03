package domain

import (
	"fmt"
	errorUtils "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils/error_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Order_Shop_Items_Required_Error(t *testing.T) {
	err := NewOrderShopItemsRequiredError("order items required")
	assert.True(t, IsOrderShopItemsRequiredError(err))
	fmt.Println(errorUtils.ErrorsWithStack(err))
}

func Test_Order_Not_Found_Error(t *testing.T) {
	err := NewOrderNotFoundError(1)
	assert.True(t, IsOrderNotFoundError(err))
	fmt.Println(errorUtils.ErrorsWithStack(err))
}

func Test_Invalid_Delivery_Address_Error(t *testing.T) {
	err := NewInvalidDeliveryAddressError("address is not valid")
	assert.True(t, IsInvalidDeliveryAddressError(err))
	fmt.Println(errorUtils.ErrorsWithStack(err))
}

func Test_InvalidEmail_Address_Error(t *testing.T) {
	err := NewInvalidEmailAddressError("email address is not valid")
	assert.True(t, IsInvalidEmailAddressError(err))
	fmt.Println(errorUtils.ErrorsWithStack(err))
}
