package main

import (
	"fmt"
	"go_code/project01/ch13/chatroom/server/model"
	"net"
	"time"
)

// 初始化函数
func init() {
	// 服务器启动时初始化连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	// 初始化UserDao
	model.MyUserDao = model.NewUserDao(pool)
}

// 处理客户端通讯
func process(conn net.Conn) {
	// 需要关闭conn
	defer conn.Close()
	processor := &Processor{Conn: conn}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误err=", err)
		return
	}
}

func main() {
	fmt.Println("服务器在8889端口监听....")
	listen, err := net.Listen("tcp", "0.0.0.0:8889") // IP地址?
	defer listen.Close()
	if err != nil {
		fmt.Println("监听错误,err=", err)
		return
	}

	for {
		fmt.Println("等待客户端连接服务器")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		// 连接成功,启动一个协程和客户端通讯
		go process(conn)
	}
}
