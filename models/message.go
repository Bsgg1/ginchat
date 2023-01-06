package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	FromId   int64  //发送者
	TargetId int64  //接收者
	Type     int    //消息类型 群聊私聊广播
	Media    int    //消息类型 文字 图片 音频
	Content  string //消息内容
	Pic      string //图片
	Url      string
	Desc     string //描述
	Amount   int    //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, req *http.Request) {
	//获取参数并且校验token等
	//token:=query.Get("token")
	query := req.URL.Query()
	Id := query.Get("userID")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	//msgType := query.Get("type")
	//targetId := query.Get("targetID")
	//context := query.Get("context")
	isvalida := true //check token
	conn, err := (&websocket.Upgrader{
		//token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, req, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//获取连接
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe), //线程安全
	}
	//todo 用户关系
	//userid 跟node 绑定并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//todo 发送逻辑
	go sendProc(node)
	//todo 接收逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws]<<<<< ", data)
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 31, 190),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case data := <-udpsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1:
		sendMsg(msg.FromId, data) //私信
		/*case 2:sendGroupMsg()//群发
		case 3:sendAllMsg()//广播
		case 4:
		*/
	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.Unlock()
	if ok {
		node.DataQueue <- msg
	}
}
