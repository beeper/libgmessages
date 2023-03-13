package pblite

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func typeError(name, typ string, value interface{}) error {
	return fmt.Errorf("field %q is not a %s, found %T", name, typ, value)
}

func unmarshalScalar(field protoreflect.FieldDescriptor, value interface{}) (protoreflect.Value, error) {
	var prValue protoreflect.Value

	name := string(field.Name())

	switch kind := field.Kind(); kind {
	case protoreflect.BoolKind:
		tmp, ok := value.(bool)
		if !ok {
			return prValue, typeError(name, "boolean", value)
		}
		prValue = protoreflect.ValueOfBool(tmp)
	case protoreflect.StringKind:
		tmp, ok := value.(string)
		if !ok {
			return prValue, typeError(name, "string", value)
		}
		prValue = protoreflect.ValueOfString(tmp)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		tmp, ok := value.(float64)
		if !ok {
			return prValue, typeError(name, "number", value)
		}
		prValue = protoreflect.ValueOfInt32(int32(tmp))
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
		protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		// 64-bit integers are treated as strings
		tmps, ok := value.(string)
		if !ok {
			return prValue, typeError(name, "int64", value)
		}

		tmpi, err := strconv.ParseInt(tmps, 10, 64)
		if err != nil {
			return prValue, err
		}

		prValue = protoreflect.ValueOfInt64(tmpi)
	case protoreflect.BytesKind:
		tmps, ok := value.(string)
		if !ok {
			return prValue, typeError(name, "string", value)
		}

		data, err := base64.StdEncoding.DecodeString(tmps)
		if err != nil {
			return prValue, err
		}

		prValue = protoreflect.ValueOfBytes(data)
	case protoreflect.EnumKind:
		tmp, ok := value.(float64)
		if !ok {
			return prValue, typeError(name, "number", value)
		}

		prValue = protoreflect.ValueOfEnum(protoreflect.EnumNumber(tmp))
	default:
		return prValue, fmt.Errorf("field %s has unhandled field of type %s", field.Name(), kind)
	}

	return prValue, nil
}

func unmarshalSingular(m protoreflect.Message, field protoreflect.FieldDescriptor, value interface{}) error {
	var prValue protoreflect.Value
	var err error

	switch field.Kind() {
	case protoreflect.MessageKind:
		tmp, ok := value.([]interface{})
		if !ok {
			return typeError(string(field.Name()), "array", value)
		}

		prValue = m.NewField(field)
		if err := unmarshalMessage(prValue.Message(), tmp); err != nil {
			return err
		}
	default:
		prValue, err = unmarshalScalar(field, value)
	}

	if err != nil {
		return err
	}

	m.Set(field, prValue)

	return nil
}

func unmarshalList(list protoreflect.List, field protoreflect.FieldDescriptor, value interface{}) error {
	array, ok := value.([]interface{})
	if !ok {
		return typeError(string(field.Name()), "array", value)
	}

	switch field.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		for _, arrayItem := range array {
			val := list.NewElement()
			if err := unmarshalMessage(val.Message(), arrayItem.([]interface{})); err != nil {
				return err
			}
			list.Append(val)
		}
	default:
		for _, arrayItem := range array {
			val, err := unmarshalScalar(field, arrayItem)
			if err != nil {
				return err
			}
			list.Append(val)
		}
	}

	return nil
}

func unmarshalMessage(m protoreflect.Message, data []interface{}) error {
	fields := m.Descriptor().Fields()

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldNumber := int(field.Number())

		if fieldNumber > len(data) {
			return fmt.Errorf(
				"field %s number is out of bounds: %d > %d)",
				field.Name(), fieldNumber, len(data))
		}

		value := data[fieldNumber-1]
		switch {
		case field.IsList():
			list := m.Mutable(field).List()
			if err := unmarshalList(list, field, value); err != nil {
				return err
			}
		case field.IsMap():
			return fmt.Errorf("maps are not currently supported")
		default:
			if err := unmarshalSingular(m, field, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func Unmarshal(data []byte, m proto.Message) error {
	raw := []interface{}{}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	return unmarshalMessage(m.ProtoReflect(), raw)
}
