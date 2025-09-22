# Cards of Hope – Jogo de Cartas Multiplayer (TEC502)

Este repositório apresenta a solução completa para o Problema 1 da disciplina TEC502 – Concorrência e Conectividade (UEFS), superando os requisitos do barema com uma arquitetura robusta, escalável e documentada.

## 📝 Descrição Geral

O Cards of Hope é um jogo de cartas online multiplayer, centrado em duelos táticos 1v1 que combinam estratégia, blefe e gerenciamento de recursos. Inspirado no clássico pedra-papel-tesoura, cada jogador possui cartas de três tipos (rock, paper, scissors), mas cada carta também possui um valor de estrelas (★) que determina sua força dentro do tipo.

**Regras e Mecânicas do Jogo:**

- Cada jogador começa com um conjunto de cartas, podendo adquirir mais através do comando `/buy`, que entrega um pacote aleatório de cartas de cada tipo.
- As partidas ocorrem em salas privadas, criadas e acessadas pelos próprios jogadores.
- Em cada rodada, ambos os jogadores escolhem secretamente uma carta de sua mão para jogar.
- O vencedor da rodada é determinado primeiro pelo tipo (rock > scissors > paper > rock), e, em caso de empate de tipo, vence quem tiver a carta com mais estrelas. Se ambos jogarem o mesmo tipo e valor de estrelas, a rodada empata.
- O sistema de estoque global garante que cada carta só possa ser adquirida por um jogador, promovendo justiça e competição pelo recurso.

O sistema é composto por três componentes principais: servidor centralizado, cliente interativo e cliente de estresse, todos containerizados para máxima portabilidade e reprodutibilidade.


## 📐 Arquitetura

A arquitetura do Cards of Hope foi cuidadosamente modularizada para garantir clareza, manutenibilidade e separação de responsabilidades, tanto no servidor quanto no cliente.

### Cliente (`client-of-hope`)

- **`cmd/app/`**: Ponto de entrada da aplicação cliente.
- **`internal/api/`**: Implementa o cliente TCP, protocolos de requisição/resposta e handlers para comandos (autenticação, chat, jogo, etc).
- **`internal/application/`**: Gerencia o roteamento de comandos do usuário para os handlers apropriados.
- **`internal/state/`**: Mantém o estado local do cliente, como cartas, jogadas e informações do usuário.
- **`internal/ui/`**: Interface de usuário baseada em terminal, utilizando Bubble Tea para uma experiência interativa e responsiva.
- **`internal/utils/`**: Utilitários genéricos para manipulação de dados e estruturas auxiliares.
- **`data/`**: Armazena arquivos de log gerados pelo cliente.

### Servidor (`server-of-hope`)

- **`cmd/app/`**: Ponto de entrada da aplicação servidor.
- **`internal/api/`**: Implementa o servidor TCP, roteamento de comandos, protocolos de comunicação e handlers para autenticação, chat, jogo, etc.
- **`internal/application/`**: Lógica de alto nível para autenticação, gerenciamento de salas, partidas e armazenamento.
- **`internal/data/`**: Repositórios e abstrações para persistência e manipulação de dados do jogo.
- **`internal/domain/`**: Define as entidades centrais do domínio, como cartas, usuários, salas e jogos.
- **`internal/state/`**: Gerencia o estado global do servidor, serviços, ambiente e logging.
- **`internal/utils/`**: Utilitários para manipulação de dados e estruturas auxiliares.

Essa organização permite que cada componente do sistema seja desenvolvido, testado e evoluído de forma independente, promovendo alta coesão e baixo acoplamento. O uso de pacotes internos reforça o encapsulamento e a segurança do código.


## 🔗 Comunicação: Estrutura dos Protocolos

Toda comunicação entre cliente e servidor é feita via TCP, utilizando mensagens JSON com a seguinte estrutura:

### Estruturas Gerais

**Request (cliente → servidor):**
```json
{
    "method": "<nome_do_comando>",
    "data": { ... }
}
```

**Response (servidor → cliente):**
```json
{
    "method": "<nome_do_comando>",
    "status": "<ok|error>",
    "data": { ... }
}
```

---

### Exemplos de Comandos

#### 1. PING
- **REQUEST:**
    ```json
    {
        "method": "ping",
        "data": { "user_id": "<id_do_usuario>" }
    }
    ```
- **RESPONSE:**
    ```json
    {
        "method": "ping",
        "status": "ok",
        "data": {}
    }
    ```

#### 2. REGISTER
- **REQUEST:**
    ```json
    {
        "method": "register",
        "data": { "username": "<usuario>", "password": "<senha>" }
    }
    ```
- **RESPONSE:**
    ```json
    {
        "method": "register",
        "status": "ok",
        "data": { "message": "User registered successfully", "username": "<usuario>" }
    }
    ```
    (Em caso de erro, `status: "error"` e mensagem explicativa em `data.message`.)

#### 3. LOGIN
- **REQUEST:**
    ```json
    {
        "method": "login",
        "data": { "username": "<usuario>", "password": "<senha>" }
    }
    ```
- **RESPONSE:**
    ```json
    {
        "method": "login",
        "status": "ok",
        "data": { "user_id": "<id_do_usuario>" }
    }
    ```

#### 4. CRIAR SALA
- **REQUEST:**
    ```json
    {
        "method": "create",
        "data": { "user_id": "<id_do_usuario>" }
    }
    ```
- **RESPONSE:**
    ```json
    {
        "method": "create",
        "status": "ok",
        "data": { "room_id": "<id_da_sala>" }
    }
    ```

#### 5. ENTRAR EM SALA
- **REQUEST:**
    ```json
    {
        "method": "join",
        "data": { "user_id": "<id_do_usuario>", "room_id": "<id_da_sala>" }
    }
    ```
- **RESPONSE:**
    ```json
    {
        "method": "join",
        "status": "ok",
        "data": {}
    }
    ```

#### 6. SAIR DA SALA
- **REQUEST:**
    ```json
    {
        "method": "leave",
        "data": { "user_id": "<id_do_usuario>", "room_id": "<id_da_sala>" }
    }
    ```
- **RESPONSE:**
    ```json
    {
        "method": "leave",
        "status": "ok",
        "data": {}
    }
    ```

#### 7. ENVIAR MENSAGEM (CHAT)
- **REQUEST:**
    ```json
    {
        "method": "send",
        "data": { "user_id": "<id_do_usuario>", "room_id": "<id_da_sala>", "message": "<texto>" }
    }
    ```
- **RESPONSE:**
    ```json
    {
        "method": "send",
        "status": "ok",
        "data": {}
    }
    ```

#### 8. COMPRAR PACOTE DE CARTAS
- **REQUEST:**
    ```json
    {
        "method": "buy",
        "data": null
    }
    ```
- **RESPONSE:**
    ```json
    {
        "method": "buy",
        "status": "ok",
        "data": {
            "package": {
                "rock": <int>,
                "paper": <int>,
                "scissors": <int>
            }
        }
    }
    ```

---

## 🛡️ API Remota & Encapsulamento

- Todas as interações (login, registro, chat, compra de pacotes, jogada, etc.) são comandos explícitos, documentados e validados.
- Dados encapsulados em structs Go, serializados/deserializados via JSON.
- Validação rigorosa de entrada/saída e tratamento de erros para garantir integridade e segurança.

## ⚡ Concorrência & Desempenho

- O servidor emprega Goroutines para cada conexão, com sincronização via `sync.Mutex`, `sync.Map` e pools para recursos críticos.
- Worker pools otimizam tarefas pesadas (ex: compra de pacotes), evitando gargalos e garantindo justiça.
- Testes de estresse automatizados comprovam a escalabilidade e ausência de race conditions.

## ⏱️ Latência & Responsividade

- Comando `/ping` disponível a qualquer momento para medir latência real entre cliente e servidor.
- Estrutura de mensagens e lógica de processamento minimizam delays, mesmo sob alta carga.


## 🥇 Partidas & Pareamento

- O sistema permite que os próprios jogadores criem e entrem manualmente em salas para disputar partidas 1v1.
- Cada jogador só pode estar em uma sala por vez, garantindo que não haja múltiplos pareamentos simultâneos.
- O isolamento entre partidas é garantido pela separação lógica das salas, evitando interferência entre jogos distintos.

## 🎴 Pacotes & Estoque Global

- Mecânica de compra de pacotes implementada como "estoque" global, protegido por locks para garantir atomicidade.
- Distribuição justa: cada carta só pode ser adquirida por um jogador, mesmo sob concorrência extrema.

## 🧪 Testes & Emulação

- Cliente de estresse simula milhares de jogadores, validando justiça, desempenho e ausência de falhas.
- Todos os componentes são executados em contêineres Docker, permitindo testes reprodutíveis e escaláveis em qualquer ambiente.

## 🛠️ Tecnologias Utilizadas

- **Linguagem:** Go 1.24+
- **Concorrência:** Goroutines, Channels, Mutex, Pools
- **Comunicação:** Sockets TCP nativos
- **Serialização:** JSON
- **Containerização:** Docker & Docker Compose

## 🚀 Como Executar o Projeto

> **Importante:**
>
> Se você for rodar servidor e cliente em máquinas diferentes, ou em containers independentes, defina a variável de ambiente `SERVER_ADDR` no cliente com o IP real ou hostname da máquina onde o servidor está rodando. Exemplo:
>
> ```bash
> docker-compose run --rm -e SERVER_ADDR=10.0.1.2:8080 client-of-hope
> ```
>
> O IP pode ser obtido com o comando `hostname --all-ip-addresses` na máquina do servidor. Não use nomes de serviço Docker (ex: `server-of-hope`) se não estiver usando a mesma rede Docker.

### Pré-requisitos
- Git
- Docker (>= 20.10)
- Docker Compose

### Passos para Execução

1. **Clone o repositório:**
    ```bash
    git clone <URL_DO_REPOSITORIO>
    cd cards-of-hope
    ```

2. **Build das imagens:**
    ```bash
    docker-compose build
    ```

---

### Como rodar o servidor

Em um terminal, execute:
```bash
docker-compose run --rm -p 8080:8080 server-of-hope
```
Isso irá iniciar o servidor ouvindo na porta 8080 do host.

---

### Como rodar o cliente

Em outro terminal, execute (ajuste o IP para o endereço do servidor):
```bash
docker-compose run --rm -e SERVER_ADDR=<IP_DO_SERVIDOR>:8080 client-of-hope
```
Exemplo prático:
```bash
docker-compose run --rm -e SERVER_ADDR=10.0.1.2:8080 client-of-hope
```
Repita em outros terminais para múltiplos jogadores.

---

### Como rodar o cliente de estresse

Em outro terminal, execute (ajuste o IP para o endereço do servidor):

```bash
docker-compose run --rm stress-client-of-hope -addr <IP_DO_SERVIDOR>:8080
```

> **Importante:**
>
> O stress-client **não utiliza** a variável de ambiente `SERVER_ADDR`. Sempre informe o endereço do servidor usando o argumento `-addr`.

Exemplo prático:
```bash
docker-compose run --rm stress-client-of-hope -addr 10.0.1.2:8080 -clients 10 -onlyconn
```

#### Argumentos e Modos de Execução do Cliente de Estresse

O cliente de estresse (`stress-client-of-hope`) aceita argumentos via linha de comando para personalizar o teste. Os principais argumentos são:

- `-addr` — Endereço do servidor no formato `host:porta` (padrão: `localhost:8080`)
- `-clients` — Número de conexões simultâneas (padrão: `100`)
- `-interval` — Intervalo entre pings em milissegundos (padrão: `100`)
- `-duration` — Duração do teste em segundos (padrão: `10`)
- `-onlyconn` — Se definido, testa apenas o limite de conexões simultâneas, sem enviar comandos (padrão: `false`)

### Comandos do Jogo

- `/register <usuario> <senha>` – Registrar novo usuário
- `/login <usuario> <senha>` – Fazer login
- `/logout` – Fazer logout da sessão atual
- `/create <nome_da_sala>` – Criar uma nova sala de chat/jogo
- `/join <nome_da_sala>` – Entrar em uma sala existente
- `/leave` – Sair da sala atual
- `/send <mensagem>` – Enviar mensagem para a sala atual (ou apenas digite a mensagem sem `/`)
- `/play <carta>` – Jogar uma carta (`rock`, `paper` ou `scissors`)
- `/cards` – Mostrar suas cartas atuais
- `/buy` – Comprar um novo pacote de cartas
- `/whoami` – Exibir informações do usuário logado
- `/whereami` – Exibir a sala em que você está
- `/ping` – Verificar a conexão com o servidor
- `/clear` – Limpar o histórico de mensagens do chat
- `/help` – Exibir a lista de comandos disponíveis
- `/quit` – Sair do aplicativo

---

Este projeto foi desenvolvido como solução individual para a disciplina TEC502 – Concorrência e Conectividade (UEFS), superando os requisitos mínimos e demonstrando domínio prático de sistemas concorrentes, comunicação em rede e engenharia de software moderna.
