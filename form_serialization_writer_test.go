package formserialization

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/kiota-serialization-form-go/internal"

	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	assert.Equal(t, "key=15%3A04%3A05", string(result[:]))
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

func TestWriteMultipleCollections(t *testing.T) {
	serializer := NewFormSerializationWriter()
	value := "value"
	serializer.WriteStringValue("key", &value)

	adlData := map[string]interface{}{
		"string":      []string{"str1", "str2"},
		"int32":       []int32{34, 56},
		"bool":        []bool{true, false},
		"uint8":       []uint8{1, 2},
		"float32":     []float32{32.0, 34.1},
		"float64":     []float64{100.1, 400.1},
		"int64":       []int64{567, 765},
		"time":        []time.Time{time.Now()},
		"timeonly":    []*absser.ISODuration{absser.NewDuration(0, 0, 1, 0, 0, 0, 0)},
		"dateonly":    []*absser.DateOnly{absser.NewDateOnly(time.Now())},
		"isoduration": []*absser.ISODuration{absser.NewDuration(0, 0, 1, 0, 0, 0, 0)},
		"uuid":        []uuid.UUID{uuid.New()},
		"base64":      []int64{567, 765},
		"int8":        []int8{2, 3},
	}
	serializer.WriteAdditionalData(adlData)

	result, err := serializer.GetSerializedContent()

	assert.Nil(t, err)
	resultVal := string(result[:])
	assert.Contains(t, resultVal, "key=value&")
	assert.Contains(t, resultVal, "string=str1&")
	assert.Contains(t, resultVal, "string=str2")
	assert.Contains(t, resultVal, "int32=34")
	assert.Contains(t, resultVal, "int32=56")
	assert.Contains(t, resultVal, "bool=true")
	assert.Contains(t, resultVal, "bool=false")
	assert.Contains(t, resultVal, "uint8=1")
	assert.Contains(t, resultVal, "uint8=2")
	assert.Contains(t, resultVal, "float32=32")
	assert.Contains(t, resultVal, "float32=34")
	assert.Contains(t, resultVal, "float64=100.1")
	assert.Contains(t, resultVal, "float64=400.1")
	assert.Contains(t, resultVal, "int64=567")
	assert.Contains(t, resultVal, "int64=765")
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
	err = serializer.WriteObjectValue("", test)
	assert.NoError(t, err)

	assert.Equal(t, 1, countBefore)
	assert.Equal(t, 1, countAfter)
	assert.Equal(t, 1, countStart)
}

func TestFormSerializationWriterNilGuards(t *testing.T) {
	var nilWriter *FormSerializationWriter
	assert.NotPanics(t, func() {
		content, err := nilWriter.GetSerializedContent()
		assert.EqualError(t, err, "writer cannot be nil")
		assert.Nil(t, content)

		err = nilWriter.Close()
		assert.NoError(t, err)

		strVal := "test"
		assert.EqualError(t, nilWriter.WriteStringValue("k", &strVal), "writer cannot be nil")

		boolVal := true
		assert.EqualError(t, nilWriter.WriteBoolValue("k", &boolVal), "writer cannot be nil")

		byteVal := byte(42)
		assert.EqualError(t, nilWriter.WriteByteValue("k", &byteVal), "writer cannot be nil")

		int8Val := int8(8)
		assert.EqualError(t, nilWriter.WriteInt8Value("k", &int8Val), "writer cannot be nil")

		int32Val := int32(32)
		assert.EqualError(t, nilWriter.WriteInt32Value("k", &int32Val), "writer cannot be nil")

		int64Val := int64(64)
		assert.EqualError(t, nilWriter.WriteInt64Value("k", &int64Val), "writer cannot be nil")

		float32Val := float32(3.14)
		assert.EqualError(t, nilWriter.WriteFloat32Value("k", &float32Val), "writer cannot be nil")

		float64Val := float64(6.28)
		assert.EqualError(t, nilWriter.WriteFloat64Value("k", &float64Val), "writer cannot be nil")

		timeVal := time.Now()
		assert.EqualError(t, nilWriter.WriteTimeValue("k", &timeVal), "writer cannot be nil")

		durationVal, _ := absser.ParseISODuration("PT1H")
		assert.EqualError(t, nilWriter.WriteISODurationValue("k", durationVal), "writer cannot be nil")

		timeOnlyVal, _ := absser.ParseTimeOnly("12:30:00")
		assert.EqualError(t, nilWriter.WriteTimeOnlyValue("k", timeOnlyVal), "writer cannot be nil")

		dateOnlyVal, _ := absser.ParseDateOnly("2026-10-23")
		assert.EqualError(t, nilWriter.WriteDateOnlyValue("k", dateOnlyVal), "writer cannot be nil")

		uuidVal := uuid.New()
		assert.EqualError(t, nilWriter.WriteUUIDValue("k", &uuidVal), "writer cannot be nil")

		assert.EqualError(t, nilWriter.WriteByteArrayValue("k", []byte{1, 2, 3}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteNullValue("k"), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteObjectValue("k", nil), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfObjectValues("k", nil), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteAnyValue("k", "any"), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteAdditionalData(map[string]interface{}{"k": "v"}), "writer cannot be nil")

		assert.EqualError(t, nilWriter.WriteCollectionOfStringValues("k", []string{"a"}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfBoolValues("k", []bool{true}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfByteValues("k", []byte{1}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfInt8Values("k", []int8{1}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfInt32Values("k", []int32{1}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfInt64Values("k", []int64{1}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfFloat32Values("k", []float32{1.0}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfFloat64Values("k", []float64{1.0}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfTimeValues("k", []time.Time{timeVal}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfISODurationValues("k", []absser.ISODuration{*durationVal}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfTimeOnlyValues("k", []absser.TimeOnly{*timeOnlyVal}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfDateOnlyValues("k", []absser.DateOnly{*dateOnlyVal}), "writer cannot be nil")
		assert.EqualError(t, nilWriter.WriteCollectionOfUUIDValues("k", []uuid.UUID{uuidVal}), "writer cannot be nil")

		assert.Nil(t, nilWriter.GetOnBeforeSerialization())
		assert.EqualError(t, nilWriter.SetOnBeforeSerialization(func(absser.Parsable) error { return nil }), "writer cannot be nil")
		assert.Nil(t, nilWriter.GetOnAfterObjectSerialization())
		assert.EqualError(t, nilWriter.SetOnAfterObjectSerialization(func(absser.Parsable) error { return nil }), "writer cannot be nil")
		assert.Nil(t, nilWriter.GetOnStartObjectSerialization())
		assert.EqualError(t, nilWriter.SetOnStartObjectSerialization(func(absser.Parsable, absser.SerializationWriter) error { return nil }), "writer cannot be nil")
	})
}

type mockTestParsable struct {
	serializeFunc func(absser.SerializationWriter) error
}

func (m *mockTestParsable) Serialize(writer absser.SerializationWriter) error {
	if m.serializeFunc != nil {
		return m.serializeFunc(writer)
	}
	return nil
}

func (m *mockTestParsable) GetFieldDeserializers() map[string]func(absser.ParseNode) error {
	return make(map[string]func(absser.ParseNode) error)
}

func TestWriteObjectValueDepthLifecycleAndReusability(t *testing.T) {
	serializer := NewFormSerializationWriter()

	// 1. First top-level object write
	entity1 := internal.NewTestEntity()
	id1 := "id-1"
	entity1.SetId(&id1)
	err := serializer.WriteObjectValue("", entity1)
	require.NoError(t, err)

	// 2. Second top-level object write on the same writer without Close()
	entity2 := internal.NewTestEntity()
	id2 := "id-2"
	entity2.SetId(&id2)
	err = serializer.WriteObjectValue("", entity2)
	require.NoError(t, err)

	// 3. Nested object serialization attempt inside Serialize() must fail
	var nestedErr error
	nestedEntity := &mockTestParsable{
		serializeFunc: func(w absser.SerializationWriter) error {
			child := internal.NewTestEntity()
			nestedErr = w.WriteObjectValue("child", child)
			return nestedErr
		},
	}
	err = serializer.WriteObjectValue("", nestedEntity)
	assert.Error(t, err)
	assert.Error(t, nestedErr)
	assert.Contains(t, nestedErr.Error(), "nested objects serialization is not supported")

	// 4. After nested error, subsequent top-level writes must still succeed
	entity3 := internal.NewTestEntity()
	id3 := "id-3"
	entity3.SetId(&id3)
	err = serializer.WriteObjectValue("", entity3)
	require.NoError(t, err)

	// 5. Close() must reset writer state and allow new writes
	err = serializer.Close()
	require.NoError(t, err)
	err = serializer.WriteObjectValue("", entity1)
	require.NoError(t, err)
}

func TestWriteObjectValueSerializationErrorDoesNotInvokeOnAfter(t *testing.T) {
	serializer := NewFormSerializationWriter()

	countAfter := 0
	err := serializer.SetOnAfterObjectSerialization(func(parsable absser.Parsable) error {
		countAfter++
		return nil
	})
	require.NoError(t, err)

	expectedErr := errors.New("failed serialization")
	failingEntity := &mockTestParsable{
		serializeFunc: func(w absser.SerializationWriter) error {
			return expectedErr
		},
	}

	err = serializer.WriteObjectValue("", failingEntity)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 0, countAfter, "OnAfterObjectSerialization must not be called when Serialize fails")

	// Verify writer remains usable after serialization error
	entity := internal.NewTestEntity()
	id := "test-id"
	entity.SetId(&id)
	err = serializer.WriteObjectValue("", entity)
	require.NoError(t, err)
	assert.Equal(t, 1, countAfter, "OnAfterObjectSerialization should be called on successful write")
}

func TestWriteAdditionalDataErrorPropagationAndIntSupport(t *testing.T) {
	serializer := NewFormSerializationWriter()

	// Test Go int, *int, and []int serialization
	intVal := 42
	intSlice := []int{10, 20}
	additionalData := map[string]interface{}{
		"plainInt": intVal,
		"ptrInt":   &intVal,
		"sliceInt": intSlice,
	}
	err := serializer.WriteAdditionalData(additionalData)
	require.NoError(t, err)

	content, err := serializer.GetSerializedContent()
	require.NoError(t, err)
	contentStr := string(content)
	assert.Contains(t, contentStr, "plainInt=42")
	assert.Contains(t, contentStr, "ptrInt=42")
	assert.Contains(t, contentStr, "sliceInt=10")
	assert.Contains(t, contentStr, "sliceInt=20")

	// Test error propagation: error must not be swallowed by subsequent map entries
	errSerializer := NewFormSerializationWriter()
	failErr := errors.New("additional data parsable failed")
	failingItem := &mockTestParsable{
		serializeFunc: func(w absser.SerializationWriter) error {
			return failErr
		},
	}
	errData := map[string]interface{}{
		"a_fail": failingItem,
		"z_str":  "must_not_mask_error",
	}
	err = errSerializer.WriteAdditionalData(errData)
	assert.Error(t, err)
}

func TestWriteAnyValueAndCollectionOfObjectValuesNilSucceeds(t *testing.T) {
	serializer := NewFormSerializationWriter()
	defer serializer.Close()

	// Nil value should safely succeed without error
	err := serializer.WriteAnyValue("key", nil)
	assert.NoError(t, err)

	// Non-nil value should return unsupported error
	err = serializer.WriteAnyValue("key", "test")
	assert.EqualError(t, err, "serialization of any value is not supported with FormSerializationWriter")

	// Nil collection should safely succeed without error
	err = serializer.WriteCollectionOfObjectValues("key", nil)
	assert.NoError(t, err)

	// Non-nil collection should return unsupported error
	err = serializer.WriteCollectionOfObjectValues("key", []absser.Parsable{})
	assert.EqualError(t, err, "collections serialization is not supported with FormSerializationWriter")
}

func TestWriteObjectValueWithNilElementsAndCallbacks(t *testing.T) {
	// 1. Calling WriteObjectValue with item == nil and all additionalValues nil should not write property or panic
	serializer := NewFormSerializationWriter()
	err := serializer.WriteObjectValue("emptyProp", nil, nil, nil)
	require.NoError(t, err)
	content, err := serializer.GetSerializedContent()
	require.NoError(t, err)
	assert.Empty(t, content, "must not emit property when all objects are nil")

	// 2. Calling WriteObjectValue with item == nil but valid additional object should serialize valid object
	serializer = NewFormSerializationWriter()
	var beforeCalled []string
	var afterCalled []string
	err = serializer.SetOnBeforeSerialization(func(parsable absser.Parsable) error {
		require.NotNil(t, parsable, "OnBeforeSerialization must never receive nil parsable")
		if entity, ok := parsable.(internal.TestEntityable); ok && entity.GetId() != nil {
			beforeCalled = append(beforeCalled, *entity.GetId())
		}
		return nil
	})
	require.NoError(t, err)
	err = serializer.SetOnAfterObjectSerialization(func(parsable absser.Parsable) error {
		require.NotNil(t, parsable, "OnAfterObjectSerialization must never receive nil parsable")
		if entity, ok := parsable.(internal.TestEntityable); ok && entity.GetId() != nil {
			afterCalled = append(afterCalled, *entity.GetId())
		}
		return nil
	})
	require.NoError(t, err)

	entity1 := internal.NewTestEntity()
	id1 := "id-1"
	entity1.SetId(&id1)

	entity2 := internal.NewTestEntity()
	id2 := "id-2"
	entity2.SetId(&id2)

	// Pass nil item, valid entity1, nil additional, and valid entity2
	err = serializer.WriteObjectValue("", nil, entity1, nil, entity2)
	require.NoError(t, err)

	content, err = serializer.GetSerializedContent()
	require.NoError(t, err)
	contentStr := string(content)
	assert.Contains(t, contentStr, "id=id-1")
	assert.Contains(t, contentStr, "id=id-2")
	assert.Equal(t, []string{"id-1", "id-2"}, beforeCalled)
	assert.Equal(t, []string{"id-1", "id-2"}, afterCalled)

	// 3. Calling WriteObjectValue with typed-nil parsables should not invoke callbacks or panic
	var typedNilEntity *internal.TestEntity
	serializer = NewFormSerializationWriter()
	err = serializer.SetOnBeforeSerialization(func(parsable absser.Parsable) error {
		t.Fatal("OnBeforeSerialization must not be called for typed-nil parsable")
		return nil
	})
	require.NoError(t, err)
	err = serializer.SetOnAfterObjectSerialization(func(parsable absser.Parsable) error {
		t.Fatal("OnAfterObjectSerialization must not be called for typed-nil parsable")
		return nil
	})
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		err = serializer.WriteObjectValue("typedNilProp", typedNilEntity, typedNilEntity, nil)
	})
	require.NoError(t, err)
	content, err = serializer.GetSerializedContent()
	require.NoError(t, err)
	assert.Empty(t, content, "must not emit property when all objects are typed-nil")

	// 4. Calling WriteObjectValue with typed-nil item and additional values mixed with valid entity
	serializer = NewFormSerializationWriter()
	assert.NotPanics(t, func() {
		err = serializer.WriteObjectValue("", typedNilEntity, entity1, typedNilEntity)
	})
	require.NoError(t, err)
	content, err = serializer.GetSerializedContent()
	require.NoError(t, err)
	assert.Contains(t, string(content), "id=id-1")
}

func TestWriteAnyValue_TypedNilSafety(t *testing.T) {
	serializer := NewFormSerializationWriter()

	var pStr *string = nil
	err := serializer.WriteAnyValue("nilStr", pStr)
	assert.NoError(t, err)

	var pInt *int32 = nil
	err = serializer.WriteAnyValue("nilInt", pInt)
	assert.NoError(t, err)

	var pEntity *internal.TestEntity = nil
	err = serializer.WriteAnyValue("nilEntity", pEntity)
	assert.NoError(t, err)

	var pSlice []string = nil
	err = serializer.WriteAnyValue("nilSlice", pSlice)
	assert.NoError(t, err)

	var pMap map[string]string = nil
	err = serializer.WriteAnyValue("nilMap", pMap)
	assert.NoError(t, err)

	err = serializer.WriteAnyValue("rawNil", nil)
	assert.NoError(t, err)

	// Non-nil value must return error as FormSerializationWriter doesn't support any-value serialization
	err = serializer.WriteAnyValue("strVal", "not-nil")
	assert.Error(t, err)

	content, err := serializer.GetSerializedContent()
	require.NoError(t, err)
	assert.Empty(t, content)
}

func TestFormSerialization_WriteNullValueRoundTrip(t *testing.T) {
	writer := NewFormSerializationWriter()

	require.NoError(t, writer.WriteStringValue("normalStr", ref("hello")))
	require.NoError(t, writer.WriteNullValue("nullBool"))
	require.NoError(t, writer.WriteNullValue("nullInt8"))
	require.NoError(t, writer.WriteNullValue("nullByte"))
	require.NoError(t, writer.WriteNullValue("nullInt32"))
	require.NoError(t, writer.WriteNullValue("nullInt64"))
	require.NoError(t, writer.WriteNullValue("nullFloat32"))
	require.NoError(t, writer.WriteNullValue("nullFloat64"))
	require.NoError(t, writer.WriteNullValue("nullTime"))
	require.NoError(t, writer.WriteNullValue("nullDuration"))
	require.NoError(t, writer.WriteNullValue("nullTimeOnly"))
	require.NoError(t, writer.WriteNullValue("nullDateOnly"))
	require.NoError(t, writer.WriteNullValue("nullUUID"))
	require.NoError(t, writer.WriteByteArrayValue("nilBytes", nil))
	require.NoError(t, writer.WriteByteArrayValue("emptyBytes", []byte{}))

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.NotEmpty(t, content)

	// Round-trip parse using FormParseNode
	node, err := NewFormParseNode(content)
	require.NoError(t, err)

	// Verify normal string parses correctly
	normalChild, err := node.GetChildNode("normalStr")
	require.NoError(t, err)
	normalVal, err := normalChild.GetStringValue()
	require.NoError(t, err)
	require.NotNil(t, normalVal)
	assert.Equal(t, "hello", *normalVal)

	// Verify every typed primitive field parses to nil without syntax, range, or decode errors
	boolChild, err := node.GetChildNode("nullBool")
	require.NoError(t, err)
	bVal, err := boolChild.GetBoolValue()
	require.NoError(t, err)
	assert.Nil(t, bVal)

	i8Child, err := node.GetChildNode("nullInt8")
	require.NoError(t, err)
	i8Val, err := i8Child.GetInt8Value()
	require.NoError(t, err)
	assert.Nil(t, i8Val)

	byteChild, err := node.GetChildNode("nullByte")
	require.NoError(t, err)
	byteVal, err := byteChild.GetByteValue()
	require.NoError(t, err)
	assert.Nil(t, byteVal)

	i32Child, err := node.GetChildNode("nullInt32")
	require.NoError(t, err)
	i32Val, err := i32Child.GetInt32Value()
	require.NoError(t, err)
	assert.Nil(t, i32Val)

	i64Child, err := node.GetChildNode("nullInt64")
	require.NoError(t, err)
	i64Val, err := i64Child.GetInt64Value()
	require.NoError(t, err)
	assert.Nil(t, i64Val)

	f32Child, err := node.GetChildNode("nullFloat32")
	require.NoError(t, err)
	f32Val, err := f32Child.GetFloat32Value()
	require.NoError(t, err)
	assert.Nil(t, f32Val)

	f64Child, err := node.GetChildNode("nullFloat64")
	require.NoError(t, err)
	f64Val, err := f64Child.GetFloat64Value()
	require.NoError(t, err)
	assert.Nil(t, f64Val)

	timeChild, err := node.GetChildNode("nullTime")
	require.NoError(t, err)
	timeVal, err := timeChild.GetTimeValue()
	require.NoError(t, err)
	assert.Nil(t, timeVal)

	durChild, err := node.GetChildNode("nullDuration")
	require.NoError(t, err)
	durVal, err := durChild.GetISODurationValue()
	require.NoError(t, err)
	assert.Nil(t, durVal)

	toChild, err := node.GetChildNode("nullTimeOnly")
	require.NoError(t, err)
	toVal, err := toChild.GetTimeOnlyValue()
	require.NoError(t, err)
	assert.Nil(t, toVal)

	doChild, err := node.GetChildNode("nullDateOnly")
	require.NoError(t, err)
	doVal, err := doChild.GetDateOnlyValue()
	require.NoError(t, err)
	assert.Nil(t, doVal)

	uuidChild, err := node.GetChildNode("nullUUID")
	require.NoError(t, err)
	uuidVal, err := uuidChild.GetUUIDValue()
	require.NoError(t, err)
	assert.Nil(t, uuidVal)

	nilBytesChild, err := node.GetChildNode("nilBytes")
	require.NoError(t, err)
	assert.Nil(t, nilBytesChild)

	emptyBytesChild, err := node.GetChildNode("emptyBytes")
	require.NoError(t, err)
	require.NotNil(t, emptyBytesChild)
	emptyBytesVal, err := emptyBytesChild.GetByteArrayValue()
	require.NoError(t, err)
	assert.Nil(t, emptyBytesVal)
}

func TestFormSerialization_ByteArrayRoundTrip_Base64NullPayload(t *testing.T) {
	writer := NewFormSerializationWriter()
	payload := []byte{158, 233, 101}
	err := writer.WriteByteArrayValue("data", payload)
	require.NoError(t, err)

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "data=null", string(content))

	node, err := NewFormParseNode(content)
	require.NoError(t, err)

	child, err := node.GetChildNode("data")
	require.NoError(t, err)
	require.NotNil(t, child)

	parsedBytes, err := child.GetByteArrayValue()
	require.NoError(t, err)
	assert.Equal(t, payload, parsedBytes)
}

func TestFormSerialization_StringRoundTrip_LiteralNullPayloads(t *testing.T) {
	for _, literal := range []string{"null", "NULL", " Null ", "nullify", "NULL_VALUE"} {
		writer := NewFormSerializationWriter()
		err := writer.WriteStringValue("str", ref(literal))
		require.NoError(t, err)

		content, err := writer.GetSerializedContent()
		require.NoError(t, err)

		node, err := NewFormParseNode(content)
		require.NoError(t, err)

		child, err := node.GetChildNode("str")
		require.NoError(t, err)
		require.NotNil(t, child)

		val, err := child.GetStringValue()
		require.NoError(t, err)
		require.NotNil(t, val)
		assert.Equal(t, literal, *val)
	}

	// String collections containing literal nulls
	writer := NewFormSerializationWriter()
	err := writer.WriteCollectionOfStringValues("items", []string{"null", "NULL", "hello"})
	require.NoError(t, err)

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)

	node, err := NewFormParseNode(content)
	require.NoError(t, err)

	child, err := node.GetChildNode("items")
	require.NoError(t, err)
	require.NotNil(t, child)

	parsed, err := child.GetCollectionOfPrimitiveValues("string")
	require.NoError(t, err)
	require.Len(t, parsed, 3)
	require.NotNil(t, parsed[0])
	assert.Equal(t, "null", *parsed[0].(*string))
	require.NotNil(t, parsed[1])
	assert.Equal(t, "NULL", *parsed[1].(*string))
	require.NotNil(t, parsed[2])
	assert.Equal(t, "hello", *parsed[2].(*string))
}

func TestWriteObjectValue_UnionTypeMockPrimary(t *testing.T) {
	mock := internal.NewUnionTypeMock()
	entity := internal.NewTestEntity()
	id := "12345"
	office := "Redmond"
	entity.SetId(&id)
	entity.SetOfficeLocation(&office)
	mock.SetComposedType1(entity)

	writer := NewFormSerializationWriter()
	err := writer.WriteObjectValue("", mock)
	require.NoError(t, err)

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	result := string(content)
	assert.Contains(t, result, "id=12345")
	assert.Contains(t, result, "officeLocation=Redmond")
}

func TestWriteObjectValue_UnionTypeMockSecondary(t *testing.T) {
	mock := internal.NewUnionTypeMock()
	entity := internal.NewSecondTestEntity()
	displayName := "SecondModel"
	var id int64 = 98765
	entity.SetDisplayName(&displayName)
	entity.SetId(&id)
	mock.SetComposedType2(entity)

	writer := NewFormSerializationWriter()
	err := writer.WriteObjectValue("", mock)
	require.NoError(t, err)

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	result := string(content)
	assert.Contains(t, result, "displayName=SecondModel")
	assert.Contains(t, result, "id=98765")
}

func TestWriteObjectValue_IntersectionTypeMock(t *testing.T) {
	mock := internal.NewIntersectionTypeMock()
	entity1 := internal.NewTestEntity()
	id1 := "test-id-1"
	office := "Building-A"
	entity1.SetId(&id1)
	entity1.SetOfficeLocation(&office)
	mock.SetComposedType1(entity1)

	entity2 := internal.NewSecondTestEntity()
	displayName := "ComposedDisplayName"
	var id2 int64 = 55443
	entity2.SetDisplayName(&displayName)
	entity2.SetId(&id2)
	mock.SetComposedType2(entity2)

	writer := NewFormSerializationWriter()
	err := writer.WriteObjectValue("", mock)
	require.NoError(t, err)

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	result := string(content)
	assert.Contains(t, result, "id=test-id-1")
	assert.Contains(t, result, "officeLocation=Building-A")
	assert.Contains(t, result, "displayName=ComposedDisplayName")
}

func TestWriteObjectValue_UnionTypeMockStringValue(t *testing.T) {
	mock := internal.NewUnionTypeMock()
	strVal := "hello world"
	mock.SetStringValue(&strVal)

	writer := NewFormSerializationWriter()
	err := writer.WriteObjectValue("", mock)
	require.NoError(t, err)

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "hello+world", string(content))
}

func TestWriteObjectValue_UnionTypeMockEmpty(t *testing.T) {
	mock := internal.NewUnionTypeMock()

	writer := NewFormSerializationWriter()
	err := writer.WriteObjectValue("", mock)
	require.NoError(t, err)

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Empty(t, content)
}

func TestWriteObjectValue_IntersectionTypeMockStringValue(t *testing.T) {
	mock := internal.NewIntersectionTypeMock()
	strVal := "intersection text"
	mock.SetStringValue(&strVal)

	writer := NewFormSerializationWriter()
	err := writer.WriteObjectValue("", mock)
	require.NoError(t, err)

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "intersection+text", string(content))
}

func TestWriteObjectValue_MultiLevelComposedDelegation(t *testing.T) {
	innerUnion := internal.NewUnionTypeMock()
	entity := internal.NewTestEntity()
	id := "deep-id"
	entity.SetId(&id)
	innerUnion.SetComposedType1(entity)

	outerWrapper := &mockTestParsable{
		serializeFunc: func(w absser.SerializationWriter) error {
			return w.WriteObjectValue("", innerUnion)
		},
	}

	writer := NewFormSerializationWriter()
	err := writer.WriteObjectValue("", outerWrapper)
	require.NoError(t, err)

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "id=deep-id", string(content))
}

func TestWriteObjectValue_NestedParsableWithNonEmptyKeyFails(t *testing.T) {
	innerEntityWithNested := &mockTestParsable{
		serializeFunc: func(w absser.SerializationWriter) error {
			child := internal.NewTestEntity()
			return w.WriteObjectValue("nestedProp", child)
		},
	}

	outerUnion := &mockTestParsable{
		serializeFunc: func(w absser.SerializationWriter) error {
			// Composed delegation with empty key
			return w.WriteObjectValue("", innerEntityWithNested)
		},
	}

	writer := NewFormSerializationWriter()
	err := writer.WriteObjectValue("", outerUnion)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nested objects serialization is not supported")
}

func TestFormSerialization_UnionTypeMockStringRoundTrip(t *testing.T) {
	mock := internal.NewUnionTypeMock()
	strVal := "union round trip"
	mock.SetStringValue(&strVal)

	writer := NewFormSerializationWriter()
	err := writer.WriteObjectValue("", mock)
	require.NoError(t, err)

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)

	node, err := NewFormParseNode(content)
	require.NoError(t, err)

	parsedObj, err := node.GetObjectValue(internal.CreateUnionTypeMockFromDiscriminator)
	require.NoError(t, err)
	require.NotNil(t, parsedObj)

	parsedUnion, ok := parsedObj.(internal.UnionTypeMockable)
	require.True(t, ok)
	require.NotNil(t, parsedUnion.GetStringValue())
	assert.Equal(t, "union round trip", *parsedUnion.GetStringValue())
}

func TestFormSerialization_IntersectionTypeMockRoundTrip(t *testing.T) {
	// 1. String value round-trip
	mock := internal.NewIntersectionTypeMock()
	strVal := "intersection round trip"
	mock.SetStringValue(&strVal)

	writer := NewFormSerializationWriter()
	err := writer.WriteObjectValue("", mock)
	require.NoError(t, err)

	content, err := writer.GetSerializedContent()
	require.NoError(t, err)

	node, err := NewFormParseNode(content)
	require.NoError(t, err)

	parsedObj, err := node.GetObjectValue(internal.CreateIntersectionTypeMockFromDiscriminator)
	require.NoError(t, err)
	require.NotNil(t, parsedObj)

	parsedIntersection, ok := parsedObj.(internal.IntersectionTypeMockable)
	require.True(t, ok)
	require.NotNil(t, parsedIntersection.GetStringValue())
	assert.Equal(t, "intersection round trip", *parsedIntersection.GetStringValue())

	// 2. Composed entities serialization and parsing of merged fields
	mockComposed := internal.NewIntersectionTypeMock()
	entity1 := internal.NewTestEntity()
	id1 := "12345"
	office := "Studio B"
	entity1.SetId(&id1)
	entity1.SetOfficeLocation(&office)
	mockComposed.SetComposedType1(entity1)

	entity2 := internal.NewSecondTestEntity()
	displayName := "SecondComposed"
	rate := 12.5
	entity2.SetDisplayName(&displayName)
	entity2.SetFailureRate(&rate)
	mockComposed.SetComposedType2(entity2)

	writer2 := NewFormSerializationWriter()
	err = writer2.WriteObjectValue("", mockComposed)
	require.NoError(t, err)

	content2, err := writer2.GetSerializedContent()
	require.NoError(t, err)

	node2, err := NewFormParseNode(content2)
	require.NoError(t, err)

	parsedEntity1, err := node2.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	require.NoError(t, err)
	parsed1, ok := parsedEntity1.(internal.TestEntityable)
	require.True(t, ok)
	assert.Equal(t, id1, *parsed1.GetId())
	assert.Equal(t, office, *parsed1.GetOfficeLocation())

	parsedEntity2, err := node2.GetObjectValue(internal.CreateSecondTestEntityFromDiscriminator)
	require.NoError(t, err)
	parsed2, ok := parsedEntity2.(internal.SecondTestEntityable)
	require.True(t, ok)
	assert.Equal(t, displayName, *parsed2.GetDisplayName())
	assert.Equal(t, rate, *parsed2.GetFailureRate())
	assert.Equal(t, int64(12345), *parsed2.GetId())
}

func TestWriteObjectValue_KeyedObjectAndCompositionRejected(t *testing.T) {
	entity := internal.NewTestEntity()
	id := "123"
	entity.SetId(&id)

	t.Run("rejects keyed union at top level without emitting partial output", func(t *testing.T) {
		union := internal.NewUnionTypeMock()
		union.SetComposedType1(entity)

		writer := NewFormSerializationWriter()
		err := writer.WriteObjectValue("manager", union)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nested objects serialization is not supported")

		content, contentErr := writer.GetSerializedContent()
		require.NoError(t, contentErr)
		assert.Empty(t, content)
	})

	t.Run("rejects keyed entity at top level", func(t *testing.T) {
		writer := NewFormSerializationWriter()
		err := writer.WriteObjectValue("manager", entity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nested objects serialization is not supported")

		content, contentErr := writer.GetSerializedContent()
		require.NoError(t, contentErr)
		assert.Empty(t, content)
	})

	t.Run("rejects keyed intersection at top level", func(t *testing.T) {
		intersection := internal.NewIntersectionTypeMock()
		intersection.SetComposedType1(entity)

		writer := NewFormSerializationWriter()
		err := writer.WriteObjectValue("manager", intersection)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nested objects serialization is not supported")
	})

	t.Run("rejects keyed union delegated from outer wrapper", func(t *testing.T) {
		innerUnion := internal.NewUnionTypeMock()
		innerUnion.SetComposedType1(entity)
		outerWrapper := &mockTestParsable{
			serializeFunc: func(w absser.SerializationWriter) error {
				return w.WriteObjectValue("manager", innerUnion)
			},
		}

		writer := NewFormSerializationWriter()
		err := writer.WriteObjectValue("", outerWrapper)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nested objects serialization is not supported")
	})

	t.Run("rejects parsable in WriteAdditionalData", func(t *testing.T) {
		writer := NewFormSerializationWriter()
		err := writer.WriteAdditionalData(map[string]interface{}{
			"manager": entity,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nested objects serialization is not supported")
	})
}

func TestFormSerialization_SequentialGetSerializedContentAndWrites(t *testing.T) {
	writer := NewFormSerializationWriter()
	strVal1 := "first"
	strVal2 := "second"
	intVal := int32(42)

	require.NoError(t, writer.WriteStringValue("prop1", &strVal1))

	content1, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "prop1=first", string(content1))

	// Verify idempotence of multiple GetSerializedContent calls
	content1Repeat, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "prop1=first", string(content1Repeat))

	// Write subsequent properties after reading serialized content
	require.NoError(t, writer.WriteStringValue("prop2", &strVal2))
	require.NoError(t, writer.WriteInt32Value("prop3", &intVal))

	content2, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "prop1=first&prop2=second&prop3=42", string(content2))
}

func TestFormSerialization_EndToEndIncrementalRoundTrip(t *testing.T) {
	writer := NewFormSerializationWriter()

	// 1. Initial write
	name := "Alice"
	require.NoError(t, writer.WriteStringValue("name", &name))
	bytes1, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "name=Alice", string(bytes1))

	// 2. Incremental writes with collections and special characters
	rate := float64(9.99)
	require.NoError(t, writer.WriteFloat64Value("rate%", &rate))
	require.NoError(t, writer.WriteCollectionOfStringValues("tags", []string{"go", "kiota", "sdk"}))

	finalBytes, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "name=Alice&rate%25=9.99&tags=go&tags=kiota&tags=sdk", string(finalBytes))

	// 3. Parse and verify round trip
	node, err := NewFormParseNode(finalBytes)
	require.NoError(t, err)

	nameChild, err := node.GetChildNode("name")
	require.NoError(t, err)
	require.NotNil(t, nameChild)
	parsedName, err := nameChild.GetStringValue()
	require.NoError(t, err)
	assert.Equal(t, "Alice", *parsedName)

	// Verify direct lookup with unescaped key containing '%'
	rateChild, err := node.GetChildNode("rate%")
	require.NoError(t, err)
	require.NotNil(t, rateChild)
	parsedRate, err := rateChild.GetFloat64Value()
	require.NoError(t, err)
	assert.Equal(t, 9.99, *parsedRate)

	// Verify collection parsing
	tagsChild, err := node.GetChildNode("tags")
	require.NoError(t, err)
	require.NotNil(t, tagsChild)
	parsedTags, err := tagsChild.GetCollectionOfPrimitiveValues("string")
	require.NoError(t, err)
	assert.Equal(t, []interface{}{ref("go"), ref("kiota"), ref("sdk")}, parsedTags)
}

func serializeUntyped(t *testing.T, key string, item absser.Parsable, extra ...absser.Parsable) string {
	t.Helper()
	w := NewFormSerializationWriter()
	defer w.Close()
	require.NoError(t, w.WriteObjectValue(key, item, extra...))
	content, err := w.GetSerializedContent()
	require.NoError(t, err)
	return string(content)
}

func serializeUntypedError(t *testing.T, key string, item absser.Parsable, extra ...absser.Parsable) error {
	t.Helper()
	w := NewFormSerializationWriter()
	defer w.Close()
	return w.WriteObjectValue(key, item, extra...)
}

func TestWriteObjectValue_UntypedNodes(t *testing.T) {
	t.Run("serializes scalar untyped primitives with key", func(t *testing.T) {
		w := NewFormSerializationWriter()
		defer w.Close()
		require.NoError(t, w.WriteObjectValue("str", absser.NewUntypedString("hello world")))
		require.NoError(t, w.WriteObjectValue("bool", absser.NewUntypedBoolean(true)))
		require.NoError(t, w.WriteObjectValue("int", absser.NewUntypedInteger(42)))
		require.NoError(t, w.WriteObjectValue("long", absser.NewUntypedLong(9007199254740993)))
		require.NoError(t, w.WriteObjectValue("float", absser.NewUntypedFloat(3.14)))
		require.NoError(t, w.WriteObjectValue("double", absser.NewUntypedDouble(2.71828)))
		require.NoError(t, w.WriteObjectValue("nilVal", absser.NewUntypedNull()))
		content, err := w.GetSerializedContent()
		require.NoError(t, err)
		str := string(content)
		for _, exp := range []string{"str=hello+world", "bool=true", "int=42", "long=9007199254740993", "float=3.14", "double=2.71828", "nilVal=null"} {
			assert.Contains(t, str, exp)
		}
	})

	t.Run("serializes root UntypedObject", func(t *testing.T) {
		str := serializeUntyped(t, "", absser.NewUntypedObject(map[string]absser.UntypedNodeable{
			"title":    absser.NewUntypedString("Form Spec"),
			"priority": absser.NewUntypedInteger(1),
			"enabled":  absser.NewUntypedBoolean(true),
		}))
		assert.Contains(t, str, "title=Form+Spec")
		assert.Contains(t, str, "priority=1")
		assert.Contains(t, str, "enabled=true")
	})

	t.Run("rejects nested UntypedObject with non-empty key", func(t *testing.T) {
		err := serializeUntypedError(t, "nested", absser.NewUntypedObject(map[string]absser.UntypedNodeable{"inner": absser.NewUntypedString("val")}))
		assert.ErrorContains(t, err, "nested objects serialization is not supported")
	})

	t.Run("serializes keyed UntypedArray of primitives as repeated parameters", func(t *testing.T) {
		arr := absser.NewUntypedArray([]absser.UntypedNodeable{absser.NewUntypedString("first"), absser.NewUntypedString("second")})
		assert.Equal(t, "tags=first&tags=second", serializeUntyped(t, "tags", arr))
	})

	t.Run("rejects root UntypedArray in WriteObjectValue", func(t *testing.T) {
		err := serializeUntypedError(t, "", absser.NewUntypedArray([]absser.UntypedNodeable{absser.NewUntypedString("item")}))
		assert.ErrorContains(t, err, "collections serialization is not supported")
	})

	t.Run("rejects keyed UntypedArray containing UntypedObject", func(t *testing.T) {
		nested := absser.NewUntypedObject(map[string]absser.UntypedNodeable{"inner": absser.NewUntypedString("val")})
		arr := absser.NewUntypedArray([]absser.UntypedNodeable{nested, absser.NewUntypedString("item2")})
		assert.ErrorContains(t, serializeUntypedError(t, "items", arr), "nested objects serialization is not supported")
	})

	t.Run("WriteAdditionalData with untyped nodes", func(t *testing.T) {
		w := NewFormSerializationWriter()
		defer w.Close()
		require.NoError(t, w.WriteAdditionalData(map[string]interface{}{
			"extraString":  absser.NewUntypedString("extra"),
			"extraBoolean": absser.NewUntypedBoolean(false),
			"extraInteger": absser.NewUntypedInteger(99),
			"extraNull":    absser.NewUntypedNull(),
		}))
		content, err := w.GetSerializedContent()
		require.NoError(t, err)
		for _, exp := range []string{"extraString=extra", "extraBoolean=false", "extraInteger=99", "extraNull=null"} {
			assert.Contains(t, string(content), exp)
		}
	})
}

func TestFormSerialization_UntypedObjectAndCollectionEndToEndRoundTrip(t *testing.T) {
	t.Run("untyped object full round-trip", func(t *testing.T) {
		obj := absser.NewUntypedObject(map[string]absser.UntypedNodeable{
			"strVal":  absser.NewUntypedString("Kiota Form Serialization"),
			"boolVal": absser.NewUntypedBoolean(true),
			"intVal":  absser.NewUntypedInteger(12345),
			"longVal": absser.NewUntypedLong(9007199254740993),
			"fltVal":  absser.NewUntypedDouble(12.75),
			"nullVal": absser.NewUntypedNull(),
		})
		str := serializeUntyped(t, "", obj)
		pn, err := NewFormParseNode([]byte(str))
		require.NoError(t, err)
		res, err := pn.GetObjectValue(absser.CreateUntypedNodeFromDiscriminatorValue)
		require.NoError(t, err)
		props := res.(*absser.UntypedObject).GetValue()
		assert.Equal(t, "Kiota Form Serialization", *props["strVal"].(*absser.UntypedString).GetValue())
		assert.Equal(t, true, *props["boolVal"].(*absser.UntypedBoolean).GetValue())
		assert.Equal(t, int32(12345), *props["intVal"].(*absser.UntypedInteger).GetValue())
		assert.Equal(t, int64(9007199254740993), *props["longVal"].(*absser.UntypedLong).GetValue())
		assert.Equal(t, 12.75, *props["fltVal"].(*absser.UntypedDouble).GetValue())
		assert.NotNil(t, props["nullVal"].(*absser.UntypedNull))
	})

	t.Run("entity with AdditionalData containing untyped nodes round-trip", func(t *testing.T) {
		w := NewFormSerializationWriter()
		defer w.Close()
		entity := internal.NewTestEntity()
		id := "test-id-123"
		entity.SetId(&id)
		require.NoError(t, w.WriteObjectValue("", entity))
		require.NoError(t, w.WriteAdditionalData(map[string]interface{}{
			"extraString":  absser.NewUntypedString("extraVal"),
			"extraBoolean": absser.NewUntypedBoolean(true),
			"extraInteger": absser.NewUntypedInteger(42),
			"extraNull":    absser.NewUntypedNull(),
		}))
		content, err := w.GetSerializedContent()
		require.NoError(t, err)
		str := string(content)
		for _, exp := range []string{"id=test-id-123", "extraString=extraVal", "extraBoolean=true", "extraInteger=42", "extraNull=null"} {
			assert.Contains(t, str, exp)
		}
	})

	t.Run("entity deserialization populates additional data from form fields", func(t *testing.T) {
		pn, err := NewFormParseNode([]byte("id=entity-789&customKey=customValue&extraNum=99"))
		require.NoError(t, err)
		res, err := pn.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
		require.NoError(t, err)
		entity := res.(*internal.TestEntity)
		assert.Equal(t, "entity-789", *entity.GetId())
		assert.Equal(t, "customValue", entity.GetAdditionalData()["customKey"])
		assert.Equal(t, "99", entity.GetAdditionalData()["extraNum"])
	})
}

func TestFormSerialization_MergeBehavior(t *testing.T) {
	t.Run("UntypedObject followed by regular entity", func(t *testing.T) {
		id := "entity-123"
		entity := internal.NewTestEntity()
		entity.SetId(&id)
		obj := absser.NewUntypedObject(map[string]absser.UntypedNodeable{"untypedProp": absser.NewUntypedString("untypedVal")})
		str := serializeUntyped(t, "", obj, entity)
		assert.Contains(t, str, "untypedProp=untypedVal")
		assert.Contains(t, str, "id=entity-123")
	})

	t.Run("regular entity followed by UntypedObject", func(t *testing.T) {
		id := "entity-456"
		entity := internal.NewTestEntity()
		entity.SetId(&id)
		obj := absser.NewUntypedObject(map[string]absser.UntypedNodeable{"extraField": absser.NewUntypedString("extraVal")})
		str := serializeUntyped(t, "", entity, obj)
		assert.Contains(t, str, "id=entity-456")
		assert.Contains(t, str, "extraField=extraVal")
	})

	t.Run("multiple UntypedObjects merged", func(t *testing.T) {
		obj1 := absser.NewUntypedObject(map[string]absser.UntypedNodeable{"prop1": absser.NewUntypedString("val1")})
		obj2 := absser.NewUntypedObject(map[string]absser.UntypedNodeable{"prop2": absser.NewUntypedInteger(100)})
		obj3 := absser.NewUntypedObject(map[string]absser.UntypedNodeable{"prop3": absser.NewUntypedBoolean(true)})
		str := serializeUntyped(t, "", obj1, obj2, obj3)
		assert.Contains(t, str, "prop1=val1")
		assert.Contains(t, str, "prop2=100")
		assert.Contains(t, str, "prop3=true")
	})
}

func TestFormSerialization_SerializationCallbacks(t *testing.T) {
	t.Run("invokes callbacks in order on UntypedObject", func(t *testing.T) {
		w := NewFormSerializationWriter()
		defer w.Close()
		events := []string{}
		require.NoError(t, w.SetOnBeforeSerialization(func(p absser.Parsable) error {
			events = append(events, "before")
			return nil
		}))
		require.NoError(t, w.SetOnStartObjectSerialization(func(p absser.Parsable, sw absser.SerializationWriter) error {
			events = append(events, "start")
			val := "callbackVal"
			return sw.WriteStringValue("callbackProp", &val)
		}))
		require.NoError(t, w.SetOnAfterObjectSerialization(func(p absser.Parsable) error {
			events = append(events, "after")
			return nil
		}))
		require.NoError(t, w.WriteObjectValue("", absser.NewUntypedObject(map[string]absser.UntypedNodeable{"title": absser.NewUntypedString("Doc")})))
		assert.Equal(t, []string{"before", "start", "after"}, events)
		content, _ := w.GetSerializedContent()
		assert.Contains(t, string(content), "callbackProp=callbackVal")
		assert.Contains(t, string(content), "title=Doc")
	})

	t.Run("callback error propagates and prevents subsequent writes", func(t *testing.T) {
		w := NewFormSerializationWriter()
		defer w.Close()
		expectedErr := errors.New("callback start failure")
		require.NoError(t, w.SetOnStartObjectSerialization(func(p absser.Parsable, sw absser.SerializationWriter) error { return expectedErr }))
		assert.ErrorIs(t, w.WriteObjectValue("", absser.NewUntypedObject(map[string]absser.UntypedNodeable{"k": absser.NewUntypedString("v")})), expectedErr)
	})

	t.Run("scalar untyped nodes do not invoke object callbacks", func(t *testing.T) {
		w := NewFormSerializationWriter()
		defer w.Close()
		startCalled := false
		require.NoError(t, w.SetOnStartObjectSerialization(func(p absser.Parsable, sw absser.SerializationWriter) error { startCalled = true; return nil }))
		require.NoError(t, w.WriteObjectValue("key", absser.NewUntypedString("hello")))
		assert.False(t, startCalled)
	})
}

func TestFormSerialization_RepeatedFieldsRoundTrip(t *testing.T) {
	t.Run("repeated keys round trip faithfully", func(t *testing.T) {
		obj := absser.NewUntypedObject(map[string]absser.UntypedNodeable{
			"tags":        absser.NewUntypedArray([]absser.UntypedNodeable{absser.NewUntypedString("golang"), absser.NewUntypedString("kiota"), absser.NewUntypedString("serialization")}),
			"singleComma": absser.NewUntypedString("hello,world"),
			"emptyStr":    absser.NewUntypedString(""),
			"nullVal":     absser.NewUntypedNull(),
			"numZero":     absser.NewUntypedInteger(0),
			"numOne":      absser.NewUntypedInteger(1),
			"zipCode":     absser.NewUntypedString("00123"),
			"sciExp":      absser.NewUntypedDouble(0.01),
		})
		str := serializeUntyped(t, "", obj)
		assert.Contains(t, str, "tags=golang")
		assert.Contains(t, str, "tags=kiota")
		assert.Contains(t, str, "tags=serialization")

		pn, err := NewFormParseNode([]byte(str))
		require.NoError(t, err)
		res, err := pn.GetObjectValue(absser.CreateUntypedNodeFromDiscriminatorValue)
		require.NoError(t, err)
		props := res.(*absser.UntypedObject).GetValue()

		tagsArr := props["tags"].(*absser.UntypedArray).GetValue()
		require.Len(t, tagsArr, 3)
		assert.Equal(t, "golang", *tagsArr[0].(*absser.UntypedString).GetValue())
		assert.Equal(t, "kiota", *tagsArr[1].(*absser.UntypedString).GetValue())
		assert.Equal(t, "serialization", *tagsArr[2].(*absser.UntypedString).GetValue())
		assert.Equal(t, "hello,world", *props["singleComma"].(*absser.UntypedString).GetValue())
		assert.Equal(t, "", *props["emptyStr"].(*absser.UntypedString).GetValue())
		assert.NotNil(t, props["nullVal"].(*absser.UntypedNull))
		assert.Equal(t, int32(0), *props["numZero"].(*absser.UntypedInteger).GetValue())
		assert.Equal(t, int32(1), *props["numOne"].(*absser.UntypedInteger).GetValue())
		assert.Equal(t, "00123", *props["zipCode"].(*absser.UntypedString).GetValue())
		assert.Equal(t, 0.01, *props["sciExp"].(*absser.UntypedDouble).GetValue())
	})
}

func TestFormSerialization_UntypedWriterEdgeCases(t *testing.T) {
	t.Run("formats whole-valued untyped floats with decimal point to prevent type loss", func(t *testing.T) {
		str := serializeUntyped(t, "", absser.NewUntypedObject(map[string]absser.UntypedNodeable{
			"flt": absser.NewUntypedFloat(42.0),
			"dbl": absser.NewUntypedDouble(100.0),
		}))
		assert.Contains(t, str, "flt=42.0")
		assert.Contains(t, str, "dbl=100.0")

		pn, err := NewFormParseNode([]byte(str))
		require.NoError(t, err)
		res, err := pn.GetObjectValue(absser.CreateUntypedNodeFromDiscriminatorValue)
		require.NoError(t, err)
		props := res.(*absser.UntypedObject).GetValue()
		assert.Equal(t, float64(42.0), *props["flt"].(*absser.UntypedDouble).GetValue())
		assert.Equal(t, float64(100.0), *props["dbl"].(*absser.UntypedDouble).GetValue())
	})

	t.Run("rejects UntypedArray with fewer than 2 elements", func(t *testing.T) {
		assert.ErrorContains(t, serializeUntypedError(t, "items", absser.NewUntypedArray([]absser.UntypedNodeable{})), "untyped arrays with fewer than 2 elements are not supported")
		assert.ErrorContains(t, serializeUntypedError(t, "items", absser.NewUntypedArray([]absser.UntypedNodeable{absser.NewUntypedString("onlyOne")})), "untyped arrays with fewer than 2 elements are not supported")
	})

	t.Run("serializes nil array elements as explicit null", func(t *testing.T) {
		arr := absser.NewUntypedArray([]absser.UntypedNodeable{absser.NewUntypedString("first"), nil, absser.NewUntypedString("third")})
		assert.Equal(t, "items=first&items=null&items=third", serializeUntyped(t, "items", arr))
	})

	t.Run("rejects nested UntypedArray inside UntypedArray", func(t *testing.T) {
		inner := absser.NewUntypedArray([]absser.UntypedNodeable{absser.NewUntypedString("a"), absser.NewUntypedString("b")})
		arr := absser.NewUntypedArray([]absser.UntypedNodeable{inner, absser.NewUntypedString("c")})
		assert.ErrorContains(t, serializeUntypedError(t, "items", arr), "nested collections serialization is not supported")
	})

	t.Run("rejects empty or whitespace-padded property names in UntypedObject", func(t *testing.T) {
		for _, name := range []string{"", "  ", " padded", "padded "} {
			err := serializeUntypedError(t, "", absser.NewUntypedObject(map[string]absser.UntypedNodeable{name: absser.NewUntypedString("val")}))
			assert.ErrorContains(t, err, "empty or whitespace-padded property names are not supported")
		}
	})

	t.Run("rejects non-finite untyped floats and doubles", func(t *testing.T) {
		for _, node := range []absser.UntypedNodeable{
			absser.NewUntypedDouble(math.NaN()),
			absser.NewUntypedDouble(math.Inf(1)),
			absser.NewUntypedDouble(math.Inf(-1)),
			absser.NewUntypedFloat(float32(math.NaN())),
			absser.NewUntypedFloat(float32(math.Inf(1))),
		} {
			assert.ErrorContains(t, serializeUntypedError(t, "", absser.NewUntypedObject(map[string]absser.UntypedNodeable{"val": node})), "non-finite float values are not supported")
		}
	})

	t.Run("omits nil typed objects while emitting explicit UntypedNull", func(t *testing.T) {
		w := NewFormSerializationWriter()
		defer w.Close()
		require.NoError(t, w.WriteObjectValue("nilTyped", nil))
		require.NoError(t, w.WriteObjectValue("nilUntyped", absser.NewUntypedNull()))
		content, err := w.GetSerializedContent()
		require.NoError(t, err)
		assert.Equal(t, "nilUntyped=null", string(content))
	})

	t.Run("rejects merging root untyped scalars with entities or other values", func(t *testing.T) {
		id := "e1"
		entity := internal.NewTestEntity()
		entity.SetId(&id)
		assert.ErrorContains(t, serializeUntypedError(t, "", absser.NewUntypedString("str"), entity), "merging root untyped scalars is not supported")
	})

	t.Run("serializes UntypedObject and AdditionalData keys in deterministic sorted order", func(t *testing.T) {
		obj := absser.NewUntypedObject(map[string]absser.UntypedNodeable{
			"z": absser.NewUntypedString("val1"),
			"a": absser.NewUntypedString("val2"),
			"m": absser.NewUntypedString("val3"),
		})
		assert.Equal(t, "a=val2&m=val3&z=val1", serializeUntyped(t, "", obj))

		w := NewFormSerializationWriter()
		defer w.Close()
		require.NoError(t, w.WriteAdditionalData(map[string]interface{}{"z": "val1", "a": "val2", "m": "val3"}))
		content, _ := w.GetSerializedContent()
		assert.Equal(t, "a=val2&m=val3&z=val1", string(content))
	})

	t.Run("rejects ambiguous untyped strings that resemble primitives", func(t *testing.T) {
		ambiguous := []string{"null", "NULL", "true", "TRUE", "false", "FALSE", "42", "-17", "1.20", "0.0", "-0.0"}
		for _, s := range ambiguous {
			w := NewFormSerializationWriter()
			err := w.WriteObjectValue("key", absser.NewUntypedString(s))
			assert.ErrorContains(t, err, "ambiguous untyped string values resembling primitives are not supported")
			w.Close()

			obj := absser.NewUntypedObject(map[string]absser.UntypedNodeable{
				"key": absser.NewUntypedString(s),
			})
			assert.ErrorContains(t, serializeUntypedError(t, "", obj), "ambiguous untyped string values resembling primitives are not supported")
		}

		nonAmbiguous := []string{"hello", "nullify", "true_val", " 42 ", "001.20", "+1.5"}
		for _, s := range nonAmbiguous {
			w := NewFormSerializationWriter()
			err := w.WriteObjectValue("key", absser.NewUntypedString(s))
			assert.NoError(t, err)
			w.Close()
		}
	})

	t.Run("rejects unrecognized untyped nodes with unsupported-type error", func(t *testing.T) {
		dummy := &customUntypedNode{UntypedNodeable: absser.NewUntypedString("sample")}
		w := NewFormSerializationWriter()
		defer w.Close()
		err := w.WriteObjectValue("", dummy)
		assert.ErrorContains(t, err, "unsupported untyped node type for FormSerializationWriter")
	})

	t.Run("rejects untyped arrays in WriteAdditionalData", func(t *testing.T) {
		arr := absser.NewUntypedArray([]absser.UntypedNodeable{
			absser.NewUntypedString("val1"),
			absser.NewUntypedString("val2"),
		})
		w := NewFormSerializationWriter()
		defer w.Close()
		err := w.WriteAdditionalData(map[string]interface{}{"tags": arr})
		assert.ErrorContains(t, err, "untyped arrays in additional data are not supported with FormSerializationWriter")
	})

	t.Run("serializes nil and typed-nil properties in UntypedObject as null tokens", func(t *testing.T) {
		var typedNilStr *absser.UntypedString
		obj := absser.NewUntypedObject(map[string]absser.UntypedNodeable{
			"a": nil,
			"b": typedNilStr,
			"c": absser.NewUntypedString("val"),
		})
		assert.Equal(t, "a=null&b=null&c=val", serializeUntyped(t, "", obj))
	})
}

type customUntypedNode struct {
	absser.UntypedNodeable
}
