package formserialization

import (
	"encoding/base64"
	"errors"
	"math"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// FormSerializationWriter implements SerializationWriter for URI form encoded.
type FormSerializationWriter struct {
	writer                     []string
	onBeforeSerialization      absser.ParsableAction
	onAfterObjectSerialization absser.ParsableAction
	onStartObjectSerialization absser.ParsableWriter
}

// NewFormSerializationWriter creates a new instance of the FormSerializationWriter.
func NewFormSerializationWriter() *FormSerializationWriter {
	return &FormSerializationWriter{
		writer: make([]string, 0),
	}
}
func (w *FormSerializationWriter) writeRawValue(value string) {
	if w == nil {
		return
	}
	w.writer = append(w.writer, url.QueryEscape(value))
}
func (w *FormSerializationWriter) writeStringValue(value string) {
	w.writeRawValue(value)
}
func (w *FormSerializationWriter) writePropertyName(key string) {
	if w == nil {
		return
	}
	w.writer = append(w.writer, url.QueryEscape(key)+"=")
}
func (w *FormSerializationWriter) writePropertySeparator() {
	if w == nil {
		return
	}
	w.writer = append(w.writer, "&")
}

// WriteStringValue writes a String value to underlying the byte array.
func (w *FormSerializationWriter) WriteStringValue(key string, value *string) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue(*value)
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteBoolValue writes a Bool value to underlying the byte array.
func (w *FormSerializationWriter) WriteBoolValue(key string, value *bool) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeRawValue(strconv.FormatBool(*value))
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteByteValue writes a Byte value to underlying the byte array.
func (w *FormSerializationWriter) WriteByteValue(key string, value *byte) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt8Value writes a int8 value to underlying the byte array.
func (w *FormSerializationWriter) WriteInt8Value(key string, value *int8) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt32Value writes a Int32 value to underlying the byte array.
func (w *FormSerializationWriter) WriteInt32Value(key string, value *int32) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt64Value writes a Int64 value to underlying the byte array.
func (w *FormSerializationWriter) WriteInt64Value(key string, value *int64) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeRawValue(strconv.FormatInt(*value, 10))
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteFloat32Value writes a Float32 value to underlying the byte array.
func (w *FormSerializationWriter) WriteFloat32Value(key string, value *float32) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if value != nil {
		cast := float64(*value)
		return w.WriteFloat64Value(key, &cast)
	}
	return nil
}

// WriteFloat64Value writes a Float64 value to underlying the byte array.
func (w *FormSerializationWriter) WriteFloat64Value(key string, value *float64) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeRawValue(strconv.FormatFloat(*value, 'f', -1, 64))
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteTimeValue writes a Time value to underlying the byte array.
func (w *FormSerializationWriter) WriteTimeValue(key string, value *time.Time) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue((*value).Format(time.RFC3339))
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteISODurationValue writes a ISODuration value to underlying the byte array.
func (w *FormSerializationWriter) WriteISODurationValue(key string, value *absser.ISODuration) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue((*value).String())
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteTimeOnlyValue writes a TimeOnly value to underlying the byte array.
func (w *FormSerializationWriter) WriteTimeOnlyValue(key string, value *absser.TimeOnly) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue((*value).String())
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteDateOnlyValue writes a DateOnly value to underlying the byte array.
func (w *FormSerializationWriter) WriteDateOnlyValue(key string, value *absser.DateOnly) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue((*value).String())
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteUUIDValue writes a UUID value to underlying the byte array.
func (w *FormSerializationWriter) WriteUUIDValue(key string, value *uuid.UUID) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue((*value).String())
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteByteArrayValue writes a ByteArray value to underlying the byte array.
func (w *FormSerializationWriter) WriteByteArrayValue(key string, value []byte) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue(base64.StdEncoding.EncodeToString(value))
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

func isNil(val any) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func isAmbiguousUntypedString(s string) bool {
	if strings.TrimSpace(s) != s || s == "" {
		return false
	}
	if strings.EqualFold(s, "null") {
		return true
	}
	if strings.EqualFold(s, "true") || strings.EqualFold(s, "false") {
		return true
	}
	if isCanonicalInteger(s) {
		if _, err := strconv.ParseInt(s, 10, 64); err == nil {
			return true
		}
	}
	if isCanonicalFloat(s) {
		if f64, err := strconv.ParseFloat(s, 64); err == nil && !math.IsNaN(f64) && !math.IsInf(f64, 0) {
			if f64 != 0 || isZeroMantissa(s) {
				return true
			}
		}
	}
	return false
}

func formatUntypedFloat(f float64, bitSize int) (string, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return "", errors.New("non-finite float values are not supported in form URL-encoded serialization")
	}
	s := strconv.FormatFloat(f, 'f', -1, bitSize)
	if math.Abs(f) >= 1e15 || (math.Abs(f) > 0 && math.Abs(f) <= 1e-4) {
		s = strconv.FormatFloat(f, 'g', -1, bitSize)
	}
	if !strings.ContainsAny(s, ".eE") {
		s += ".0"
	}
	return s, nil
}

func (w *FormSerializationWriter) writeRawProperty(key string, rawVal string) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" {
		w.writePropertyName(key)
	}
	w.writeRawValue(rawVal)
	if key != "" {
		w.writePropertySeparator()
	}
	return nil
}

// WriteObjectValue writes a Parsable value to underlying the byte array.
func (w *FormSerializationWriter) WriteObjectValue(key string, item absser.Parsable, additionalValuesToMerge ...absser.Parsable) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	var values []absser.Parsable
	if !isNil(item) {
		values = append(values, item)
	}
	for _, additionalValue := range additionalValuesToMerge {
		if !isNil(additionalValue) {
			values = append(values, additionalValue)
		}
	}
	if len(values) == 0 {
		return nil
	}

	if key == "" && len(values) > 1 {
		for _, val := range values {
			if untyped, ok := val.(absser.UntypedNodeable); ok {
				if _, isObj := untyped.(*absser.UntypedObject); !isObj {
					return errors.New("merging root untyped scalars is not supported with FormSerializationWriter")
				}
			}
		}
	}

	for _, val := range values {
		if untyped, ok := val.(absser.UntypedNodeable); ok {
			var writeErr error
			switch value := untyped.(type) {
			case *absser.UntypedBoolean:
				writeErr = w.WriteBoolValue(key, value.GetValue())
			case *absser.UntypedFloat:
				if v := value.GetValue(); v != nil {
					var s string
					s, writeErr = formatUntypedFloat(float64(*v), 32)
					if writeErr == nil {
						writeErr = w.writeRawProperty(key, s)
					}
				} else {
					writeErr = w.WriteNullValue(key)
				}
			case *absser.UntypedDouble:
				if v := value.GetValue(); v != nil {
					var s string
					s, writeErr = formatUntypedFloat(*v, 64)
					if writeErr == nil {
						writeErr = w.writeRawProperty(key, s)
					}
				} else {
					writeErr = w.WriteNullValue(key)
				}
			case *absser.UntypedInteger:
				writeErr = w.WriteInt32Value(key, value.GetValue())
			case *absser.UntypedLong:
				writeErr = w.WriteInt64Value(key, value.GetValue())
			case *absser.UntypedNull:
				writeErr = w.WriteNullValue(key)
			case *absser.UntypedString:
				if v := value.GetValue(); v != nil {
					if isAmbiguousUntypedString(*v) {
						return errors.New("ambiguous untyped string values resembling primitives are not supported with FormSerializationWriter")
					}
					writeErr = w.WriteStringValue(key, v)
				} else {
					writeErr = w.WriteNullValue(key)
				}
			case *absser.UntypedArray:
				if key == "" {
					return errors.New("collections serialization is not supported with FormSerializationWriter")
				}
				arr := value.GetValue()
				if len(arr) < 2 {
					return errors.New("untyped arrays with fewer than 2 elements are not supported with FormSerializationWriter")
				}
				for _, elem := range arr {
					if elem == nil || isNil(elem) {
						err := w.WriteNullValue(key)
						if err != nil {
							return err
						}
						continue
					}
					if _, isNestedArr := elem.(*absser.UntypedArray); isNestedArr {
						return errors.New("nested collections serialization is not supported with FormSerializationWriter")
					}
					err := w.WriteObjectValue(key, elem)
					if err != nil {
						return err
					}
				}
			case *absser.UntypedObject:
				if key != "" {
					return errors.New("nested objects serialization is not supported with FormSerializationWriter")
				}
				abstractions.InvokeParsableAction(w.GetOnBeforeSerialization(), val)
				err := abstractions.InvokeParsableWriter(w.GetOnStartObjectSerialization(), val, w)
				if err != nil {
					return err
				}
				props := value.GetValue()
				if props != nil {
					keys := make([]string, 0, len(props))
					for propKey := range props {
						keys = append(keys, propKey)
					}
					sort.Strings(keys)
					for _, propKey := range keys {
						propVal := props[propKey]
						if propKey == "" || strings.TrimSpace(propKey) != propKey {
							return errors.New("empty or whitespace-padded property names are not supported with FormSerializationWriter")
						}
						if propVal == nil || isNil(propVal) {
							err := w.WriteNullValue(propKey)
							if err != nil {
								return err
							}
							continue
						}
						err := w.WriteObjectValue(propKey, propVal)
						if err != nil {
							return err
						}
					}
				}
				abstractions.InvokeParsableAction(w.GetOnAfterObjectSerialization(), val)
			default:
				return errors.New("unsupported untyped node type for FormSerializationWriter")
			}
			if writeErr != nil {
				return writeErr
			}
			continue
		}

		if key != "" {
			return errors.New("nested objects serialization is not supported with FormSerializationWriter")
		}

		abstractions.InvokeParsableAction(w.GetOnBeforeSerialization(), val)
		err := abstractions.InvokeParsableWriter(w.GetOnStartObjectSerialization(), val, w)
		if err != nil {
			return err
		}
		err = val.Serialize(w)
		if err != nil {
			return err
		}
		abstractions.InvokeParsableAction(w.GetOnAfterObjectSerialization(), val)
	}

	return nil
}

// WriteCollectionOfObjectValues writes a collection of Parsable values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfObjectValues(key string, collection []absser.Parsable) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if collection == nil {
		return nil
	}
	return errors.New("collections serialization is not supported with FormSerializationWriter")
}

func writeCollectionOfPrimitiveValues[T interface{}](key string, writer func(string, *T) error, collection []T) error {
	if collection != nil && len(collection) > 0 {
		for _, item := range collection {
			err := writer(key, &item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteCollectionOfStringValues writes a collection of String values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfStringValues(key string, collection []string) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteStringValue, collection)
}

// WriteCollectionOfInt32Values writes a collection of Int32 values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfInt32Values(key string, collection []int32) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteInt32Value, collection)
}

// WriteCollectionOfInt64Values writes a collection of Int64 values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfInt64Values(key string, collection []int64) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteInt64Value, collection)
}

// WriteCollectionOfFloat32Values writes a collection of Float32 values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfFloat32Values(key string, collection []float32) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteFloat32Value, collection)
}

// WriteCollectionOfFloat64Values writes a collection of Float64 values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfFloat64Values(key string, collection []float64) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteFloat64Value, collection)
}

// WriteCollectionOfTimeValues writes a collection of Time values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfTimeValues(key string, collection []time.Time) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteTimeValue, collection)
}

// WriteCollectionOfISODurationValues writes a collection of ISODuration values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfISODurationValues(key string, collection []absser.ISODuration) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteISODurationValue, collection)
}

// WriteCollectionOfTimeOnlyValues writes a collection of TimeOnly values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfTimeOnlyValues(key string, collection []absser.TimeOnly) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteTimeOnlyValue, collection)
}

// WriteCollectionOfDateOnlyValues writes a collection of DateOnly values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfDateOnlyValues(key string, collection []absser.DateOnly) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteDateOnlyValue, collection)
}

// WriteCollectionOfUUIDValues writes a collection of UUID values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfUUIDValues(key string, collection []uuid.UUID) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteUUIDValue, collection)
}

// WriteCollectionOfBoolValues writes a collection of Bool values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfBoolValues(key string, collection []bool) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteBoolValue, collection)
}

// WriteCollectionOfByteValues writes a collection of Byte values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfByteValues(key string, collection []byte) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteByteValue, collection)
}

// WriteCollectionOfInt8Values writes a collection of int8 values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfInt8Values(key string, collection []int8) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	return writeCollectionOfPrimitiveValues(key, w.WriteInt8Value, collection)
}

// GetSerializedContent returns the resulting byte array from the serialization writer.
func (w *FormSerializationWriter) GetSerializedContent() ([]byte, error) {
	if w == nil {
		return nil, errors.New("writer cannot be nil")
	}
	resultStr := strings.TrimSuffix(strings.Join(w.writer, ""), "&")
	return []byte(resultStr), nil
}

// WriteAnyValue an unknown value as a parameter.
func (w *FormSerializationWriter) WriteAnyValue(key string, value interface{}) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if isNil(value) {
		return nil
	}
	return errors.New("serialization of any value is not supported with FormSerializationWriter")
}

// WriteAdditionalData writes additional data to underlying the byte array.
func (w *FormSerializationWriter) WriteAdditionalData(value map[string]interface{}) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if len(value) != 0 {
		keys := make([]string, 0, len(value))
		for k := range value {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			input := value[key]
			var err error
			switch value := input.(type) {
			case *absser.UntypedArray:
				return errors.New("untyped arrays in additional data are not supported with FormSerializationWriter")
			case absser.Parsable:
				if _, isArr := value.(*absser.UntypedArray); isArr {
					return errors.New("untyped arrays in additional data are not supported with FormSerializationWriter")
				}
				err = w.WriteObjectValue(key, value)
			case []absser.Parsable:
				err = w.WriteCollectionOfObjectValues(key, value)
			case []string:
				err = w.WriteCollectionOfStringValues(key, value)
			case []bool:
				err = w.WriteCollectionOfBoolValues(key, value)
			case []byte:
				err = w.WriteCollectionOfByteValues(key, value)
			case []int8:
				err = w.WriteCollectionOfInt8Values(key, value)
			case []int32:
				err = w.WriteCollectionOfInt32Values(key, value)
			case []int64:
				err = w.WriteCollectionOfInt64Values(key, value)
			case []int:
				if value != nil {
					vals := make([]int64, len(value))
					for i, v := range value {
						vals[i] = int64(v)
					}
					err = w.WriteCollectionOfInt64Values(key, vals)
				}
			case []float32:
				err = w.WriteCollectionOfFloat32Values(key, value)
			case []float64:
				err = w.WriteCollectionOfFloat64Values(key, value)
			case []uuid.UUID:
				err = w.WriteCollectionOfUUIDValues(key, value)
			case []time.Time:
				err = w.WriteCollectionOfTimeValues(key, value)
			case []absser.ISODuration:
				err = w.WriteCollectionOfISODurationValues(key, value)
			case []absser.TimeOnly:
				err = w.WriteCollectionOfTimeOnlyValues(key, value)
			case []absser.DateOnly:
				err = w.WriteCollectionOfDateOnlyValues(key, value)
			case *string:
				err = w.WriteStringValue(key, value)
			case string:
				err = w.WriteStringValue(key, &value)
			case *bool:
				err = w.WriteBoolValue(key, value)
			case bool:
				err = w.WriteBoolValue(key, &value)
			case *byte:
				err = w.WriteByteValue(key, value)
			case byte:
				err = w.WriteByteValue(key, &value)
			case *int8:
				err = w.WriteInt8Value(key, value)
			case int8:
				err = w.WriteInt8Value(key, &value)
			case *int32:
				err = w.WriteInt32Value(key, value)
			case int32:
				err = w.WriteInt32Value(key, &value)
			case *int64:
				err = w.WriteInt64Value(key, value)
			case int64:
				err = w.WriteInt64Value(key, &value)
			case *int:
				if value != nil {
					val64 := int64(*value)
					err = w.WriteInt64Value(key, &val64)
				}
			case int:
				val64 := int64(value)
				err = w.WriteInt64Value(key, &val64)
			case *float32:
				err = w.WriteFloat32Value(key, value)
			case float32:
				err = w.WriteFloat32Value(key, &value)
			case *float64:
				err = w.WriteFloat64Value(key, value)
			case float64:
				err = w.WriteFloat64Value(key, &value)
			case *uuid.UUID:
				err = w.WriteUUIDValue(key, value)
			case uuid.UUID:
				err = w.WriteUUIDValue(key, &value)
			case *time.Time:
				err = w.WriteTimeValue(key, value)
			case time.Time:
				err = w.WriteTimeValue(key, &value)
			case *absser.ISODuration:
				err = w.WriteISODurationValue(key, value)
			case absser.ISODuration:
				err = w.WriteISODurationValue(key, &value)
			case *absser.TimeOnly:
				err = w.WriteTimeOnlyValue(key, value)
			case absser.TimeOnly:
				err = w.WriteTimeOnlyValue(key, &value)
			case *absser.DateOnly:
				err = w.WriteDateOnlyValue(key, value)
			case absser.DateOnly:
				err = w.WriteDateOnlyValue(key, &value)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Close clears the internal buffer.
func (w *FormSerializationWriter) Close() error {
	if w == nil {
		return nil
	}
	w.writer = make([]string, 0)
	return nil
}

func (w *FormSerializationWriter) WriteNullValue(key string) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	if key != "" {
		w.writePropertyName(key)
	}

	w.writeRawValue("null")

	if key != "" {
		w.writePropertySeparator()
	}

	return nil
}

func (w *FormSerializationWriter) GetOnBeforeSerialization() absser.ParsableAction {
	if w == nil {
		return nil
	}
	return w.onBeforeSerialization
}

func (w *FormSerializationWriter) SetOnBeforeSerialization(action absser.ParsableAction) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	w.onBeforeSerialization = action
	return nil
}

func (w *FormSerializationWriter) GetOnAfterObjectSerialization() absser.ParsableAction {
	if w == nil {
		return nil
	}
	return w.onAfterObjectSerialization
}

func (w *FormSerializationWriter) SetOnAfterObjectSerialization(action absser.ParsableAction) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	w.onAfterObjectSerialization = action
	return nil
}

func (w *FormSerializationWriter) GetOnStartObjectSerialization() absser.ParsableWriter {
	if w == nil {
		return nil
	}
	return w.onStartObjectSerialization
}

func (w *FormSerializationWriter) SetOnStartObjectSerialization(writer absser.ParsableWriter) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}
	w.onStartObjectSerialization = writer
	return nil
}
