package formserialization

import (
	"testing"

	assert "github.com/stretchr/testify/assert"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

func TestFormParseNodeFactoryHonoursInterface(t *testing.T) {
	instance := NewFormParseNodeFactory()
	assert.Implements(t, (*absser.ParseNodeFactory)(nil), instance)
}
