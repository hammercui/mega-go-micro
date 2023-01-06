package mongo

import (
	"context"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo()  map[string]*mongo.Client {
	_map := make(map[string]*mongo.Client)
	if conf.GetConf().MongoMap == nil || len(conf.GetConf().MongoMap) == 0 {
		return _map
	}

	log.Logger().Infof("-------mongo init console-------")
	for k, v := range conf.GetConf().MongoMap {
		if v.Enable {
			_map[k] = newMongoClient(v)
		}
	}
	return _map
}

func newMongoClient(c *conf.MongoConf) *mongo.Client {
	credential := options.Credential{
		Username: c.Username,
		Password: c.Password,
	}
	clientOpts := options.Client().ApplyURI(c.Addr).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Logger().Errorf("mongo connect failed! uri: %s, err: %v", c.Addr,err)
		panic(err)
	}
	log.Logger().Infof("mongo connect success! uri: %s",c.Addr)
	return client
}



func clearMongoClient(mongoMap map[string]*mongo.Client)  {
	for _, value := range mongoMap{
		if value != nil {
			 value.Disconnect(context.TODO())
		}
	}
}