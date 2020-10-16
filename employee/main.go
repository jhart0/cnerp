package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	employee_grpc "../proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

func main() {
	server := grpc.NewServer(grpc.Creds(LoadKeyPair()), grpc.UnaryInterceptor(interceptor))
	employeeServer := &EmployeeServer{}

	employee_grpc.RegisterEmployeeServer(server, employeeServer)

	go func() {
		s, err := net.Listen("tcp", ":5002")
		if err != nil {
			panic(err)
		}
		if err := server.Serve(s); err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	server.GracefulStop()

}

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if p, ok := peer.FromContext(ctx); ok {
		if _, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			log.Println("Valid TLS Cert")
		}
	}
	return handler(ctx, req)
}

type EmployeeServer struct {
	employee_grpc.UnimplementedEmployeeServer
}

func GetManagerName(empName string) string {
	return "Sanjay"
}

func (g *EmployeeServer) GetManager(ctx context.Context, req *employee_grpc.ManagerRequest) (*employee_grpc.ManagerReply, error) {
	empName := req.GetName()
	respdata := "Manager is: " + GetManagerName(empName)
	return &employee_grpc.ManagerReply{Message: respdata}, nil
}

func LoadKeyPair() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair("../crt/cert.pem", "../crt/key.pem")
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile("../crt/cert.pem")
	if err != nil {
		panic(err)
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(data) {
		panic(err)
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    capool,
	}
	return credentials.NewTLS(tlsConfig)
}
