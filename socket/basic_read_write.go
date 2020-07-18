package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8082")
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cliConn, err := net.Dial("tcp", "127.0.0.1:8082")
		if err != nil {
			panic(err)
		}
		var dataBuf bytes.Buffer
		b := make([]byte, 100)
		data := []byte("3\n")
		for {
			n, err := cliConn.Write(data)
			if err != nil {
				fmt.Printf("Write error: %v\n", err)
				break
			}
			fmt.Printf("Write %d bytes\n", n)
			if n == len(data) {
				break
			} else {
				data = data[n+1:]
			}
		}
		for {
			n, err := cliConn.Read(b)
			if err != nil {
				if err == io.EOF {
					fmt.Println("conn closed")
					cliConn.Close()
				} else {
					fmt.Printf("Read error: %v\n", err)
				}
				break
			}
			dataBuf.Write(b[:n])
			if strings.Contains(dataBuf.String(), "\n") {
				fmt.Printf("read resp %s", dataBuf.String())
				break
			}
		}
		wg.Done()
	}()

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(conn)
	data, err := reader.ReadBytes(byte('\n'))
	if err != nil {
		if err == io.EOF {
			fmt.Println("conn closed")
			conn.Close()
		} else {
			fmt.Printf("Read error: %v", err)
		}
	}
	fmt.Printf("Read bytes: %s\n", data)
	i, err := strconv.ParseInt(string(data[0:len(data)-1]), 10, 64)
	if err != nil {
		fmt.Printf("ParseInt error: %v\n", err)
	}
	resp := []byte(strconv.FormatInt(i*i*i, 10) + "\n")

	writer := bufio.NewWriter(conn)
	for {
		n, err := writer.Write(resp)
		if err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		fmt.Printf("Write %d bytes, len(resp)=%d\n", n, len(resp))
		if n == len(resp) {
			break
		} else {
			resp = resp[n+1:]
		}
	}
	writer.Flush()

	wg.Wait()
	fmt.Println("demo done")
}
