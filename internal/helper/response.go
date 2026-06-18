package helper

import "go-fiber-svelte/internal/lang"

var Res res

type res struct{}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (res) Success(msg string, status ...int) *Response {
	return &Response{Message: msg}
}

func (res) SuccessData(data any, msg string, status ...int) *Response {
	return &Response{Message: msg, Data: data}
}

func (res) Error(msg string, errors any, status ...int) *Response {
	return &Response{Message: msg, Errors: errors}
}

func (res) Paginate(data any, meta *Meta, msg string) *Response {
	return &Response{Message: msg, Data: data, Meta: meta}
}

func (res) Catch(err error) *Response {
	return &Response{Message: lang.T.Get().SOMETHING_WENT_WRONG}
}

func (res) Validate(err error) *Response {
	return &Response{Message: "Validation error", Errors: err.Error()}
}
