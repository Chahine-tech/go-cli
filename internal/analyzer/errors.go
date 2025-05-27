package parser

import (
	"fmt"
)

type FileNotFoundError struct {
	Path string
	Err  error
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("file not found or inaccessible: %s", e.Path)
}

func (e *FileNotFoundError) Unwrap() error {
	return e.Err
}

func NewFileNotFoundError(path string, err error) *FileNotFoundError {
	return &FileNotFoundError{
		Path: path,
		Err:  err,
	}
}

type ParseError struct {
	LogID   string
	Message string
	Err     error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error for log %s: %s", e.LogID, e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

func NewParseError(logID, message string, err error) *ParseError {
	return &ParseError{
		LogID:   logID,
		Message: message,
		Err:     err,
	}
} 