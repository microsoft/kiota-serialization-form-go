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

func TestFormSerializationWriterFactoryGetSerializationWriter(t *testing.T) {
	factory := NewFormSerializationWriterFactory()

	writer, err := factory.GetSerializationWriter("application/x-www-form-urlencoded")
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	writer, err = factory.GetSerializationWriter("application/x-www-form-urlencoded; charset=utf-8")
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	writer, err = factory.GetSerializationWriter("  APPLICATION/X-WWW-FORM-URLENCODED  ")
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	_, err = factory.GetSerializationWriter("")
	assert.Error(t, err)

	_, err = factory.GetSerializationWriter("application/json")
	assert.Error(t, err)

	var nilFactory *FormSerializationWriterFactory
	_, err = nilFactory.GetSerializationWriter("application/x-www-form-urlencoded")
	assert.Error(t, err)
}
