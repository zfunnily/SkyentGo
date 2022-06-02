package mongoproxy

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"pro2d/common"
	"pro2d/common/logger"
	"reflect"
	"sort"
	"strings"
	"time"
)

func DB() *mongo.Database {
	return mongoDatabase
}

func ConnectMongo(conf *common.MongoConf, ID int64) error {
	var uri string
	if conf.User != "" {
		//uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?w=majority", conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/?w=majority", conf.User, conf.Password, conf.Host, conf.Port)
	} else {
		//uri = fmt.Sprintf("mongodb://%s:%d/%s?w=majority", conf.Host, conf.Port, conf.DBName)
		uri = fmt.Sprintf("mongodb://%s:%d/?w=majority", conf.Host, conf.Port)
	}
	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.TimeOut)*time.Second)
	defer cancel()
	// 通过传进来的uri连接相关的配置
	o := options.Client().ApplyURI(uri)
	// 设置最大连接数 - 默认是100 ，不设置就是最大 max 64
	o.SetMaxPoolSize(uint64(conf.MaxNum))
	// 发起链接
	var err error
	mongoClient, err = mongo.Connect(ctx, o)
	if err != nil {
		return err
	}
	// 判断服务是不是可用
	if err = mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
		return err
	}

	if conf.DBName != "account" {
		mongoDatabase = mongoClient.Database(fmt.Sprintf("%s_%d", conf.DBName, ID))
	} else {
		mongoDatabase = mongoClient.Database(conf.DBName)
	}
	return nil
}

func CloseMongo() {
	mongoClient.Disconnect(context.TODO())
}

func CreateTable(tb string) error {
	colls, _ := DB().ListCollectionNames(context.TODO(), bson.D{})
	pos := sort.SearchStrings(colls, tb)
	if pos != len(colls) {
		if tb == colls[pos] {
			return DB().CreateCollection(context.TODO(), tb)
		}
	}
	return DB().CreateCollection(context.TODO(), tb)
}

func FindOne(coll string, pri interface{}, schema interface{}) error {
	r := mongoDatabase.Collection(coll).FindOne(context.TODO(), pri)
	return r.Decode(schema)
}

func FindMany(coll string, key string, val interface{}, schema interface{}) error {
	r, err := mongoDatabase.Collection(coll).Find(context.TODO(), bson.M{key: val})
	if err != nil {
		return err
	}
	return r.All(context.TODO(), schema)
}

func DelOne(coll string, key string, value interface{}) error {
	filter := bson.D{{key, value}}
	_, err := mongoDatabase.Collection(coll).DeleteOne(context.TODO(), filter, nil)
	return err
}

func DelMany(coll string, filter bson.D) error {
	r, err := mongoDatabase.Collection(coll).DeleteMany(context.TODO(), filter, nil)
	logger.Debug(r.DeletedCount)
	return err
}

func GetBsonD(key string, value interface{}) interface{} {
	return bson.D{{key, value}}
}

func GetBsonM(key string, value interface{}) interface{} {
	return bson.M{key: value}
}

func GetSchemaType(schema interface{}) reflect.Type {
	s := reflect.TypeOf(schema)
	if s.Kind() == reflect.Ptr {
		s = reflect.TypeOf(schema).Elem()
	}
	return s
}

func GetCollName(schema interface{}) string {
	return strings.ToLower(GetSchemaType(schema).Name())
}

func GetPriKey(schema interface{}) string {
	s := GetSchemaType(schema)

	var pri string
	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).Tag.Get("pri") == "1" {
			pri = strings.ToLower(s.Field(i).Name)
			break
		}
	}
	return pri
}

func FindIndex(schema interface{}) (string, []string) {
	s := GetSchemaType(schema)

	var index []string
	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).Tag.Get("index") != "" {
			js := strings.Split(s.Field(i).Tag.Get("json"), ",")
			if len(js) == 0 {
				continue
			}
			index = append(index, js[0])
		}
	}
	return strings.ToLower(s.Name()), index
}

func SetUnique(coll, key string) (string, error) {
	return DB().Collection(coll).Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{key, bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)
}

func InitDoc(schema ...interface{}) {
	for _, s := range schema {
		coll, keys := FindIndex(s)
		CreateTable(coll)
		for _, index := range keys {

			logger.Debug("InitDoc collect: %v, createIndex: %s", coll, index)
			res, err := SetUnique(coll, index)
			if err != nil {
				logger.Error("InitDoc unique: %s, err: %v", res, err)
				continue
			}
		}
	}

}
