package tests

import (
	"context"
	"fmt"
	tokengenerator "github.com/ascenmmo/token-generator/token_generator"
	tokentype "github.com/ascenmmo/token-generator/token_type"
	"github.com/ascenmmo/udp-server/env"
	"github.com/ascenmmo/udp-server/pkg/api/types"
	"github.com/ascenmmo/udp-server/pkg/clients/udpGameServer"
	"github.com/ascenmmo/udp-server/pkg/start"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	clients = 20
	msgs    = 100

	baseURl = "http://" + env.ServerAddress + ":" + env.TCPPort
	udpAddr = env.ServerAddress + ":" + env.UDPPort
	token   = env.TokenKey
)

var ctx, cancel = context.WithCancel(context.Background())
var min, max time.Duration
var maxMsgs int

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
		logger,
		false)
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
	for j := 0; j < 10; j++ {
		msg := createToken(t, i)
		_, err := connection.Write([]byte(msg))
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 1)
	for j := 0; j < msgs; j++ {
		if ctx.Err() != nil {
			return
		}
		msg := buildMessage(t, i, j)

		_, err := connection.Write(msg)
		assert.NoError(t, err)

		msg = buildMessageWithTime(t)
		_, err = connection.Write(msg)
		assert.NoError(t, err)

		time.Sleep(time.Millisecond * 1)
	}
}

func Listener(t *testing.T, i int) {
	defer cancel()
	connection := newConnection(t, i)
	for j := 0; j < 10; j++ {
		if ctx.Err() != nil {
			return
		}
		msg := createToken(t, i)
		_, err := connection.Write([]byte(msg))
		assert.NoError(t, err)
	}
	response := listen(t, connection)
	fmt.Println("done pubSub", i, "with msgs", response)
	time.Sleep(time.Second * 5)
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
	}, tokengenerator.JWT)
	assert.Nil(t, err, "gen token expected nil")

	return token
}

func createRoom(t *testing.T, token string) {
	cli := udpGameServer.New(baseURl)
	err := cli.ServerSettings().CreateRoom(context.Background(), token, types.CreateRoomRequest{})
	assert.Nil(t, err, "client.do expected nil")
}

func buildMessage(t *testing.T, i, j int) (msg []byte) {
	msg = []byte(fmt.Sprintf("token:%sотправила горутина:%d \t номер сообщения: %d", createToken(t, i), i, j))
	if j > msgs/2 {
		msg = []byte("close")
	}
	return msg
}

func buildMessageWithTime(t *testing.T) (msg []byte) {
	msg = []byte(time.Now().Format(time.RFC3339Nano))
	return msg
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
		if len(msg) == 36 {
			continue
		}
		counter++
		maxMsgs++
		if string(msg) == "close" {
			return counter
		}
		if strings.Contains(string(msg), "close") {
			return counter
		}

		parse, err := time.Parse(time.RFC3339Nano, string(msg))
		if err != nil {
			continue
		}

		timeNow, err := time.Parse(time.RFC3339Nano, time.Now().UTC().Format(time.RFC3339Nano))
		if err != nil {
			continue
		}

		sub := timeNow.Sub(parse)

		if min == 0 {
			min = time.Duration(time.Now().Unix())
		}
		if min > sub {
			min = sub
		}

		if max < sub {
			max = sub
		}

	}
}
