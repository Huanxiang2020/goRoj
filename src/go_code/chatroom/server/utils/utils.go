package utils

import (
	"encoding/json"
	"fmt"
	"go_code/project01/ch13/chatroom/common/message"
	"net"
	"strconv"
)

// Transfer 将方法关联到结构体
type Transfer struct {
	Conn net.Conn
	Buf  []byte // 传输时使用的缓冲
}

func (tran *Transfer) ReadPkg() (mes message.Message, err error) {
	tran.Buf = make([]byte, 8192)
	fmt.Println("读取客户端数据...")
	var n int
	// conn.Read只有在conn没有被关闭的情况下阻塞
	n, err = tran.Conn.Read(tran.Buf) // ??????????????????????
	if err != nil {
		// err = errors.New("readPkg header error")
		fmt.Println("read1 err=", err)
		return
	}

	pkgLen, _ := strconv.Atoi(string(tran.Buf[0:n])) // 83

	n, err = tran.Conn.Read(tran.Buf[:pkgLen])
	//n, err = conn.Read(buf)

	if n != pkgLen || err != nil {
		// err = errors.New("readPkg header error")
		fmt.Println("read2 err=", err)
	}

	err = json.Unmarshal(tran.Buf[:n], &mes) // ******注意传地址******
	if err != nil {
		fmt.Println("Unmarshal err=", err)
		return
	}
	return
}

func (tran *Transfer) WritePkg(data []byte) (err error) {
	// 先发送长度,后发送本体
	pkgLen := []byte(strconv.Itoa(len(data)))
	var n int
	n, err = tran.Conn.Write(pkgLen)
	if n != len(pkgLen) || err != nil {
		fmt.Println("Write err=", err)
		return
	}

	// 发送消息本身
	_, err = tran.Conn.Write(data)
	if err != nil {
		fmt.Println("Write err=", err)
		return
	}
	return
}
