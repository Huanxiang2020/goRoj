package test

import (
	"encoding/json"
	"fmt"
	"go_code/project01/ch13/chatroom/common/message"
	"go_code/project01/ch13/chatroom/server/utils"
	"net"
	"testing"
)

func TestServ(t *testing.T) {
	var conn net.Conn
	var err error
	// 连接服务器
	conn, err = net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("Dial err=", err)
		return
	}
	defer conn.Close()

	var userId int
	var userPwd string
	fmt.Println("id")
	fmt.Scanf("%d\n", &userId)
	fmt.Println("pwd")
	fmt.Scanf("%s\n", &userPwd)

	// 通过conn向服务器发送消息
	var mes message.Message
	mes.Type = message.LoginMesType
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	// loginMes 序列化
	var data []byte
	data, err = json.Marshal(loginMes)
	if err != nil {
		fmt.Println("Marshal err=", err)
		return
	}
	mes.Data = string(data)
	// mes 序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("Marshal err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("WritePkg err=", err)
	}

	// 服务器返回的消息
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("ReadPkg err=", err)
	}

	// 将mes的Data部分反序列化
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	if loginResMes.Code == 200 {
		//fmt.Println("登录成功")
		// 客户端启动一个协程,保持和服务器通讯,接收服务器数据推送,并显示在终端
		fmt.Println("成功")

	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}

func TestConn(t *testing.T) {
	fmt.Println("服务器在8889端口监听....")
	listen, err := net.Listen("tcp", "localhost:8889") // IP地址?
	defer listen.Close()
	if err != nil {
		fmt.Println("监听错误,err=", err)
		return
	}
}

func TestMySQL(t *testing.T) {

}
