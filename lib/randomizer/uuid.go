package seed

import (
	"fmt"
	"math"
)

type Generator interface {
	Intn(n int) int
}

// UUID4 Генерирует UUID v4 без использования crypto/rand
func UUID4(generator Generator) string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(generator.Intn(math.MaxUint8 + 1))
	}
	return fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
