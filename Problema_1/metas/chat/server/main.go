package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	clientes  = make(map[net.Conn]bool)
	mensagens = make(chan string)
	mutex     = &sync.Mutex{}
)

func main() {
	ouvinte, erro := net.Listen("tcp", "localhost:8080")
	if erro != nil {
		log.Fatal("ERRO AO OUVIR: ", erro)
	}
	defer ouvinte.Close()
	fmt.Println("OUVINDO EM: localhost:8080")

	go broadcastMessages()

	for {
		conexao, erro := ouvinte.Accept()
		if erro != nil {
			log.Println("ERRO AO ACEITAR: ", erro)
			continue
		}

		go handleConnection(conexao)
	}
}

func broadcastMessages() {
	for {
		mensagem := <-mensagens

		mutex.Lock()
		for cliente := range clientes {
			_, erro := fmt.Fprintln(cliente, mensagem)
			if erro != nil {
				log.Println("ERRO AO ESCREVER PARA O CLIENTE: ", erro)
				cliente.Close()
				delete(clientes, cliente)
			}
		}
		mutex.Unlock()
	}
}

func handleConnection(conexao net.Conn) {
	log.Println("Novo cliente conectado:", conexao.RemoteAddr().String())

	mutex.Lock()
	clientes[conexao] = true
	mutex.Unlock()

	defer func() {
		log.Println("Cliente desconectado:", conexao.RemoteAddr().String())
		mutex.Lock()
		delete(clientes, conexao)
		mutex.Unlock()
		conexao.Close()
	}()

	leitor := bufio.NewReader(conexao)
	for {
		mensagem, erro := leitor.ReadString('\n')
		if erro != nil {
			return
		}
		mensagemCompleta := conexao.RemoteAddr().String() + ": " + mensagem
		mensagens <- mensagemCompleta
	}
}
