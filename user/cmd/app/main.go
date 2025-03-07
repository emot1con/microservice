package main

import (
	"sync"

	"user_service/transport/grpc"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Starting Application")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		logrus.Info("Calling GRPCListen()")
		grpc.GRPCListen()
		logrus.Info("GRPCListen() exited")
	}()

	wg.Wait()
}
