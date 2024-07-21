package repository

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/data"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	defaultLogger "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mongodb"
	mongocontainer "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/containers/testcontainer/mongo"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type mongoGenericRepositoryTest struct {
	suite.Suite
	databaseName                   string
	collectionName                 string
	mongoClient                    *mongo.Client
	productRepository              data.GenericRepository[*ProductMongo]
	productRepositoryWithDataModel data.GenericRepositoryWithDataModel[*ProductMongo, *Product]
	products                       []*ProductMongo
}

func TestMongoGenericRepository(t *testing.T) {
	suite.Run(
		t,
		&mongoGenericRepositoryTest{
			databaseName:   "catalogs_write",
			collectionName: "products",
		},
	)
}

func (c *mongoGenericRepositoryTest) SetupSuite() {
	opts, err := mongocontainer.NewMongoTestContainers(defaultLogger.GetLogger()).
		PopulateContainerOptions(context.Background(), c.T())
	c.Require().NoError(err)

	mongoClient, err := mongodb.NewMongoDB(opts)
	c.Require().NoError(err)
	c.mongoClient = mongoClient

	c.productRepository = NewGenericMongoRepository[*ProductMongo](
		mongoClient,
		c.databaseName,
		c.collectionName,
	)
	c.productRepositoryWithDataModel = NewGenericMongoRepositoryWithDataModel[*ProductMongo, *Product](
		mongoClient,
		c.databaseName,
		c.collectionName,
	)

	err = mapper.CreateMap[*ProductMongo, *Product]()
	c.Require().NoError(err)

	err = mapper.CreateMap[*Product, *ProductMongo]()
	c.Require().NoError(err)
}

func (c *mongoGenericRepositoryTest) SetupTest() {
	p, err := c.seedData(context.Background())
	c.Require().NoError(err)
	c.products = p
}

func (c *mongoGenericRepositoryTest) TearDownTest() {
	err := c.cleanupMongoData()
	c.Require().NoError(err)
}

func (c *mongoGenericRepositoryTest) Test_Add() {
	ctx := context.Background()

	product := &ProductMongo{
		// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
		ID:          uuid.NewV4().String(),
		Name:        gofakeit.Name(),
		Weight:      gofakeit.Number(100, 1000),
		IsAvailable: true,
	}

	err := c.productRepository.Add(ctx, product)
	c.Require().NoError(err)

	id, err := uuid.FromString(product.ID)
	c.Require().NoError(err)

	p, err := c.productRepository.GetById(ctx, id)
	c.Require().NoError(err)

	c.Assert().NotNil(p)
	c.Assert().Equal(product.ID, p.ID)
}

func (c *mongoGenericRepositoryTest) Test_Add_With_Data_Model() {
	ctx := context.Background()

	product := &ProductMongo{
		// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
		ID:          uuid.NewV4().String(),
		Name:        gofakeit.Name(),
		Weight:      gofakeit.Number(100, 1000),
		IsAvailable: true,
	}

	err := c.productRepository.Add(ctx, product)
	c.Require().NoError(err)

	id, err := uuid.FromString(product.ID)
	c.Require().NoError(err)

	p, err := c.productRepository.GetById(ctx, id)
	c.Require().NoError(err)

	c.Assert().NotNil(p)
	c.Assert().Equal(product.ID, p.ID)
}

func (c *mongoGenericRepositoryTest) Test_Get_By_Id() {
	ctx := context.Background()

	all, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	p := all.Items[0]
	id, err := uuid.FromString(p.ID)
	name := p.Name

	testCases := []struct {
		Name         string
		ProductId    uuid.UUID
		ExpectResult *ProductMongo
	}{
		{
			Name:         name,
			ProductId:    id,
			ExpectResult: p,
		},
		{
			Name:         "NonExistingProduct",
			ProductId:    uuid.NewV4(),
			ExpectResult: nil,
		},
	}

	for _, s := range testCases {
		c.T().Run(s.Name, func(t *testing.T) {
			t.Parallel()
			res, err := c.productRepository.GetById(ctx, s.ProductId)
			if s.ExpectResult == nil {
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

func (c *mongoGenericRepositoryTest) Test_Get_By_Id_With_Data_Model() {
	ctx := context.Background()

	all, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	p := all.Items[0]
	id, err := uuid.FromString(p.ID)
	name := p.Name

	testCases := []struct {
		Name         string
		ProductId    uuid.UUID
		ExpectResult *Product
	}{
		{
			Name:         name,
			ProductId:    id,
			ExpectResult: p,
		},
		{
			Name:         "NonExistingProduct",
			ProductId:    uuid.NewV4(),
			ExpectResult: nil,
		},
	}

	for _, s := range testCases {
		c.T().Run(s.Name, func(t *testing.T) {
			t.Parallel()
			res, err := c.productRepositoryWithDataModel.GetById(
				ctx,
				s.ProductId,
			)
			if s.ExpectResult == nil {
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

func (c *mongoGenericRepositoryTest) Test_First_Or_Default() {
	ctx := context.Background()

	all, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	p := all.Items[0]

	single, err := c.productRepository.FirstOrDefault(
		ctx,
		map[string]interface{}{"_id": p.ID},
	)
	c.Require().NoError(err)
	c.Assert().NotNil(single)
}

func (c *mongoGenericRepositoryTest) Test_First_Or_Default_With_Data_Model() {
	ctx := context.Background()

	all, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	p := all.Items[0]

	single, err := c.productRepositoryWithDataModel.FirstOrDefault(
		ctx,
		map[string]interface{}{"_id": p.ID},
	)

	c.Require().NoError(err)
	c.Assert().NotNil(single)
}

func (c *mongoGenericRepositoryTest) Test_Get_All() {
	ctx := context.Background()

	models, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
}

func (c *mongoGenericRepositoryTest) Test_Get_All_With_Data_Model() {
	ctx := context.Background()

	models, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
}

func (c *mongoGenericRepositoryTest) Test_Search() {
	ctx := context.Background()

	models, err := c.productRepository.Search(
		ctx,
		c.products[0].Name,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
	c.Assert().Equal(len(models.Items), 1)
}

func (c *mongoGenericRepositoryTest) Test_Search_With_Data_Model() {
	ctx := context.Background()

	models, err := c.productRepositoryWithDataModel.Search(
		ctx,
		c.products[0].Name,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
	c.Assert().Equal(len(models.Items), 1)
}

func (c *mongoGenericRepositoryTest) Test_GetByFilter() {
	ctx := context.Background()

	models, err := c.productRepository.GetByFilter(
		ctx,
		map[string]interface{}{"name": c.products[0].Name},
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models)
	c.Assert().Equal(len(models), 1)
}

func (c *mongoGenericRepositoryTest) Test_GetByFilter_With_Data_Model() {
	ctx := context.Background()

	models, err := c.productRepositoryWithDataModel.GetByFilter(
		ctx,
		map[string]interface{}{"name": c.products[0].Name},
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models)
	c.Assert().Equal(len(models), 1)
}

func (c *mongoGenericRepositoryTest) Test_Update() {
	ctx := context.Background()

	products, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	product := products.Items[0]

	product.Name = "product2_updated"
	err = c.productRepository.Update(ctx, product)
	c.Require().NoError(err)

	id, err := uuid.FromString(product.ID)
	c.Require().NoError(err)

	single, err := c.productRepository.GetById(ctx, id)
	c.Require().NoError(err)

	c.Assert().NotNil(single)
	c.Assert().Equal("product2_updated", single.Name)
}

func (c *mongoGenericRepositoryTest) Test_Update_With_Data_Model() {
	ctx := context.Background()

	products, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	product := products.Items[0]

	product.Name = "product2_updated"
	err = c.productRepositoryWithDataModel.Update(ctx, product)
	c.Require().NoError(err)

	id, err := uuid.FromString(product.ID)
	c.Require().NoError(err)

	single, err := c.productRepositoryWithDataModel.GetById(ctx, id)
	c.Require().NoError(err)

	c.Assert().NotNil(single)
	c.Assert().Equal("product2_updated", single.Name)
}

func (c *mongoGenericRepositoryTest) Test_Delete() {
	ctx := context.Background()

	products, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	product := products.Items[0]

	id, err := uuid.FromString(product.ID)
	c.Require().NoError(err)

	err = c.productRepository.Delete(ctx, id)
	c.Require().NoError(err)

	single, err := c.productRepository.GetById(ctx, id)
	c.Assert().Nil(single)
}

//func Test_Delete_With_Data_Model(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
//
//	products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
//	if err != nil {
//		t.Fatal(err)
//	}
//	product := products.Items[0]
//
//	id, err := uuid.FromString(product.ID)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	err = repository.Delete(ctx, id)
//	if err != nil {
//		return
//	}
//
//	single, err := repository.GetById(ctx, id)
//	assert.Nil(t, single)
//}
//
//func Test_Count(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericMongoRepository(ctx, t)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	count := repository.Count(ctx)
//	assert.Equal(t, count, int64(2))
//}
//
//func Test_Count_With_Data_Model(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	count := repository.Count(ctx)
//	assert.Equal(t, count, int64(2))
//}
//
//func Test_Skip_Take(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericMongoRepository(ctx, t)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	entities, err := repository.SkipTake(ctx, 1, 1)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	assert.Equal(t, len(entities), 1)
//}
//
//func Test_Skip_Take_With_Data_Model(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	entities, err := repository.SkipTake(ctx, 1, 1)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	assert.Equal(t, len(entities), 1)
//}
//
//func Test_Find(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericMongoRepository(ctx, t)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	entities, err := repository.Find(
//		ctx,
//		specification.And(
//			specification.Equal("is_available", true),
//			specification.Equal("name", "seed_product1"),
//		),
//	)
//	if err != nil {
//		return
//	}
//	assert.Equal(t, len(entities), 1)
//}
//
//func Test_Find_With_Data_Model(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	entities, err := repository.Find(
//		ctx,
//		specification.And(
//			specification.Equal("is_available", true),
//			specification.Equal("name", "seed_product1"),
//		),
//	)
//	if err != nil {
//		return
//	}
//	assert.Equal(t, len(entities), 1)
//}
//
//func setupGenericMongoRepositoryWithDataModel(
//	ctx context.Context,
//	t *testing.T,
//) (data.GenericRepositoryWithDataModel[*ProductMongo, *Product], error) {
//	db, err := mongocontainer.NewMongoTestContainers(defaultLogger.GetLogger()).
//		Start(ctx, t)
//	if err != nil {
//		return nil, err
//	}
//
//	err = seedData(ctx, db)
//	if err != nil {
//		return nil, err
//	}
//
//	return NewGenericMongoRepositoryWithDataModel[*ProductMongo, *Product](
//		db,
//		DatabaseName,
//		CollectionName,
//	), nil
//}

func (c *mongoGenericRepositoryTest) cleanupMongoData() error {
	collections := []string{c.collectionName}
	err := cleanupCollections(
		c.mongoClient,
		collections,
		c.databaseName,
	)

	return err
}

func cleanupCollections(
	db *mongo.Client,
	collections []string,
	databaseName string,
) error {
	database := db.Database(databaseName)
	ctx := context.Background()

	// Iterate over the collections and delete all collections
	for _, collection := range collections {
		collection := database.Collection(collection)

		err := collection.Drop(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *mongoGenericRepositoryTest) seedData(
	ctx context.Context,
) ([]*ProductMongo, error) {
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

	collection := c.mongoClient.Database(c.databaseName).
		Collection(c.collectionName)
	_, err := collection.InsertMany(ctx, data, &options.InsertManyOptions{})
	if err != nil {
		return nil, err
	}

	return seedProducts, nil
}
