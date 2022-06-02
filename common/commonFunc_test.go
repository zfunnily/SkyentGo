package common

import (
	"fmt"
	"pro2d/common/db/redisproxy"
	"pro2d/common/logger"
	"testing"
)

func TestGetNextRoleId(t *testing.T) {
	if err := redisproxy.ConnectRedis(GlobalConf.GameConf.RedisConf.DB, GlobalConf.GameConf.RedisConf.Auth, GlobalConf.GameConf.RedisConf.Address); err != nil {
		logger.Error(err)
		return
	}
	GlobalSconf = GlobalConf.GameConf

	fmt.Println(GetNextRoleId())
	fmt.Println(GetNextUId())
}
