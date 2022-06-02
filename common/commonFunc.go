package common

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"pro2d/common/db/redisproxy"
	"strconv"
	"strings"
)

func GetNextRoleId() (string, error) {
	relay, err := redisproxy.HGET(fmt.Sprintf(AutoIncrement, GlobalSconf.ID), "role")
	if err != nil {
		return "", err
	}
	ID, err := redis.Int64(relay, err)
	if err != nil {
		return "", err
	}

	//roleID的范围 [GlobalSconf.ID*MaxRoleNum, GlobalSconf.ID*MaxRoleNum + MaxRoleNum]
	if ID-GlobalSconf.ID*MaxRoleNum >= MaxCommNum-1 {
		return "", errors.New("DB_FULL")
	}

	relay, err = redisproxy.HINCRBY(fmt.Sprintf(AutoIncrement, GlobalSconf.ID), "role", 1)
	ID, err = redis.Int64(relay, err)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", ID), nil
}

func GetNextUId() (string, error) {
	relay, err := redisproxy.HGET(fmt.Sprintf(AutoIncrement, GlobalSconf.ID), "uid")
	if err != nil {
		return "", err
	}

	var ID int64 = 0
	if relay == nil {
		ID = 90000
		redisproxy.HSET(fmt.Sprintf(AutoIncrement, GlobalSconf.ID), "uid", ID)
	} else {
		relay, err = redisproxy.HINCRBY(fmt.Sprintf(AutoIncrement, GlobalSconf.ID), "uid", 1)
		ID, err = redis.Int64(relay, err)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%d", ID), nil
}

type IMapString map[string]interface{}
type IMapStringNum map[string]int32

func MapToString(params IMapString) string {
	var items bytes.Buffer
	for k, v := range params {
		items.WriteString(k)
		items.WriteString("=")
		items.WriteString(fmt.Sprintf("%v", v))
		items.WriteString(" ")
	}
	return items.String()
}

func StringToMap(items string, num bool) IMapString {
	backPack := make(map[string]interface{})
	for _, v := range strings.Split(items, " ") {
		ii := strings.Split(v, "=")
		if len(ii) < 2 {
			continue
		}
		if num {
			c, err := strconv.Atoi(ii[1])
			if err != nil {
				continue
			}
			backPack[ii[0]] = uint32(c)
		} else {
			backPack[ii[0]] = ii[1]
		}
	}
	return backPack
}

func MapNumToString(params IMapStringNum) string {
	var items bytes.Buffer
	for k, v := range params {
		items.WriteString(k)
		items.WriteString("=")
		items.WriteString(fmt.Sprintf("%v", v))
		items.WriteString(" ")
	}
	return items.String()
}

func StringToMapNum(items string) IMapStringNum {
	backPack := make(IMapStringNum)
	for _, v := range strings.Split(items, " ") {
		ii := strings.Split(v, "=")
		if len(ii) < 2 {
			continue
		}
		c, err := strconv.Atoi(ii[1])
		if err != nil {
			continue
		}
		backPack[ii[0]] = int32(c)
	}
	return backPack
}
