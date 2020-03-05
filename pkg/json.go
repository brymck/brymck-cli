package pkg

import (
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func PrintAsJson(message proto.Message) {
	m := jsonpb.Marshaler{}
	result, _ := m.MarshalToString(message)
	fmt.Println(result)
}
