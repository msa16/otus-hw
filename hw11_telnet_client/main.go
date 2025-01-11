package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	connectTimeout time.Duration
	serverAddress  string
)

func init() {
	flag.DurationVar(&connectTimeout, "timeout", time.Second*10, "connect to server timeout")
}

func parseArgs() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}
	serverAddress = net.JoinHostPort(args[0], args[1])
}

func main() {
	parseArgs()

	client := NewTelnetClient(serverAddress, connectTimeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", serverAddress)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// прием от сервера, вывод в консоль
		defer wg.Done()
		if client.Receive() == nil {
			fmt.Fprintf(os.Stderr, "...Connection was closed by peer\n")
		}
		stop()
	}()

	go func() {
		// ввод с консоли, отправка на сервер
		if client.Send() == nil {
			fmt.Fprintf(os.Stderr, "...EOF\n")
		}
		stop()
	}()

	// ждём одно из событий:
	//   - сигнал os.Interrupt (Ctrl-C)
	//   - окончание горутины приема данных от сервера по причине закрытия соединения или ошибки, будет вызвана stop()
	//   - окончания горутины отправки данных на сервер по причине EOF (Ctrl-D) или другой ошибки, будет вызвана stop()
	<-ctx.Done()
	// явно закрываем соединение, чтобы горутина приема данных от сервера завершилась, и мы ее корректно дождались
	client.Close()
	// ждём горутину приема данных
	wg.Wait()
	// простого способа корректного прерывания и ожидания завершения горутины отправки данных,
	// когда она находится внутри функции Read из os.Stdin, я не нашел
	// сейчас она завершится при завершении основной горутины
}
