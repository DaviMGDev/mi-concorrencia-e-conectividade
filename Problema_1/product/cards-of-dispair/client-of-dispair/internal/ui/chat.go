package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Chat struct {
	Inputs   chan string
	Outputs  chan string
	reader   *bufio.Reader
}

func (chat *Chat) Read() {
	for {
		line, err := chat.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				close(chat.Inputs)
				return
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line != "" {
			chat.Inputs <- line
		}
	}
}

func (chat *Chat) Write() {
	for msg := range chat.Outputs {
		fmt.Printf("\r\033[K%s\n> ", msg)
	}
}

func NewChat() *Chat {
	return &Chat{
		Inputs:   make(chan string, 1),
		Outputs:  make(chan string, 1),
		reader:   bufio.NewReader(os.Stdin),
	}
}

func (chat *Chat) Start() {
	go chat.Read()
	go chat.Write()
	fmt.Print("> ")
	go func() {
		chat.Inputs <- "/get"
	}()
}
