// Package formserialization is the default Kiota serialization implementation for URI form encoded.
package formserialization

import (
	"encoding/base64"
	"errors"
	"math"
	"math/big"
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

// FormParseNode is a ParseNode implementation for URI form encoded.
type FormParseNode struct {
	value                     string
	values                    []string
	fields                    map[string]string
	rawFieldValues            map[string][]string
	onBeforeAssignFieldValues absser.ParsableAction
	onAfterAssignFieldValues  absser.ParsableAction
}

// NewFormParseNode creates a new FormParseNode.
func NewFormParseNode(content []byte) (*FormParseNode, error) {
	if len(content) == 0 {
		return nil, errors.New("content is empty")
	}
	rawValue := string(content)
	fields, rawFieldValues, err := loadFields(rawValue)
	if err != nil {
		return nil, err
	}
	return &FormParseNode{
		value:          rawValue,
		fields:         fields,
		rawFieldValues: rawFieldValues,
	}, nil
}
func loadFields(value string) (map[string]string, map[string][]string, error) {
	result := make(map[string]string)
	rawResult := make(map[string][]string)
	if len(value) == 0 {
		return result, rawResult, nil
	}
	parts := strings.Split(value, "&")
	for _, part := range parts {
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) == 2 {
			key, err := sanitizeKey(keyValue[0])
			if err != nil {
				return nil, nil, err
			}
			if key == "" {
				continue
			}
			rawResult[key] = append(rawResult[key], keyValue[1])
			if _, exists := result[key]; !exists {
				result[key] = keyValue[1]
			} else {
				result[key] += "," + keyValue[1]
			}
		}
	}
	return result, rawResult, nil
}
func sanitizeKey(key string) (string, error) {
	if key == "" {
		return "", nil
	}
	res, err := url.QueryUnescape(key)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res), nil
}

// GetChildNode returns a new parse node for the given identifier.
func (n *FormParseNode) GetChildNode(index string) (absser.ParseNode, error) {
	if n == nil {
		return nil, errors.New("node is nil")
	}
	trimmedIndex := strings.TrimSpace(index)
	if trimmedIndex == "" {
		return nil, errors.New("index is empty")
	}
	matchedKey := trimmedIndex
	fieldValue, ok := n.fields[matchedKey]
	if !ok {
		key, err := sanitizeKey(trimmedIndex)
		if err != nil {
			return nil, err
		}
		if key == "" {
			return nil, errors.New("index is empty")
		}
		matchedKey = key
		fieldValue, ok = n.fields[matchedKey]
	}
	if !ok {
		return nil, nil
	}

	var rawVals []string
	if n.rawFieldValues != nil {
		rawVals = n.rawFieldValues[matchedKey]
	}

	node := &FormParseNode{
		value:                     fieldValue,
		values:                    rawVals,
		onBeforeAssignFieldValues: n.GetOnBeforeAssignFieldValues(),
		onAfterAssignFieldValues:  n.GetOnAfterAssignFieldValues(),
	}
	return node, nil
}

func isUntypedDiscriminator(ctor absser.ParsableFactory) bool {
	if ctor == nil {
		return false
	}
	return reflect.ValueOf(ctor).Pointer() == reflect.ValueOf(absser.CreateUntypedNodeFromDiscriminatorValue).Pointer()
}

// GetObjectValue returns the Parsable value from the node.
func (n *FormParseNode) GetObjectValue(ctor absser.ParsableFactory) (absser.Parsable, error) {
	if ctor == nil {
		return nil, errors.New("constructor is nil")
	}
	if n == nil {
		return nil, nil
	}
	if isNullOrEmptyNode(n.value) {
		if isUntypedDiscriminator(ctor) {
			return n.getUntypedValue()
		}
		var isUntyped bool
		func() {
			defer func() {
				_ = recover()
			}()
			if res, err := ctor(n); err == nil && res != nil {
				if _, ok := res.(absser.UntypedNodeable); ok {
					isUntyped = true
				}
			}
		}()
		if isUntyped {
			return n.getUntypedValue()
		}
		return nil, nil
	}
	result, err := ctor(n)
	if err != nil {
		return nil, err
	}
	if _, isUntypedNode := result.(absser.UntypedNodeable); isUntypedNode {
		return n.getUntypedValue()
	}
	if result == nil {
		return nil, nil
	}
	abstractions.InvokeParsableAction(n.GetOnBeforeAssignFieldValues(), result)
	fields := result.GetFieldDeserializers()
	if len(n.fields) != 0 {
		itemAsHolder, isHolder := result.(absser.AdditionalDataHolder)
		var itemAdditionalData map[string]interface{}
		if isHolder {
			itemAdditionalData = itemAsHolder.GetAdditionalData()
			if itemAdditionalData == nil {
				itemAdditionalData = make(map[string]interface{})
				itemAsHolder.SetAdditionalData(itemAdditionalData)
			}
		}

		keys := make([]string, 0, len(n.fields))
		for k := range n.fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := n.fields[key]
			field := fields[key]
			if field == nil {
				if value != "" && isHolder {
					decodedValue, err := url.QueryUnescape(value)
					if err != nil {
						return nil, err
					}
					itemAdditionalData[key] = decodedValue
				}
			} else {
				childNode, err := n.GetChildNode(key)
				if err != nil {
					return nil, err
				}
				if childNode != nil {
					err = field(childNode)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	abstractions.InvokeParsableAction(n.GetOnAfterAssignFieldValues(), result)
	return result, nil
}

func isCanonicalInteger(s string) bool {
	if s == "0" {
		return true
	}
	if len(s) == 0 {
		return false
	}
	raw := s
	if raw[0] == '-' {
		raw = raw[1:]
	}
	if len(raw) == 0 || raw[0] < '1' || raw[0] > '9' {
		return false
	}
	for i := 1; i < len(raw); i++ {
		if raw[i] < '0' || raw[i] > '9' {
			return false
		}
	}
	return true
}

func isCanonicalFloat(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s[0] == '+' {
		return false
	}
	raw := s
	if raw[0] == '-' {
		raw = raw[1:]
	}
	if len(raw) == 0 {
		return false
	}

	i := 0
	if raw[0] == '0' {
		i = 1
	} else if raw[0] >= '1' && raw[0] <= '9' {
		i = 1
		for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
			i++
		}
	} else {
		return false
	}

	hasDot := false
	hasExp := false

	if i < len(raw) && raw[i] == '.' {
		hasDot = true
		i++
		if i >= len(raw) || raw[i] < '0' || raw[i] > '9' {
			return false
		}
		for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
			i++
		}
	}

	if i < len(raw) && (raw[i] == 'e' || raw[i] == 'E') {
		hasExp = true
		i++
		if i < len(raw) && (raw[i] == '+' || raw[i] == '-') {
			i++
		}
		if i >= len(raw) || raw[i] < '0' || raw[i] > '9' {
			return false
		}
		for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
			i++
		}
	}

	return i == len(raw) && (hasDot || hasExp)
}

func isZeroMantissa(s string) bool {
	mantissa := s
	if idx := strings.IndexAny(s, "eE"); idx != -1 {
		mantissa = s[:idx]
	}
	for i := 0; i < len(mantissa); i++ {
		c := mantissa[i]
		if c >= '1' && c <= '9' {
			return false
		}
	}
	return true
}

func (n *FormParseNode) getUntypedValue() (absser.Parsable, error) {
	if len(n.values) > 1 {
		elements := make([]absser.UntypedNodeable, len(n.values))
		for i, rawVal := range n.values {
			child := &FormParseNode{
				value:                     rawVal,
				onBeforeAssignFieldValues: n.GetOnBeforeAssignFieldValues(),
				onAfterAssignFieldValues:  n.GetOnAfterAssignFieldValues(),
			}
			untypedChild, err := child.getUntypedValue()
			if err != nil {
				return nil, err
			}
			if un, ok := untypedChild.(absser.UntypedNodeable); ok {
				elements[i] = un
			} else {
				elements[i] = absser.NewUntypedNull()
			}
		}
		arr := absser.NewUntypedArray(elements)
		abstractions.InvokeParsableAction(n.GetOnBeforeAssignFieldValues(), arr)
		abstractions.InvokeParsableAction(n.GetOnAfterAssignFieldValues(), arr)
		return arr, nil
	}

	if len(n.fields) > 0 {
		properties := make(map[string]absser.UntypedNodeable)
		untypedObj := absser.NewUntypedObject(properties)
		abstractions.InvokeParsableAction(n.GetOnBeforeAssignFieldValues(), untypedObj)
		keys := make([]string, 0, len(n.fields))
		for k := range n.fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			childNode, err := n.GetChildNode(key)
			if err != nil {
				return nil, err
			}
			if childNode == nil {
				properties[key] = absser.NewUntypedNull()
				continue
			}
			cfn, ok := childNode.(*FormParseNode)
			if !ok {
				properties[key] = absser.NewUntypedNull()
				continue
			}
			untypedChild, err := cfn.getUntypedValue()
			if err != nil {
				return nil, err
			}
			if un, ok := untypedChild.(absser.UntypedNodeable); ok {
				properties[key] = un
			} else {
				properties[key] = absser.NewUntypedNull()
			}
		}
		abstractions.InvokeParsableAction(n.GetOnAfterAssignFieldValues(), untypedObj)
		return untypedObj, nil
	}

	scalarNode, err := n.getUntypedScalarValue()
	if err != nil || scalarNode == nil {
		return scalarNode, err
	}
	abstractions.InvokeParsableAction(n.GetOnBeforeAssignFieldValues(), scalarNode)
	abstractions.InvokeParsableAction(n.GetOnAfterAssignFieldValues(), scalarNode)
	return scalarNode, nil
}

func (n *FormParseNode) getUntypedScalarValue() (absser.Parsable, error) {
	decoded, err := url.QueryUnescape(n.value)
	if err != nil {
		return nil, err
	}

	// Preserve decoded string when it has leading or trailing whitespace instead of coercing to bool, number, or null
	if strings.TrimSpace(decoded) != decoded {
		return absser.NewUntypedString(decoded), nil
	}

	if isNullValue(decoded) {
		return absser.NewUntypedNull(), nil
	}

	trimmed := decoded
	if strings.EqualFold(trimmed, "true") {
		b := true
		return absser.NewUntypedBoolean(b), nil
	}
	if strings.EqualFold(trimmed, "false") {
		b := false
		return absser.NewUntypedBoolean(b), nil
	}

	if isCanonicalInteger(trimmed) {
		if i32, err := strconv.ParseInt(trimmed, 10, 32); err == nil {
			v32 := int32(i32)
			return absser.NewUntypedInteger(v32), nil
		}
		if i64, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
			longNode := absser.NewUntypedLong(i64)
			return longNode, nil
		}
		// Integer token out of int64 range: preserve as UntypedString without float precision loss
		return absser.NewUntypedString(decoded), nil
	}

	if isCanonicalFloat(trimmed) {
		if f64, err := strconv.ParseFloat(trimmed, 64); err == nil && !math.IsNaN(f64) && !math.IsInf(f64, 0) {
			if f64 != 0 || isZeroMantissa(trimmed) {
				return absser.NewUntypedDouble(f64), nil
			}
		}
	}

	return absser.NewUntypedString(decoded), nil
}

// GetCollectionOfObjectValues returns the collection of Parsable values from the node.
func (n *FormParseNode) GetCollectionOfObjectValues(ctor absser.ParsableFactory) ([]absser.Parsable, error) {
	if n == nil || isNullOrEmptyNode(n.value) {
		return nil, nil
	}
	return nil, errors.New("collections are not supported in form serialization")
}

// GetCollectionOfPrimitiveValues returns the collection of primitive values from the node.
func (n *FormParseNode) GetCollectionOfPrimitiveValues(targetType string) ([]interface{}, error) {
	if n == nil || isNullOrEmptyNode(n.value) {
		return nil, nil
	}
	trimmedType := strings.TrimSpace(targetType)
	if trimmedType == "" {
		return nil, errors.New("targetType is empty")
	}
	valueList := strings.Split(n.value, ",")

	result := make([]interface{}, len(valueList))
	for i, element := range valueList {
		val, err := getPrimitiveValue(element, trimmedType)
		if err != nil {
			return nil, err
		}
		if isNil(val) {
			result[i] = nil
		} else {
			result[i] = val
		}
	}
	return result, nil
}

func getPrimitiveValue(value string, targetType string) (interface{}, error) {
	switch strings.ToLower(strings.TrimSpace(targetType)) {
	case "string":
		return getStringValue(value)
	case "bool", "boolean":
		return getBoolValue(value)
	case "int8", "uint8", "sbyte":
		return getInt8Value(value)
	case "byte":
		return getByteValue(value)
	case "float32", "float":
		return getFloat32Value(value)
	case "float64", "double", "decimal":
		return getFloat64Value(value)
	case "int32":
		return getInt32Value(value)
	case "int64", "int", "integer":
		return getInt64Value(value)
	case "time", "datetime", "date-time":
		return getTimeValue(value)
	case "timeonly", "time-only":
		return getTimeOnlyValue(value)
	case "dateonly", "date", "date-only":
		return getDateOnlyValue(value)
	case "isoduration", "duration":
		return getISODurationValue(value)
	case "uuid", "guid":
		return getUUIDValue(value)
	case "base64", "[]byte", "bytearray", "binary":
		return getByteArrayValue(value)
	default:
		return nil, errors.New("targetType is not supported")
	}
}

// GetCollectionOfEnumValues returns the collection of Enum values from the node.
func (n *FormParseNode) GetCollectionOfEnumValues(parser absser.EnumFactory) ([]interface{}, error) {
	if n == nil || isNullOrEmptyNode(n.value) {
		return nil, nil
	}
	if parser == nil {
		return nil, errors.New("parser is nil")
	}
	rawValues := strings.Split(n.value, ",")
	result := make([]interface{}, len(rawValues))
	for i, rawValue := range rawValues {
		node := &FormParseNode{
			value: rawValue,
		}
		val, err := node.GetEnumValue(parser)
		if err != nil {
			return nil, err
		}
		if isNil(val) {
			result[i] = nil
		} else {
			result[i] = val
		}
	}
	return result, nil
}

// GetStringValue returns a String value from the nodes.
func (n *FormParseNode) GetStringValue() (*string, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getStringValue(n.value)
}

func isNullOrEmptyNode(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || isNullValue(trimmed) {
		return true
	}
	if !strings.ContainsAny(value, "%+") {
		return false
	}
	res, err := url.QueryUnescape(value)
	if err != nil {
		return false
	}
	trimmed = strings.TrimSpace(res)
	return trimmed == "" || isNullValue(trimmed)
}

func isNullValue(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "null")
}

func getStringValue(value string) (*string, error) {
	res, err := url.QueryUnescape(value)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func getTrimmedValue(value string) (*string, error) {
	v, err := getStringValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(*v)
	if trimmed == "" || isNullValue(trimmed) {
		return nil, nil
	}
	return &trimmed, nil
}

// GetBoolValue returns a Bool value from the nodes.
func (n *FormParseNode) GetBoolValue() (*bool, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getBoolValue(n.value)
}

func getBoolValue(value string) (*bool, error) {
	v, err := getTrimmedValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	res, err := strconv.ParseBool(*v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetInt8Value returns a int8 value from the nodes.
func (n *FormParseNode) GetInt8Value() (*int8, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getInt8Value(n.value)
}

func getInt8Value(value string) (*int8, error) {
	v, err := getTrimmedValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	val, err := parseNumberAsInt64(*v)
	if err != nil {
		return nil, err
	}
	if val > math.MaxInt8 || val < math.MinInt8 {
		return nil, &strconv.NumError{Func: "ParseInt", Num: value, Err: strconv.ErrRange}
	}
	cast := int8(val)
	return &cast, nil
}

// GetByteValue returns a Byte value from the nodes.
func (n *FormParseNode) GetByteValue() (*byte, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getByteValue(n.value)
}

func getByteValue(value string) (*byte, error) {
	v, err := getTrimmedValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	val, err := parseNumberAsInt64(*v)
	if err != nil {
		return nil, err
	}
	if val > math.MaxUint8 || val < 0 {
		return nil, &strconv.NumError{Func: "ParseUint", Num: value, Err: strconv.ErrRange}
	}
	cast := byte(val)
	return &cast, nil
}

// GetFloat32Value returns a Float32 value from the nodes.
func (n *FormParseNode) GetFloat32Value() (*float32, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getFloat32Value(n.value)
}

func getFloat32Value(value string) (*float32, error) {
	v, err := getFloat64Value(value)
	if err != nil || v == nil {
		return nil, err
	}
	cast := float32(*v)
	return &cast, nil
}

// GetFloat64Value returns a Float64 value from the nodes.
func (n *FormParseNode) GetFloat64Value() (*float64, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getFloat64Value(n.value)
}

func getFloat64Value(value string) (*float64, error) {
	v, err := getTrimmedValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	res, err := strconv.ParseFloat(*v, 64)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetInt32Value returns a Int32 value from the nodes.
func (n *FormParseNode) GetInt32Value() (*int32, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getInt32Value(n.value)
}

func parseNumberAsInt64(value string) (int64, error) {
	if strings.ContainsAny(value, ".eE") {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, err
		}
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return 0, &strconv.NumError{Func: "ParseInt", Num: value, Err: strconv.ErrRange}
		}
		r, ok := new(big.Rat).SetString(value)
		if !ok {
			return 0, &strconv.NumError{Func: "ParseInt", Num: value, Err: strconv.ErrSyntax}
		}
		var intPart big.Int
		intPart.Quo(r.Num(), r.Denom())
		if !intPart.IsInt64() {
			return 0, &strconv.NumError{Func: "ParseInt", Num: value, Err: strconv.ErrRange}
		}
		return intPart.Int64(), nil
	}
	// Parse integer tokens directly before floating-point coercion can lose precision or accept non-integers.
	return strconv.ParseInt(value, 10, 64)
}

func getInt32Value(value string) (*int32, error) {
	v, err := getTrimmedValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	val, err := parseNumberAsInt64(*v)
	if err != nil {
		return nil, err
	}
	if val > math.MaxInt32 || val < math.MinInt32 {
		return nil, &strconv.NumError{Func: "ParseInt", Num: value, Err: strconv.ErrRange}
	}
	cast := int32(val)
	return &cast, nil
}

// GetInt64Value returns a Int64 value from the nodes.
func (n *FormParseNode) GetInt64Value() (*int64, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getInt64Value(n.value)
}

func getInt64Value(value string) (*int64, error) {
	v, err := getTrimmedValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	val, err := parseNumberAsInt64(*v)
	if err != nil {
		return nil, err
	}
	return &val, nil
}

// GetTimeValue returns a Time value from the nodes.
func (n *FormParseNode) GetTimeValue() (*time.Time, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getTimeValue(n.value)
}

func getTimeValue(value string) (*time.Time, error) {
	v, err := getTrimmedValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	parsed, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		parsed, err = time.Parse(time.RFC3339Nano, *v)
		if err != nil {
			return nil, err
		}
	}
	return &parsed, nil
}

// GetISODurationValue returns a ISODuration value from the nodes.
func (n *FormParseNode) GetISODurationValue() (*absser.ISODuration, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getISODurationValue(n.value)
}

func getISODurationValue(value string) (*absser.ISODuration, error) {
	v, err := getTrimmedValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	return absser.ParseISODuration(*v)
}

// GetTimeOnlyValue returns a TimeOnly value from the nodes.
func (n *FormParseNode) GetTimeOnlyValue() (*absser.TimeOnly, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getTimeOnlyValue(n.value)
}

func getTimeOnlyValue(value string) (*absser.TimeOnly, error) {
	v, err := getTrimmedValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	return absser.ParseTimeOnly(*v)
}

// GetDateOnlyValue returns a DateOnly value from the nodes.
func (n *FormParseNode) GetDateOnlyValue() (*absser.DateOnly, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getDateOnlyValue(n.value)
}

func getDateOnlyValue(value string) (*absser.DateOnly, error) {
	v, err := getTrimmedValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	return absser.ParseDateOnly(*v)
}

// GetUUIDValue returns a UUID value from the nodes.
func (n *FormParseNode) GetUUIDValue() (*uuid.UUID, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getUUIDValue(n.value)
}

func getUUIDValue(value string) (*uuid.UUID, error) {
	v, err := getTrimmedValue(value)
	if err != nil || v == nil {
		return nil, err
	}
	parsed, err := uuid.Parse(*v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// GetEnumValue returns a Enum value from the nodes.
func (n *FormParseNode) GetEnumValue(parser absser.EnumFactory) (interface{}, error) {
	if parser == nil {
		return nil, errors.New("parser is nil")
	}
	s, err := n.GetStringValue()
	if err != nil || s == nil {
		return nil, err
	}
	if isNullValue(*s) {
		return nil, nil
	}
	return parser(*s)
}

// GetByteArrayValue returns a ByteArray value from the nodes.
func (n *FormParseNode) GetByteArrayValue() ([]byte, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	return getByteArrayValue(n.value)
}

func getByteArrayValue(value string) ([]byte, error) {
	res, err := url.QueryUnescape(value)
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(res)
	if trimmed == "" {
		return nil, nil
	}
	return base64.StdEncoding.DecodeString(trimmed)
}

// GetRawValue returns a pointer to the raw string value from the node, or nil if the node is empty or null.
func (n *FormParseNode) GetRawValue() (interface{}, error) {
	if n == nil || isNullOrEmptyNode(n.value) {
		return nil, nil
	}
	res := n.value
	return &res, nil
}

func (n *FormParseNode) GetOnBeforeAssignFieldValues() absser.ParsableAction {
	if n == nil {
		return nil
	}
	return n.onBeforeAssignFieldValues
}

func (n *FormParseNode) SetOnBeforeAssignFieldValues(action absser.ParsableAction) error {
	if n == nil {
		return errors.New("node is nil")
	}
	n.onBeforeAssignFieldValues = action
	return nil
}

func (n *FormParseNode) GetOnAfterAssignFieldValues() absser.ParsableAction {
	if n == nil {
		return nil
	}
	return n.onAfterAssignFieldValues
}

func (n *FormParseNode) SetOnAfterAssignFieldValues(action absser.ParsableAction) error {
	if n == nil {
		return errors.New("node is nil")
	}
	n.onAfterAssignFieldValues = action
	return nil
}
