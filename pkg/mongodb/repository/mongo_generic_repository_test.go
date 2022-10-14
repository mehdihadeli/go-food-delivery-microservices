package repository

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

const (
	DatabaseName   = "catalogs"
	CollectionName = "products"
)

// Product is a domain entity
type Product struct {
	ID          string
	Name        string
	Weight      int
	IsAvailable bool
}

type ProductMongo struct {
	ID          string `json:"id" bson:"_id,omitempty"` //https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/#the-_id-field
	Name        string `json:"name" bson:"name"`
	Weight      int    `json:"weight" bson:"weight"`
	IsAvailable bool   `json:"isAvailable" bson:"isAvailable"`
}

//	func Test_Add(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepository(ctx)
//
//		product := &ProductMongo{
//			Name:        "added_product",
//			Weight:      100,
//			IsAvailable: true,
//		}
//
//		err = repository.Add(ctx, product)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		p, err := repository.GetById(ctx, product.ID)
//		if err != nil {
//			return
//		}
//
//		assert.NotNil(t, p)
//		assert.Equal(t, product.ID, p.ID)
//	}
//
//	func Test_Add_With_Data_Model(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepositoryWithDataModel(ctx)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		product := &Product{
//			ID:          uuid.NewV4(),
//			Name:        "added_product",
//			Weight:      100,
//			IsAvailable: true,
//		}
//
//		err = repository.Add(ctx, product)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		p, err := repository.GetById(ctx, product.ID)
//		if err != nil {
//			return
//		}
//
//		assert.NotNil(t, p)
//		assert.Equal(t, product.ID, p.ID)
//	}
//
//	func Test_Get_By_Id(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepository(ctx)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		all, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
//		if err != nil {
//			return
//		}
//		p := all.Items[0]
//
//		single, err := repository.GetById(ctx, p.ID)
//		if err != nil {
//			t.Fatal(err)
//		}
//		assert.NotNil(t, single)
//
//		nilResult, err := repository.GetById(ctx, uuid.NewV4())
//		assert.Nil(t, nilResult)
//	}
//
//	func Test_Get_By_Id_With_Data_Model(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepositoryWithDataModel(ctx)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		all, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
//		if err != nil {
//			return
//		}
//		p := all.Items[0]
//
//		single, err := repository.GetById(ctx, p.ID)
//		if err != nil {
//			t.Fatal(err)
//		}
//		assert.NotNil(t, single)
//
//		nilResult, err := repository.GetById(ctx, uuid.NewV4())
//		assert.Nil(t, nilResult)
//	}
//
//	func Test_Get_All(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepository(ctx)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		models, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		assert.NotEmpty(t, models.Items)
//	}
//
//	func Test_Get_All_With_Data_Model(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepositoryWithDataModel(ctx)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		models, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		assert.NotEmpty(t, models.Items)
//	}
//
//	func Test_Search(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepository(ctx)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		models, err := repository.Search(ctx, "seed_product1", utils.NewListQuery(10, 1))
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		assert.NotEmpty(t, models.Items)
//		assert.Equal(t, len(models.Items), 1)
//	}
//
//	func Test_Search_With_Data_Model(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepositoryWithDataModel(ctx)
//
//		models, err := repository.Search(ctx, "seed_product1", utils.NewListQuery(10, 1))
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		assert.NotEmpty(t, models.Items)
//		assert.Equal(t, len(models.Items), 1)
//	}
//
//	func Test_Where(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepository(ctx)
//
//		models, err := repository.Where(ctx, map[string]interface{}{"Name": "seed_product1"})
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		assert.NotEmpty(t, models)
//		assert.Equal(t, len(models), 1)
//	}
//
//	func Test_Where_With_Data_Model(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepositoryWithDataModel(ctx)
//
//		models, err := repository.Where(ctx, map[string]interface{}{"Name": "seed_product1"})
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		assert.NotEmpty(t, models)
//		assert.Equal(t, len(models), 1)
//	}
//
//	func Test_Update(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepository(ctx)
//
//		products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
//		if err != nil {
//			t.Fatal(err)
//		}
//		product := products.Items[0]
//
//		product.Name = "product2_updated"
//		err = repository.Update(ctx, product)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		single, err := repository.GetById(ctx, product.ID)
//		if err != nil {
//			t.Fatal(err)
//		}
//		assert.NotNil(t, single)
//		assert.Equal(t, "product2_updated", single.Name)
//	}
//
//	func Test_Update_With_Data_Model(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepositoryWithDataModel(ctx)
//
//		products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
//		if err != nil {
//			t.Fatal(err)
//		}
//		product := products.Items[0]
//
//		product.Name = "product2_updated"
//		err = repository.Update(ctx, product)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		single, err := repository.GetById(ctx, product.ID)
//		if err != nil {
//			t.Fatal(err)
//		}
//		assert.NotNil(t, single)
//		assert.Equal(t, "product2_updated", single.Name)
//	}
//
//	func Test_Delete(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepository(ctx)
//
//		products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
//		if err != nil {
//			t.Fatal(err)
//		}
//		product := products.Items[0]
//
//		err = repository.Delete(ctx, product.ID)
//		if err != nil {
//			return
//		}
//
//		single, err := repository.GetById(ctx, product.ID)
//		assert.Nil(t, single)
//	}
//
//	func Test_Delete_With_Data_Model(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepositoryWithDataModel(ctx)
//
//		products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
//		if err != nil {
//			t.Fatal(err)
//		}
//		product := products.Items[0]
//
//		err = repository.Delete(ctx, product.ID)
//		if err != nil {
//			return
//		}
//
//		single, err := repository.GetById(ctx, product.ID)
//		assert.Nil(t, single)
//	}
//
//	func Test_Count(t *testing.T) {
//		ctx := context.Background()
//		repository, err := setupGenericGormRepository(ctx)
//
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		count := repository.Count(ctx)
//
//		assert.Equal(t, count, int64(2))
//	}
func Test_Count_With_Data_Model(t *testing.T) {
	ctx := context.Background()
	_, err := setupGenericMongoRepositoryWithDataModel(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	//count := repository.Count(ctx)
	//
	//assert.Equal(t, count, int64(2))
}

func setupGenericMongoRepositoryWithDataModel(ctx context.Context, t *testing.T) (data.GenericRepositoryWithDataModel[*ProductMongo, *Product], error) {
	db, err := dockertest.NewDockerTestMongoContainer().Start(ctx, t)
	if err != nil {
		return nil, err
	}

	err = seedData(ctx, db)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//	func setupGenericGormRepository(ctx context.Context, t *testing.T) (data.GenericRepository[*ProductGorm], error) {
//		db, err := testcontainer.NewTestContainerGormContainer().Start(ctx, t)
//
//		err = seedData(ctx, db)
//		if err != nil {
//			return nil, err
//		}
//
//		return NewGenericGormRepository[*ProductGorm](db), nil
//	}
func seedData(ctx context.Context, db *mongo.Client) error {
	var seedProducts = []*ProductMongo{
		{
			Name:        "seed_product1",
			Weight:      100,
			IsAvailable: true,
		},
		{
			Name:        "seed_product2",
			Weight:      100,
			IsAvailable: true,
		},
	}

	//https://go.dev/doc/faq#convert_slice_of_interface
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
