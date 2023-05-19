package model

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"go_code/project01/ch13/chatroom/common/message"
)

type UserDao struct {
	pool *redis.Pool
}

var MyUserDao *UserDao

// 工厂模式,创建UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

func (u *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	var res string
	res, err = redis.String(conn.Do("HGet", "users", id)) // {"userId":123,"userPwd":"qwe","userName":"Tom"}
	if err != nil {
		if err == redis.ErrNil { // 未查到
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &User{}
	// res反序列化
	err = json.Unmarshal([]byte(res), user) // 反序列化失败   &{0 qwe Tom} userId为什么变成0?  反序列化时,无法调用userId
	if err != nil {
		fmt.Println("Unmarshal err=", err)
		return
	}
	return
}

// 用户登录校验
func (u *UserDao) LoginVerification(userId int, userPwd string) (user *User, err error) {
	conn := u.pool.Get()
	defer conn.Close()

	user, err = u.getUserById(conn, userId)
	if err != nil {
		return
	}
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (u *UserDao) Register(user *message.User) (err error) {
	conn := u.pool.Get()
	defer conn.Close()
	_, err = u.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXITS
		return
	}
	var data []byte
	data, err = json.Marshal(user)
	if err != nil {
		fmt.Println("Marshal err", err)
		return
	}

	// 存入数据库
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("数据库存储错误! err=", err)
		return
	}
	return
}
