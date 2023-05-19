package main

import (
	"fmt"
	"go_code/project01/ch13/chatroom/client/process"
	"os"
)

var userId int
var userPwd string
var userName string

func main() {
	var key int
	for {
		fmt.Println("-----------------欢迎登录多人聊天系统-----------------")
		fmt.Println("                    1 登录聊天室")
		fmt.Println("                    2 注册用户")
		fmt.Println("                    3 退出系统")
		fmt.Println("请选择(1-3):")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入密码:")
			fmt.Scanf("%s\n", &userPwd)

			up := &process.UserProcess{}
			err := up.Login(userId, userPwd)
			if err != nil {
				fmt.Println("Login err=", err)
			}
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户昵称:")
			fmt.Scanf("%s\n", &userName)
			up := &process.UserProcess{}
			err := up.Register(userId, userPwd, userName)
			if err != nil {
				fmt.Println("Register err=", err)
			}
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("重新选择")
		}
	}
}
