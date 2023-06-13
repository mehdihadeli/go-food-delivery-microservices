package fxapp

import "fmt"

type FxErrorHandler struct{}

func NewFxErrorHandler() *FxErrorHandler {
	return &FxErrorHandler{}
}

func (h *FxErrorHandler) HandleError(e error) {
	fmt.Println(e.Error())
}
