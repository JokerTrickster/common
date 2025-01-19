package request

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/labstack/echo/v4"
)

// RequestData represents the structured request data
type RequestData struct {
	Body  map[string]interface{} `json:"body,omitempty"`  // JSON Body
	Query map[string][]string    `json:"query,omitempty"` // Query Parameters
	Path  map[string]string      `json:"path,omitempty"`  // Path Parameters
}

// ParseRequest extracts JSON Body, Query Parameters, and Path Parameters from the request
func ParseRequest(c echo.Context) (RequestData, error) {
	var requestData RequestData

	// Extract JSON Body if the method is POST, PUT, or PATCH
	if c.Request().Method == "POST" || c.Request().Method == "PUT" || c.Request().Method == "PATCH" {
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil && err != io.EOF {
			return requestData, err
		}
		// Reset the Body to allow re-reading later if necessary
		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Unmarshal JSON Body
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &requestData.Body); err != nil {
				return requestData, err
			}
		}
	}

	// Extract Query Parameters
	requestData.Query = c.QueryParams()

	// Extract Path Parameters
	pathParams := c.ParamNames()
	pathValues := make(map[string]string)
	for _, param := range pathParams {
		pathValues[param] = c.Param(param)
	}
	requestData.Path = pathValues
	

	return requestData, nil
}
