package redisproxy

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var RedisPool *redis.Pool

//conf *conf.ServerConf
func ConnectRedis(db int, auth, address string) error {
	RedisPool = &redis.Pool{
		//最大活跃连接数，0代表无限
		MaxActive: 888,
		MaxIdle:   20,
		//闲置连接的超时时间
		IdleTimeout: time.Second * 100,
		//定义拨号获得连接的函数
		Dial: func() (redis.Conn, error) {
			option := []redis.DialOption{redis.DialDatabase(db)}
			if auth != "" {
				option = append(option, redis.DialPassword(auth))
			}
			return redis.Dial("tcp", address, option...)
		},
	}
	return nil
}

func CloseRedis() {
	RedisPool.Close()
}

func redisCommand(command string, args ...interface{}) (reply interface{}, err error) {
	conn := RedisPool.Get()
	defer conn.Close()
	return conn.Do(command, args...)
}

func ExpireKey(key interface{}, ttl interface{}) (reply interface{}, err error) {
	return redisCommand("expire", key, ttl)
}

//redis 管道操作
func PipLine(f func(conn redis.Conn)) {
	conn := RedisPool.Get()
	defer conn.Close()
	f(conn)
}

func PipLineTest() {
	PipLine(func(c redis.Conn) {
		c.Send("SET", "foo", "bar")
		c.Send("GET", "foo")
		c.Flush()
		//receive一次只从结果中拿出一个send的命令进行处理
		c.Receive()        // reply from SET
		_, _ = c.Receive() // reply from GET
	})
}

func SETNX(args ...interface{}) (reply interface{}, err error) {
	return redisCommand("SETNX", args...)
}

func SET(args ...interface{}) (reply interface{}, err error) {
	return redisCommand("SET", args...)
}

func GET(args ...interface{}) (reply interface{}, err error) {
	return redisCommand("GET", args...)
}

func DEL(args ...interface{}) (reply interface{}, err error) {
	return redisCommand("DEL", args...)
}

func HKEYS(args ...interface{}) (reply interface{}, err error) {
	return redisCommand("HKEYS", args...)
}

func HMSET(args ...interface{}) (reply interface{}, err error) {
	return redisCommand("HMSET", args...)
}

func HMGET(args ...interface{}) (reply interface{}, err error) {
	return redisCommand("HMGET", args...)
}

func HSET(args ...interface{}) (reply interface{}, err error) {
	return redisCommand("HSET", args...)
}

func HGET(args ...interface{}) (reply interface{}, err error) {
	return redisCommand("HGET", args...)
}

func HINCRBY(args ...interface{}) (reply interface{}, err error) {
	return redisCommand("HINCRBY", args...)
}
