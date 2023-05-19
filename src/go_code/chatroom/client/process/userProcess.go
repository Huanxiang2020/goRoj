package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/project01/ch13/chatroom/client/utils"
	"go_code/project01/ch13/chatroom/common/message"
	"net"
	"os"
)

type UserProcess struct {
}

func (userPro *UserProcess) Login(userId int, userPwd string) (err error) {

	// 1.连接服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("Dial err=", err)
		return
	}
	defer conn.Close()

	// 2.通过conn向服务器发送消息
	var mes message.Message
	mes.Type = message.LoginMesType
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	// loginMes 序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("Marshal err=", err)
		return
	}
	mes.Data = string(data)
	// 4.mes 序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("Marshal err=", err)
		return
	}

	/*
		// 先发送长度,后发送本体
		pkgLen := []byte(strconv.Itoa(len(data)))
		var n int
		n, err = conn.Write(pkgLen)
		if n != len(pkgLen) || err != nil {
			fmt.Println("Write err=", err)
			return
		}
		//fmt.Printf("客户端发送的消息的长度=%d 内容=%s", len(data), string(data))
		// 发送消息本身
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Write err=", err)
			return
		}
	*/

	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)

	// 发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}

	fmt.Printf("客户端，发送消息的长度=%d 内容=%s", len(data), string(data))

	// 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}

	tf := &utils.Transfer{
		Conn: conn,
	}
	// 服务器返回的消息
	mes, err = tf.ReadPkg() //???????????????????????????????
	if err != nil {
		fmt.Println("ReadPkg err=", err)
		return err
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
		go serverProcessMes(conn)
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}

func (userPro *UserProcess) Register(userId int, userPwd, userName string) (err error) {
	// 连接服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("Dial err=", err)
		return
	}
	defer conn.Close()

	// 通过conn向服务器发送消息
	var mes message.Message
	mes.Type = message.RegisterMesType
	var registerMes message.RegisterMes

	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	data, err := json.Marshal(registerMes)
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
		return err
	}

	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	if registerResMes.Code == 200 {
		fmt.Println("注册成功,请重新登录")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}
