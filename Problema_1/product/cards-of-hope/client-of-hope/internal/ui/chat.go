// Pacote ui implementa a interface de usuário baseada em terminal para o Cards of Hope, utilizando Bubble Tea.
package ui

import (
	"client-of-hope/internal/state"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// logToFile registra mensagens de depuração no arquivo de log.
//
// Parâmetros:
//   - message: string a ser registrada no log.
func logToFile(message string) {
	f, err := os.OpenFile("data/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer f.Close()
	log.New(f, "ui-debug: ", log.LstdFlags).Println(message)
}

// chatMsg representa uma mensagem recebida do servidor ou lógica do app.
type chatMsg string

// errMsg encapsula um erro para ser tratado pela interface.
type errMsg struct{ err error }

// Chat é a estrutura principal da interface de usuário do Cards of Hope.
//
// Campos:
//   - Outputs: canal para mensagens enviadas da lógica para a UI.
//   - Inputs: canal para comandos/mensagens do usuário.
//   - program: instância do programa Bubble Tea.
//   - ctx, cancel: contexto para controle de execução.
//   - Done: canal para sinalizar o término da UI.
type Chat struct {
	Outputs chan string
	Inputs  chan string
	program *tea.Program
	ctx     context.Context
	cancel  context.CancelFunc
	Done    chan struct{} // Canal para sinalizar o término
}

// NewChat cria e inicializa uma nova instância de Chat, configurando os canais de comunicação.
//
// Retorno:
//   - *Chat: ponteiro para a nova instância de Chat.
func NewChat() *Chat {
	logToFile("NewChat called")
	return &Chat{
		Outputs: make(chan string, 100),
		Inputs:  make(chan string, 100),
		Done:    make(chan struct{}), // Inicializa o canal Done
	}
}

// Start inicializa e executa o programa Bubble Tea de forma não bloqueante.
//
// Efeitos colaterais:
//   - Inicia goroutines para escutar saídas e buscar mensagens periodicamente.
//   - Sinaliza término pelo canal Done.
func (c *Chat) Start() {
	logToFile("Chat.Start called")
	c.ctx, c.cancel = context.WithCancel(context.Background())
	m := newModel(c.Inputs, c.Outputs)
	c.program = tea.NewProgram(m, tea.WithAltScreen())

	go c.listenToOutputs()
	go c.periodicFetch()

	// Roda o programa bubbletea em sua própria goroutine para não bloquear a main.
	go func() {
		logToFile("Bubbletea program starting.")
		defer func() {
			logToFile("Closing Done channel.")
			close(c.Done) // Sinaliza para a main que a UI terminou
		}()

		if err := c.program.Start(); err != nil {
			logToFile(fmt.Sprintf("Bubbletea program exited with error: %v", err))
		}
		logToFile("Bubbletea program finished.")
	}()
}

// Close encerra o programa Bubble Tea e cancela o contexto associado.
func (c *Chat) Close() {
	logToFile("Chat.Close called")
	if c.cancel != nil {
		c.cancel()
	}
	if c.program != nil {
		c.program.Quit()
	}
}

func (c *Chat) Clear() {
	if c.program != nil {
		c.program.Send(clearHistoryMsg{}) // Envia a mensagem para limpar o histórico
	}
}

// listenToOutputs escuta o canal Outputs e envia mensagens para a interface Bubble Tea.
func (c *Chat) listenToOutputs() {
	logToFile("listenToOutputs goroutine started.")
	for {
		select {
		case <-c.ctx.Done():
			logToFile("listenToOutputs goroutine stopping.")
			return
		case msg := <-c.Outputs:
			logToFile(fmt.Sprintf("Received message from Outputs channel: '%s'", msg))
			if c.program != nil {
				c.program.Send(chatMsg(msg))
			}
		}
	}
}

// periodicFetch envia o comando /fetch periodicamente enquanto o usuário estiver em uma sala.
func (c *Chat) periodicFetch() {
	logToFile("periodicFetch goroutine started.")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if state.UserID != "" && state.RoomID != "" {
				c.Inputs <- "/fetch"
			}
		}
	}
}

// --- Bubble Tea Model ---

type model struct {
	// viewport é o componente de exibição de mensagens.
	viewport viewport.Model
	// textarea é a área de entrada de texto.
	textarea textarea.Model
	// inputsChan é o canal para comandos/mensagens do usuário.
	inputsChan chan<- string
	// outputsChan é o canal para mensagens recebidas.
	outputsChan <-chan string
	// history armazena o histórico de mensagens exibidas.
	history []string
}

// newModel cria e configura o modelo Bubble Tea para a interface de chat.
//
// Parâmetros:
//   - inputs: canal para comandos/mensagens do usuário.
//   - outputs: canal para mensagens recebidas.
//
// Retorno:
//   - model: instância configurada do modelo Bubble Tea.
func newModel(inputs chan<- string, outputs <-chan string) model {
	logToFile("newModel called")
	ta := textarea.New()
	ta.Placeholder = "Digite uma mensagem ou /help para ver os comandos..."
	ta.Focus()

	ta.Prompt = state.Username + ": "
	ta.CharLimit = 280
	ta.SetHeight(1)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(0, 0)

	return model{
		textarea:    ta,
		viewport:    vp,
		inputsChan:  inputs,
		outputsChan: outputs,
		history:     []string{},
	}
}

// Init inicializa o modelo Bubble Tea.
//
// Retorno:
//   - tea.Cmd: comando inicial do Bubble Tea.
func (m model) Init() tea.Cmd {
	logToFile("model.Init called")
	return textarea.Blink
}

type clearHistoryMsg struct{}

// Update processa eventos e atualiza o estado do modelo Bubble Tea.
//
// Parâmetros:
//   - msg: mensagem/evento recebido.
//
// Retorno:
//   - tea.Model: modelo atualizado.
//   - tea.Cmd: comando a ser executado.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tiCmd, vpCmd tea.Cmd

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			input := strings.TrimSpace(m.textarea.Value())
			if input != "" {
				logToFile(fmt.Sprintf("Sending to Inputs channel: '%s'", input))
				m.inputsChan <- input
				m.addHistory(fmt.Sprintf("%s: %s", state.Username, input))
			}
			m.textarea.Reset()
		}

	case chatMsg:
		logToFile(fmt.Sprintf("Update received chatMsg: '%s'", string(msg)))
		m.addHistory(string(msg))
		return m, nil

	case errMsg:
		logToFile(fmt.Sprintf("Update received errMsg: '%s'", msg.err.Error()))
		m.addHistory("Erro: " + msg.err.Error())
		return m, nil

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - m.textarea.Height()
		m.textarea.SetWidth(msg.Width)

	case clearHistoryMsg:
		m.clearHistory()
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

// addHistory adiciona uma mensagem ao histórico e atualiza o conteúdo do viewport.
//
// Parâmetros:
//   - msg: mensagem a ser adicionada ao histórico.
func (m *model) addHistory(msg string) {
	m.history = append(m.history, msg)
	m.viewport.SetContent(strings.Join(m.history, "\n"))
	m.viewport.GotoBottom()
}

func (m *model) clearHistory() {
	m.history = []string{}
	m.viewport.SetContent("")
}

// View retorna a representação textual da interface de chat.
//
// Retorno:
//   - string: interface formatada para exibição no terminal.
func (m model) View() string {
	if state.Username != "" {
		m.textarea.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(state.Username + ": ")
	} else {
		m.textarea.Prompt = ": "
	}
	return fmt.Sprintf("%s\n%s", m.viewport.View(), m.textarea.View())
}
