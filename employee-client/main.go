package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	employee "../proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func LoadKeyPair() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair("../crt/cert.pem", "../crt/key.pem")
	if err != nil {
		panic(err)
	}

	ca, err := ioutil.ReadFile("../crt/cert.pem")
	if err != nil {
		panic(err)
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(ca) {
		panic("can't add CA cert")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      capool,
	}

	return credentials.NewTLS(tlsConfig)
}

func main() {
	conn, err := grpc.Dial("localhost:5002", grpc.WithTransportCredentials(LoadKeyPair()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := employee.NewEmployeeClient(conn)
	resp, err := client.GetManager(context.Background(), &employee.ManagerRequest{Name: "John"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
