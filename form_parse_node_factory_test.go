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

func TestFormParseNodeFactoryGetRootParseNode(t *testing.T) {
	factory := NewFormParseNodeFactory()

	node, err := factory.GetRootParseNode("application/x-www-form-urlencoded", []byte("key=val"))
	assert.NoError(t, err)
	assert.NotNil(t, node)

	node, err = factory.GetRootParseNode("application/x-www-form-urlencoded; charset=utf-8", []byte("key=val"))
	assert.NoError(t, err)
	assert.NotNil(t, node)

	node, err = factory.GetRootParseNode("  APPLICATION/X-WWW-FORM-URLENCODED  ", []byte("key=val"))
	assert.NoError(t, err)
	assert.NotNil(t, node)

	_, err = factory.GetRootParseNode("", []byte("key=val"))
	assert.Error(t, err)

	_, err = factory.GetRootParseNode("application/json", []byte("key=val"))
	assert.Error(t, err)

	var nilFactory *FormParseNodeFactory
	_, err = nilFactory.GetRootParseNode("application/x-www-form-urlencoded", []byte("key=val"))
	assert.Error(t, err)
}
