package main

import (
	"fmt"
	"go_code/project01/ch13/chatroom/common/message"
	"go_code/project01/ch13/chatroom/server/process"
	"go_code/project01/ch13/chatroom/server/utils"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// serverProcess 函数,根据客户端发送消息的种类,调用不同函数进行处理
func (pro *Processor) serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		// 处理登录
		up := &process2.UserProcess{Conn: pro.Conn}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		up := &process2.UserProcess{Conn: pro.Conn}
		err = up.ServerProcessRegister(mes)
	default:
		fmt.Println("消息有误")
	}
	return
}

func (pro *Processor) process2() (err error) {
	//读取客户端消息
	for {
		tf := &utils.Transfer{
			Conn: pro.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}
		fmt.Println("mes=", mes)
		err = pro.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}
