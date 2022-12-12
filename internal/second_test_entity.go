package internal

import absser "github.com/microsoft/kiota-abstractions-go/serialization"

type SecondTestEntity struct {
	displayName *string
	id          *int64
	failureRate *float64
}

type SecondTestEntityable interface {
	absser.Parsable
	GetDisplayName() *string
	SetDisplayName(value *string)
	GetId() *int64
	SetId(value *int64)
	GetFailureRate() *float64
	SetFailureRate(value *float64)
}

func NewSecondTestEntity() *SecondTestEntity {
	return &SecondTestEntity{}
}

func (e *SecondTestEntity) GetDisplayName() *string {
	return e.displayName
}
func (e *SecondTestEntity) SetDisplayName(value *string) {
	e.displayName = value
}
func (e *SecondTestEntity) GetId() *int64 {
	return e.id
}
func (e *SecondTestEntity) SetId(value *int64) {
	e.id = value
}
func (e *SecondTestEntity) GetFailureRate() *float64 {
	return e.failureRate
}
func (e *SecondTestEntity) SetFailureRate(value *float64) {
	e.failureRate = value
}

func CreateSecondTestEntityFromDiscriminator(parseNode absser.ParseNode) (absser.Parsable, error) {
	return NewSecondTestEntity(), nil
}

func (e *SecondTestEntity) GetFieldDeserializers() map[string]func(absser.ParseNode) error {
	res := make(map[string]func(absser.ParseNode) error)
	res["id"] = func(n absser.ParseNode) error {
		val, err := n.GetInt64Value()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetId(val)
		}
		return nil
	}
	res["displayName"] = func(n absser.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetDisplayName(val)
		}
		return nil
	}
	res["failureRate"] = func(n absser.ParseNode) error {
		val, err := n.GetFloat64Value()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetFailureRate(val)
		}
		return nil
	}
	return res
}
func (m *SecondTestEntity) Serialize(writer absser.SerializationWriter) error {
	{
		err := writer.WriteInt64Value("id", m.GetId())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteStringValue("displayName", m.GetDisplayName())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteFloat64Value("failureRate", m.GetFailureRate())
		if err != nil {
			return err
		}
	}
	return nil
}
