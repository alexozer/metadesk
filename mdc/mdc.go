package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/alexozer/metadesk/server"
)

func main() {
	conn, err := net.Dial("unix", server.SockAddr)
	if err != nil {
		fail("No metadesk server found")
	}

	if err = writeArgs(conn); err != nil {
		fail(err.Error())
	}

	exitCode, err := printResponse(conn)
	if err != nil {
		fail(err.Error())
	}

	os.Exit(exitCode)
}

func writeArgs(conn net.Conn) error {
	writer := bufio.NewWriter(conn)
	args := os.Args[1:]

	_, err := writer.WriteString(fmt.Sprintf("%d\n", len(args)))
	if err != nil {
		return errors.New("Failed to write argument count to server")
	}

	for _, arg := range args {
		_, err := writer.WriteString(arg + "\n")
		if err != nil {
			return errors.New("Failed to write argument to server")
		}
	}

	if writer.Flush() != nil {
		return errors.New("Failed to flush arguments to server")
	}

	return nil
}

func printResponse(conn net.Conn) (exitCode int, err error) {
	reader := bufio.NewReader(conn)
	exitCodeStr, err := reader.ReadString('\n')
	if err != nil {
		err = errors.New("Failed to read exit code from server")
		return
	}

	exitCode, err = strconv.Atoi(exitCodeStr[:len(exitCodeStr)-1])
	if err != nil {
		err = errors.New("Non-integer exit code received from server")
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // connection closed
		}

		fmt.Print(line)
	}

	return
}

func fail(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
