package repository

import (
	"context"
	"log"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/data"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/data/specification"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	mongo2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/mongo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName   = "catalogs_write"
	CollectionName = "products"
)

// Product is a domain_events entity
type Product struct {
	ID          string
	Name        string
	Weight      int
	IsAvailable bool
}

type ProductMongo struct {
	ID          string `json:"id"          bson:"_id,omitempty"` // https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/#the-_id-field
	Name        string `json:"name"        bson:"name"`
	Weight      int    `json:"weight"      bson:"weight"`
	IsAvailable bool   `json:"isAvailable" bson:"isAvailable"`
}

func init() {
	err := mapper.CreateMap[*ProductMongo, *Product]()
	if err != nil {
		log.Fatal(err)
	}

	err = mapper.CreateMap[*Product, *ProductMongo]()
	if err != nil {
		log.Fatal(err)
	}
}

func Test_Add(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepository(ctx, t)

	product := &ProductMongo{
		ID: uuid.NewV4().
			String(),
		// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
		Name:        "added_product",
		Weight:      100,
		IsAvailable: true,
	}

	err = repository.Add(ctx, product)
	if err != nil {
		t.Fatal(err)
	}

	id, err := uuid.FromString(product.ID)
	if err != nil {
		return
	}

	p, err := repository.GetById(ctx, id)
	if err != nil {
		return
	}

	assert.NotNil(t, p)
	assert.Equal(t, product.ID, p.ID)
}

func Test_Add_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	product := &Product{
		ID: uuid.NewV4().
			String(),
		// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
		Name:        "added_product",
		Weight:      100,
		IsAvailable: true,
	}

	err = repository.Add(ctx, product)
	if err != nil {
		t.Fatal(err)
	}

	id, err := uuid.FromString(product.ID)
	if err != nil {
		return
	}

	p, err := repository.GetById(ctx, id)
	if err != nil {
		return
	}

	assert.NotNil(t, p)
	assert.Equal(t, product.ID, p.ID)
}

func Test_Get_By_Id(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepository(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	all, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		return
	}
	p := all.Items[0]
	id, err := uuid.FromString(p.ID)

	testCases := []struct {
		Name         string
		ProductId    uuid.UUID
		ExpectResult *ProductMongo
	}{
		{
			Name:         "ExistingProduct",
			ProductId:    id,
			ExpectResult: p,
		},
		{
			Name:         "NonExistingProduct",
			ProductId:    uuid.NewV4(),
			ExpectResult: nil,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			res, err := repository.GetById(ctx, c.ProductId)
			if c.ExpectResult == nil {
				assert.Error(t, err)
				assert.True(t, customErrors.IsNotFoundError(err))
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, p.ID, res.ID)
			}
		})
	}
}

func Test_Get_By_Id_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	all, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		return
	}

	p := all.Items[0]
	id, err := uuid.FromString(p.ID)

	testCases := []struct {
		Name         string
		ProductId    uuid.UUID
		ExpectResult *Product
	}{
		{
			Name:         "ExistingProduct",
			ProductId:    id,
			ExpectResult: p,
		},
		{
			Name:         "NonExistingProduct",
			ProductId:    uuid.NewV4(),
			ExpectResult: nil,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			res, err := repository.GetById(ctx, c.ProductId)
			if c.ExpectResult == nil {
				assert.Error(t, err)
				assert.True(t, customErrors.IsNotFoundError(err))
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, p.ID, res.ID)
			}
		})
	}
}

func Test_First_Or_Default(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepository(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	all, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		return
	}
	p := all.Items[0]

	single, err := repository.FirstOrDefault(ctx, map[string]interface{}{"_id": p.ID})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, single)
}

func Test_First_Or_Default_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	all, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		return
	}
	p := all.Items[0]

	single, err := repository.FirstOrDefault(ctx, map[string]interface{}{"_id": p.ID})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, single)
}

func Test_Get_All(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepository(ctx, t)
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
	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
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
	repository, err := setupGenericMongoRepository(ctx, t)
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
	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)

	models, err := repository.Search(ctx, "seed_product1", utils.NewListQuery(10, 1))
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, models.Items)
	assert.Equal(t, len(models.Items), 1)
}

func Test_GetByFilter(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepository(ctx, t)

	models, err := repository.GetByFilter(ctx, map[string]interface{}{"name": "seed_product1"})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, models)
	assert.Equal(t, len(models), 1)
}

func Test_GetByFilter_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)

	models, err := repository.GetByFilter(ctx, map[string]interface{}{"name": "seed_product1"})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, models)
	assert.Equal(t, len(models), 1)
}

func Test_Update(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepository(ctx, t)

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

	id, err := uuid.FromString(product.ID)
	if err != nil {
		t.Fatal(err)
	}

	single, err := repository.GetById(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, single)
	assert.Equal(t, "product2_updated", single.Name)
}

func Test_Update_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)

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

	id, err := uuid.FromString(product.ID)
	if err != nil {
		t.Fatal(err)
	}

	single, err := repository.GetById(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, single)
	assert.Equal(t, "product2_updated", single.Name)
}

func Test_Delete(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepository(ctx, t)

	products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		t.Fatal(err)
	}
	product := products.Items[0]

	id, err := uuid.FromString(product.ID)
	if err != nil {
		t.Fatal(err)
	}

	err = repository.Delete(ctx, id)
	if err != nil {
		return
	}

	single, err := repository.GetById(ctx, id)
	assert.Nil(t, single)
}

func Test_Delete_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)

	products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
	if err != nil {
		t.Fatal(err)
	}
	product := products.Items[0]

	id, err := uuid.FromString(product.ID)
	if err != nil {
		t.Fatal(err)
	}

	err = repository.Delete(ctx, id)
	if err != nil {
		return
	}

	single, err := repository.GetById(ctx, id)
	assert.Nil(t, single)
}

func Test_Count(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepository(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	count := repository.Count(ctx)
	assert.Equal(t, count, int64(2))
}

func Test_Count_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	count := repository.Count(ctx)
	assert.Equal(t, count, int64(2))
}

func Test_Skip_Take(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepository(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	entities, err := repository.SkipTake(ctx, 1, 1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(entities), 1)
}

func Test_Skip_Take_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	entities, err := repository.SkipTake(ctx, 1, 1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(entities), 1)
}

func Test_Find(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepository(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	entities, err := repository.Find(
		ctx,
		specification.And(specification.Equal("is_available", true), specification.Equal("name", "seed_product1")),
	)
	if err != nil {
		return
	}
	assert.Equal(t, len(entities), 1)
}

func Test_Find_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	entities, err := repository.Find(
		ctx,
		specification.And(specification.Equal("is_available", true), specification.Equal("name", "seed_product1")),
	)
	if err != nil {
		return
	}
	assert.Equal(t, len(entities), 1)
}

func setupGenericMongoRepositoryWithDataModel(
	ctx context.Context,
	t *testing.T,
) (data.GenericRepositoryWithDataModel[*ProductMongo, *Product], error) {
	defaultLogger.SetupDefaultLogger()

	db, err := mongo2.NewMongoTestContainers(defaultLogger.Logger).Start(ctx, t)
	if err != nil {
		return nil, err
	}

	err = seedAndMigration(ctx, db)
	if err != nil {
		return nil, err
	}

	return NewGenericMongoRepositoryWithDataModel[*ProductMongo, *Product](db, DatabaseName, CollectionName), nil
}

func setupGenericMongoRepository(ctx context.Context, t *testing.T) (data.GenericRepository[*ProductMongo], error) {
	defaultLogger.SetupDefaultLogger()

	db, err := mongo2.NewMongoTestContainers(defaultLogger.Logger).Start(ctx, t)
	if err != nil {
		return nil, err
	}

	err = seedAndMigration(ctx, db)
	if err != nil {
		return nil, err
	}

	return NewGenericMongoRepository[*ProductMongo](db, DatabaseName, CollectionName), nil
}

func seedAndMigration(ctx context.Context, db *mongo.Client) error {
	seedProducts := []*ProductMongo{
		{
			ID: uuid.NewV4().
				String(),
			// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
			Name:        "seed_product1",
			Weight:      100,
			IsAvailable: true,
		},
		{
			ID: uuid.NewV4().
				String(),
			// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
			Name:        "seed_product2",
			Weight:      100,
			IsAvailable: true,
		},
	}

	// https://go.dev/doc/faq#convert_slice_of_interface
	data := make([]interface{}, len(seedProducts))
	for i, v := range seedProducts {
		data[i] = v
	}

	collection := db.Database(DatabaseName).Collection(CollectionName)
	_, err := collection.InsertMany(ctx, data, &options.InsertManyOptions{})
	if err != nil {
		return err
	}

	return nil
}
