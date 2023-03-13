package pblite

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func marshalList(list protoreflect.List, fd protoreflect.FieldDescriptor) (interface{}, error) {
	r := make([]interface{}, list.Len())

	for i := 0; i < list.Len(); i++ {
		item := list.Get(i)

		result, err := marshalSingular(item, fd)
		if err != nil {
			return nil, err
		}

		r[i] = result
	}

	return r, nil
}

func marshalMap(list protoreflect.Map, fd protoreflect.FieldDescriptor) (interface{}, error) {
	return nil, fmt.Errorf("marshalMap is not implemented")
}

func marshalSingular(value protoreflect.Value, fd protoreflect.FieldDescriptor) (interface{}, error) {
	if !value.IsValid() {
		fmt.Println("invalid")
		return nil, nil
	}

	switch kind := fd.Kind(); kind {
	case protoreflect.BoolKind:
		return value.Bool(), nil
	case protoreflect.StringKind:
		return value.String(), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return value.Int(), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return value.Uint(), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
		protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		// 64-bit integers are written out as strings
		return value.String(), nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return value.Float(), nil
	case protoreflect.BytesKind:
		return value.Bytes(), nil
	case protoreflect.EnumKind:
		return int64(value.Enum()), nil
	case protoreflect.MessageKind:
		return marshalMessage(value.Message())
	default:
		return nil, fmt.Errorf("field %v has unknown kind %v", fd.FullName(), kind)
	}

	return nil, nil
}

func marshalMessage(m protoreflect.Message) ([]interface{}, error) {
	// A nil message should be an empty array.
	if m == nil {
		return []interface{}{}, nil
	}

	// Grab the fields
	fields := m.Descriptor().Fields()

	// Get the last field and use it's Number() to determine how many fields
	// we need in our array.
	last := fields.Get(fields.Len() - 1)

	// Allocate our output array based on the number of fields we have.
	r := make([]interface{}, last.Number())

	// Now walk through the fields we have and marshal them into the proper
	// index in our result slice.
	var err error
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		var result interface{}

		switch {
		case fd.IsList():
			result, err = marshalList(v.List(), fd)
		case fd.IsMap():
			result, err = marshalMap(v.Map(), fd)
		default:
			result, err = marshalSingular(v, fd)
		}

		// This will break the loop and continue below with error set as it's
		// a scope higher than this.
		if err != nil {
			return false
		}

		r[fd.Number()-1] = result

		return true
	})

	if err != nil {
		return nil, err
	}

	return r, nil
}

func Marshal(m proto.Message) ([]byte, error) {
	if m == nil {
		return []byte("[]"), nil
	}

	raw, err := marshalMessage(m.ProtoReflect())
	if err != nil {
		return []byte("[]"), err
	}

	return json.Marshal(raw)
}
