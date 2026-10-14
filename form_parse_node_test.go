package formserialization

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"testing"

	"github.com/microsoft/kiota-serialization-form-go/internal"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/stretchr/testify/assert"
)

// testColor is a simple enum type used in tests.
type testColor int

const (
	testColorRed testColor = iota
	testColorBlue
	testColorGreen
)

// testColorFactory parses "red", "blue", "green" into testColor values.
func testColorFactory(value string) (interface{}, error) {
	switch value {
	case "red":
		return testColorRed, nil
	case "blue":
		return testColorBlue, nil
	case "green":
		return testColorGreen, nil
	default:
		return nil, errors.New("unknown color: " + value)
	}
}

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
	require.NoError(t, err)
	someProp, err := parseNode.GetChildNode("item")
	require.NoError(t, err)

	value, err := someProp.GetCollectionOfPrimitiveValues("int32")
	require.NoError(t, err)

	expected := []interface{}{ref(int32(1)), ref(int32(2)), ref(int32(3))}
	assert.Equal(t, expected, value)
}

func TestGetCollectionOfPrimitiveValuesTypes(t *testing.T) {
	assert.Equal(t,
		[]interface{}{ref("milk"), ref("soda")},
		getCollectionValues("id=2&item=milk&item=soda", "item", "string"),
	)
	assert.Equal(t,
		[]interface{}{ref(true), ref(false)},
		getCollectionValues("id=2&item=true&item=false", "item", "bool"),
	)
	assert.Equal(t,
		[]interface{}{ref(int8(1)), ref(int8(2)), ref(int8(3))},
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "uint8"),
	)
	assert.Equal(t,
		[]interface{}{ref(byte(1)), ref(byte(2)), ref(byte(3))},
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "byte"),
	)
	assert.Equal(t,
		[]interface{}{ref(float32(1)), ref(float32(2)), ref(float32(3))},
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "float32"),
	)
	assert.Equal(t,
		[]interface{}{ref(float64(1)), ref(float64(2)), ref(float64(3))},
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "float64"),
	)
	assert.Equal(t,
		[]interface{}{ref(float64(1)), ref(float64(2)), ref(float64(3))},
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "float64"),
	)
	assert.Equal(t,
		[]interface{}{ref(int32(1)), ref(int32(2)), ref(int32(3))},
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "int32"),
	)
	assert.Equal(t,
		[]interface{}{ref(int64(1)), ref(int64(2)), ref(int64(3))},
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "int64"),
	)
	assert.Nil(t,
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "time"),
	)
	assert.Nil(t,
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "timeonly"),
	)
	assert.Nil(t,
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "dateonly"),
	)
	assert.Nil(t,
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "isoduration"),
	)
	assert.Nil(t,
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "uuid"),
	)
	assert.Nil(t,
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "base64"),
	)
}

func getCollectionValues(source string, indexName string, targetType string) []interface{} {
	sourceArray := []byte(source)
	parseNode, _ := NewFormParseNode(sourceArray)
	someProp, _ := parseNode.GetChildNode(indexName)

	value, _ := someProp.GetCollectionOfPrimitiveValues(targetType)
	return value
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

func TestGetEnumValue(t *testing.T) {
	t.Run("returns error when parser is nil", func(t *testing.T) {
		source := `color=red`
		parseNode, err := NewFormParseNode([]byte(source))
		require.NoError(t, err)
		node, err := parseNode.GetChildNode("color")
		require.NoError(t, err)

		val, err := node.GetEnumValue(nil)
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("returns enum value for valid string", func(t *testing.T) {
		source := `color=red`
		parseNode, err := NewFormParseNode([]byte(source))
		require.NoError(t, err)
		node, err := parseNode.GetChildNode("color")
		require.NoError(t, err)

		val, err := node.GetEnumValue(testColorFactory)
		require.NoError(t, err)
		assert.Equal(t, testColorRed, val)
	})

	t.Run("returns nil when node value is empty", func(t *testing.T) {
		node := &FormParseNode{value: ""}

		val, err := node.GetEnumValue(testColorFactory)
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("returns error when factory cannot parse value", func(t *testing.T) {
		source := `color=unknown`
		parseNode, err := NewFormParseNode([]byte(source))
		require.NoError(t, err)
		node, err := parseNode.GetChildNode("color")
		require.NoError(t, err)

		val, err := node.GetEnumValue(testColorFactory)
		assert.Error(t, err)
		assert.Nil(t, val)
	})
}

func TestGetCollectionOfEnumValues(t *testing.T) {
	t.Run("returns nil for nil node", func(t *testing.T) {
		var node *FormParseNode

		val, err := node.GetCollectionOfEnumValues(testColorFactory)
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("returns nil for empty value", func(t *testing.T) {
		node := &FormParseNode{value: ""}

		val, err := node.GetCollectionOfEnumValues(testColorFactory)
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("returns error when parser is nil", func(t *testing.T) {
		source := `colors=red,blue`
		parseNode, err := NewFormParseNode([]byte(source))
		require.NoError(t, err)
		node, err := parseNode.GetChildNode("colors")
		require.NoError(t, err)

		val, err := node.GetCollectionOfEnumValues(nil)
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("returns single-element slice for one value", func(t *testing.T) {
		source := `colors=green`
		parseNode, err := NewFormParseNode([]byte(source))
		require.NoError(t, err)
		node, err := parseNode.GetChildNode("colors")
		require.NoError(t, err)

		val, err := node.GetCollectionOfEnumValues(testColorFactory)
		require.NoError(t, err)
		assert.Equal(t, []interface{}{testColorGreen}, val)
	})

	t.Run("returns all enum values from comma-separated string", func(t *testing.T) {
		source := `colors=red,blue,green`
		parseNode, err := NewFormParseNode([]byte(source))
		require.NoError(t, err)
		node, err := parseNode.GetChildNode("colors")
		require.NoError(t, err)

		val, err := node.GetCollectionOfEnumValues(testColorFactory)
		require.NoError(t, err)
		assert.Equal(t, []interface{}{testColorRed, testColorBlue, testColorGreen}, val)
	})

	t.Run("returns error when one value in collection is invalid", func(t *testing.T) {
		source := `colors=red,unknown,green`
		parseNode, err := NewFormParseNode([]byte(source))
		require.NoError(t, err)
		node, err := parseNode.GetChildNode("colors")
		require.NoError(t, err)

		val, err := node.GetCollectionOfEnumValues(testColorFactory)
		assert.Error(t, err)
		assert.Nil(t, val)
	})
}

func TestInt64ValuesRoundTripWithoutPrecisionLoss(t *testing.T) {
	values := []int64{0, 42, 9007199254740991, 9007199254740993, -9007199254740993, math.MaxInt64, math.MinInt64}
	for _, value := range values {
		t.Run(fmt.Sprint(value), func(t *testing.T) {
			writer := NewFormSerializationWriter()
			require.NoError(t, writer.WriteInt64Value("id", &value))
			content, err := writer.GetSerializedContent()
			require.NoError(t, err)
			root, err := NewFormParseNodeFactory().GetRootParseNode("application/x-www-form-urlencoded", content)
			require.NoError(t, err)
			node, err := root.GetChildNode("id")
			require.NoError(t, err)
			require.NotNil(t, node)
			actual, err := node.GetInt64Value()
			require.NoError(t, err)
			require.NotNil(t, actual)
			assert.Equal(t, value, *actual)
		})
	}
	writer := NewFormSerializationWriter()
	require.NoError(t, writer.WriteCollectionOfInt64Values("ids", values))
	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	root, err := NewFormParseNode(content)
	require.NoError(t, err)
	node, err := root.GetChildNode("ids")
	require.NoError(t, err)
	actual, err := node.GetCollectionOfPrimitiveValues("int64")
	require.NoError(t, err)
	require.Len(t, actual, len(values))
	for i, value := range values {
		assert.Equal(t, value, *actual[i].(*int64))
	}
}

func TestInt64ValuePreservesNumericCoercion(t *testing.T) {
	for value, expected := range map[string]int64{"1.5": 1, "-1.5": -1, "1e3": 1000} {
		t.Run(value, func(t *testing.T) {
			node, err := NewFormParseNode([]byte(value))
			require.NoError(t, err)
			actual, err := node.GetInt64Value()
			require.NoError(t, err)
			require.NotNil(t, actual)
			assert.Equal(t, expected, *actual)
		})
	}
}

func TestByteValuesRoundTripAcrossUnsignedRange(t *testing.T) {
	values := []byte{0, 97, 127, 128, 255}
	for _, value := range values {
		t.Run(fmt.Sprint(value), func(t *testing.T) {
			writer := NewFormSerializationWriter()
			require.NoError(t, writer.WriteByteValue("value", &value))
			content, err := writer.GetSerializedContent()
			require.NoError(t, err)
			root, err := NewFormParseNode(content)
			require.NoError(t, err)
			node, err := root.GetChildNode("value")
			require.NoError(t, err)
			actual, err := node.GetByteValue()
			require.NoError(t, err)
			require.NotNil(t, actual)
			assert.Equal(t, value, *actual)
		})
	}
	writer := NewFormSerializationWriter()
	require.NoError(t, writer.WriteCollectionOfByteValues("values", values))
	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	root, err := NewFormParseNode(content)
	require.NoError(t, err)
	node, err := root.GetChildNode("values")
	require.NoError(t, err)
	actual, err := node.GetCollectionOfPrimitiveValues("byte")
	require.NoError(t, err)
	require.Len(t, actual, len(values))
	for i, value := range values {
		assert.Equal(t, value, *actual[i].(*byte))
	}
}

func TestByteValueRejectsInvalidValues(t *testing.T) {
	for _, value := range []string{"-1", "-128", "256", "invalid"} {
		t.Run(value, func(t *testing.T) {
			root, err := NewFormParseNode([]byte("value=" + value))
			require.NoError(t, err)
			node, err := root.GetChildNode("value")
			require.NoError(t, err)
			actual, err := node.GetByteValue()
			assert.Error(t, err)
			assert.Nil(t, actual)
			collection, err := node.GetCollectionOfPrimitiveValues("byte")
			assert.Error(t, err)
			assert.Nil(t, collection)
		})
	}
}

func TestInt8ValueRetainsSignedRange(t *testing.T) {
	for _, value := range []int8{-128, -1, 0, 127} {
		node, err := NewFormParseNode([]byte(fmt.Sprint(value)))
		require.NoError(t, err)
		actual, err := node.GetInt8Value()
		require.NoError(t, err)
		require.NotNil(t, actual)
		assert.Equal(t, value, *actual)
	}
}
