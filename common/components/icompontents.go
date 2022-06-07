package components

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"pro2d/common/ccnet"
	"reflect"
)

//-----------------
//----net start----
//-----------------
type (
	//网络包头
	IHead interface {
		GetDataLen() uint32  //获取消息数据段长度
		GetMsgID() uint32    //获取消息ID
		GetErrCode() int32   //获取消息错误码
		GetPreserve() uint32 //获取预留数据
	}
	// IMessage 网络包
	IMessage interface {
		GetHeader() IHead       //获取消息头
		SetHeader(header IHead) //设置消息头

		GetData() []byte //获取消息内容
		SetData([]byte)  //设置消息内容

		SetSID(id uint32) //设置连接
		GetSID() uint32   //获取连接
	}
	// ISplitter 网络拆包解包器
	ISplitter interface {
		UnPack([]byte) (IMessage, error)
		Pack(cmd uint32, data []byte, errcode int32, preserve uint32) ([]byte, error)
		ParseMsg(data []byte, atEOF bool) (advance int, token []byte, err error)
		GetHeadLen() uint32
	}
	// IEncipher 加解密
	IEncipher interface {
		Encrypt([]byte) ([]byte, error)
		Decrypt([]byte) ([]byte, error)
	}

	// 事件 Poll
	IPoll interface {
		CreatePoll()
		AddEvent()
		DelEvent()
		Poll()
		Name()
	}

	// 上下文
	IContext interface {
		Handle() uint32
	}

	// IConnection 链接
	IConnection interface {
		GetID() uint32
		Start()
		Stop()
		Send(errCode int32, cmd uint32, b []byte) error
		SendSuccess(cmd uint32, b []byte) error
		CustomChan() chan<- func()
		GetCtx() IContext

		SetConnectionCallback(ccnet.ConnectionCallback)
		SetMessageCallback(ccnet.MessageCallback)
		SetCloseCallback(ccnet.CloseCallback)
		SetTimerCallback(ccnet.TimerCallback)
	}
	// IConnManage connManage
	IConnManage interface {
		AddConn(id uint32, connection IConnection)
		GetConn(id uint32) IConnection
		DelConn(id uint32) IConnection
		Range(f func(key interface{}, value interface{}) bool)
		StopAllConns()

		AddRID(rid string, id uint32)
		DelRID(rid string)
		GetConnByRID(rid string) IConnection

		AddUID(uid string, id uint32)
		DelUID(uid string)
		GetConnByUID(uid string) IConnection
	}
	// IServer server
	IServer interface {
		Start() error
		Stop()

		GetSplitter() ISplitter
		GetPlugin() IPlugin
		GetAction(uint32) interface{}
		SetActions(map[interface{}]interface{})
		GetConnManage() IConnManage

		SetConnectionCallback(ccnet.ConnectionCallback)
		SetMessageCallback(ccnet.MessageCallback)
		SetCloseCallback(ccnet.CloseCallback)
		SetTimerCallback(ccnet.TimerCallback)
	}
	IAgent interface {
		GetSchema() ISchema
		SetSchema(schema ISchema)

		GetServer() IServer
		SetServer(server IServer)
		SendMsg(errCode int32, cmd uint32, msg proto.Message)
	}

	// IConnector Connector
	IConnector interface {
		Connect() error
		DisConnect()

		GetConn() IConnection
		Send(cmd uint32, b []byte) error
	}

	// IHttp httpserver
	IHttp interface {
		Start() error
		Stop()
		SetHandlerFuncCallback(func(tvl, obj reflect.Value) gin.HandlerFunc)
		BindHandler(interface{})
	}
	ActionHandler func(msg IMessage) (int32, interface{})
	// IPlugin 用于热更逻辑的插件接口
	IPlugin interface {
		LoadPlugin() error
		SetActions(map[interface{}]interface{})
		GetAction(uint32) interface{}
	}
)

//-----------------
//-----db start----
//-----------------
type (
	IDB interface {
		CreateTable() error

		Create() (interface{}, error)
		Save() error
		Load() error
		FindOne() error
		UpdateProperty(key string, val interface{}) error
		UpdateProperties(properties map[string]interface{}) error

		SetUnique(key string) (string, error)
	}

	ISchema interface {
		Init()
		GetDB() IDB

		GetPri() interface{}
		GetSchema() interface{}
		GetSchemaName() string
		UpdateSchema(interface{})
		SetConn(conn IConnection)
		GetConn() IConnection

		Load() error
		Create() error
		Update()

		SetProperty(key string, val interface{})
		SetProperties(properties map[string]interface{})
		IncrProperty(key string, val int64) int64
		ParseFields(message protoreflect.Message, properties map[string]interface{}) []int32
	}
)

//-----------------
//-----db end------
//-----------------
