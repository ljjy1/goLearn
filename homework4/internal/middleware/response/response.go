package response

/**
 * @Description: 封装统一返回和错误处理
 */
import (
	"net/http"

	"homework4/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	CodeSuccess      = 200 // 成功返回码
	CodeServerError  = 500 // 服务器错误返回码
	CodeBadRequest   = 400 // 参数错误或者业务错误返回码
	CodeUnauthorized = 401 // 未授权返回码
)

type BizError struct {
	Code    int         `json:"code" example:"400"`
	Message string      `json:"message" example:"参数错误"`
	Detail  interface{} `json:"detail"`
}

func (e *BizError) Error() string {
	return e.Message
}

func NewBizError(code int, message string, detail interface{}) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// 未授权异常
func NewUnauthorizedError(message string) *BizError {
	return &BizError{
		Code:    CodeUnauthorized,
		Message: message,
		Detail:  nil,
	}
}

// 参数错误异常
func NewBadRequestError(message string) *BizError {
	return &BizError{
		Code:    CodeBadRequest,
		Message: message,
		Detail:  nil,
	}
}

type Response struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type HandlerFunc func(c *gin.Context) error

func Success(data interface{}) *Response {
	return &Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	}
}

func Fail(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}

/**
 * @Description: 发送成功JSON响应
 * @param c
 * @param data
 */
func SendJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Success(data))
}

/**
 * @Description: 发送错误JSON响应
 * @param c
 * @param err
 */
func SendErrorJSON(c *gin.Context, err error) {
	logger := logger.AppLog
	var resp *Response
	var statusCode int

	if bizErr, ok := err.(*BizError); ok {
		resp = Fail(bizErr.Code, bizErr.Message)

		switch bizErr.Code {
		case CodeUnauthorized:
			statusCode = http.StatusUnauthorized
		case CodeBadRequest:
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusOK
		}

		logger.Error("业务错误",
			zap.Int("code", bizErr.Code),
			zap.String("message", bizErr.Message),
			zap.Any("detail", bizErr.Detail),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
	} else {
		resp = Fail(CodeServerError, "服务器内部错误")
		statusCode = http.StatusInternalServerError
		logger.Error("系统错误",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
	}

	c.JSON(statusCode, resp)
}

/**
 * @Description: 错误处理中间件
 * @return gin.HandlerFunc
 */
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			SendErrorJSON(c, err)
			c.Abort()
		}
	}
}

/**
 * @Description: 包装处理函数，捕获异常
 * @param h
 * @return gin.HandlerFunc
 */
func WrapHandler(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var e error
				switch v := err.(type) {
				case error:
					e = v
				default:
					e = &BizError{
						Code:    CodeServerError,
						Message: "系统内部错误",
						Detail:  v,
					}
				}
				c.Error(e)
			}
		}()

		if err := h(c); err != nil {
			c.Error(err)
		}
	}
}
