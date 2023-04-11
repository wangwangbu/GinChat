package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// 消息
type Message struct {
	gorm.Model
	FromId   uint   // 发送者
	TargetId uint   // 接收者
	Type     int    // 发送类型 群聊 私聊 广播
	Media    int    // 消息类型 文字 图片 音频
	Content  string // 消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int // 其他数据统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	// GroupSets set.Interface
}

// 映射关系
var clientMap map[uint]*Node = make(map[uint]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Chat....")
	// 1. 获取参数并校验合法性
	query := request.URL.Query()
	strUserId := query.Get("userId")
	intUserId, _ := strconv.Atoi(strUserId)
	userId := uint(intUserId)
	// strTarId := query.Get("targetId")
	// targetId, _ := strconv.ParseInt(strTarId, 10, 64)
	// context := query.Get("context")
	// msgType := query.Get("type")
	isValid := true // checkToken()
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValid
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 2. 获取 conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		// GroupSets: set.New(set.ThreadSafe),
	}
	// 3. 用户关系
	// 4. userID 和 node 绑定并且加锁
	rwLocker.RLock()
	clientMap[userId] = node
	rwLocker.RUnlock()
	// 5. 完成发送逻辑
	go sendProc(node)
	// 6. 完成接受逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天系统"))
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
		fmt.Println("[ws]<<<<<<", data)
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

// 完成 udp 数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(10, 249, 0, 255),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
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

// 完成 udp 数据接收协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
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

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		sendMsg(msg.TargetId, data)
	// case 2: // 群发
	// 	sendGroupMsg()
	// case 3: // 广播
	// 	sendAllMsg()
	}
}

func sendMsg(userId uint, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
