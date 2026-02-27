package main

import (
	"linda/core"
	"sync"
)

// Modelo de tuplas
type Tuple struct {
	Key   string
	Value string
}

type TupleSpace struct {
	// várias tuplas com a mesma chave
	data map[string][]Tuple
	mu   sync.Mutex
	cond *sync.Cond
}

// fábrica - "construtor"
func NewTupleSpace() *TupleSpace {
	ts := &TupleSpace{
		data: make(map[string][]Tuple),
	}
	ts.cond = sync.NewCond(&ts.mu) // pra saber qual mutex liberar quando alguém chamar wait
	return ts
}

/*	====================================
	WR - pertence a TupleSpace
	------------------------------------
	lock -> cria tuple -> append no map
	-> unlock -> return ok
	====================================
*/

func (ts *TupleSpace) WR(key, value string) string {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	tuple := Tuple{Key: key, Value: value}
	ts.data[key] = append(ts.data[key], tuple)

	// ts.cond.Signal()
	ts.cond.Broadcast()

	return "OK"
}

/*
====================================
RD - pertence a TupleSpace
------------------------------------
lock -> se não tem tupla, wait -> se tem,
FIFO mantém -> unlock -> return ok valor
====================================
*/
func (ts *TupleSpace) RD(key string) string {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	for len(ts.data[key]) == 0 {
		ts.cond.Wait()
	}

	tuple := ts.data[key][0]

	return "OK " + tuple.Value
}

/*
====================================
IN - pertence a TupleSpace
------------------------------------
lock -> se não tem tupla, wait -> se tem,
FIFO remove -> unlock -> return ok valor
====================================
*/
func (ts *TupleSpace) IN(key string) string {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	for len(ts.data[key]) == 0 {
		ts.cond.Wait()
	}

	tuple := ts.data[key][0]
	ts.data[key] = ts.data[key][1:] // cria nova slice ts.data[key] do íntide 1 até o fim da ts.data[key] inicial

	return "OK " + tuple.Value
}

/*
====================================
EX - pertence a TupleSpace
------------------------------------
lock -> se não tem tupla, wait -> se tem,
FIFO remove -> se não tem serviço, NO-SERVICE
-> unlock -> se tem, aplica -> lock -> insere nova tupla ->
unlock -> return ok
====================================
*/
func (ts *TupleSpace) EX(key_in, key_out string, svc_id int) string {
	ts.mu.Lock()
	// defer ts.mu.Unlock()

	for len(ts.data[key_in]) == 0 {
		ts.cond.Wait()
	}

	tuple := ts.data[key_in][0]
	ts.data[key_in] = ts.data[key_in][1:]

	service, exists := core.Services[svc_id]
	if !exists {
		ts.mu.Unlock()
		return "NO-SERVICE"
	}
	ts.mu.Unlock()

	vOut := service(tuple.Value)

	ts.mu.Lock()
	defer ts.mu.Unlock()

	newTuple := Tuple{Key: key_out, Value: vOut}
	ts.data[key_out] = append(ts.data[key_out], newTuple)

	//ts.cond.Signal()
	ts.cond.Broadcast()

	return "OK"
}

// // receiver - se quiser fazer o método ser pertencente a um tipo específico
// func (receiver) nomeDaFuncao(parametro tipo) tipoRetorno {
//     // corpo da função
//     return valor
// }
