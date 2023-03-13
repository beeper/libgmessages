# pblite

pblite is a codec for protobufs where the encoded data has keys removed and is
packed into order dependent JSON arrays. All option fields with numbers before
a set field must but included as the array position maps directly to the field
number.

For example say you have the following `.proto` file.

```protobuf
message Simple {
	optional int32 Foo = 1;
	optional int32 Bar = 2;
}
```

We'll create a message where we only set field 2.

```golang
simple := &pb.Simple {
	Bar: int32ptr(42),
}
```

This will be encoded as the following JSON array.

```json
[0, 42]
```

This obviously gets a lot more complicated with messages and lists, but the
concept is the same regardless. If you only set the first value, the second
value can be skipped as it is optional.

## errata

Currently maps are not implemented as there hasn't been a need yet.
