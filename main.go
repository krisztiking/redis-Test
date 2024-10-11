package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func ToRespArray(command string) string {
	args := strings.Fields(command)
	length := len(args)

	resp := "*" + strconv.Itoa(length) + "\r\n"

	for _, arg := range args {
		resp += "$" + strconv.Itoa(len(arg)) + "\r\n" + arg + "\r\n"
	}

	return resp
}

func main() {
	// Kapcsolódás a szerverhez a 9090-es porton
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		fmt.Println("Hiba a kapcsolat létrehozásakor:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to the server...")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			break
		}
		command := scanner.Text()

		if command == "exit" {
			fmt.Println("Kilépés...")
			break
		}
		respString := ToRespArray(command)
		// Küldd el a parancsot a szervernek
		_, err = conn.Write([]byte(respString))
		if err != nil {
			fmt.Println("Error - sending: ", err)
			return
		}

		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ == "string" {
			fmt.Println(value.str)
		}

		if value.typ == "bulk" {
			fmt.Printf("\"%s\"\n", value.bulk)
		}

	}

}
