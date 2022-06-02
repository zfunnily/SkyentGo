package mongoproxy

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"pro2d/common/components"
	"sort"
	"strings"
)

var (
	mongoClient   *mongo.Client
	mongoDatabase *mongo.Database
)

type MgoColl struct {
	components.IDB
	Schema components.ISchema

	dbname string
	coll   *mongo.Collection
}

func NewMongoColl(dbname string, schema components.ISchema) *MgoColl {
	m := &MgoColl{
		dbname: dbname,
		coll:   DB().Collection(dbname),
		Schema: schema,
	}
	return m
}

func (m *MgoColl) CreateTable() error {
	colls, _ := DB().ListCollectionNames(context.TODO(), bson.D{})
	pos := sort.SearchStrings(colls, m.dbname)
	if pos != len(colls) {
		if m.dbname == colls[pos] {
			return DB().CreateCollection(context.TODO(), m.dbname)
		}
	}
	return DB().CreateCollection(context.TODO(), m.dbname)
}

func (m *MgoColl) Create() (interface{}, error) {
	return m.coll.InsertOne(context.TODO(), m.Schema.GetSchema())
}

func (m *MgoColl) Save() error {
	_, err := m.coll.UpdateOne(context.TODO(), m.Schema.GetPri(), bson.D{{"$set", m.Schema.GetSchema()}})
	if err != nil {
		return err
	}
	return nil
}

func (m *MgoColl) Load() error {
	r := m.coll.FindOne(context.TODO(), m.Schema.GetPri())
	err := r.Decode(m.Schema.GetSchema())
	if err != nil {
		return err
	}
	return nil
}

// 查询单个
func (m *MgoColl) FindOne() error {
	singleResult := m.coll.FindOne(context.TODO(), m.Schema.GetPri())
	return singleResult.Decode(m.Schema.GetSchema())
}

func (m *MgoColl) UpdateOne(filter interface{}, update interface{}) *mongo.UpdateResult {
	res, err := m.coll.UpdateOne(context.TODO(), filter, bson.D{{"$set", update}})
	if err != nil {
		return nil
	}
	return res
}

func (m *MgoColl) UpdateProperty(key string, val interface{}) error {
	_, err := m.coll.UpdateOne(context.TODO(), m.Schema.GetPri(), bson.D{{"$set", bson.M{strings.ToLower(key): val}}})
	return err
}

func (m *MgoColl) UpdateProperties(properties map[string]interface{}) error {
	_, err := m.coll.UpdateOne(context.TODO(), m.Schema.GetPri(), bson.D{{"$set", properties}})
	return err
}

//索引
func (m *MgoColl) SetUnique(key string) (string, error) {
	return m.coll.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{key, bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)
}

func (m *MgoColl) Delete(key string, value interface{}) int64 {
	filter := bson.D{{key, value}}
	count, err := m.coll.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		fmt.Println(err)
	}
	return count.DeletedCount

}
