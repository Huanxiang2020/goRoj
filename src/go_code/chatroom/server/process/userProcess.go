package process2

import (
	"encoding/json"
	"fmt"
	"go_code/project01/ch13/chatroom/common/message"
	"go_code/project01/ch13/chatroom/server/model"
	"go_code/project01/ch13/chatroom/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn
}

// ServerProcessLogin 处理登录请求
func (userPro *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	var resMes message.Message
	resMes.Type = message.LoginResMesType
	var loginResMes message.LoginResMes

	// 用户账密校验
	//var user *model.User
	user, err := model.MyUserDao.LoginVerification(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error() // "该用户不存在,请先注册..."
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "未知错误"
		}
	} else {
		loginResMes.Code = 200
		fmt.Println(user, "登录成功")
	}

	var data []byte
	data, err = json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 使用分层模式(mvc)
	tf := &utils.Transfer{
		Conn: userPro.Conn,
	}
	err = tf.WritePkg(data)
	return
}

func (userPro *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXITS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXITS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "未知错误类型"
		}
	} else {
		registerResMes.Code = 200
	}

	var data []byte
	data, err = json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: userPro.Conn,
	}
	err = tf.WritePkg(data)
	return

}
