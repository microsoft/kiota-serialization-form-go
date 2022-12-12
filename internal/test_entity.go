package internal

import (
	"time"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TestEntity struct {
	additionalData map[string]interface{}
	id             *string
	officeLocation *string
	//TODO test numbers enum
	workDuration    *absser.ISODuration
	birthDay        *absser.DateOnly
	startWorkTime   *absser.TimeOnly
	endWorkTime     *absser.TimeOnly
	createdDateTime *time.Time
}

type TestEntityable interface {
	absser.Parsable
	absser.AdditionalDataHolder
	GetId() *string
	SetId(value *string)
	GetOfficeLocation() *string
	SetOfficeLocation(value *string)
	GetWorkDuration() *absser.ISODuration
	SetWorkDuration(value *absser.ISODuration)
	GetBirthDay() *absser.DateOnly
	SetBirthDay(value *absser.DateOnly)
	GetStartWorkTime() *absser.TimeOnly
	SetStartWorkTime(value *absser.TimeOnly)
	GetEndWorkTime() *absser.TimeOnly
	SetEndWorkTime(value *absser.TimeOnly)
	GetCreatedDateTime() *time.Time
	SetCreatedDateTime(value *time.Time)
}

func NewTestEntity() *TestEntity {
	return &TestEntity{
		additionalData: make(map[string]interface{}),
	}
}

func (e *TestEntity) GetAdditionalData() map[string]interface{} {
	return e.additionalData
}
func (e *TestEntity) SetAdditionalData(value map[string]interface{}) {
	e.additionalData = value
}
func (e *TestEntity) SetWorkDuration(value *absser.ISODuration) {
	e.workDuration = value
}
func (e *TestEntity) GetWorkDuration() *absser.ISODuration {
	return e.workDuration
}
func (e *TestEntity) GetBirthDay() *absser.DateOnly {
	return e.birthDay
}
func (e *TestEntity) SetBirthDay(value *absser.DateOnly) {
	e.birthDay = value
}
func (e *TestEntity) GetEndWorkTime() *absser.TimeOnly {
	return e.endWorkTime
}
func (e *TestEntity) SetEndWorkTime(value *absser.TimeOnly) {
	e.endWorkTime = value
}
func (e *TestEntity) GetStartWorkTime() *absser.TimeOnly {
	return e.startWorkTime
}
func (e *TestEntity) SetStartWorkTime(value *absser.TimeOnly) {
	e.startWorkTime = value
}
func (e *TestEntity) GetId() *string {
	return e.id
}
func (e *TestEntity) SetId(value *string) {
	e.id = value
}
func (e *TestEntity) GetOfficeLocation() *string {
	return e.officeLocation
}
func (e *TestEntity) SetOfficeLocation(value *string) {
	e.officeLocation = value
}
func (e *TestEntity) GetCreatedDateTime() *time.Time {
	return e.createdDateTime
}
func (e *TestEntity) SetCreatedDateTime(value *time.Time) {
	e.createdDateTime = value
}

func CreateTestEntityFromDiscriminator(parseNode absser.ParseNode) (absser.Parsable, error) {
	return NewTestEntity(), nil
}

func (e *TestEntity) GetFieldDeserializers() map[string]func(absser.ParseNode) error {
	res := make(map[string]func(absser.ParseNode) error)
	res["id"] = func(n absser.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetId(val)
		}
		return nil
	}
	res["officeLocation"] = func(n absser.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetOfficeLocation(val)
		}
		return nil
	}
	res["workDuration"] = func(n absser.ParseNode) error {
		val, err := n.GetISODurationValue()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetWorkDuration(val)
		}
		return nil
	}
	res["birthDay"] = func(n absser.ParseNode) error {
		val, err := n.GetDateOnlyValue()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetBirthDay(val)
		}
		return nil
	}
	res["startWorkTime"] = func(n absser.ParseNode) error {
		val, err := n.GetTimeOnlyValue()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetStartWorkTime(val)
		}
		return nil
	}
	res["endWorkTime"] = func(n absser.ParseNode) error {
		val, err := n.GetTimeOnlyValue()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetEndWorkTime(val)
		}
		return nil
	}
	res["createdDateTime"] = func(n absser.ParseNode) error {
		val, err := n.GetTimeValue()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetCreatedDateTime(val)
		}
		return nil
	}
	return res
}

func (m *TestEntity) Serialize(writer absser.SerializationWriter) error {
	{
		err := writer.WriteStringValue("id", m.GetId())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteStringValue("officeLocation", m.GetOfficeLocation())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteISODurationValue("workDuration", m.GetWorkDuration())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteDateOnlyValue("birthDay", m.GetBirthDay())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteTimeOnlyValue("startWorkTime", m.GetStartWorkTime())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteTimeOnlyValue("endWorkTime", m.GetEndWorkTime())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
		if err != nil {
			return err
		}
	}
	return nil
}
