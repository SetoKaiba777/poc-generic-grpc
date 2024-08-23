package genericclient

import (
	"context"
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type Client struct {
	conn  *grpc.ClientConn
	proto []*desc.FileDescriptor
	resp  *dynamicpb.Message
}

func NewClient(uri, protoPath string) (*Client, error) {
	conn, err := grpc.Dial(uri, grpc.WithInsecure())
	if err != nil {
		return &Client{}, fmt.Errorf("failed to connect in grpc client: %v", err)
	}
	parser := protoparse.Parser{}
	fds, err := parser.ParseFiles(protoPath)
	if err != nil {
		return &Client{},fmt.Errorf("failed to parse .proto file: %v", err)
	}
	return &Client{conn, fds, nil}, nil
}

func (c *Client) NewRequest(serviceName string, methodName string, data map[string]interface{}) error {
	fd := c.proto[0]
	serviceDesc := fd.FindService(serviceName)
	if serviceDesc == nil {
		return fmt.Errorf("service not found: %s", serviceName)
	}
	methodDesc := serviceDesc.FindMethodByName(methodName)
	if methodDesc == nil {
		return fmt.Errorf("method not found: %s", methodName)
	}

	requestMsg := dynamicpb.NewMessage(methodDesc.GetInputType().UnwrapMessage())
	if err := fillMessage(requestMsg, data); err != nil {
		return fmt.Errorf("failed to fill message: %v", err)
	}
	responseMessage := dynamicpb.NewMessage(methodDesc.GetOutputType().UnwrapMessage())
	err := c.conn.Invoke(context.Background(), fmt.Sprintf("/%s/%s", serviceName, methodName), requestMsg, responseMessage)
	if err != nil {
		return fmt.Errorf("failed to invoke method: %v", err)
	}
	c.resp = responseMessage
	return nil
}

func fillMessage(msg *dynamicpb.Message, data map[string]interface{}) error {
	for key, value := range data {
		field := msg.Descriptor().Fields().ByName(protoreflect.Name(key))
		if field == nil {
			return fmt.Errorf("field %s not found in message", key)
		}

		var v protoreflect.Value
		switch field.Kind() {
		case protoreflect.BoolKind:
			v = protoreflect.ValueOfBool(value.(bool))
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			v = protoreflect.ValueOfInt32(value.(int32))
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			v = protoreflect.ValueOfInt64(value.(int64))
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			v = protoreflect.ValueOfUint32(value.(uint32))
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			v = protoreflect.ValueOfUint64(value.(uint64))
		case protoreflect.FloatKind:
			v = protoreflect.ValueOfFloat32(value.(float32))
		case protoreflect.DoubleKind:
			v = protoreflect.ValueOfFloat64(value.(float64))
		case protoreflect.StringKind:
			v = protoreflect.ValueOfString(value.(string))
		case protoreflect.BytesKind:
			v = protoreflect.ValueOfBytes(value.([]byte))
		case protoreflect.MessageKind:
			nestedMsg := dynamicpb.NewMessage(field.Message())
			nestedData := value.(map[string]interface{})
			if err := fillMessage(nestedMsg, nestedData); err != nil {
				return err
			}
			v = protoreflect.ValueOfMessage(nestedMsg)
		default:
			return fmt.Errorf("unsupported field type for field %s", key)
		}

		msg.Set(field, v)
	}
	return nil
}

func (c *Client) PrintResponse() {
	fmt.Println(c.resp)
}
