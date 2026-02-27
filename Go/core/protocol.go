package core

import (
	"strconv"
	"strings"
)

/*
====================================
HandleCommand
------------------------------------
limpa string -> divide por espaço ->
acha o comando e chama função certa
-> retorna resposta
====================================
*/
func HandleCommand(input string, ts TupleSpace) string {
	input = strings.TrimSpace(input)
	parts := strings.Fields(input)

	if len(parts) == 0 {
		return "ERROR"
	}

	switch parts[0] {

	case "WR":
		if len(parts) < 3 {
			return "ERROR"
		}
		valor := strings.Join(parts[2:], " ")
		//             chave   valor
		return ts.WR(parts[1], valor)

	case "RD":
		if len(parts) != 2 {
			return "ERROR"
		}
		//           chave
		return ts.RD(parts[1])

	case "IN":
		if len(parts) != 2 {
			return "ERROR"
		}
		//           chave
		return ts.IN(parts[1])

	case "EX":
		if len(parts) != 4 {
			return "ERROR"
		}

		svcID, err := strconv.Atoi(parts[3])
		if err != nil {
			return "ERROR"
		}
		//     chave_entrada chave_saida svc_id
		return ts.EX(parts[1], parts[2], svcID)

	default:
		return "ERROR"
	}
}
