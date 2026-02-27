package core

import (
	"strconv"
	"strings"
	//"time"
)

type ServiceFunc func(string) string

// Tabela de serviços: svc_id −→ funcao(string) −→ string
var Services = map[int]ServiceFunc{
	// converter a string para maiúsculas
	1: func(s string) string {
		//time.Sleep(20 * time.Second)
		return strings.ToUpper(s)
	},
	// inverter a string;
	2: func(s string) string {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	},
	// retornar o tamanho da string
	3: func(s string) string {
		return strconv.Itoa(len(s))
	},
}
