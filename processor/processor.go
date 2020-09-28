package processor

import (
	"net/http"
)

type ProcessorInterface interface {
	RequestProcess(*http.Request) *http.Request
	ResponseProcess(*http.Response) *http.Response
	SetNext(ProcessorInterface) ProcessorInterface
}
