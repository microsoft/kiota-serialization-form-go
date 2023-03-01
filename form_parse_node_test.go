package formserialization

import (
	testing "testing"

	"github.com/stretchr/testify/require"

	"github.com/microsoft/kiota-serialization-form-go/internal"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	assert "github.com/stretchr/testify/assert"
)

func TestGetRawValue(t *testing.T) {
	source := `id=2&status=200&item`
	sourceArray := []byte(source)
	parseNode, err := NewFormParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}
	someProp, err := parseNode.GetChildNode("item")
	require.NoError(t, err)
	assert.Nil(t, someProp)

	someProp, err = parseNode.GetChildNode("status")
	value, err := someProp.GetRawValue()
	assert.Equal(t, "200", *value.(*string))
}

func TestGetCollectionOfPrimitiveValues(t *testing.T) {
	source := `id=2&status=200&item=1&item=2&item=3`
	sourceArray := []byte(source)
	parseNode, err := NewFormParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}
	someProp, err := parseNode.GetChildNode("item")
	require.NoError(t, err)

	value, err := someProp.GetCollectionOfPrimitiveValues("int32")
	require.NoError(t, err)

	expected := []interface{}{ref(int32(1)), ref(int32(2)), ref(int32(3))}
	assert.Equal(t, expected, value)
}

func ref[T interface{}](t T) *T {
	return &t
}

func TestFormParseNodeHonoursInterface(t *testing.T) {
	instance := &FormParseNode{}
	assert.Implements(t, (*absser.ParseNode)(nil), instance)
}

func TestFunctional(t *testing.T) {
	sourceArray := []byte(FunctionalTestSource)
	parseNode, err := NewFormParseNode(sourceArray)
	assert.NoError(t, err)
	assert.NotNil(t, parseNode)
	parsable, err := parseNode.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	assert.NoError(t, err)
	assert.NotNil(t, parsable)
	entity, ok := parsable.(internal.TestEntityable)
	assert.True(t, ok)
	assert.Equal(t, "48d31887-5fad-4d73-a9f5-3c356e68a038", *entity.GetId())

}

const FunctionalTestSource = "displayName=Megan+Bowen&" +
	"numbers=one,two,thirtytwo&" +
	"givenName=Megan&" +
	"accountEnabled=true&" +
	"createdDateTime=2017-07-29T03:07:25Z&" +
	"jobTitle=Auditor&" +
	"mail=MeganB@M365x214355.onmicrosoft.com&" +
	"mobilePhone=null&" +
	"officeLocation=null&" +
	"preferredLanguage=en-US&" +
	"surname=Bowen&" +
	"workDuration=PT1H&" +
	"startWorkTime=08:00:00.0000000&" +
	"endWorkTime=17:00:00.0000000&" +
	"userPrincipalName=MeganB@M365x214355.onmicrosoft.com&" +
	"birthDay=2017-09-04&" +
	"id=48d31887-5fad-4d73-a9f5-3c356e68a038"
