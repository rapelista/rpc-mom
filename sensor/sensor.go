package main

import (
	"fmt"
	"net/rpc"
	"os"
)

type DataSensor struct {
	A, B, C string
}

var port string

func main() {
	var suhu, kelembaban, karbon string
	// suhu := os.Args[1]
	// kelembaban := os.Args[2]
	// karbon := os.Args[3]
	port = os.Args[1]

	fmt.Print("Masukkan suhu\t\t: ")
	fmt.Scanf("%s", &suhu)
	fmt.Print("Masukkan kelembaban\t: ")
	fmt.Scanf("%s", &kelembaban)
	fmt.Print("Masukkan kadar CO2\t: ")
	fmt.Scanf("%s", &karbon)

	fmt.Println(os.Args)
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("127.0.0.1:%s", port))
	handleError(err)
	data_sensor := &DataSensor{suhu, kelembaban, karbon}
	var res bool
	err = client.Call("RS.Add", data_sensor, &res)
	handleError(err)
	handleResult(res)
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())
	}
}

func handleResult(res bool) {
	if !res {
		fmt.Println("Gagal kirim data")
	} else {
		fmt.Println("Berhasil kirim data ke server port", port)
	}
}
