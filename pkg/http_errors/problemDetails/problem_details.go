package problemDetails

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	ContentTypeJSON = "application/problem+json"
)

// ProblemDetailErr ProblemDetail error interface
type ProblemDetailErr interface {
	GetStatus() int
	SetStatus(status int) ProblemDetailErr
	GetTitle() string
	SetTitle(title string) ProblemDetailErr
	GetStackTrace() string
	SetStackTrace(stackTrace string) ProblemDetailErr
	GetDetail() string
	SetDetail(detail string) ProblemDetailErr
	GetType() string
	SetType(typ string) ProblemDetailErr
	Error() string
	ErrBody() error
	WriteTo(w http.ResponseWriter) (int, error)
}

// ProblemDetail error struct
type problemDetail struct {
	Status     int       `json:"status,omitempty"`
	Title      string    `json:"title,omitempty"`
	Detail     string    `json:"detail,omitempty"`
	Type       string    `json:"type,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	StackTrace string    `json:"stackTrace,omitempty"`
}

// ErrBody Error body
func (p *problemDetail) ErrBody() error {
	return p
}

// Error  Error() interface method
func (p *problemDetail) Error() string {
	return fmt.Sprintf("Error Title: %s - Error Status: %d - Error Detail: %s", p.Title, p.Status, p.Detail)
}

func (p *problemDetail) GetStatus() int {
	return p.Status
}

func (p *problemDetail) SetStatus(status int) ProblemDetailErr {
	p.Status = status

	return p
}

func (p *problemDetail) GetTitle() string {
	return p.Title
}

func (p *problemDetail) SetTitle(title string) ProblemDetailErr {
	p.Title = title

	return p
}

func (p *problemDetail) GetType() string {
	return p.Type
}

func (p *problemDetail) SetType(typ string) ProblemDetailErr {
	p.Type = typ

	return p
}

func (p *problemDetail) GetDetail() string {
	return p.Detail
}

func (p *problemDetail) SetDetail(detail string) ProblemDetailErr {
	p.Detail = detail

	return p
}

func (p *problemDetail) GetStackTrace() string {
	return p.StackTrace
}

func (p *problemDetail) SetStackTrace(stackTrace string) ProblemDetailErr {
	p.StackTrace = stackTrace

	return p
}

func (p *problemDetail) json() []byte {
	b, _ := json.Marshal(&p)
	return b
}

// WriteTo writes the JSON Problem to an HTTP Response Writer
func (p *problemDetail) WriteTo(w http.ResponseWriter) (int, error) {
	p.writeHeaderTo(w)
	return w.Write(p.json())
}

func (p *problemDetail) writeHeaderTo(w http.ResponseWriter) {
	w.Header().Set("Content-Type", ContentTypeJSON)
	status := p.GetStatus()
	if status == 0 {
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
}

// NewProblemDetail New ProblemDetail Error
func NewProblemDetail(status int, title string, detail string, stackTrace string) ProblemDetailErr {
	problemDetail := &problemDetail{
		Status:     status,
		Title:      title,
		Timestamp:  time.Now(),
		Detail:     detail,
		Type:       getDefaultType(status),
		StackTrace: stackTrace,
	}

	return problemDetail
}

// NewProblemDetailFromCode New ProblemDetail Error With Message
func NewProblemDetailFromCode(status int, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Status:     status,
		Title:      http.StatusText(status),
		Timestamp:  time.Now(),
		Type:       getDefaultType(status),
		StackTrace: stackTrace,
	}
}

// NewProblemDetailFromCodeAndDetail New ProblemDetail Error With Message
func NewProblemDetailFromCodeAndDetail(status int, detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Status:     status,
		Title:      http.StatusText(status),
		Detail:     detail,
		Timestamp:  time.Now(),
		Type:       getDefaultType(status),
		StackTrace: stackTrace,
	}
}

func getDefaultType(statusCode int) string {
	return fmt.Sprintf("https://httpstatuses.io/%d", statusCode)
}
