
// Pacote utils fornece utilitários diversos, incluindo um contador global simples.
package utils

import (
	"strconv"
	"time"
)

// NowMillis retorna o timestamp atual em milissegundos desde a época Unix.
func NowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}


// count armazena o valor atual do contador global.
var count uint64 = 0

// Count incrementa o contador global e retorna seu valor como string.
//
// Retorno:
//   - string: valor atual do contador após incremento.
func Count() string {
	count++
	return strconv.FormatUint(count, 10)
}
