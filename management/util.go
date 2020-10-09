package management

import (
	"fmt"

	"github.com/ResultadosDigitais/x9/crypto"
)

func getID(params ...interface{}) string {
	payload := fmt.Sprintf("%v", params[0])
	for _, param := range params[1:len(params)] {
		payload = fmt.Sprintf("%v,%v", payload, param)
	}
	return crypto.SHA256(payload)
}
