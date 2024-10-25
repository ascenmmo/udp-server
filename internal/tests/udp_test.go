package tests

import (
	"context"
	"encoding/json"
	"fmt"
	tokengenerator "github.com/ascenmmo/token-generator/token_generator"
	tokentype "github.com/ascenmmo/token-generator/token_type"
	"github.com/ascenmmo/udp-server/env"
	"github.com/ascenmmo/udp-server/pkg/clients/udpGameServer"
	"github.com/ascenmmo/udp-server/pkg/restconnection/types"
	"github.com/ascenmmo/udp-server/pkg/start"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	clients = 20
	msgs    = 1000

	baseURl = "http://" + env.ServerAddress + ":" + env.TCPPort
	udpAddr = env.ServerAddress + ":" + env.UDPPort
	token   = env.TokenKey
)
var hash = ""

var ctx, cancel = context.WithCancel(context.Background())
var min, max time.Duration
var maxMsgs int

type Message struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

type Request struct {
	Token string  `json:"token,omitempty"`
	Data  Message `json:"data,omitempty"`
}

type Response struct {
	Data Message
}

func TestConnection(t *testing.T) {
	//logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger := zerolog.Logger{}

	go start.StartUDP(
		ctx,
		env.ServerAddress,
		env.TCPPort,
		env.UDPPort,
		env.TokenKey,
		msgs*clients,
		2,
		3,
		logger)
	time.Sleep(time.Second * 1)

	for i := 0; i < clients; i++ {
		createRoom(t, createToken(t, i))
		go Listener(t, i)
		go Publisher(t, i)
	}
	<-ctx.Done()

	fmt.Println(max, min, maxMsgs)
}

func Publisher(t *testing.T, i int) {
	connection := newConnection(t, i)
	for j := 0; j < msgs; j++ {
		if ctx.Err() != nil {
			return
		}
		msg := buildMessage(t, i, j)
		_, err := connection.Write(msg)
		assert.NoError(t, err)
	}
}

func Listener(t *testing.T, i int) {
	defer cancel()
	connection := newConnection(t, i)
	for j := 0; j < msgs; j++ {
		if ctx.Err() != nil {
			return
		}
		msg := buildMessage(t, i, j)
		_, err := connection.Write(msg)
		assert.NoError(t, err)
	}
	response := listen(t, connection)
	fmt.Println("done pubSub", i, "with msgs", response)

}

func createToken(t *testing.T, i int) string {
	z := 0
	if i > clients/2 {
		z = 1
	}
	gameID := uuid.NewMD5(uuid.UUID{}, []byte(strconv.Itoa(i)))
	roomID := uuid.NewMD5(uuid.UUID{}, []byte(strconv.Itoa(i)+strconv.Itoa(z)))
	userID := uuid.New()

	tokenGen, err := tokengenerator.NewTokenGenerator(token)
	assert.Nil(t, err, "init gen token expected nil")

	token, err := tokenGen.GenerateToken(tokentype.Info{
		GameID: gameID,
		RoomID: roomID,
		UserID: userID,
		TTL:    time.Second * 100,
	}, tokengenerator.AESGCM)
	assert.Nil(t, err, "gen token expected nil")

	return token
}

func createRoom(t *testing.T, token string) {
	hash = token
	cli := udpGameServer.New(baseURl)
	err := cli.ServerSettings().CreateRoom(context.Background(), token, types.CreateRoomRequest{})
	assert.Nil(t, err, "client.do expected nil")
}

func buildMessage(t *testing.T, i, j int) (msg []byte) {
	data := Message{
		Text:      fmt.Sprintf("отправила горутина:%d \t номер сообщения: %d %s %s %s", i, j, token, token, token),
		CreatedAt: time.Now(),
	}
	if j > msgs/2 {
		data = Message{
			Text:      "close",
			CreatedAt: time.Now(),
		}
	}

	req := Request{
		Token: createToken(t, i),
		Data:  data,
	}
	marshal, err := json.Marshal(req)
	assert.Nil(t, err, "client.do expected nil")
	return marshal
}

func newConnection(t *testing.T, i int) *net.UDPConn {
	serverAddr, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		fmt.Println("Ошибка разрешения адреса:", err, "горутина:", i)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		os.Exit(1)
	}
	assert.Nil(t, err)
	return conn
}

func listen(t *testing.T, conn *net.UDPConn) int {
	counter := 0
	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		assert.Nil(t, err, "ReadFromUDP expected nil")

		msg := buf[:n]
		var res Response
		err = json.Unmarshal(msg, &res)
		assert.Nil(t, err, "Unmarshal expected nil")
		counter++
		maxMsgs++
		if res.Data.Text == "close" {
			return counter
		}

		sub := time.Now().Sub(res.Data.CreatedAt)
		if min == 0 {
			min = sub
		}
		if min > sub {
			min = sub
		}

		if max < sub {
			max = sub
		}

	}
}
