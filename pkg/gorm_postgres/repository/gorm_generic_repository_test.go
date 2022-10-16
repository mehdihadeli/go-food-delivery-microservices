package repository

import (
	"context"
	_ "github.com/lib/pq" // postgres driver
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data/specification"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/testcontainer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"testing"
)

// Product is a domain entity
type Product struct {
	ID          uuid.UUID
	Name        string
	Weight      int
	IsAvailable bool
}

// ProductGorm is DTO used to map Product entity to database
type ProductGorm struct {
	ID          uuid.UUID `gorm:"primaryKey;column:id"`
	Name        string    `gorm:"column:name"`
	Weight      int       `gorm:"column:weight"`
	IsAvailable bool      `gorm:"column:is_available"`
}

func init() {
	err := mapper.CreateMap[*ProductGorm, *Product]()
	if err != nil {
		log.Fatal(err)
	}

	err = mapper.CreateMap[*Product, *ProductGorm]()
	if err != nil {
		log.Fatal(err)
	}
}

func Test_Add(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepository(ctx, t)

	product := &ProductGorm{
		ID:          uuid.NewV4(),
		Name:        "added_product",
		Weight:      100,
		IsAvailable: true,
	}

	err = repository.Add(ctx, product)
	if err != nil {
		t.Fatal(err)
	}

	p, err := repository.GetById(ctx, product.ID)
	if err != nil {
		return
	}

	assert.NotNil(t, p)
	assert.Equal(t, product.ID, p.ID)
}

func Test_Add_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	product := &Product{
		ID:          uuid.NewV4(),
		Name:        "added_product",
		Weight:      100,
		IsAvailable: true,
	}

	err = repository.Add(ctx, product)
	if err != nil {
		t.Fatal(err)
	}

	p, err := repository.GetById(ctx, product.ID)
	if err != nil {
		return
	}

	assert.NotNil(t, p)
	assert.Equal(t, product.ID, p.ID)
}

func Test_Get_By_Id(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepository(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	all, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		return
	}
	p := all.Items[0]

	testCases := []struct {
		Name      string
		ProductId uuid.UUID
	}{
		{
			Name:      "ExistingProduct",
			ProductId: p.ID,
		},
		{
			Name:      "NonExistingProduct",
			ProductId: uuid.NewV4(),
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			if c.Name == "NonExistingProduct" {
				nilResult, err := repository.GetById(ctx, c.ProductId)
				if err != nil {
					t.Fatal(err)
				}
				assert.Nil(t, nilResult)
			} else {
				single, err := repository.GetById(ctx, c.ProductId)
				if err != nil {
					t.Fatal(err)
				}
				assert.NotNil(t, single)
			}
		})
	}
}

func Test_Get_By_Id_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	all, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		return
	}
	p := all.Items[0]

	testCases := []struct {
		Name      string
		ProductId uuid.UUID
	}{
		{
			Name:      "ExistingProduct",
			ProductId: p.ID,
		},
		{
			Name:      "NonExistingProduct",
			ProductId: uuid.NewV4(),
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			if c.Name == "NonExistingProduct" {
				nilResult, err := repository.GetById(ctx, c.ProductId)
				if err != nil {
					t.Fatal(err)
				}
				assert.Nil(t, nilResult)
			} else {
				single, err := repository.GetById(ctx, c.ProductId)
				if err != nil {
					t.Fatal(err)
				}
				assert.NotNil(t, single)
			}
		})
	}
}

func Test_Get_All(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepository(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	models, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, models.Items)
}

func Test_Get_All_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	models, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, models.Items)
}

func Test_Search(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepository(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	models, err := repository.Search(ctx, "seed_product1", utils.NewListQuery(10, 1))
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, models.Items)
	assert.Equal(t, len(models.Items), 1)
}

func Test_Search_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)

	models, err := repository.Search(ctx, "seed_product1", utils.NewListQuery(10, 1))
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, models.Items)
	assert.Equal(t, len(models.Items), 1)
}

func Test_Where(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepository(ctx, t)

	models, err := repository.Where(ctx, map[string]interface{}{"name": "seed_product1"})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, models)
	assert.Equal(t, len(models), 1)
}

func Test_Where_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)

	models, err := repository.Where(ctx, map[string]interface{}{"name": "seed_product1"})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, models)
	assert.Equal(t, len(models), 1)
}

func Test_Update(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepository(ctx, t)

	products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		t.Fatal(err)
	}
	product := products.Items[0]

	product.Name = "product2_updated"
	err = repository.Update(ctx, product)
	if err != nil {
		t.Fatal(err)
	}

	single, err := repository.GetById(ctx, product.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, single)
	assert.Equal(t, "product2_updated", single.Name)
}

func Test_Update_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)

	products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		t.Fatal(err)
	}
	product := products.Items[0]

	product.Name = "product2_updated"
	err = repository.Update(ctx, product)
	if err != nil {
		t.Fatal(err)
	}

	single, err := repository.GetById(ctx, product.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, single)
	assert.Equal(t, "product2_updated", single.Name)
}

func Test_Delete(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepository(ctx, t)

	products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		t.Fatal(err)
	}
	product := products.Items[0]

	err = repository.Delete(ctx, product.ID)
	if err != nil {
		return
	}

	single, err := repository.GetById(ctx, product.ID)
	assert.Nil(t, single)
}

func Test_Delete_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)

	products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		t.Fatal(err)
	}
	product := products.Items[0]

	err = repository.Delete(ctx, product.ID)
	if err != nil {
		return
	}

	single, err := repository.GetById(ctx, product.ID)
	assert.Nil(t, single)
}

func Test_Count(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepository(ctx, t)

	if err != nil {
		t.Fatal(err)
	}

	count := repository.Count(ctx)

	assert.Equal(t, count, int64(2))
}

func Test_Count_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	count := repository.Count(ctx)

	assert.Equal(t, count, int64(2))
}

func Test_Find(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepository(ctx, t)

	if err != nil {
		t.Fatal(err)
	}

	entities, err := repository.Find(ctx, specification.And(specification.Equal("is_available", true), specification.Equal("name", "seed_product1")))
	if err != nil {
		return
	}
	assert.Equal(t, len(entities), 1)
}

func Test_Find_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)

	if err != nil {
		t.Fatal(err)
	}

	entities, err := repository.Find(ctx, specification.And(specification.Equal("is_available", true), specification.Equal("name", "seed_product1")))
	if err != nil {
		return
	}
	assert.Equal(t, len(entities), 1)
}

func setupGenericGormRepositoryWithDataModel(ctx context.Context, t *testing.T) (data.GenericRepositoryWithDataModel[*ProductGorm, *Product], error) {
	db, err := testcontainer.NewGormTestContainers().Start(ctx, t)
	if err != nil {
		return nil, err
	}

	err = seedData(ctx, db)
	if err != nil {
		return nil, err
	}

	return NewGenericGormRepositoryWithDataModel[*ProductGorm, *Product](db), nil
}

func setupGenericGormRepository(ctx context.Context, t *testing.T) (data.GenericRepository[*ProductGorm], error) {
	db, err := testcontainer.NewGormTestContainers().Start(ctx, t)

	err = seedData(ctx, db)
	if err != nil {
		return nil, err
	}

	return NewGenericGormRepository[*ProductGorm](db), nil
}

func seedData(ctx context.Context, db *gorm.DB) error {
	err := db.AutoMigrate(ProductGorm{})
	if err != nil {
		return err
	}

	seedProducts := []*ProductGorm{
		{
			ID:          uuid.NewV4(),
			Name:        "seed_product1",
			Weight:      100,
			IsAvailable: true,
		},
		{
			ID:          uuid.NewV4(),
			Name:        "seed_product2",
			Weight:      100,
			IsAvailable: true,
		},
	}

	err = db.WithContext(ctx).Create(seedProducts).Error
	if err != nil {
		return err
	}
	return nil
}
