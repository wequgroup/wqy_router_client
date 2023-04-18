package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

const (
	Qos    = 0
	Server = "mqtt-hw.wequ.net:1883"
)

var rt = 0

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	cmd := ParseShellJSON(msg.Payload())
	switch cmd.ShellType {
	case 0:
		fmt.Println(cmd.ShellContent)
	default:
		fmt.Println("Unsupported Command!")
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connect success!")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
	fmt.Println("Trying to reconnect...")
	if rt > 10 {
		fmt.Println("Can not reconnect")
		os.Exit(-1)
	}
	client.Connect()
	rt = rt + 1
}

type WeQu interface {
	Join() WeQuMqtt
	Subscribe()
	Disconnect()
}

type LoginConfig struct {
	ClientID string
	Username string
	Password string
}

type WeQuMqtt struct {
	WQ       mqtt.Client
	Username string
}

func NewWeQu(clientID, userName, passWord string) LoginConfig {
	return LoginConfig{clientID, userName, passWord}
}

func (lc LoginConfig) Join() WeQuMqtt {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(Server)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetClientID(lc.ClientID)
	opts.SetUsername(lc.Username)
	opts.SetPassword(lc.Password)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(1 * time.Minute)
	fmt.Println("Starting MQTT connection...")
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(-1)
	}
	fmt.Println("Connected from MQTT server")
	return WeQuMqtt{client, lc.Username}
}

func (m WeQuMqtt) Ping() {
	for {
		token := m.WQ.Publish("0", 0, false, "1")
		if token.Wait() && token.Error() != nil {
			break
		}
		done := make(chan bool)
		go func() {
			time.Sleep(50 * time.Millisecond)
			done <- true
		}()
		<-done
	}
}

func (m WeQuMqtt) Subscribe() {
	token := m.WQ.Subscribe(fmt.Sprintf("duck/%s", m.Username), Qos, messageHandler)
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

func (m WeQuMqtt) Disconnect() {
	m.WQ.Disconnect(250)
}
