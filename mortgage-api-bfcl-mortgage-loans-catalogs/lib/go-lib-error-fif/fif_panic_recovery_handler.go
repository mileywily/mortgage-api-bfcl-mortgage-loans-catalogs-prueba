package error_fif

import (
	"fmt"
)

func FifPanicRecoveryHandler() PanicRecoveryHandler {
	return func(recovered interface{}) (int, interface{}) {
		return CreateInternalServerError(fmt.Sprintf("%s", recovered))
	}
}
