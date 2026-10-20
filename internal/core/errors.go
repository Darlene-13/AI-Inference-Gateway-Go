package core

import (
	"errors"
	"fmt"
)

type Kind int

const (
	KindInternal Kind = iota // A bug in the gateway
	KindInvalidRequest
	KindModelNotFound
	KindUpStreamRejected
	KindUpStreamError
	KindUpStreamUnavailable
	KindUpStreamTimeout
	KindCancelled
)

var kindNames = map[Kind]string{
	KindInternal:            "internal",
	KindInvalidRequest:      "invalid_request",
	KindModelNotFound:       "model_not_found",
	KindUpStreamRejected:    "up_stream_rejected",
	KindUpStreamError:       "up_stream_error",
	KindUpStreamUnavailable: "up_stream_unavailable",
	KindUpStreamTimeout:     "up_stream_timeout",
	KindCancelled:           "cancelled",
}

func (k Kind) String() string {
	if name, ok := kindNames[k]; ok {
		return name
	}

	return fmt.Sprintf("%d", int(k))
}

type Error struct {
	Kind           Kind
	Code           string
	Message        string
	Provider       string
	UpstreamStatus string
	Err            error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s : %v", e.Code, e.Err)
	}

	return fmt.Sprintf("%s : %s", e.Kind, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Retryable() bool {
	switch e.Kind {
	case KindUpStreamUnavailable, KindUpStreamTimeout:
		return true
	case KindUpStreamError:
		return e.UpstreamStatus == 0 || e.UpstreamStatus >= 500
	case KindInternal, KindInvalidRequest, KindModelNotFound, KindUpStreamRejected, KindCancelled:
		return false
	default:
		return false
	}

}

func kindOf(err error) Kind {
	var gerr *Error
	if errors.As(err, &gerr) {
		return gerr.Kind
	}

	return KindInternal
}

func InvalidRequest(code, message string) *Error {
	return &Error{Kind: KindInvalidRequest, Code: code, Message: message}
}

func ModeNotFound(alias string) *Error {
	return &Error{
		Kind:    KindModelNotFound,
		Code:    "model_not_found",
		Message: fmt.Sprintf("The model %q does not exist. Call Get /v1/models to see available models.", alias),
	}
}
