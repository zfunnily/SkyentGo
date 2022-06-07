package common

type RedisConf struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
	DB      int    `json:"db"`
}

type MongoConf struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	TimeOut  int    `yaml:"timeout"`
	MaxNum   int    `yaml:"maxnum"`
	DBName   string `yaml:"dbname"`
}

type SConf struct {
	ID             int64      `yaml:"id"`
	Name           string     `yaml:"name"`
	IP             string     `yaml:"ip"`
	Port           int        `yaml:"port"`
	DebugPort      int        `yaml:"debugport"`
	GMPort         int        `yaml:"gm"`
	Encipher       bool       `yaml:"encipher"`
	MongoConf      *MongoConf `yaml:"mongo"`
	RedisConf      *RedisConf `yaml:"redis"`
	WorkerPoolSize int        `yaml:"pool_size"`
	PluginPath     string     `yaml:"plugin_path"`
}

type LogConsole struct {
	Level string `yaml:"level" json:"level"`
	Color bool   `yaml:"color" json:"color"`
}

type LogFile struct {
	Level    string `yaml:"level" json:"level"`
	Daily    bool   `yaml:"daily" json:"daily"`
	Maxlines int    `yaml:"maxlines" json:"maxlines"`
	Maxsize  int    `yaml:"maxsize" json:"maxsize"`
	Maxdays  int    `yaml:"maxdays" json:"maxdays"`
	Append   bool   `yaml:"append" json:"append"`
	Permit   string `yaml:"permit" json:"permit"`
}

type LogConn struct {
	Net            string `yaml:"net" json:"net"`
	Addr           string `yaml:"addr" json:"addr"`
	Level          string `yaml:"level" json:"level"`
	Reconnect      bool   `yaml:"reconnect" json:"reconnect"`
	ReconnectOnMsg bool   `yaml:"reconnectOnMsg" json:"reconnectOnMsg"`
}

type LogConf struct {
	TimeFormat string      `yaml:"TimeFormat" json:"TimeFormat"`
	LogConsole *LogConsole `yaml:"Console" json:"Console"`
	LogFile    *LogFile    `yaml:"File" json:"File"`
	LogConn    *LogConn    `yaml:"Conn" json:"Conn"`
}

type TestClient struct {
	Ip    string `yaml:"ip"`
	Port  int    `yaml:"port"`
	Count int    `yaml:"count"`
}

//
//type GameService struct {
//	ServiceInfo []*pb.ServiceInfo `yaml:"server_list"`
//}

type ServerConf struct {
	ID           string      `yaml:"id"`
	Name         string      `yaml:"name"`
	WorkerID     int64       `yaml:"workerid"`
	DatacenterID int64       `yaml:"datacenterid"`
	AccountConf  *SConf      `yaml:"server_account"`
	GameConf     *SConf      `yaml:"server_game"`
	LogConf      *LogConf    `yaml:"logconf" json:"logconf"`
	TestClient   *TestClient `yaml:"test_client"`
	//GameService  *GameService `yaml:"server_list"`
}

var (
	GlobalConf  ServerConf
	GlobalSconf *SConf
)

func init() {
	//configFile, err := ioutil.ReadFile("conf/conf.yml")
	//if err != nil {
	//	fmt.Printf("conf read faild: %v", err)
	//	return
	//}

	//servList, err := ioutil.ReadFile("conf/serverlist.yml")
	//if err != nil {
	//	fmt.Printf("serverlist read faild: %v\n", err)
	//	return
	//}

	//初始化配置
	//if err = yaml.Unmarshal(configFile, &GlobalConf); err != nil {
	//	fmt.Printf("config.yml unmarshal faild: %v\n", err)
	//	return
	//}

	////游戏服务列表
	//if err = yaml.Unmarshal(servList, &GlobalConf.GameService); err != nil {
	//	fmt.Printf("serverlist.yml unmarshal faild: %v\n", err)
	//	return
	//}

	//c, err := json.Marshal(&GlobalConf.LogConf)
	//if err != nil {
	//	fmt.Errorf("log conf %v", err)
	//	return
	//}
	////初始化日志
	//err = logger.SetLogger(string(c), fmt.Sprintf("logs/%s", strings.ToLower(GlobalConf.GameConf.Name)))
	//if err != nil {
	//	fmt.Errorf("log conf %v", err)
	//	return
	//}
}
