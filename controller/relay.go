package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"one-api/common"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Role    string  `json:"role"`
	Content string  `json:"content"`
	Name    *string `json:"name,omitempty"`
}

const (
	RelayModeUnknown = iota
	RelayModeChatCompletions
	RelayModeCompletions
	RelayModeEmbeddings
	RelayModeModerations
	RelayModeImagesGenerations
	RelayModeEdits
	RelayModeImagesEdits
)

// https://platform.openai.com/docs/api-reference/chat

type GeneralOpenAIRequest struct {
	Model       string    `json:"model,omitempty"`
	Messages    []Message `json:"messages,omitempty"`
	Prompt      any       `json:"prompt,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	N           int       `json:"n,omitempty"`
	Input       any       `json:"input,omitempty"`
	Instruction string    `json:"instruction,omitempty"`
	Size        string    `json:"size,omitempty"`
}

type GeneralChatbaseRequest struct {
	Model          string    `json:"model,omitempty"`
	Messages       []Message `json:"messages,omitempty"`
	Stream         bool      `json:"stream,omitempty"`
	Temperature    float64   `json:"temperature,omitempty"`
	ChatbotId      string    `json:"chatbotId,omitempty"`
	ConversationId string    `json:"conversationId,omitempty"`
}

type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type TextRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	Prompt    string    `json:"prompt"`
	MaxTokens int       `json:"max_tokens"`
	//Stream   bool      `json:"stream"`
}

type ImageRequest struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ImageEditsRequest struct {
	Image          string `json:"image"`
	Mask           string `json:"mask,omitempty"`
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatbaseError struct {
	Message string `json:"message"`
}

type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    any    `json:"code"`
}
type OpenAIResponseError struct {
	Error OpenAIError
}

type OpenAIErrorWithStatusCode struct {
	OpenAIError
	StatusCode int `json:"status_code"`
}

type ChatbaseErrorWithStatusCode struct {
	ChatbaseError
	StatusCode int    `json:"status_code"`
	Type       string `json:"type"`
}

type TextResponse struct {
	Choices []OpenAITextResponseChoice `json:"choices"`
	Usage   `json:"usage"`
	Error   OpenAIError `json:"error"`
}

type OpenAITextResponseChoice struct {
	Index        int `json:"index"`
	Message      `json:"message"`
	FinishReason string `json:"finish_reason"`
}

type OpenAITextResponse struct {
	Id      string                     `json:"id"`
	Object  string                     `json:"object"`
	Created int64                      `json:"created"`
	Choices []OpenAITextResponseChoice `json:"choices"`
	Usage   `json:"usage"`
}
type ChatbaseTextResponse struct {
	Text string `json:"text"`
}

type OpenAIEmbeddingResponseItem struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

type OpenAIEmbeddingResponse struct {
	Object string                        `json:"object"`
	Data   []OpenAIEmbeddingResponseItem `json:"data"`
	Model  string                        `json:"model"`
	Usage  `json:"usage"`
}

type ImageResponse struct {
	Created int `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	}
}

type ChatCompletionsStreamResponseChoice struct {
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}

type ChatCompletionsStreamResponse struct {
	Id      string                                `json:"id"`
	Object  string                                `json:"object"`
	Created int64                                 `json:"created"`
	Model   string                                `json:"model"`
	Choices []ChatCompletionsStreamResponseChoice `json:"choices"`
}

type CompletionsStreamResponse struct {
	Choices []struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type CustomerInfo struct {
	Name     string `json:"name,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Question string `json:"question,omitempty"`
}

func SendEmail(c *gin.Context) {
	var customerInfo CustomerInfo
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	err = json.Unmarshal(data, &customerInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	err = common.SendEmail("华瑞24小时客服获客", "shixd@rntd.cn",
		fmt.Sprintf("姓名: %s，电话：%s,邮箱 %s,问题:%s", customerInfo.Name, customerInfo.Phone, customerInfo.Email, customerInfo.Question))
	if err != nil {
		common.SysError("failed to send email" + err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"text": "发送成功",
	})
}
func RelayChatbase(c *gin.Context) {
	err := relayChatbaseTextHelper(c)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{
			"message": err.ChatbaseError.Message,
		})
	}
}

func Relay(c *gin.Context) {
	relayMode := RelayModeUnknown
	if strings.HasPrefix(c.Request.URL.Path, "/v1/chat/completions") {
		relayMode = RelayModeChatCompletions
	} else if strings.HasPrefix(c.Request.URL.Path, "/v1/completions") {
		relayMode = RelayModeCompletions
	} else if strings.HasPrefix(c.Request.URL.Path, "/v1/embeddings") {
		relayMode = RelayModeEmbeddings
	} else if strings.HasSuffix(c.Request.URL.Path, "embeddings") {
		relayMode = RelayModeEmbeddings
	} else if strings.HasPrefix(c.Request.URL.Path, "/v1/moderations") {
		relayMode = RelayModeModerations
	} else if strings.HasPrefix(c.Request.URL.Path, "/v1/images/generations") {
		relayMode = RelayModeImagesGenerations
	} else if strings.HasPrefix(c.Request.URL.Path, "/v1/images/edits") {
		relayMode = RelayModeImagesEdits
	} else if strings.HasPrefix(c.Request.URL.Path, "/v1/edits") {
		relayMode = RelayModeEdits
	}
	var err *OpenAIErrorWithStatusCode
	switch relayMode {
	case RelayModeImagesGenerations:
		err = relayImageHelper(c, relayMode)
	case RelayModeImagesEdits:
		err = relayImageEditHelper(c, relayMode)
	default:
		err = relayTextHelper(c, relayMode)
	}
	if err != nil {
		retryTimesStr := c.Query("retry")
		retryTimes, _ := strconv.Atoi(retryTimesStr)
		if retryTimesStr == "" {
			retryTimes = common.RetryTimes
		}
		if retryTimes > 0 {
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?retry=%d", c.Request.URL.Path, retryTimes-1))
		} else {
			if err.StatusCode == http.StatusTooManyRequests {
				if strings.Contains(err.Message, "You exceeded your current quota") {
					err.OpenAIError.Message = "余额不足，通道将自动关闭"
				} else {
					err.OpenAIError.Message = "当前分组上游负载已饱和，请稍后再试"
				}

			}
			c.JSON(err.StatusCode, gin.H{
				"error": err.OpenAIError,
			})
		}
		channelId := c.GetInt("channel_id")
		common.SysError(fmt.Sprintf("relay error (channel #%d): %s", channelId, err.Message))
		// https://platform.openai.com/docs/guides/error-codes/api-errors
		if shouldDisableChannel(&err.OpenAIError) {
			channelId := c.GetInt("channel_id")
			channelName := c.GetString("channel_name")
			disableChannel(channelId, channelName, err.Message)
		}
	}
}

func RelayNotImplemented(c *gin.Context) {
	err := OpenAIError{
		Message: "API not implemented",
		Type:    "one_api_error",
		Param:   "",
		Code:    "api_not_implemented",
	}
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": err,
	})
}

func RelayNotFound(c *gin.Context) {
	err := OpenAIError{
		Message: fmt.Sprintf("Invalid URL (%s %s)", c.Request.Method, c.Request.URL.Path),
		Type:    "invalid_request_error",
		Param:   "",
		Code:    "",
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": err,
	})
}
