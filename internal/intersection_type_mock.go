package internal

import absser "github.com/microsoft/kiota-abstractions-go/serialization"

type IntersectionTypeMock struct {
	composedType1 TestEntityable
	composedType2 SecondTestEntityable
	stringValue   *string
	composedType3 []TestEntityable
}
type IntersectionTypeMockable interface {
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

func NewIntersectionTypeMock() *IntersectionTypeMock {
	return &IntersectionTypeMock{}
}
func (e *IntersectionTypeMock) GetComposedType1() TestEntityable {
	return e.composedType1
}
func (e *IntersectionTypeMock) SetComposedType1(value TestEntityable) {
	e.composedType1 = value
}
func (e *IntersectionTypeMock) GetComposedType2() SecondTestEntityable {
	return e.composedType2
}
func (e *IntersectionTypeMock) SetComposedType2(value SecondTestEntityable) {
	e.composedType2 = value
}
func (e *IntersectionTypeMock) GetStringValue() *string {
	return e.stringValue
}
func (e *IntersectionTypeMock) SetStringValue(value *string) {
	e.stringValue = value
}
func (e *IntersectionTypeMock) GetComposedType3() []TestEntityable {
	return e.composedType3
}
func (e *IntersectionTypeMock) SetComposedType3(value []TestEntityable) {
	e.composedType3 = value
}
func CreateIntersectionTypeMockFromDiscriminator(parseNode absser.ParseNode) (absser.Parsable, error) {
	result := NewIntersectionTypeMock()
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
	} else {
		result.SetComposedType1(NewTestEntity())
		result.SetComposedType2(NewSecondTestEntity())
	}
	return result, nil
}
func (e *IntersectionTypeMock) GetFieldDeserializers() map[string]func(absser.ParseNode) error {
	if e.GetComposedType1() != nil || e.GetComposedType2() != nil {
		res, err := absser.MergeDeserializersForIntersectionWrapper(e.GetComposedType1(), e.GetComposedType2())
		if err != nil {
			return nil
		}
		return res
	}
	return make(map[string]func(absser.ParseNode) error)
}
func (e *IntersectionTypeMock) Serialize(writer absser.SerializationWriter) error {
	if e.GetStringValue() != nil {
		return writer.WriteStringValue("", e.GetStringValue())
	} else if e.GetComposedType3() != nil {
		cast := make([]absser.Parsable, len(e.GetComposedType3()))
		for i, v := range e.GetComposedType3() {
			cast[i] = v.(absser.Parsable)
		}
		return writer.WriteCollectionOfObjectValues("", cast)
	} else {
		return writer.WriteObjectValue("", e.GetComposedType1(), e.GetComposedType2())
	}
}
