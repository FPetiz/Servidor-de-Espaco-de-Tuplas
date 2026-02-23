package main

import (
	"bufio"
	"fmt"
	"net"
)

func StartServer(port string, ts *TupleSpace) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	fmt.Println("Servidor escutando na porta", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}

		// Goroutine para atender o cliente
		go handleConnection(conn, ts)
	}
}

/*
====================================
handleConnection
------------------------------------
cria reader -> lê da rede -> se erro,
encerra -> senão, HandleCommand ->
escreve resposta e envia de pela rede
-> fecha conexão com o cliente
====================================
*/
func handleConnection(conn net.Conn, ts *TupleSpace) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		response := HandleCommand(line, ts)

		conn.Write([]byte(response + "\n"))
	}
}
