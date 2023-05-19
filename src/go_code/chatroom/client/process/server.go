package process

import (
	"fmt"
	"go_code/project01/ch13/chatroom/client/utils"
	"net"
	"os"
)

// 显示登录成功后界面
func ShowMenu() {
	fmt.Println("-----------------用户登录成功-----------------")
	fmt.Println("              1. 显示在线用户列表")
	fmt.Println("              2. 发送消息")
	fmt.Println("              3. 信息列表")
	fmt.Println("              4. 退出系统")
	fmt.Println("请选择(1-4)")

	var key int
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		fmt.Println("显示在线用户列表")
	case 2:
		fmt.Println("发送消息")
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("退出系统")
		os.Exit(0)
	default:
		fmt.Println("输入有误!")
	}
}

// 和服务器保持通讯
func serverProcessMes(conn net.Conn) {
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务器发送的数据")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("ReadPkg err", err)
			return
		}
		fmt.Println("mes=", mes)
	}
}
