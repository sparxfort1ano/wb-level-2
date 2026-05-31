package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			if err := handleConnection(conn); err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

func handleConnection(conn net.Conn) (err error) {
	defer func() {
		err = conn.Close()
	}()

	name := conn.RemoteAddr().String()

	log.Printf("%+v connected\n", name)
	conn.Write([]byte("Hello, " + name + "\n"))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			conn.Write([]byte("bye\n"))
			break
		} else if text != "" {
			log.Println(name, "client:", text)
			conn.Write([]byte("you enter " + text + "\n"))
		}
	}
	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("scanner error: %w", err)
	} else {
		log.Println(name, "disconnected")
	}

	return err
}
