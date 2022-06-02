package common

const (
	//最大包大
	MaxPacketLength = 10 * 1024 * 1024
	MaxMsgChan      = 100

	//心跳
	HeartTimerInterval   = 5  //s
	HeartTimeoutCountMax = 20 //最大超时次数

	//保存数据时间
	SaveDataInterval = 5 //s

	//自增id相关
	MaxCommNum = 1000000
	MaxRoleNum = MaxCommNum
	MaxHeroNum = 100
	MaxTeamNum = 10
	MaxUidNum  = 90000

	AutoIncrement     = "pro2d_autoincrement_set:%d"
	AutoIncrementHero = "pro2d_autoincrement:hero"

	//gm参数属性
	Role_ = "_role"
	Conn_ = "_conn"

	//背包容量
	LimitCommon = 100
)

//redis keys
const (
	NickNames = "name:%s"
	SMSCode   = "smscode:%s"
)
