package main

import (
	"fmt"

	"github.com/beeper/libgmessages/pb"
	"github.com/beeper/libgmessages/pblite"
)

func strptr(s string) *string {
	return &s
}

func int32ptr(x int32) *int32 {
	return &x
}

func main() {
	msg := &pb.ImageUploadRequest{
		Header: &pb.RequestHeader{
			Id:               strptr("requestId"),
			App:              strptr("requestApp"),
			AuthTokenPayload: []byte{0x01, 0x02, 0x03, 0x04},
			ClientInfo: &pb.ClientInfo{
				Major:        int32ptr(1),
				Minor:        int32ptr(2),
				Point:        int32ptr(3),
				ApiVersion:   int32ptr(100),
				PlatformType: pb.PlatformType_DESKTOP.Enum(),
			},
		},
		SenderId: &pb.SenderId{
			Type: pb.SenderId_DEVICE_ID.Enum(),
			Id:   strptr("senderId"),
			App:  strptr("senderApp"),
		},
	}

	fmt.Println("=== Input     ==========")
	fmt.Println(msg)

	fmt.Println("=== Marshal   ==========")
	data, err := pblite.Marshal(msg)
	fmt.Println("error:", err)
	fmt.Println("output:", string(data))

	fmt.Println("=== Unmarshal ==========")
	out := &pb.ImageUploadRequest{}
	err = pblite.Unmarshal(data, out)
	fmt.Println("error:", err)
	fmt.Println("output:", out)
}
