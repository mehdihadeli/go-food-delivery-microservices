package gorm_postgres

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

// Product is a domain entity
type Product struct {
	ID          int
	Name        string
	Weight      int
	IsAvailable bool
}

// ProductGorm is DTO used to map Product entity to database
type ProductGorm struct {
	ID          int    `gorm:"primaryKey;column:id"`
	Name        string `gorm:"column:name"`
	Weight      int    `gorm:"column:weight"`
	IsAvailable bool   `gorm:"column:is_available"`
}

func (g *ProductGorm) ToEntity() *Product {
	return &Product{
		ID:          g.ID,
		Name:        g.Name,
		Weight:      g.Weight,
		IsAvailable: g.IsAvailable,
	}
}

func (g *ProductGorm) FromEntity(product *Product) {
	g.ID = product.ID
	g.Name = product.Name
	g.Weight = product.Weight
	g.IsAvailable = product.IsAvailable
}

type GormGenericRepository struct {
	*testing.T
	repository data.GenericRepository[*ProductGorm, *Product]
	ctx        context.Context
}

func TestRunner(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:test?mode=memory"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = db.AutoMigrate(ProductGorm{})
	if err != nil {
		t.Fatal(err)
	}

	var repository data.GenericRepository[*ProductGorm, *Product]
	repository = NewGenericGormRepository[*ProductGorm, *Product](db)

	ctx := context.Background()

	//https://pkg.go.dev/testing@master#hdr-Subtests_and_Sub_benchmarks
	t.Run("A=request-response", func(t *testing.T) {
		test := GormGenericRepository{T: t, repository: repository, ctx: ctx}
		test.Test_Add()
		test.Test_Get_By_Id()
		test.Test_Update()
	})
}

func (t *GormGenericRepository) Test_Add() {
	product := Product{
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	err := t.repository.Add(t.ctx, &product)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(product)
	assert.NotZero(t, product.ID)
}

func (t *GormGenericRepository) Test_Get_By_Id() {
	product := Product{
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	err := t.repository.Add(t.ctx, &product)
	if err != nil {
		t.Fatal(err)
	}

	single, err := t.repository.GetById(t.ctx, product.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(single)
	assert.NotNil(t, single)
}

func (t *GormGenericRepository) Test_Update() {
	product := Product{
		Name:        "product2",
		Weight:      50,
		IsAvailable: true,
	}
	err := t.repository.Add(t.ctx, &product)
	if err != nil {
		t.Fatal(err)
	}

	product.Name = "product2_updated"
	err = t.repository.Update(t.ctx, &product)
	if err != nil {
		return
	}

	single, err := t.repository.GetById(t.ctx, product.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(single)
	assert.NotNil(t, single)
	assert.Equal(t, "product2_updated", single.Name)
}