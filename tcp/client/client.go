package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 15*time.Second, "maximum amount of time a dial will wait for a connect to complete. Default is 15s")
	flag.Parse()
	ip := flag.Args()[0]
	port := flag.Args()[1]

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	remoteAddr := ip
	if !strings.HasPrefix(port, ":") {
		remoteAddr += ":"
	}
	remoteAddr += port

	conn, err := d.DialContext(ctx, "tcp", remoteAddr)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	errChan := make(chan error, 2)

	go func() {
		errChan <- scanHandle(conn)
	}()

	go func() {
		errChan <- writeHandle(conn)
	}()

	if err := <-errChan; err != nil {
		log.Println(err)
	}
}

func scanHandle(conn net.Conn) (err error) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println("server:", text)
	}
	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("scanner error: %w", err)
	} else {
		log.Println("disconnected")
	}

	return err
}

func writeHandle(conn net.Conn) (err error) {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		text := scanner.Text()

		text += "\r\n"
		writer.WriteString(text)

		writer.Flush()
	}
	if err = scanner.Err(); err != nil || errors.Is(err, io.EOF) {
		err = fmt.Errorf("scanner error: %w", err)
	} else {
		log.Println("disconnected")
	}

	return err
}
