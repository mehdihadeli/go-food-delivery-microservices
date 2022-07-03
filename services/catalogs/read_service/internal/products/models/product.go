package models

import (
	"time"
)

type Product struct {
	ProductID   string    `json:"productId" bson:"_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Price       float64   `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// ProductsList products list response with pagination
type ProductsList struct {
	TotalCount int64      `json:"totalCount" bson:"totalCount"`
	TotalPages int64      `json:"totalPages" bson:"totalPages"`
	Page       int64      `json:"page" bson:"page"`
	Size       int64      `json:"size" bson:"size"`
	Products   []*Product `json:"products" bson:"products"`
}

//func NewProductListWithPagination(products []*Product, count int64, pagination *utils.ListResult[Product]) *ProductsList {
//	return &ProductsList{
//		TotalCount: count,
//		TotalPages: int64(pagination.TotalPage),
//		Page:       int64(pagination.Page),
//		Size:       int64(pagination.Size),
//		Products:   products,
//	}
//}
