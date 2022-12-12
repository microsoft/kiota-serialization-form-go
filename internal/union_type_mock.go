package internal

import (
	"strings"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnionTypeMock struct {
	composedType1 TestEntityable
	composedType2 SecondTestEntityable
	stringValue   *string
	composedType3 []TestEntityable
}
type UnionTypeMockable interface {
	absser.Parsable
	GetComposedType1() TestEntityable
	SetComposedType1(value TestEntityable)
	GetComposedType2() SecondTestEntityable
	SetComposedType2(value SecondTestEntityable)
	GetStringValue() *string
	SetStringValue(value *string)
	GetComposedType3() []TestEntityable
	SetComposedType3(value []TestEntityable)
}

func NewUnionTypeMock() *UnionTypeMock {
	return &UnionTypeMock{}
}
func (e *UnionTypeMock) GetComposedType1() TestEntityable {
	return e.composedType1
}
func (e *UnionTypeMock) SetComposedType1(value TestEntityable) {
	e.composedType1 = value
}
func (e *UnionTypeMock) GetComposedType2() SecondTestEntityable {
	return e.composedType2
}
func (e *UnionTypeMock) SetComposedType2(value SecondTestEntityable) {
	e.composedType2 = value
}
func (e *UnionTypeMock) GetStringValue() *string {
	return e.stringValue
}
func (e *UnionTypeMock) SetStringValue(value *string) {
	e.stringValue = value
}
func (e *UnionTypeMock) GetComposedType3() []TestEntityable {
	return e.composedType3
}
func (e *UnionTypeMock) SetComposedType3(value []TestEntityable) {
	e.composedType3 = value
}
func CreateUnionTypeMockFromDiscriminator(parseNode absser.ParseNode) (absser.Parsable, error) {
	result := NewUnionTypeMock()
	mappingValueNode, err := parseNode.GetChildNode("@odata.type")
	if err != nil {
		return nil, err
	}
	if mappingValueNode != nil {
		mappingValue, err := mappingValueNode.GetStringValue()
		if err != nil {
			return nil, err
		}
		if mappingValue != nil {
			if strings.EqualFold(*mappingValue, "#microsoft.graph.testEntity") {
				result.SetComposedType1(NewTestEntity())
			} else if strings.EqualFold(*mappingValue, "#microsoft.graph.secondTestEntity") {
				result.SetComposedType2(NewSecondTestEntity())
			}
		}
	}
	if val, err := parseNode.GetStringValue(); val != nil {
		if err != nil {
			return nil, err
		}
		result.SetStringValue(val)
	} else if val, err := parseNode.GetCollectionOfObjectValues(CreateTestEntityFromDiscriminator); val != nil {
		if err != nil {
			return nil, err
		}
		cast := make([]TestEntityable, len(val))
		for i, v := range val {
			cast[i] = v.(TestEntityable)
		}
		result.SetComposedType3(cast)
	}
	return result, nil
}
func (e *UnionTypeMock) GetFieldDeserializers() map[string]func(absser.ParseNode) error {
	if e.GetComposedType1() != nil {
		return e.GetComposedType1().GetFieldDeserializers()
	} else if e.GetComposedType2() != nil {
		return e.GetComposedType2().GetFieldDeserializers()
	} else {
		return make(map[string]func(absser.ParseNode) error)
	}
}
func (e *UnionTypeMock) Serialize(writer absser.SerializationWriter) error {
	if e.GetComposedType1() != nil {
		return writer.WriteObjectValue("", e.GetComposedType1())
	} else if e.GetComposedType2() != nil {
		return writer.WriteObjectValue("", e.GetComposedType2())
	} else if e.GetComposedType3() != nil {
		cast := make([]absser.Parsable, len(e.GetComposedType3()))
		for i, v := range e.GetComposedType3() {
			cast[i] = v.(absser.Parsable)
		}
		return writer.WriteCollectionOfObjectValues("", cast)
	} else if e.GetStringValue() != nil {
		return writer.WriteStringValue("", e.GetStringValue())
	}
	return nil
}
