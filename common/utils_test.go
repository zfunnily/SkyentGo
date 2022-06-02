package common

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRandomString(t *testing.T) {
	rand.Seed(Timex())
	fmt.Println(RandomName(DefaultName))
}

func TestRandomName(t *testing.T) {
	//if err := redisproxy.ConnectRedis(GlobalConf.GameConf.RedisConf.DB, GlobalConf.GameConf.RedisConf.Auth, GlobalConf.GameConf.RedisConf.Address); err != nil {
	//	logger.Error(err)
	//	return
	//}
	fmt.Println(FirstCharToUpper("name"))
}
