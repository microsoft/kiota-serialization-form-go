package formserialization

import (
	"fmt"
	"github.com/microsoft/kiota-serialization-form-go/internal"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

func referenceTime() (value time.Time) {
	value, _ = time.Parse(time.Layout, time.Layout)
	return
}

func TestItDoesntWriteAnythingForNilAdditionalData(t *testing.T) {
	serializer := NewFormSerializationWriter()
	serializer.WriteAdditionalData(nil)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestItDoesntWriteAnythingForEmptyAdditionalData(t *testing.T) {
	serializer := NewFormSerializationWriter()
	serializer.WriteAdditionalData(make(map[string]interface{}))
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestItTrimsTheLastAmp(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := "value"
	serializer.WriteStringValue("key", &value)
	serializer.WriteAdditionalData(make(map[string]interface{}))
	value2 := "value2"
	serializer.WriteStringValue("key2", &value2)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, "key=value&key2=value2", string(result[:]))
}

func TestWriteTimeValue(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := referenceTime()
	serializer.WriteTimeValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, "key=2006-01-02T15%3A04%3A05-07%3A00", string(result[:]))
}

func TestWriteISODurationValue(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := absser.NewDuration(1, 0, 2, 3, 4, 5, 6)
	serializer.WriteISODurationValue("key", value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("key=%v", value), string(result[:]))
}

func TestWriteTimeOnlyValue(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := absser.NewTimeOnly(referenceTime())
	serializer.WriteTimeOnlyValue("key", value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, "key=15%3A04%3A05.000000000", string(result[:]))
}

func TestWriteDateOnlyValue(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := absser.NewDateOnly(referenceTime())
	serializer.WriteDateOnlyValue("key", value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("key=%v", value), string(result[:]))
}

func TestWriteBoolValue(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := true
	serializer.WriteBoolValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("key=%t", value), string(result[:]))
}

func TestWriteInt8Value(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := int8(125)
	serializer.WriteInt8Value("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("key=%d", value), string(result[:]))
}

func TestWriteByteValue(t *testing.T) {
	serializer := NewFormSerializationWriter()
	var value byte = 97
	serializer.WriteByteValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("key=%d", value), string(result[:]))
}

// ByteArray values are encoded to Base64 when stored
func TestWriteByteArrayValue(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := []byte("SerialWriter")
	serializer.WriteByteArrayValue("key", value)
	expected := "U2VyaWFsV3JpdGVy"
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("key=%s", expected), string(result[:]))
}

func TestBufferClose(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := "W/\"CQAAABYAAAAs+XSiyjZdS4Rhtwk0v1pGAAC5bsJ2\""
	serializer.WriteStringValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.True(t, len(result) > 0)
	serializer.Close()
	assert.True(t, len(result) > 0)
	empty, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.True(t, len(empty) == 0)
	dateOnly := absser.NewDateOnly(referenceTime())
	serializer.WriteDateOnlyValue("today", dateOnly)
	notEmpty, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.True(t, len(notEmpty) > 0)
}

func TestFormSerializationWriterHonoursInterface(t *testing.T) {
	instance := NewFormSerializationWriter()
	assert.Implements(t, (*absser.SerializationWriter)(nil), instance)
}

func TestWriteMultipleTypes(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := "value"
	serializer.WriteStringValue("key", &value)
	pointer := "pointer"
	adlData := map[string]interface{}{
		"add1": "string",
		"add2": &pointer,
	}
	serializer.WriteAdditionalData(adlData)
	value2 := "value2"
	serializer.WriteStringValue("key2", &value2)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Contains(t, string(result[:]), "key=value&")
	assert.Contains(t, string(result[:]), "add1=string&")
	assert.Contains(t, string(result[:]), "add2=pointer&")
	assert.Contains(t, string(result[:]), "key2=value2")
	assert.Equal(t, len("key=value&add1=string&add2=pointer&key2=value2"), len(string(result[:])))
}

func TestEscapesNewLinesInStrings(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := "value\nwith\nnew\nlines"
	serializer.WriteStringValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, "key=value%0Awith%0Anew%0Alines", string(result[:]))
}

func TestJsonSerializationWriter_WriteNullValue(t *testing.T) {
	serializer := NewFormSerializationWriter()

	err := serializer.WriteNullValue("name")
	assert.NoError(t, err)
	result, err := serializer.GetSerializedContent()
	assert.NoError(t, err)
	converted := string(result)

	assert.Equal(t, "name=null", converted)
}

func TestJsonSerializationWriter(t *testing.T) {
	serializer := NewFormSerializationWriter()
	countBefore := 0
	onBefore := func(parsable absser.Parsable) error {
		countBefore++
		return nil
	}
	err := serializer.SetOnBeforeSerialization(onBefore)
	assert.NoError(t, err)

	countAfter := 0
	onAfter := func(parsable absser.Parsable) error {
		countAfter++
		return nil
	}
	err = serializer.SetOnAfterObjectSerialization(onAfter)
	assert.NoError(t, err)

	countStart := 0
	onStart := func(absser.Parsable, absser.SerializationWriter) error {
		countStart++
		return nil
	}

	err = serializer.SetOnStartObjectSerialization(onStart)
	assert.NoError(t, err)

	assert.Equal(t, 0, countBefore)
	assert.Equal(t, 0, countAfter)
	assert.Equal(t, 0, countStart)

	test := internal.NewTestEntity()
	err = serializer.WriteObjectValue("name", test)
	assert.NoError(t, err)

	assert.Equal(t, 1, countBefore)
	assert.Equal(t, 1, countAfter)
	assert.Equal(t, 1, countStart)
}
