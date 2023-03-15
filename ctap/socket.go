package ctap

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type message struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	X    string `json:"X"`
	Y    string `json:"Y"`
}

var conn *websocket.Conn

func parseString(s string) ([]uint8, error) {
	var result []uint8
	parts := strings.Split(s, "x")
	for _, part := range parts[1:] { // Ignore first element since it's empty
		num, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return nil, err
		}
		result = append(result, uint8(num))
	}
	return result, nil
}

func requestAuth(allowed *string, rp string, wg *sync.WaitGroup) {
	defer wg.Done()

	socketUrl := "ws://127.0.0.1:4041" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	fmt.Println("Listening..")
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("auth"))
	if err != nil {
		log.Println("Error during writing to websocket:", err)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			break
		}
		msgStr := string(msg)

		fmt.Println(msgStr)

		*allowed = msgStr

		break

	}

}

func requestSignature(dataVar *[]byte, data []byte, wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println("Signature...")
	socketUrl := "ws://localhost:4041" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	var dataStr string
	// encode as string thing
	for _, elem := range data {
		dataStr += "x" + strconv.Itoa(int(elem))
	}

	print("encoded signature")

	err = conn.WriteMessage(websocket.TextMessage, []byte(dataStr))
	if err != nil {
		log.Println("Error during writing to websocket:", err)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			break
		}
		bytes, err := parseString(string(msg))

		if err != nil {
			log.Println("Error parsing message signature")
			break
		}

		*dataVar = bytes
		break
	}
}

func requestCert(certVar *[]byte, wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println("cert..")
	socketUrl := "ws://localhost:4041" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("cert"))
	if err != nil {
		log.Println("Error during writing to websocket:", err)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			break
		}
		bytes, err := parseString(string(msg))

		fmt.Println("CERT: ")
		fmt.Println(msg)

		if err != nil {
			log.Println("Error parsing message certificate")
			break
		}

		*certVar = bytes
		break

	}

}

func requestX(x *[]byte, wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println("getting X")

	socketUrl := "ws://localhost:4041" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("X_KEY"))
	if err != nil {
		log.Println("Error during writing to websocket:", err)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			break
		}
		bytes, err := parseString(string(msg))

		if err != nil {
			log.Println("Error parsing message X")
			break
		}

		*x = bytes
		break
	}

}

func requestY(y *[]byte, wg *sync.WaitGroup) {

	fmt.Println("getting Y")

	defer wg.Done()

	socketUrl := "ws://localhost:4041" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("Y_KEY"))
	if err != nil {
		log.Println("Error during writing to websocket:", err)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			break
		}
		bytes, err := parseString(string(msg))

		if err != nil {
			log.Println("Error parsing message Y")
			break
		}

		*y = bytes
		break
	}
}
