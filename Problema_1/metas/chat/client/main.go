package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conexao, erro := net.Dial("tcp", "localhost:8080")
	if erro != nil {
		log.Fatal("NÃO FOI POSSÍVEL CONECTAR AO SERVIDOR: ", erro)
	}
	defer conexao.Close()

	go func() {
		io.Copy(os.Stdout, conexao)
	}()

	io.Copy(conexao, os.Stdin)
}
