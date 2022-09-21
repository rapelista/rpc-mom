package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var data = map[string][]float64{
	"suhu": {},
	"kelembaban": {},
	"karbon": {},
}

type DataSensor struct {
	A, B, C string
}

type RS struct {
}

func (t *RS) Add(data_sensor *DataSensor, result *bool) error {
	fmt.Println("Menerima data dari sensor.")
	handlePub(data_sensor)
	*result = true
	return nil
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Berhasil terhubung ke broker!")
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("Menerima data sensor dari broker.")
	data_sensor := string(msg.Payload())
	handleAddData(data_sensor)
}

var opts *mqtt.ClientOptions = mqtt.NewClientOptions()
var client_mqtt mqtt.Client
var token mqtt.Token
var client_rpc *rpc.Client
var err error

func main() {
	port := os.Args[1]
	port_client := os.Args[2]

	// mqtt
	opts.AddBroker("tcp://127.0.0.1:1883")
	opts.OnConnect = connectHandler
	opts.SetDefaultPublishHandler(messageHandler)
	client_mqtt = mqtt.NewClient(opts)
	token = client_mqtt.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println("Error", token.Error())
	}
	fmt.Println("Berhasil inisialisasi mqtt")

	// rpc as client (untuk mengirim data ke client)
	client_rpc, err = rpc.DialHTTP("tcp", fmt.Sprintf("127.0.0.1:%s", port_client))
	handleError(err)

	// rpc as server (untuk menerima data dari node sensor)
	rs := &RS{}
	rpc.Register(rs)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	handleError(err)

	handleSub()
	http.Serve(listener, nil)
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())
	}
}

func handleSub() {
	topic := "/sensor"
	client_mqtt.Subscribe(topic, 1, nil)
	fmt.Println("Berhasil subscribe di", topic)
}

func handlePub(data_sensor *DataSensor) {
	topic := "/sensor"
	message := fmt.Sprintf("%s|%s|%s", data_sensor.A, data_sensor.B, data_sensor.C)
	client_mqtt.Publish(topic, 1, false, message)
	fmt.Println("Publish data sensor ke broker.")
}

func handleAddData(data_sensor string) {
	data_sensor_splitted := strings.Split(data_sensor, "|")
	handleSendDataToClient(data_sensor_splitted)

	suhu, err := strconv.ParseFloat(data_sensor_splitted[0], 64);
	handleError(err)
	kelembaban, err := strconv.ParseFloat(data_sensor_splitted[1], 64);
	handleError(err)
	karbon, err := strconv.ParseFloat(data_sensor_splitted[2], 64);
	handleError(err)
	
	data["suhu"] = append(data["suhu"], suhu)
	data["kelembaban"] = append(data["kelembaban"], kelembaban)
	data["karbon"] = append(data["karbon"], karbon)

	fmt.Println("Update data sensor di server.")
}

func handleSendDataToClient(data_sensor []string) {
	fmt.Println("Mengirim data sensor ke client.")
	var res bool
	err = client_rpc.Call("Client.Update", data_sensor, &res)
	handleError(err)
}