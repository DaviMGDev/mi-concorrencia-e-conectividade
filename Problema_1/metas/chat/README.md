# Chat Simples em Go com Docker

Este projeto implementa um chat simples cliente-servidor em Go, utilizando Docker para a conteinerização das aplicações.

## Pré-requisitos

- Docker

## Como usar

O projeto utiliza um `Makefile` para simplificar o processo de build e execução dos contêineres.

### 1. Iniciar o servidor

Em um terminal, execute o seguinte comando para construir a imagem do servidor e iniciar o contêiner:

```bash
make server
```

O servidor estará ouvindo na porta `8080`.

### 2. Iniciar o cliente

Abra um novo terminal e execute o seguinte comando para construir a imagem do cliente e iniciar o contêiner:

```bash
make client
```

Você pode iniciar múltiplos clientes, cada um em seu próprio terminal. As mensagens enviadas por um cliente serão transmitidas para todos os outros clientes conectados.
