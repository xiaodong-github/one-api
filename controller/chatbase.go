package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func chatbaseHandler(c *gin.Context, resp *http.Response) *ChatbaseErrorWithStatusCode {
	var textResponse ChatbaseTextResponse
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errorWrapper2(err, "read_response_body_failed", http.StatusInternalServerError)
	}
	err = resp.Body.Close()
	if err != nil {
		return errorWrapper2(err, "close_response_body_failed", http.StatusInternalServerError)
	}
	err = json.Unmarshal(responseBody, &textResponse)
	if err != nil {
		return errorWrapper2(err, "unmarshal_response_body_failed", http.StatusInternalServerError)
	}

	// Reset response body
	resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))

	for k, v := range resp.Header {
		c.Writer.Header().Set(k, v[0])
	}
	c.Writer.WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		return errorWrapper2(err, "copy_response_body_failed", http.StatusInternalServerError)
	}
	err = resp.Body.Close()
	if err != nil {
		return errorWrapper2(err, "close_response_body_failed", http.StatusInternalServerError)
	}
	return nil
}
