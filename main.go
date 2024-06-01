package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	Device string
}

var (
	endpoint string
	device   string
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func connect(url *url.URL) mqtt.Client {
	host := url.Host
	username := url.User.Username()
	password, _ := url.User.Password()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", host))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	return mqtt.NewClient(opts)

}

func (cfg *Config) Sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 0, cfg.OnMessageReceived)
	token.Wait()
	log.Printf("Subscribed to topic %s", topic)
}

func (cfg *Config) OnMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", message.Payload(), message.Topic())
	display, err := strconv.Atoi(string(message.Payload()))
	if err != nil {
		log.Fatal(err)
	}
	switchDisplay(cfg.Device, display)
}

func switchDisplay(device string, display int) {
	log.Printf("device: %s | display: %d", device, display)
	path, err := exec.LookPath("ddccontrol")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("path is: ", path)
	_, err = exec.Command(path, "-r 0x60", fmt.Sprintf("-w %d", display), fmt.Sprintf("dev:%s", device)).Output()
	if err != nil {
		log.Fatal(err)
	}
	//log.Print("command result: ", cmd)
}

func main() {
	flag.StringVar(&endpoint, "url", "", "example: mqtt://<user>:<pass>@<server>.cloudmqtt.com:<port>/<topic>")
	flag.StringVar(&device, "device", "", "device path. example: /dev/i2c-2")

	flag.Parse()

	log.Println("endpoint: ", endpoint)
	log.Println("device: ", device)
	cfg := &Config{
		Device: device,
	}

	url, err := url.Parse(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	topic := url.Path[1:len(url.Path)]
	if topic == "" {
		topic = "test"
	}

	client := connect(url)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	cfg.Sub(client, topic)

	keepAlive := make(chan os.Signal)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)
	<-keepAlive

	log.Printf("unsub from topic: %s", topic)
	client.Unsubscribe(topic)
	log.Print("disconnecting from broker")
	client.Disconnect(250)
	log.Println("bye bye")

}
