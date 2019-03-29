package main // package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"zhouwei/goprice/lib"

	"github.com/axgle/mahonia"
)

// Queue Queue
var Queue *lib.Queue

func getUserKey() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	userKey := fmt.Sprintf("%08v", rnd.Int31n(1000000))
	return userKey
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func convertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func sendMsg(conn net.Conn) {
	words := make(map[string]string)
	words["ApiName"] = "onRecvSpotQuotation"
	words["DataType"] = "TSpotQuotation"
	words["term_type"] = "05"
	words["user_key"] = getUserKey()
	words["user_type"] = "3"
	words["user_id"] = "8888"
	words["branch_id"] = "8888"
	words["lan_ip"] = "58.132.211.84"
	strs := make([]string, 0)
	for k, v := range words {
		strs = append(strs, k+"="+v)
	}
	mes := strings.Join(strs, "#")
	message := fmt.Sprintf("00000%d%s#", len(mes), mes)
	fmt.Println(message)
	data := convertToString(message, "utf-8", "gbk")
	n, errDial := conn.Write([]byte(data))
	checkError(errDial)
	fmt.Println("sended:", n)
}

func parseHeader(str string) int {
	result, _ := strconv.Atoi(str)
	return result
}

func main() {
	host := "117.141.138.101:41701"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	checkError(err)

	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	checkError(err)

	fmt.Println("connect success")
	sendMsg(conn)

	Queue := lib.QueueInstance()

	// 处理消息
	go ProcessMessage(Queue)

	for {
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err == io.EOF {
			conn.Close()
			break
		}
		Queue.Push(buf[:n])
	}
}

// ProcessMessage ProcessMessage
func ProcessMessage(Queue *lib.Queue) {
	msgLength := 8
	for {
		message := Queue.Pop(msgLength)
		if msgLength == 8 {
			header := parseHeader(string(message))
			msgLength = header
		} else {
			body := convertToString(string(message), "gbk", "utf-8")
			fmt.Println("-------------------")
			fmt.Println(body)
			msgLength = 8
		}
	}
}
