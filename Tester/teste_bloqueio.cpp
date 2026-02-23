#include <arpa/inet.h>
#include <netdb.h>
#include <sys/socket.h>
#include <unistd.h>
#include <cstring>
#include <iostream>
#include <sstream>
#include <string>

// Envia uma linha (cmd + '\n') e lê uma linha de resposta
std::string send_command(int sock, const std::string &cmd) {
    std::string to_send = cmd + "\n";

    ssize_t total_sent = 0;

    while (total_sent < static_cast<ssize_t>(to_send.size())) {
        ssize_t sent = ::send(sock, to_send.data() + total_sent,
        to_send.size()- total_sent, 0);

        if (sent <= 0) {
            throw std::runtime_error("Erro ao enviar comando ao servidor");
        }
        total_sent += sent;
    }

    std::string response;
    char ch;

    while (true) {
        ssize_t rec = ::recv(sock, &ch, 1, 0);
        
        if (rec <= 0) {
            throw std::runtime_error("Conexao encerrada pelo servidor");
        }
        if (ch == '\n') {
            break;
        }
        if (ch != '\r') {
            response.push_back(ch);
        }
    }
    return response;
}

void expect_prefix(const std::string &resp, const std::string &prefix,
const std::string &context) {

    if (resp.rfind(prefix, 0) != 0) {
        std::cerr << "[FALHA] " << context
        << " resposta inesperada: \"" << resp << "\"\n";
    } else {
        std::cout << "[OK] " << context
        << " resposta: \"" << resp << "\"\n";
    }
}

int main(int argc, char *argv[]) {
    if (argc != 3) {
        std::cerr << "Uso: " << argv[0] << " <host> <porta>\n";
        std::cerr << "Exemplo: " << argv[0] << " 127.0.0.1 12345\n";
        return 1;
    }

    std::string host = argv[1];
    std::string port = argv[2];

    // Cria socket e conecta
    addrinfo hints{};
    hints.ai_family = AF_INET;
    hints.ai_socktype = SOCK_STREAM;
    addrinfo *result;
    int ret = ::getaddrinfo(host.c_str(), port.c_str(), &hints, &result);
    
    if (ret != 0) {
        std::cerr << "getaddrinfo: " << gai_strerror(ret) << "\n";
        return 1;
    }

    int sock =-1;
    for (addrinfo *rp = result; rp != nullptr; rp = rp->ai_next) {
        sock = ::socket(rp->ai_family, rp->ai_socktype, rp->ai_protocol);
        if (sock ==-1)
            continue;
        if (::connect(sock, rp->ai_addr, rp->ai_addrlen) == 0) {
            break; // conectou
        }
        ::close(sock);
        sock =-1;
    }
    freeaddrinfo(result);

    if (sock ==-1) {
        std::cerr << "Nao foi possivel conectar ao servidor\n";
        return 1;
    }

    try {
        std::cout << "Conectado a " << host << ":" << port << "\n";

        // Teste de EX com servico inexistente
        {
            std::string cmd = "WR in3 xyz";
            std::string resp = send_command(sock, cmd);
            expect_prefix(resp, "OK", "WR in3");
            cmd = "EX in3 out3 99"; // supondo que 99 nao exista
            resp = send_command(sock, cmd);
            expect_prefix(resp, "NO-SERVICE", "EX 99");
            // opcional: tentar RD out3, deve bloquear se implementado corretamente
            // entao NAO fazemos RD out3 aqui para evitar travar o tester
        }

        // Teste de bloqueio sugerido acima
        {
            std::cout << "Aguardando 5 segundos para ver se o servidor trava...\n";

            // Timeout de 5 segundos no socket para ele não travar pra sempre
            struct timeval tv;
            tv.tv_sec = 5;
            tv.tv_usec = 0;
            setsockopt(sock, SOL_SOCKET, SO_RCVTIMEO, (const char*)&tv, sizeof tv);

            std::string cmd = "RD out3";
            
            try {
                std::string resp = send_command(sock, cmd);
                
                // Se responder alguma coisa
                std::cerr << "[FALHA] O servidor nao bloqueou! Ele respondeu: \"" << resp << "\"\n";
                
            } catch (const std::exception &e) {
                // Se deu erro (timeout), recv() desistiu de esperar -> bloqueou 
                std::cout << "[OK] Servidor bloqueou a leitura. Timeout atingido!\n";
            }
        }

        std::cout << "Teste bloqueio concluido.\n";
    } catch (const std::exception &e) {
        std::cerr << "Erro: " << e.what() << "\n";
        ::close(sock);
        return 1;
    }
    ::close(sock);
    return 0;
}