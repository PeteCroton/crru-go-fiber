package helpers

import (
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Success(code int, data any) IResponse
	Error(code int, tractId, msg string) IResponse
	Res() error
}

type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	Context    *fiber.Ctx
	IsError    bool
}

type ErrorResponse struct {
	TraceId string `json:"trace_id"`
	Msg     string `json:"message"`
}

func NewResponse(c *fiber.Ctx) IResponse {
	return &Response{
		Context: c,
	}
}

func (r *Response) Success(code int, data any) IResponse {
	r.StatusCode = code
	r.Data = data
	r.IsError = false
	InitLogger(r.Context, &r.Data, code).Print()
	return r
}

func (r *Response) Error(code int, tractId, msg string) IResponse {
	r.StatusCode = code
	r.ErrorRes = &ErrorResponse{
		TraceId: tractId,
		Msg:     msg,
	}
	r.IsError = true
	InitLogger(r.Context, &r.ErrorRes, code).Print().Save()
	return r
}
func (r *Response) Res() error {
	// return r.Context.Status(r.StatusCode).JSON(func() any {
	// 	if r.IsError {
	// 		return &r.ErrorRes
	// 	}
	// 	return &r.Data
	// }())

	return r.Context.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    r.Data,
		"message": &r.ErrorRes,
		"code":    r.StatusCode,
		"success": !r.IsError,
	})
}

type PaginateRes struct {
	Data      any `json:"data"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
	TotalItem int `json:"total_item"`
}
