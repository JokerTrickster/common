package error

import (
	"fmt"
	"strings"
)

/*
	에러 핸들링 로직
*/

// HandleError combines error messages and additional parameters
func HandleError(errMsg string, args ...interface{}) string {
	// Convert arguments to strings
	var params []string
	for _, arg := range args {
		params = append(params, fmt.Sprintf("%v", arg))
	}

	// Combine error message and parameters
	return fmt.Sprintf("Error: %s | Parameters: %s", errMsg, strings.Join(params, ", "))
}
