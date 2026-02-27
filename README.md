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

O projeto implementa o Tuple Space de duas formas distintas:

### 1. Versão Mutex (/mutex)
Utiliza os recursos clássicos de sistemas operacionais:

sync.Mutex para proteger o mapa de dados contra leituras e escritas simultâneas.

sync.Cond (Variáveis de Condição) para gerenciar o bloqueio das operações RD e IN sem consumir CPU (sem busy-wait), acordando as goroutines apenas quando um WR ocorre (Broadcast).

### 2. Versão Channels (/channels)
Utiliza a filosofia nativa do Go: "Do not communicate by sharing memory; instead, share memory by communicating".

Single-Threaded Loop: Apenas uma goroutine possui acesso direto ao mapa de dados, eliminando a necessidade de cadeados (Locks).

Filas de Espera (Waiters): Operações bloqueantes têm seus canais de resposta (replyCh) armazenados em uma fila. Quando o dado chega, ele é enviado diretamente pelo canal, destravando o cliente.

---

## Estrutura do Projeto

O repositório está organizado da seguinte forma:

```text
LINDA/
│
├── Go/                      # Código fonte do Servidor em Go
│   ├── Dockerfile           # Instruções unificadas com build-args para gerar as imagens
│   ├── go.mod               # Gerenciador de dependências e módulo do Go
│   │
│   ├── core/                # Lógica compartilhada do servidor
│   │   ├── protocol.go      # Interpretador (Parser) dos comandos do protocolo (WR, RD, etc.)
│   │   ├── server.go        # Lógica de conexão TCP e ciclo de vida dos clientes
│   │   └── services.go      # Tabela e implementação das funções do comando EX
│   │
│   ├── mutex/               # Implementação Clássica
│   │   ├── main.go          # Ponto de entrada da versão Mutex
│   │   └── tuplespace.go    # Estrutura de dados usando sync.Mutex e sync.Cond
│   │
│   └── channels/            # Implementação Idiomática Go
│       ├── main.go          # Ponto de entrada da versão Channels
│       └── tuplespace.go    # Estrutura de dados usando goroutines únicas e chan
│
└── Tester/                  # Scripts de teste automatizados em C++
    ├── prof_tests.cpp       # Casos de teste fornecidos pelo professor
    └── teste_bloqueio.cpp   # Script adicional para validar o bloqueio (concorrência)
```

---

### Inicie o Servidor

Primeiro, inicie o servidor (via Docker ou execução nativa). Ele deve permanecer executando durante toda a execução dos testes.

Exemplo com Docker:

No diretório `Go/`:
```bash
docker run -p 54321:54321 linda-go 54321
Servidor escutando na porta 54321
```

## Execução

O servidor recebe a **porta como argumento**.

### Porta utilizada nos testes:

Porta: `54321` (Localhost / `127.0.0.1`)

---

### Opção 1 — Execução via Docker (não precisa ter Go instalado)

É possível compilar qualquer uma das versões usando a variável `ARG VERSION`. No diretório `LINDA/Go/`, construa a imagem:

- Versão Mutex
```bash
docker build --build-arg VERSION=mutex -t linda-mutex .
docker run -p 54321:54321 linda-mutex 54321
```

- Versão Channels
```bash
docker build --build-arg VERSION=channels -t linda-channels .
docker run -p 54321:54321 linda-channels 54321
```

### Opção 2 — Execução Nativa (Linux ou WSL)

Requer Go 1.24 instalado. No diretório `Go/` execute:

```bash
# Executa a versão Mutex
go run ./mutex 54321

# OU executa a versão Channels
go run ./channels 54321
```

## Testando via TCP (Netcat)

Conecte ao servidor:

```bash
nc 127.0.0.1 54321
```

## Exemplos de interação

Abra dois ou três terminais. Sempre tente ler a partir do terminal que não escreveu a chave.

### Escrever uma tupla - Terminal 1
```bash
> WR chave_teste ola_mundo
< OK
```

### Ler sem remover -Terminal 2
```bash
> RD chave_teste
< OK ola_mundo
```

### Executar serviço - Terminal 2
```bash
> EX chave_teste chave_resultado 1
< OK
```

### Ler resultado - Terminal 1
```bash
> RD chave_resultado
< OK OLA_MUNDO 
```

### Consumir tupla - Terminal 1
```bash
> IN chave_resultado
< OK OLA_MUNDO
```

### Verificar consumo - Terminal 2
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

Abra um novo terminal, entre no diretório `Tester/`, compile o código C++ e execute:

```bash
cd Tester
g++ -std=c++17 prof_tests.cpp -o prof_tests
./prof_tests 127.0.0.1 54321
```

Se todas as mensagens aparecerem com [OK], o servidor está funcionando corretamente.

## Autora

Fernanda Petiz