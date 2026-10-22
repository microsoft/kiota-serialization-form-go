package formserialization

import (
	"errors"
	"strings"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// FormParseNodeFactory is a ParseNodeFactory implementation for URI form encoded
type FormParseNodeFactory struct {
}

// Creates a new FormParseNodeFactory
func NewFormParseNodeFactory() *FormParseNodeFactory {
	return &FormParseNodeFactory{}
}

// GetValidContentType returns the content type this factory's parse nodes can deserialize.
func (f *FormParseNodeFactory) GetValidContentType() (string, error) {
	return "application/x-www-form-urlencoded", nil
}

// GetRootParseNode return a new ParseNode instance that is the root of the content
func (f *FormParseNodeFactory) GetRootParseNode(contentType string, content []byte) (absser.ParseNode, error) {
	if f == nil {
		return nil, errors.New("factory cannot be nil")
	}
	validType, err := f.GetValidContentType()
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(contentType)
	if trimmed == "" {
		return nil, errors.New("contentType is empty")
	}

	mediaType, _, _ := strings.Cut(trimmed, ";")
	if strings.ToLower(strings.TrimSpace(mediaType)) != validType {
		return nil, errors.New("contentType is not valid")
	}

	return NewFormParseNode(content)
}
