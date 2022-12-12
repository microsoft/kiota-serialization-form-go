package formserialization

import (
	"testing"

	assert "github.com/stretchr/testify/assert"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

func TestFormSerializationFactoryWriterHonoursInterface(t *testing.T) {
	instance := NewFormSerializationWriterFactory()
	assert.Implements(t, (*absser.SerializationWriterFactory)(nil), instance)
}
