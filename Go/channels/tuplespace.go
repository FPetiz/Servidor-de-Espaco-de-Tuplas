package main

import "linda/core"

type request struct {
	op      string
	keyIn   string
	keyOut  string
	value   string
	svcID   int
	replyCh chan string
}

type TupleSpace struct {
	reqChan chan request
}

func NewTupleSpace() *TupleSpace {
	ts := &TupleSpace{
		reqChan: make(chan request),
	}
	go ts.loop()
	return ts
}

func (ts *TupleSpace) loop() {
	// Apenas esta goroutine toca nestas variáveis
	data := make(map[string][]string)
	waiters := make(map[string][]request) // Filas de espera

	for req := range ts.reqChan { // range em um canal pausa a execução da goroutine se o canal estiver vazio
		switch req.op {

		case "WR":
			k := req.keyIn
			v := req.value

			// Se houver alguém esperando por esta chave, atende o primeiro da fila
			if len(waiters[k]) > 0 {
				pending := waiters[k][0]
				waiters[k] = waiters[k][1:]

				// Processa a requisição pendente com o valor que acabou de chegar
				ts.processRequest(pending, v, data)
			} else {
				// Ninguém esperando, guarda no mapa
				data[k] = append(data[k], v)
			}
			req.replyCh <- "OK"

		case "RD", "IN", "EX":
			k := req.keyIn

			if len(data[k]) > 0 {
				val := data[k][0]
				if req.op != "RD" {
					data[k] = data[k][1:]
				}
				ts.processRequest(req, val, data)
			} else {
				// Se não existe, coloca na fila de espera
				waiters[k] = append(waiters[k], req)
			}
		}
	}
}

// processRequest auxilia o loop a responder ou disparar serviços
func (ts *TupleSpace) processRequest(req request, val string, data map[string][]string) {
	switch req.op {
	case "RD", "IN":
		req.replyCh <- "OK " + val

	case "EX":
		service, exists := core.Services[req.svcID]
		if !exists {
			req.replyCh <- "NO-SERVICE"
			return
		}
		// Executa o serviço em background para não travar o loop principal
		go func() {
			vOut := service(val)
			// Envia um novo WR para o sistema inserir o resultado
			ts.WR(req.keyOut, vOut)
			req.replyCh <- "OK"
		}()
	}
}

// --- Métodos Públicos (Interface) ---

func (ts *TupleSpace) WR(key, value string) string {
	reply := make(chan string)
	ts.reqChan <- request{op: "WR", keyIn: key, value: value, replyCh: reply}
	return <-reply
}

func (ts *TupleSpace) RD(key string) string {
	reply := make(chan string)
	ts.reqChan <- request{op: "RD", keyIn: key, replyCh: reply}
	return <-reply
}

func (ts *TupleSpace) IN(key string) string {
	reply := make(chan string)
	ts.reqChan <- request{op: "IN", keyIn: key, replyCh: reply}
	return <-reply
}

func (ts *TupleSpace) EX(keyIn, keyOut string, svcID int) string {
	reply := make(chan string)
	ts.reqChan <- request{op: "EX", keyIn: keyIn, keyOut: keyOut, svcID: svcID, replyCh: reply}
	return <-reply
}
