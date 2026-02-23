# Servidor de Espaço de Tuplas (Tuple Space - Linda)

Implementação de um servidor concorrente de Espaço de Tuplas baseado no modelo **Linda**, desenvolvido em Go.

O servidor permite que múltiplos clientes se conectem via TCP e executem operações concorrentes de leitura, escrita, remoção e processamento de tuplas.

---

## Funcionalidades

O servidor suporta os seguintes comandos:

- `WR chave valor` → escreve uma tupla
- `RD chave` → lê sem remover (bloqueante)
- `IN chave` → remove e retorna (bloqueante)
- `EX chave_entrada chave_saida id_servico` → executa serviço sobre tupla

---

## Concorrência

O servidor utiliza:

- Goroutines para múltiplas conexões simultâneas
- Mutex + variáveis de condição para sincronização
- Bloqueio adequado nas operações `RD` e `IN`

Garantindo segurança em ambiente multi-thread.

---

## Estrutura do Projeto

O repositório está organizado da seguinte forma:

```text
LINDA/
│
├── Go/                      # Código fonte do Servidor em Go
|   ├── Dockerfile           # Instruções para construir a imagem do container
|   ├── go.mod               # Gerenciador de dependências e módulo do Go
|   ├── main.go              # Ponto de entrada (Entry point) do servidor
|   ├── initial_tests.go     # Testes locais iniciais de desenvolvimento
|   ├── protocol.go          # Interpretador (Parser) dos comandos do protocolo (WR, RD, etc.)
|   ├── server.go            # Lógica de conexão TCP e ciclo de vida dos clientes
|   ├── services.go          # Implementação das funções de transformação do comando EX
|   └── tuplespace.go        # Núcleo de concorrência: Estrutura de dados, Mutex e Cond
│
└── Tester/                  # Scripts de teste automatizados em C++
│   ├── prof_tests.cpp       # Casos de teste fornecidos pelo professor
│   └── teste_bloqueio.cpp   # Script adicional para validar o bloqueio (concorrência)
```

---

## Execução

O servidor recebe a **porta como argumento**.

Exemplo: `./linda 54321`

---

### Opção 1 — Execução via Docker (não precisa ter Go instalado)

No diretório onde está Dockerfile, construa a imagem:

```bash
docker build -t linda-go .
docker run -p 54321:54321 linda-go 54321
```

### Opção 2 — Execução Nativa (Linux ou WSL)

Requer Go instalado. No diretório contendo os arquivos `.go` execute:

```bash
go build -o linda .
./linda 54321
```

## Porta utilizada nos testes:

Porta: `54321` (Localhost / `127.0.0.1`)

## Testando via TCP (Netcat)

Conecte ao servidor:

```bash
nc 127.0.0.1 54321
```

## Exemplos de interação

### Escrever uma tupla
```bash
> WR chave_teste ola_mundo
< OK
```

### Ler sem remover
```bash
> RD chave_teste
< OK ola_mundo
```

### Executar serviço
```bash
> EX chave_teste chave_resultado 1
< OK
```

### Ler resultado
```bash
> RD chave_resultado
< OK OLA_MUNDO 
```

### Consumir tupla
```bash
> IN chave_resultado
< OK OLA_MUNDO
```

### Verificar consumo
```bash
RD chave_resultado
```

Resposta esperada:
O cliente ficará bloqueado até que outra conexão execute:
```bash
WR chave_resultado valor
```

Para sair:
`ctrl + c` 

## Serviços Disponíveis (EX)

|  ID  | Serviço    | Descrição                       | Função Interna |
| :--: | ---------- | ------------------------------- | ---------------|
|  1   | Maiúsculas | Converte texto para UPPERCASE   | `ToUpper()`    |
|  2   | Inversão   | Inverte a string                | `Reverse()`    |
|  3   | Tamanho    | Retorna o comprimento da string | `Length()`     |

Caso o serviço não exista, a resposta esperada é `NO-SERVICE`.

## Testes do diretório Tester

Os testes em C++ requerem ambiente Linux ou WSL devido ao uso de bibliotecas POSIX (arpa/inet.h, sys/socket.h).

### Inicie o Servidor

Primeiro, inicie o servidor (via Docker ou execução nativa). Ele deve permanecer em execução durante toda a execução dos testes.

Exemplo com Docker:

No diretório `Go/`:
```bash
docker run -p 54321:54321 linda-go 54321
Servidor escutando na porta 54321
```

### Execute os testes

Abra outro terminal, entre no diretório `Tester/`, compile o código e execute os testes:
```bash
cd Tester
g++ -std=c++17 prof_tests.cpp -o prof_tests
./prof_tests 127.0.0.1 54321
```

Se todas as mensagens aparecerem com [OK], o servidor está funcionando corretamente.

## Autora

Fernanda Petiz