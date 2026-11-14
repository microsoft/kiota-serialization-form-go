package formserialization

import (
	"errors"
	"github.com/google/uuid"
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"strconv"
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
		[]interface{}{ref(int8(1)), ref(int8(2)), ref(int8(3))},
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "int8"),
	)
	assert.Equal(t,
		[]interface{}{ref(int32(1)), ref(int32(2)), ref(int32(3))},
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "int32"),
	)
	assert.Equal(t,
		[]interface{}{ref(int32(42)), ref(int32(-42))},
		getCollectionValues("id=2&status=200&item=42000000000e-9&item=-42000000000e-9", "item", "int32"),
	)
	assert.Equal(t,
		[]interface{}{ref(int64(1)), ref(int64(2)), ref(int64(3))},
		getCollectionValues("id=2&status=200&item=1&item=2&item=3", "item", "int64"),
	)
	assert.Equal(t,
		[]interface{}{ref(int64(-9223372036854775808)), ref(int64(9223372036854775807))},
		getCollectionValues("id=2&status=200&item=-9223372036854775808e0&item=9223372036854775807e0", "item", "int64"),
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

func TestInt64ValuesRejectOutOfRangeIntegers(t *testing.T) {
	for _, value := range []string{
		"9223372036854775808",
		"-9223372036854775809",
		"18446744073709551616",
		"9223372036854775808e0",
		"-9223372036854775809e0",
		"420000000000000000000e-1",
		"-420000000000000000000e-1",
		"1e20",
		"-1e20",
	} {
		t.Run(value, func(t *testing.T) {
			t.Run("scalar", func(t *testing.T) {
				node, err := NewFormParseNode([]byte(value))
				require.NoError(t, err)
				actual, err := node.GetInt64Value()
				require.ErrorIs(t, err, strconv.ErrRange)
				assert.Nil(t, actual)
			})
			t.Run("collection", func(t *testing.T) {
				node, err := NewFormParseNode([]byte("42," + value))
				require.NoError(t, err)
				values, err := node.GetCollectionOfPrimitiveValues("int64")
				require.ErrorIs(t, err, strconv.ErrRange)
				assert.Nil(t, values)
			})
		})
	}
}

func TestInt32ValuesRejectOutOfRangeIntegers(t *testing.T) {
	for _, value := range []string{
		"2147483648",
		"-2147483649",
		"2147483648e0",
		"-2147483649e0",
		"42000000000e-1",
		"-42000000000e-1",
		"1e11",
		"-1e11",
	} {
		t.Run(value, func(t *testing.T) {
			t.Run("scalar", func(t *testing.T) {
				node, err := NewFormParseNode([]byte(value))
				require.NoError(t, err)
				actual, err := node.GetInt32Value()
				require.ErrorIs(t, err, strconv.ErrRange)
				assert.Nil(t, actual)
			})
			t.Run("collection", func(t *testing.T) {
				node, err := NewFormParseNode([]byte("42," + value))
				require.NoError(t, err)
				values, err := node.GetCollectionOfPrimitiveValues("int32")
				require.ErrorIs(t, err, strconv.ErrRange)
				assert.Nil(t, values)
			})
		})
	}
}

func TestInt64ValuePreservesNumericCoercion(t *testing.T) {
	for value, expected := range map[string]int64{
		"1.5": 1,
		"-1.5": -1,
		"1e3": 1000,
		"-9223372036854775808e0": -9223372036854775808,
		"9223372036854775807e0": 9223372036854775807,
		"420000000000000000000e-10": 42000000000,
		"-420000000000000000000e-10": -42000000000,
	} {
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

func TestAdditionalDataStringsRoundTrip(t *testing.T) {
	for _, value := range []string{"a&b", "hello world", "a+b", "%2F", "café", "one,two", "plain"} {
		t.Run(value, func(t *testing.T) {
			entity := internal.NewTestEntity()
			entity.SetOfficeLocation(&value)
			entity.SetAdditionalData(map[string]interface{}{"custom": value})
			writer := NewFormSerializationWriter()
			require.NoError(t, writer.WriteObjectValue("", entity))
			require.NoError(t, writer.WriteAdditionalData(entity.GetAdditionalData()))
			content, err := writer.GetSerializedContent()
			require.NoError(t, err)
			node, err := NewFormParseNodeFactory().GetRootParseNode("application/x-www-form-urlencoded", content)
			require.NoError(t, err)
			result, err := node.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
			require.NoError(t, err)
			actual := result.(internal.TestEntityable)
			require.NotNil(t, actual.GetOfficeLocation())
			assert.Equal(t, value, *actual.GetOfficeLocation())
			assert.Equal(t, value, actual.GetAdditionalData()["custom"])
		})
	}
}

func TestFieldValuesPreserveEqualsSigns(t *testing.T) {
	for _, value := range []string{"a=b", "a=b=c", "=leading", "trailing=", "="} {
		t.Run(value, func(t *testing.T) {
			node, err := NewFormParseNode([]byte("officeLocation=" + value + "&custom=" + value + "&other=plain"))
			require.NoError(t, err)
			result, err := node.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
			require.NoError(t, err)
			actual := result.(internal.TestEntityable)
			require.NotNil(t, actual.GetOfficeLocation())
			assert.Equal(t, value, *actual.GetOfficeLocation())
			assert.Equal(t, value, actual.GetAdditionalData()["custom"])
			assert.Equal(t, "plain", actual.GetAdditionalData()["other"])
		})
	}
}

func TestEncodedEqualsSignsDecodeOnce(t *testing.T) {
	node, err := NewFormParseNode([]byte("officeLocation=a%3Db&custom=a%3Db&literal=%253D"))
	require.NoError(t, err)
	result, err := node.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	actual := result.(internal.TestEntityable)
	require.NotNil(t, actual.GetOfficeLocation())
	assert.Equal(t, "a=b", *actual.GetOfficeLocation())
	assert.Equal(t, "a=b", actual.GetAdditionalData()["custom"])
	assert.Equal(t, "%3D", actual.GetAdditionalData()["literal"])
}

func TestAdditionalDataRepeatedValuesDecodeOnce(t *testing.T) {
	node, err := NewFormParseNode([]byte("custom%2Bfield=one%2Btwo&custom%2Bfield=three%2Cfour&custom%2Bfield=%252F"))
	require.NoError(t, err)
	result, err := node.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Equal(t, "one+two,three,four,%2F", result.(internal.TestEntityable).GetAdditionalData()["custom+field"])
}

func TestAdditionalDataRejectsInvalidEscapes(t *testing.T) {
	for _, value := range []string{"%", "%ZZ"} {
		t.Run(value, func(t *testing.T) {
			node, err := NewFormParseNode([]byte("custom=" + value))
			require.NoError(t, err)
			result, err := node.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestFieldParsingPreservesCollectionAndEmptyValueBehavior(t *testing.T) {
	node, err := NewFormParseNode([]byte("items=a%2Cb&items=%252F&items=one%2Btwo&custom=&other=plain"))
	require.NoError(t, err)
	child, err := node.GetChildNode("items")
	require.NoError(t, err)
	values, err := child.GetCollectionOfPrimitiveValues("string")
	require.NoError(t, err)
	assert.Equal(t, []interface{}{ref("a,b"), ref("%2F"), ref("one+two")}, values)
	result, err := node.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.NotContains(t, result.(internal.TestEntityable).GetAdditionalData(), "custom")
}

func TestGetInt32Value(t *testing.T) {
	t.Run("parses valid int32", func(t *testing.T) {
		node, err := NewFormParseNode([]byte("value=2147483647"))
		require.NoError(t, err)
		child, err := node.GetChildNode("value")
		require.NoError(t, err)
		val, err := child.GetInt32Value()
		require.NoError(t, err)
		assert.Equal(t, int32(2147483647), *val)
	})

	t.Run("parses negative int32", func(t *testing.T) {
		node, err := NewFormParseNode([]byte("value=-2147483648"))
		require.NoError(t, err)
		child, err := node.GetChildNode("value")
		require.NoError(t, err)
		val, err := child.GetInt32Value()
		require.NoError(t, err)
		assert.Equal(t, int32(-2147483648), *val)
	})

	t.Run("rejects overflow int32", func(t *testing.T) {
		node, err := NewFormParseNode([]byte("value=2147483648"))
		require.NoError(t, err)
		child, err := node.GetChildNode("value")
		require.NoError(t, err)
		val, err := child.GetInt32Value()
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("coerces floating point integer format", func(t *testing.T) {
		node, err := NewFormParseNode([]byte("value=42.0"))
		require.NoError(t, err)
		child, err := node.GetChildNode("value")
		require.NoError(t, err)
		val, err := child.GetInt32Value()
		require.NoError(t, err)
		assert.Equal(t, int32(42), *val)
	})

	t.Run("parses scientific notation with large mantissa", func(t *testing.T) {
		node, err := NewFormParseNode([]byte("value=42000000000e-9"))
		require.NoError(t, err)
		child, err := node.GetChildNode("value")
		require.NoError(t, err)
		val, err := child.GetInt32Value()
		require.NoError(t, err)
		assert.Equal(t, int32(42), *val)
	})

	t.Run("parses negative scientific notation with large mantissa", func(t *testing.T) {
		node, err := NewFormParseNode([]byte("value=-42000000000e-9"))
		require.NoError(t, err)
		child, err := node.GetChildNode("value")
		require.NoError(t, err)
		val, err := child.GetInt32Value()
		require.NoError(t, err)
		assert.Equal(t, int32(-42), *val)
	})
}


func TestGetTimeValueWithSubsecondsAndWhitespace(t *testing.T) {
	node, err := NewFormParseNode([]byte("t1=2020-12-31T23%3A59%3A59.123456789Z&t2=2020-12-31T23%3A59%3A59.999Z&t3=%202020-12-31T23%3A59%3A59Z%20"))
	assert.NoError(t, err)

	child1, err := node.GetChildNode("t1")
	assert.NoError(t, err)
	t1, err := child1.GetTimeValue()
	assert.NoError(t, err)
	assert.NotNil(t, t1)
	assert.Equal(t, 123456789, t1.Nanosecond())

	child2, err := node.GetChildNode("t2")
	assert.NoError(t, err)
	t2, err := child2.GetTimeValue()
	assert.NoError(t, err)
	assert.NotNil(t, t2)
	assert.Equal(t, 999000000, t2.Nanosecond())

	child3, err := node.GetChildNode("t3")
	assert.NoError(t, err)
	t3, err := child3.GetTimeValue()
	assert.NoError(t, err)
	assert.NotNil(t, t3)
	assert.Equal(t, 2020, t3.Year())
}

func TestFormParseNodeNilReceiverGuards(t *testing.T) {
	var nilNode *FormParseNode

	assert.NotPanics(t, func() {
		child, err := nilNode.GetChildNode("key")
		assert.Error(t, err)
		assert.Nil(t, child)

		objs, err := nilNode.GetCollectionOfObjectValues(nil)
		assert.NoError(t, err)
		assert.Nil(t, objs)

		assert.Nil(t, nilNode.GetOnBeforeAssignFieldValues())
		assert.Error(t, nilNode.SetOnBeforeAssignFieldValues(func(absser.Parsable) error { return nil }))

		assert.Nil(t, nilNode.GetOnAfterAssignFieldValues())
		assert.Error(t, nilNode.SetOnAfterAssignFieldValues(func(absser.Parsable) error { return nil }))
	})
}

func TestLoadFieldsPreservesInitialEmptyValuesAndSkipsEmptyKeys(t *testing.T) {
	// Leading empty value for items collection and empty key in query string
	node, err := NewFormParseNode([]byte("=emptykey&items=&items=second&foo=bar"))
	require.NoError(t, err)

	// Child node for items must contain both empty string and "second"
	child, err := node.GetChildNode("items")
	require.NoError(t, err)
	require.NotNil(t, child)

	values, err := child.GetCollectionOfPrimitiveValues("string")
	require.NoError(t, err)
	require.Len(t, values, 2)
	assert.Equal(t, "", *values[0].(*string))
	assert.Equal(t, "second", *values[1].(*string))

	// Empty key "=emptykey" must not be present in fields
	assert.NotContains(t, node.fields, "")
	_, hasEmptyKey := node.fields[""]
	assert.False(t, hasEmptyKey, "fields must not contain an empty key")
	emptyKeyChild, err := node.GetChildNode("")
	assert.Error(t, err)
	assert.Nil(t, emptyKeyChild)
}

func TestGetUUIDValueNilOnErrorAndTrim(t *testing.T) {
	// Invalid UUID string should return (nil, error) without zero-UUID pointer
	invalidNode, err := NewFormParseNode([]byte("id=not-a-valid-uuid"))
	require.NoError(t, err)
	idNode, err := invalidNode.GetChildNode("id")
	require.NoError(t, err)
	uuidVal, err := idNode.GetUUIDValue()
	assert.Error(t, err)
	assert.Nil(t, uuidVal, "returned pointer must be nil on error")

	// Whitespace-padded UUID should parse correctly
	validUUIDStr := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	paddedNode, err := NewFormParseNode([]byte("id=%20%20" + validUUIDStr + "%20%20"))
	require.NoError(t, err)
	paddedIdNode, err := paddedNode.GetChildNode("id")
	require.NoError(t, err)
	uuidVal, err = paddedIdNode.GetUUIDValue()
	require.NoError(t, err)
	require.NotNil(t, uuidVal)
	expectedUUID, _ := uuid.Parse(validUUIDStr)
	assert.Equal(t, expectedUUID, *uuidVal)

	// Whitespace-only UUID should return (nil, nil)
	whitespaceNode, err := NewFormParseNode([]byte("id=%20%20%20"))
	require.NoError(t, err)
	wsIdNode, err := whitespaceNode.GetChildNode("id")
	require.NoError(t, err)
	uuidVal, err = wsIdNode.GetUUIDValue()
	require.NoError(t, err)
	assert.Nil(t, uuidVal)
}

func TestTemporalParsersTrimWhitespace(t *testing.T) {
	node, err := NewFormParseNode([]byte("duration=%20PT1H%20&time=%2012%3A30%3A00%20&date=%202026-11-01%20&emptyDuration=%20%20&emptyTime=%20%20&emptyDate=%20%20"))
	require.NoError(t, err)

	durNode, err := node.GetChildNode("duration")
	require.NoError(t, err)
	durVal, err := durNode.GetISODurationValue()
	require.NoError(t, err)
	require.NotNil(t, durVal)

	timeNode, err := node.GetChildNode("time")
	require.NoError(t, err)
	timeVal, err := timeNode.GetTimeOnlyValue()
	require.NoError(t, err)
	require.NotNil(t, timeVal)
	assert.Equal(t, "12:30:00", timeVal.String())

	dateNode, err := node.GetChildNode("date")
	require.NoError(t, err)
	dateVal, err := dateNode.GetDateOnlyValue()
	require.NoError(t, err)
	require.NotNil(t, dateVal)
	assert.Equal(t, "2026-11-01", dateVal.String())

	// Empty / whitespace-only cases must return (nil, nil)
	durNodeEmpty, err := node.GetChildNode("emptyDuration")
	require.NoError(t, err)
	durValEmpty, err := durNodeEmpty.GetISODurationValue()
	require.NoError(t, err)
	assert.Nil(t, durValEmpty)

	timeNodeEmpty, err := node.GetChildNode("emptyTime")
	require.NoError(t, err)
	timeValEmpty, err := timeNodeEmpty.GetTimeOnlyValue()
	require.NoError(t, err)
	assert.Nil(t, timeValEmpty)

	dateNodeEmpty, err := node.GetChildNode("emptyDate")
	require.NoError(t, err)
	dateValEmpty, err := dateNodeEmpty.GetDateOnlyValue()
	require.NoError(t, err)
	assert.Nil(t, dateValEmpty)
}

func TestGetCollectionOfPrimitiveValuesCaseInsensitiveAndWhitespace(t *testing.T) {
	assert.Equal(t,
		[]interface{}{ref("milk"), ref("soda")},
		getCollectionValues("id=2&item=milk&item=soda", "item", "String"),
	)
	assert.Equal(t,
		[]interface{}{ref("milk"), ref("soda")},
		getCollectionValues("id=2&item=milk&item=soda", "item", "STRING"),
	)
	assert.Equal(t,
		[]interface{}{ref(true), ref(false)},
		getCollectionValues("id=2&item=true&item=false", "item", "Bool"),
	)
	assert.Equal(t,
		[]interface{}{ref(int32(1)), ref(int32(2))},
		getCollectionValues("id=2&item=1&item=2", "item", "Int32"),
	)
	assert.Equal(t,
		[]interface{}{ref(float64(1.5)), ref(float64(2.5))},
		getCollectionValues("id=2&item=1.5&item=2.5", "item", "Float64"),
	)

	node, err := NewFormParseNode([]byte("item=milk,soda"))
	assert.NoError(t, err)

	_, err = node.GetCollectionOfPrimitiveValues("")
	assert.EqualError(t, err, "targetType is empty")

	_, err = node.GetCollectionOfPrimitiveValues("   ")
	assert.EqualError(t, err, "targetType is empty")
}


func TestEmptyFieldDeserializationDoesNotPanic(t *testing.T) {
	node, err := NewFormParseNode([]byte("officeLocation=&id=48d31887-5fad-4d73-a9f5-3c356e68a038"))
	require.NoError(t, err)
	result, err := node.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	require.NotNil(t, result)
	entity := result.(internal.TestEntityable)
	assert.Nil(t, entity.GetOfficeLocation())
	assert.Equal(t, "48d31887-5fad-4d73-a9f5-3c356e68a038", *entity.GetId())
}

func TestGetChildNodeForEmptyAndMissingFields(t *testing.T) {
	node, err := NewFormParseNode([]byte("emptyField=&populatedField=hello"))
	require.NoError(t, err)

	// Missing field returns (nil, nil)
	missingChild, err := node.GetChildNode("missingField")
	require.NoError(t, err)
	assert.Nil(t, missingChild)

	// Empty field returns a valid non-nil child node whose getters safely return (nil, nil)
	emptyChild, err := node.GetChildNode("emptyField")
	require.NoError(t, err)
	require.NotNil(t, emptyChild)

	strVal, err := emptyChild.GetStringValue()
	require.NoError(t, err)
	assert.Nil(t, strVal)

	intVal, err := emptyChild.GetInt32Value()
	require.NoError(t, err)
	assert.Nil(t, intVal)

	boolVal, err := emptyChild.GetBoolValue()
	require.NoError(t, err)
	assert.Nil(t, boolVal)

	byteVal, err := emptyChild.GetByteArrayValue()
	require.NoError(t, err)
	assert.Nil(t, byteVal)
}

func TestGetByteArrayValueTrimsWhitespace(t *testing.T) {
	// Base64 for "Hello World" is "SGVsbG8gV29ybGQ="
	node, err := NewFormParseNode([]byte("data=%20%20SGVsbG8gV29ybGQ%3D%0A%20&emptyData=%20%20%0A%0D%20"))
	require.NoError(t, err)

	child, err := node.GetChildNode("data")
	require.NoError(t, err)
	require.NotNil(t, child)

	val, err := child.GetByteArrayValue()
	require.NoError(t, err)
	assert.Equal(t, []byte("Hello World"), val)

	emptyChild, err := node.GetChildNode("emptyData")
	require.NoError(t, err)
	require.NotNil(t, emptyChild)

	emptyVal, err := emptyChild.GetByteArrayValue()
	require.NoError(t, err)
	assert.Nil(t, emptyVal)
}

func TestGetCollectionOfPrimitiveValuesIntAndByteArrayAliases(t *testing.T) {
	// Test "int" alias mapping to int64
	node, err := NewFormParseNode([]byte("numbers=10,20,30"))
	require.NoError(t, err)
	child, err := node.GetChildNode("numbers")
	require.NoError(t, err)
	require.NotNil(t, child)

	intValues, err := child.GetCollectionOfPrimitiveValues("int")
	require.NoError(t, err)
	require.Len(t, intValues, 3)
	assert.Equal(t, int64(10), *intValues[0].(*int64))
	assert.Equal(t, int64(20), *intValues[1].(*int64))
	assert.Equal(t, int64(30), *intValues[2].(*int64))

	// Test "[]byte" and "bytearray" aliases
	// "Hello" -> SGVsbG8=, "World" -> V29ybGQ=
	byteNode, err := NewFormParseNode([]byte("blobs=SGVsbG8%3D,V29ybGQ%3D"))
	require.NoError(t, err)
	blobChild, err := byteNode.GetChildNode("blobs")
	require.NoError(t, err)
	require.NotNil(t, blobChild)

	sliceByteValues, err := blobChild.GetCollectionOfPrimitiveValues("[]byte")
	require.NoError(t, err)
	require.Len(t, sliceByteValues, 2)
	assert.Equal(t, []byte("Hello"), sliceByteValues[0].([]byte))
	assert.Equal(t, []byte("World"), sliceByteValues[1].([]byte))

	byteArrayValues, err := blobChild.GetCollectionOfPrimitiveValues("bytearray")
	require.NoError(t, err)
	require.Len(t, byteArrayValues, 2)
	assert.Equal(t, []byte("Hello"), byteArrayValues[0].([]byte))
	assert.Equal(t, []byte("World"), byteArrayValues[1].([]byte))
}

func TestPrimitiveNumericAndBoolParsersHandleWhitespaceAndEscapes(t *testing.T) {
	// URL-encoded and whitespace-padded primitives
	node, err := NewFormParseNode([]byte("b=%20true%20&i8=%20-120%20&u8=%20250%20&f32=%2012.5%20&f64=%20123.456%20&i32=%2042%20&i64=%209007199254740991%20&emptyB=%20&emptyNum=%20%20"))
	require.NoError(t, err)

	bNode, err := node.GetChildNode("b")
	require.NoError(t, err)
	bVal, err := bNode.GetBoolValue()
	require.NoError(t, err)
	require.NotNil(t, bVal)
	assert.True(t, *bVal)

	i8Node, err := node.GetChildNode("i8")
	require.NoError(t, err)
	i8Val, err := i8Node.GetInt8Value()
	require.NoError(t, err)
	require.NotNil(t, i8Val)
	assert.Equal(t, int8(-120), *i8Val)

	u8Node, err := node.GetChildNode("u8")
	require.NoError(t, err)
	u8Val, err := u8Node.GetByteValue()
	require.NoError(t, err)
	require.NotNil(t, u8Val)
	assert.Equal(t, byte(250), *u8Val)

	f32Node, err := node.GetChildNode("f32")
	require.NoError(t, err)
	f32Val, err := f32Node.GetFloat32Value()
	require.NoError(t, err)
	require.NotNil(t, f32Val)
	assert.Equal(t, float32(12.5), *f32Val)

	f64Node, err := node.GetChildNode("f64")
	require.NoError(t, err)
	f64Val, err := f64Node.GetFloat64Value()
	require.NoError(t, err)
	require.NotNil(t, f64Val)
	assert.Equal(t, 123.456, *f64Val)

	i32Node, err := node.GetChildNode("i32")
	require.NoError(t, err)
	i32Val, err := i32Node.GetInt32Value()
	require.NoError(t, err)
	require.NotNil(t, i32Val)
	assert.Equal(t, int32(42), *i32Val)

	i64Node, err := node.GetChildNode("i64")
	require.NoError(t, err)
	i64Val, err := i64Node.GetInt64Value()
	require.NoError(t, err)
	require.NotNil(t, i64Val)
	assert.Equal(t, int64(9007199254740991), *i64Val)

	// Empty and whitespace-only return (nil, nil)
	emptyBNode, err := node.GetChildNode("emptyB")
	require.NoError(t, err)
	emptyBVal, err := emptyBNode.GetBoolValue()
	require.NoError(t, err)
	assert.Nil(t, emptyBVal)

	emptyNumNode, err := node.GetChildNode("emptyNum")
	require.NoError(t, err)
	emptyI32Val, err := emptyNumNode.GetInt32Value()
	require.NoError(t, err)
	assert.Nil(t, emptyI32Val)

	emptyI64Val, err := emptyNumNode.GetInt64Value()
	require.NoError(t, err)
	assert.Nil(t, emptyI64Val)

	emptyF64Val, err := emptyNumNode.GetFloat64Value()
	require.NoError(t, err)
	assert.Nil(t, emptyF64Val)

	// Comma-separated primitive collection with whitespace between elements
	wsCollectionNode, err := NewFormParseNode([]byte("items=1,%202,%203&flags=true,%20false&floats=1.25,%202.5"))
	require.NoError(t, err)

	itemsChild, err := wsCollectionNode.GetChildNode("items")
	require.NoError(t, err)
	parsedItems, err := itemsChild.GetCollectionOfPrimitiveValues("int32")
	require.NoError(t, err)
	require.Len(t, parsedItems, 3)
	assert.Equal(t, int32(1), *parsedItems[0].(*int32))
	assert.Equal(t, int32(2), *parsedItems[1].(*int32))
	assert.Equal(t, int32(3), *parsedItems[2].(*int32))

	flagsChild, err := wsCollectionNode.GetChildNode("flags")
	require.NoError(t, err)
	parsedFlags, err := flagsChild.GetCollectionOfPrimitiveValues("bool")
	require.NoError(t, err)
	require.Len(t, parsedFlags, 2)
	assert.True(t, *parsedFlags[0].(*bool))
	assert.False(t, *parsedFlags[1].(*bool))

	floatsChild, err := wsCollectionNode.GetChildNode("floats")
	require.NoError(t, err)
	parsedFloats, err := floatsChild.GetCollectionOfPrimitiveValues("float64")
	require.NoError(t, err)
	require.Len(t, parsedFloats, 2)
	assert.Equal(t, 1.25, *parsedFloats[0].(*float64))
	assert.Equal(t, 2.5, *parsedFloats[1].(*float64))
}

func TestFormParseNode_EncodedPositiveExponent(t *testing.T) {
	// value=1e%2B3 should decode to 1000 for int32, int64, and collection variants
	node, err := NewFormParseNode([]byte("val=1e%2B3&coll=1e%2B2,%202e%2B3"))
	require.NoError(t, err)

	valNode, err := node.GetChildNode("val")
	require.NoError(t, err)

	i64Val, err := valNode.GetInt64Value()
	require.NoError(t, err)
	require.NotNil(t, i64Val)
	assert.Equal(t, int64(1000), *i64Val)

	i32Val, err := valNode.GetInt32Value()
	require.NoError(t, err)
	require.NotNil(t, i32Val)
	assert.Equal(t, int32(1000), *i32Val)

	f64Val, err := valNode.GetFloat64Value()
	require.NoError(t, err)
	require.NotNil(t, f64Val)
	assert.Equal(t, float64(1000), *f64Val)

	collNode, err := node.GetChildNode("coll")
	require.NoError(t, err)

	i64Coll, err := collNode.GetCollectionOfPrimitiveValues("int64")
	require.NoError(t, err)
	require.Len(t, i64Coll, 2)
	assert.Equal(t, int64(100), *i64Coll[0].(*int64))
	assert.Equal(t, int64(2000), *i64Coll[1].(*int64))

	i32Coll, err := collNode.GetCollectionOfPrimitiveValues("int32")
	require.NoError(t, err)
	require.Len(t, i32Coll, 2)
	assert.Equal(t, int32(100), *i32Coll[0].(*int32))
	assert.Equal(t, int32(2000), *i32Coll[1].(*int32))
}

func TestFormParseNode_NullMarkerAcrossPrimitiveTypes(t *testing.T) {
	// Form encoded query string with "null" values across various primitive fields
	source := "bVal=null&i8Val=null&byteVal=null&i32Val=null&i64Val=null&f32Val=null&f64Val=null&tVal=null&durVal=null&toVal=null&doVal=null&uuidVal=null&bytesVal=&b64Bytes=null&enumVal=null"
	node, err := NewFormParseNode([]byte(source))
	require.NoError(t, err)

	bNode, err := node.GetChildNode("bVal")
	require.NoError(t, err)
	bVal, err := bNode.GetBoolValue()
	require.NoError(t, err)
	assert.Nil(t, bVal)

	i8Node, err := node.GetChildNode("i8Val")
	require.NoError(t, err)
	i8Val, err := i8Node.GetInt8Value()
	require.NoError(t, err)
	assert.Nil(t, i8Val)

	byteNode, err := node.GetChildNode("byteVal")
	require.NoError(t, err)
	byteVal, err := byteNode.GetByteValue()
	require.NoError(t, err)
	assert.Nil(t, byteVal)

	i32Node, err := node.GetChildNode("i32Val")
	require.NoError(t, err)
	i32Val, err := i32Node.GetInt32Value()
	require.NoError(t, err)
	assert.Nil(t, i32Val)

	i64Node, err := node.GetChildNode("i64Val")
	require.NoError(t, err)
	i64Val, err := i64Node.GetInt64Value()
	require.NoError(t, err)
	assert.Nil(t, i64Val)

	f32Node, err := node.GetChildNode("f32Val")
	require.NoError(t, err)
	f32Val, err := f32Node.GetFloat32Value()
	require.NoError(t, err)
	assert.Nil(t, f32Val)

	f64Node, err := node.GetChildNode("f64Val")
	require.NoError(t, err)
	f64Val, err := f64Node.GetFloat64Value()
	require.NoError(t, err)
	assert.Nil(t, f64Val)

	tNode, err := node.GetChildNode("tVal")
	require.NoError(t, err)
	tVal, err := tNode.GetTimeValue()
	require.NoError(t, err)
	assert.Nil(t, tVal)

	durNode, err := node.GetChildNode("durVal")
	require.NoError(t, err)
	durVal, err := durNode.GetISODurationValue()
	require.NoError(t, err)
	assert.Nil(t, durVal)

	toNode, err := node.GetChildNode("toVal")
	require.NoError(t, err)
	toVal, err := toNode.GetTimeOnlyValue()
	require.NoError(t, err)
	assert.Nil(t, toVal)

	doNode, err := node.GetChildNode("doVal")
	require.NoError(t, err)
	doVal, err := doNode.GetDateOnlyValue()
	require.NoError(t, err)
	assert.Nil(t, doVal)

	uuidNode, err := node.GetChildNode("uuidVal")
	require.NoError(t, err)
	uuidVal, err := uuidNode.GetUUIDValue()
	require.NoError(t, err)
	assert.Nil(t, uuidVal)

	bytesNode, err := node.GetChildNode("bytesVal")
	require.NoError(t, err)
	bytesVal, err := bytesNode.GetByteArrayValue()
	require.NoError(t, err)
	assert.Nil(t, bytesVal)

	b64BytesNode, err := node.GetChildNode("b64Bytes")
	require.NoError(t, err)
	b64BytesVal, err := b64BytesNode.GetByteArrayValue()
	require.NoError(t, err)
	assert.Equal(t, []byte{158, 233, 101}, b64BytesVal)

	enumNode, err := node.GetChildNode("enumVal")
	require.NoError(t, err)
	enumVal, err := enumNode.GetEnumValue(testColorFactory)
	require.NoError(t, err)
	assert.Nil(t, enumVal)
}

func TestFormParseNode_NullMarkerCaseInsensitiveAndCollections(t *testing.T) {
	source := "upper=NULL&coll=1,null,2"
	node, err := NewFormParseNode([]byte(source))
	require.NoError(t, err)

	upperNode, err := node.GetChildNode("upper")
	require.NoError(t, err)
	i32Val, err := upperNode.GetInt32Value()
	require.NoError(t, err)
	assert.Nil(t, i32Val)

	bVal, err := upperNode.GetBoolValue()
	require.NoError(t, err)
	assert.Nil(t, bVal)

	collNode, err := node.GetChildNode("coll")
	require.NoError(t, err)
	parsed, err := collNode.GetCollectionOfPrimitiveValues("int32")
	require.NoError(t, err)
	require.Len(t, parsed, 3)
	assert.Equal(t, int32(1), *parsed[0].(*int32))
	assert.Nil(t, parsed[1])
	assert.True(t, parsed[1] == nil)
	assert.Equal(t, int32(2), *parsed[2].(*int32))
}

func TestFormParseNode_Int8AndByteExponentAndBounds(t *testing.T) {
	// Valid int8 exponents
	node, err := NewFormParseNode([]byte("pos=1e2&neg=-1e2&max=127&min=-128"))
	require.NoError(t, err)

	posNode, err := node.GetChildNode("pos")
	require.NoError(t, err)
	posVal, err := posNode.GetInt8Value()
	require.NoError(t, err)
	require.NotNil(t, posVal)
	assert.Equal(t, int8(100), *posVal)

	negNode, err := node.GetChildNode("neg")
	require.NoError(t, err)
	negVal, err := negNode.GetInt8Value()
	require.NoError(t, err)
	require.NotNil(t, negVal)
	assert.Equal(t, int8(-100), *negVal)

	// Int8 boundary overflow/underflow
	overflowNode, err := NewFormParseNode([]byte("over=128&under=-129&large=1e3"))
	require.NoError(t, err)

	overNode, err := overflowNode.GetChildNode("over")
	require.NoError(t, err)
	_, err = overNode.GetInt8Value()
	assert.Error(t, err)

	underNode, err := overflowNode.GetChildNode("under")
	require.NoError(t, err)
	_, err = underNode.GetInt8Value()
	assert.Error(t, err)

	largeNode, err := overflowNode.GetChildNode("large")
	require.NoError(t, err)
	_, err = largeNode.GetInt8Value()
	assert.Error(t, err)

	// Valid byte exponents and boundary checks
	byteNode, err := NewFormParseNode([]byte("bExp=1e2&bFloatExp=2.5e2&bMax=255&bOver=256&bNeg=-1&bLarge=1e3"))
	require.NoError(t, err)

	bExpChild, err := byteNode.GetChildNode("bExp")
	require.NoError(t, err)
	bExpVal, err := bExpChild.GetByteValue()
	require.NoError(t, err)
	require.NotNil(t, bExpVal)
	assert.Equal(t, byte(100), *bExpVal)

	bFloatChild, err := byteNode.GetChildNode("bFloatExp")
	require.NoError(t, err)
	bFloatVal, err := bFloatChild.GetByteValue()
	require.NoError(t, err)
	require.NotNil(t, bFloatVal)
	assert.Equal(t, byte(250), *bFloatVal)

	bMaxChild, err := byteNode.GetChildNode("bMax")
	require.NoError(t, err)
	bMaxVal, err := bMaxChild.GetByteValue()
	require.NoError(t, err)
	require.NotNil(t, bMaxVal)
	assert.Equal(t, byte(255), *bMaxVal)

	bOverChild, err := byteNode.GetChildNode("bOver")
	require.NoError(t, err)
	_, err = bOverChild.GetByteValue()
	assert.Error(t, err)

	bNegChild, err := byteNode.GetChildNode("bNeg")
	require.NoError(t, err)
	_, err = bNegChild.GetByteValue()
	assert.Error(t, err)

	bLargeChild, err := byteNode.GetChildNode("bLarge")
	require.NoError(t, err)
	_, err = bLargeChild.GetByteValue()
	assert.Error(t, err)
}

func TestFormParseNode_PrimitiveCollectionSchemaAliases(t *testing.T) {
	source := "bools=true,false&doubles=1.5,2.5&floats=3.5,4.5&ints=10,20&bins=AQID,BAUG&durs=PT1H,PT2M&times=2022-01-01T00:00:00Z,2022-01-02T00:00:00Z"
	node, err := NewFormParseNode([]byte(source))
	require.NoError(t, err)

	boolsChild, err := node.GetChildNode("bools")
	require.NoError(t, err)
	bools, err := boolsChild.GetCollectionOfPrimitiveValues("boolean")
	require.NoError(t, err)
	require.Len(t, bools, 2)
	assert.True(t, *bools[0].(*bool))
	assert.False(t, *bools[1].(*bool))

	doublesChild, err := node.GetChildNode("doubles")
	require.NoError(t, err)
	doubles, err := doublesChild.GetCollectionOfPrimitiveValues("double")
	require.NoError(t, err)
	require.Len(t, doubles, 2)
	assert.Equal(t, 1.5, *doubles[0].(*float64))
	assert.Equal(t, 2.5, *doubles[1].(*float64))

	floatsChild, err := node.GetChildNode("floats")
	require.NoError(t, err)
	floats, err := floatsChild.GetCollectionOfPrimitiveValues("float")
	require.NoError(t, err)
	require.Len(t, floats, 2)
	assert.Equal(t, float32(3.5), *floats[0].(*float32))
	assert.Equal(t, float32(4.5), *floats[1].(*float32))

	intsChild, err := node.GetChildNode("ints")
	require.NoError(t, err)
	ints, err := intsChild.GetCollectionOfPrimitiveValues("integer")
	require.NoError(t, err)
	require.Len(t, ints, 2)
	assert.Equal(t, int64(10), *ints[0].(*int64))
	assert.Equal(t, int64(20), *ints[1].(*int64))

	binsChild, err := node.GetChildNode("bins")
	require.NoError(t, err)
	bins, err := binsChild.GetCollectionOfPrimitiveValues("binary")
	require.NoError(t, err)
	require.Len(t, bins, 2)
	assert.Equal(t, []byte{1, 2, 3}, bins[0].([]byte))
	assert.Equal(t, []byte{4, 5, 6}, bins[1].([]byte))

	dursChild, err := node.GetChildNode("durs")
	require.NoError(t, err)
	durs, err := dursChild.GetCollectionOfPrimitiveValues("duration")
	require.NoError(t, err)
	require.Len(t, durs, 2)

	timesChild, err := node.GetChildNode("times")
	require.NoError(t, err)
	timesDt, err := timesChild.GetCollectionOfPrimitiveValues("datetime")
	require.NoError(t, err)
	require.Len(t, timesDt, 2)

	timesDash, err := timesChild.GetCollectionOfPrimitiveValues("date-time")
	require.NoError(t, err)
	require.Len(t, timesDash, 2)
}

func TestFormParseNode_ByteArrayValidBase64CasePayloads(t *testing.T) {
	node, err := NewFormParseNode([]byte("upper=NULL&mixed=Null&lower=null"))
	require.NoError(t, err)

	upperChild, err := node.GetChildNode("upper")
	require.NoError(t, err)
	upperBytes, err := upperChild.GetByteArrayValue()
	require.NoError(t, err)
	assert.Equal(t, []byte{53, 66, 203}, upperBytes)

	mixedChild, err := node.GetChildNode("mixed")
	require.NoError(t, err)
	mixedBytes, err := mixedChild.GetByteArrayValue()
	require.NoError(t, err)
	assert.Equal(t, []byte{54, 233, 101}, mixedBytes)

	lowerChild, err := node.GetChildNode("lower")
	require.NoError(t, err)
	lowerBytes, err := lowerChild.GetByteArrayValue()
	require.NoError(t, err)
	assert.Equal(t, []byte{158, 233, 101}, lowerBytes)
}

func TestFormParseNode_GetEnumValue_PreservesWhitespace(t *testing.T) {
	received := ""
	customFactory := func(val string) (interface{}, error) {
		received = val
		return val, nil
	}

	node, err := NewFormParseNode([]byte("status=%20ready%20"))
	require.NoError(t, err)
	child, err := node.GetChildNode("status")
	require.NoError(t, err)

	val, err := child.GetEnumValue(customFactory)
	require.NoError(t, err)
	assert.Equal(t, " ready ", val)
	assert.Equal(t, " ready ", received)

	// Ensure null enum still returns (nil, nil)
	nullNode, err := NewFormParseNode([]byte("status=null"))
	require.NoError(t, err)
	nullChild, err := nullNode.GetChildNode("status")
	require.NoError(t, err)

	nullVal, err := nullChild.GetEnumValue(customFactory)
	require.NoError(t, err)
	assert.Nil(t, nullVal)
}

func TestFormParseNode_GetObjectValue_NullAndWhitespace(t *testing.T) {
	// Root null node
	rootNull, err := NewFormParseNode([]byte("null"))
	require.NoError(t, err)
	obj, err := rootNull.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, obj)

	// Root case-insensitive null node
	rootUpperNull, err := NewFormParseNode([]byte("NULL"))
	require.NoError(t, err)
	obj, err = rootUpperNull.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, obj)

	// Root whitespace node
	rootWhitespace, err := NewFormParseNode([]byte("   "))
	require.NoError(t, err)
	obj, err = rootWhitespace.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, obj)

	// Child null node
	payload := "entity=null&upperEntity=NULL&wsEntity=%20%20&encodedNull=%6e%75%6c%6c&plusWs=++"
	node, err := NewFormParseNode([]byte(payload))
	require.NoError(t, err)

	childEntity, err := node.GetChildNode("entity")
	require.NoError(t, err)
	obj, err = childEntity.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, obj)

	childUpper, err := node.GetChildNode("upperEntity")
	require.NoError(t, err)
	obj, err = childUpper.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, obj)

	childWs, err := node.GetChildNode("wsEntity")
	require.NoError(t, err)
	obj, err = childWs.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, obj)

	childEncodedNull, err := node.GetChildNode("encodedNull")
	require.NoError(t, err)
	obj, err = childEncodedNull.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, obj)

	childPlusWs, err := node.GetChildNode("plusWs")
	require.NoError(t, err)
	obj, err = childPlusWs.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, obj)
}

func TestFormParseNode_CollectionGetters_NullAndWhitespace(t *testing.T) {
	dummyParser := func(val string) (interface{}, error) {
		return val, nil
	}

	// Root null nodes return (nil, nil)
	rootNull, err := NewFormParseNode([]byte("null"))
	require.NoError(t, err)
	primColl, err := rootNull.GetCollectionOfPrimitiveValues("int32")
	require.NoError(t, err)
	assert.Nil(t, primColl)

	enumColl, err := rootNull.GetCollectionOfEnumValues(dummyParser)
	require.NoError(t, err)
	assert.Nil(t, enumColl)

	// Root whitespace nodes return (nil, nil)
	rootWs, err := NewFormParseNode([]byte("   "))
	require.NoError(t, err)
	primColl, err = rootWs.GetCollectionOfPrimitiveValues("string")
	require.NoError(t, err)
	assert.Nil(t, primColl)

	enumColl, err = rootWs.GetCollectionOfEnumValues(dummyParser)
	require.NoError(t, err)
	assert.Nil(t, enumColl)

	// Child null collections return (nil, nil)
	payload := "numbers=null&upperRoles=NULL&wsItems=%20%20&mixed=10,null,20"
	node, err := NewFormParseNode([]byte(payload))
	require.NoError(t, err)

	numChild, err := node.GetChildNode("numbers")
	require.NoError(t, err)
	primColl, err = numChild.GetCollectionOfPrimitiveValues("int32")
	require.NoError(t, err)
	assert.Nil(t, primColl)

	rolesChild, err := node.GetChildNode("upperRoles")
	require.NoError(t, err)
	enumColl, err = rolesChild.GetCollectionOfEnumValues(dummyParser)
	require.NoError(t, err)
	assert.Nil(t, enumColl)

	wsChild, err := node.GetChildNode("wsItems")
	require.NoError(t, err)
	primColl, err = wsChild.GetCollectionOfPrimitiveValues("string")
	require.NoError(t, err)
	assert.Nil(t, primColl)

	// Mixed collections still preserve non-null elements
	mixedChild, err := node.GetChildNode("mixed")
	require.NoError(t, err)
	mixedColl, err := mixedChild.GetCollectionOfPrimitiveValues("int32")
	require.NoError(t, err)
	require.Len(t, mixedColl, 3)
	assert.Equal(t, int32(10), *mixedColl[0].(*int32))
	assert.Nil(t, mixedColl[1])
	assert.Equal(t, int32(20), *mixedColl[2].(*int32))
}

func TestFormParseNode_CollectionOfObjectsAndRawValue_NullAndWhitespace(t *testing.T) {
	// Root null node for object collections
	rootNull, err := NewFormParseNode([]byte("null"))
	require.NoError(t, err)
	objColl, err := rootNull.GetCollectionOfObjectValues(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, objColl)

	rawVal, err := rootNull.GetRawValue()
	require.NoError(t, err)
	assert.Nil(t, rawVal)

	// Root case-insensitive null node
	rootUpper, err := NewFormParseNode([]byte("NULL"))
	require.NoError(t, err)
	objColl, err = rootUpper.GetCollectionOfObjectValues(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, objColl)

	rawVal, err = rootUpper.GetRawValue()
	require.NoError(t, err)
	assert.Nil(t, rawVal)

	// Root whitespace node
	rootWs, err := NewFormParseNode([]byte("   "))
	require.NoError(t, err)
	objColl, err = rootWs.GetCollectionOfObjectValues(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, objColl)

	rawVal, err = rootWs.GetRawValue()
	require.NoError(t, err)
	assert.Nil(t, rawVal)

	// Child nodes
	payload := "users=null&upperUsers=NULL&wsNode=%20%20&activeUser=john"
	node, err := NewFormParseNode([]byte(payload))
	require.NoError(t, err)

	nullUsers, err := node.GetChildNode("users")
	require.NoError(t, err)
	objColl, err = nullUsers.GetCollectionOfObjectValues(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, objColl)

	rawVal, err = nullUsers.GetRawValue()
	require.NoError(t, err)
	assert.Nil(t, rawVal)

	upperUsers, err := node.GetChildNode("upperUsers")
	require.NoError(t, err)
	objColl, err = upperUsers.GetCollectionOfObjectValues(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	assert.Nil(t, objColl)

	wsChild, err := node.GetChildNode("wsNode")
	require.NoError(t, err)
	rawVal, err = wsChild.GetRawValue()
	require.NoError(t, err)
	assert.Nil(t, rawVal)

	// Non-null non-empty node returns raw value and rejects collection of objects
	activeUser, err := node.GetChildNode("activeUser")
	require.NoError(t, err)
	rawVal, err = activeUser.GetRawValue()
	require.NoError(t, err)
	require.NotNil(t, rawVal)
	assert.Equal(t, "john", *rawVal.(*string))

	_, err = activeUser.GetCollectionOfObjectValues(internal.CreateTestEntityFromDiscriminator)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "collections are not supported in form serialization")
}

func TestFormParseNode_KeySanitizationAndIndexValidation(t *testing.T) {
	node, err := NewFormParseNode([]byte("user=alice"))
	require.NoError(t, err)

	// Whitespace-only and encoded whitespace indices must return index is empty error
	for _, wsIndex := range []string{"", " ", "   ", "\t", "\r\n", "%20", "%09", "+", "%0D%0A", "  %20  ", "%20%20"} {
		child, err := node.GetChildNode(wsIndex)
		assert.Nil(t, child)
		assert.EqualError(t, err, "index is empty")
	}

	// Keys with tabs, newlines, and spaces must be properly trimmed during field loading
	whitespaceKeyPayload := "%09user%09=alice&%0D%0Arole%0D%0A=admin&%20tag%20=golang"
	wsNode, err := NewFormParseNode([]byte(whitespaceKeyPayload))
	require.NoError(t, err)

	userChild, err := wsNode.GetChildNode("user")
	require.NoError(t, err)
	require.NotNil(t, userChild)
	userVal, err := userChild.GetStringValue()
	require.NoError(t, err)
	assert.Equal(t, "alice", *userVal)

	roleChild, err := wsNode.GetChildNode("role")
	require.NoError(t, err)
	require.NotNil(t, roleChild)
	roleVal, err := roleChild.GetStringValue()
	require.NoError(t, err)
	assert.Equal(t, "admin", *roleVal)

	tagChild, err := wsNode.GetChildNode("tag")
	require.NoError(t, err)
	require.NotNil(t, tagChild)
	tagVal, err := tagChild.GetStringValue()
	require.NoError(t, err)
	assert.Equal(t, "golang", *tagVal)
}

func TestFormParseNode_GetChildNode_LiteralUnescapedKeyCharacters(t *testing.T) {
	// Form payload containing an escaped key that decodes to a literal percent sign
	payload := "discount%25=15&rate%25=3.5"
	node, err := NewFormParseNode([]byte(payload))
	require.NoError(t, err)

	// Direct lookup with unescaped key containing '%' should succeed without URL escape error
	child, err := node.GetChildNode("discount%")
	require.NoError(t, err)
	require.NotNil(t, child)
	val, err := child.GetInt32Value()
	require.NoError(t, err)
	require.NotNil(t, val)
	assert.Equal(t, int32(15), *val)

	// Escaped lookup should also succeed
	childEscaped, err := node.GetChildNode("discount%25")
	require.NoError(t, err)
	require.NotNil(t, childEscaped)
	valEscaped, err := childEscaped.GetInt32Value()
	require.NoError(t, err)
	require.NotNil(t, valEscaped)
	assert.Equal(t, int32(15), *valEscaped)

	// Non-existent key with invalid escape syntax should still return error
	invalidChild, err := node.GetChildNode("nonexistent%invalid")
	require.Error(t, err)
	assert.Nil(t, invalidChild)
}

