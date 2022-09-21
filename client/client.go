package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Client struct {}

func (t *Client) Update(data_sensor []string, res *bool) error {
	fmt.Println("Menerima data terbaru.")
	fmt.Println("|| Suhu\t\t:", data_sensor[0])
	fmt.Println("|| Kelembaban\t:", data_sensor[1])
	fmt.Println("|| Kadar CO2\t:", data_sensor[2])
	*res = true
	return nil
}

func main() {
	port := os.Args[1]
	client := &Client{}

	rpc.Register(client)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	handleError(err)
	
	fmt.Println("Client berjalan di port", port)
	http.Serve(listener, nil)
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())
	}
}
