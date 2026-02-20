package domain

import "fmt"

type ErrorCode string

const (
	ErrValidation   ErrorCode = "VALIDATION_ERROR"
	ErrPolicyBlock  ErrorCode = "POLICY_BLOCKED"
	ErrStateConflct ErrorCode = "STATE_CONFLICT"
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrRateLimited  ErrorCode = "RATE_LIMITED"
)

type DomainError struct {
	Code    ErrorCode
	Message string
	Detail  string
}

func (e *DomainError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewValidationError(msg, detail string) *DomainError {
	return &DomainError{Code: ErrValidation, Message: msg, Detail: detail}
}

func NewPolicyBlockedError(msg, detail string) *DomainError {
	return &DomainError{Code: ErrPolicyBlock, Message: msg, Detail: detail}
}

func NewStateConflictError(msg, detail string) *DomainError {
	return &DomainError{Code: ErrStateConflct, Message: msg, Detail: detail}
}

func NewNotFoundError(msg, detail string) *DomainError {
	return &DomainError{Code: ErrNotFound, Message: msg, Detail: detail}
}

func NewRateLimitedError(msg, detail string) *DomainError {
	return &DomainError{Code: ErrRateLimited, Message: msg, Detail: detail}
}
