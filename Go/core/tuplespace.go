package core

// Interface para poder usar a mesma estrtutura nas duas implementações
type TupleSpace interface {
	WR(key, value string) string
	RD(key string) string
	IN(key string) string
	EX(keyIn, keyOut string, serviceID int) string
}
