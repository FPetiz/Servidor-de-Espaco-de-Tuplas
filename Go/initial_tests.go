package main

import (
	"fmt"
	"sync"
	"time"
)

func RunDevTests() {
	fmt.Println("----------------------------------------")
	fmt.Println("Iniciando testes...\n")

	ts := NewTupleSpace()
	var wg sync.WaitGroup

	// =====================================================
	fmt.Println("Teste 1: RD BLOQUEANTE")
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		fmt.Println("RD aguardando chave A...")
		resp := ts.RD("A")
		fmt.Println("RD recebeu:", resp)
		fmt.Println("Tempo bloqueado:", time.Since(start))
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("WR inserindo (A, valor1)")
	ts.WR("A", "valor1")

	wg.Wait()
	fmt.Println()

	// =====================================================
	fmt.Println("Teste 2: IN FIFO")
	ts.WR("B", "primeiro")
	ts.WR("B", "segundo")

	fmt.Println("IN 1:", ts.IN("B"))
	fmt.Println("IN 2:", ts.IN("B"))
	fmt.Println()

	// =====================================================
	fmt.Println("Teste 3: EX COM SERVIÇO VÁLIDO")

	ts.WR("C", "ola")

	resp := ts.EX("C", "D", 1)
	fmt.Println("Resposta EX:", resp)
	fmt.Println("RD D:", ts.RD("D"))
	fmt.Println()

	// =====================================================
	fmt.Println("Teste 4: EX COM SERVIÇO INEXISTENTE")

	ts.WR("E", "teste")
	resp = ts.EX("E", "F", 999)
	fmt.Println("Resposta EX:", resp)
	fmt.Println("Tentando RD(F) (deve bloquear)...")

	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		fmt.Println(ts.RD("F"))
		fmt.Println("Tempo bloqueado:", time.Since(start))
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("Inserindo manualmente F para liberar...")
	ts.WR("F", "liberado")

	wg.Wait()
	fmt.Println()

	// =====================================================
	fmt.Println("Teste 5: EX BLOQUEANTE")

	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		fmt.Println("EX aguardando chave Z...")
		fmt.Println("Resposta EX:", ts.EX("Z", "Y", 1))
		fmt.Println("Tempo bloqueado:", time.Since(start))
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("WR inserindo (Z, teste)")
	ts.WR("Z", "teste")

	wg.Wait()

	fmt.Println("RD Y:", ts.RD("Y"))
	fmt.Println()

	// =====================================================
	fmt.Println("Teste 6: CONCORRÊNCIA MÚLTIPLA WR")

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ts.WR("X", fmt.Sprintf("valor-%d", id))
		}(i)
	}

	wg.Wait()

	for i := 0; i < 5; i++ {
		fmt.Println("IN X:", ts.IN("X"))
	}

	fmt.Println("\nTestes finalizados.")
	fmt.Println("--------------------------------------\n")
}
