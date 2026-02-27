package grpc

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// THBValidationInterceptor ensures that any field related to currency is strictly "THB" or empty.
func THBValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		msg, ok := req.(proto.Message)
		if ok {
			if err := checkTHBConstraints(msg.ProtoReflect()); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "THB Constraint Error: %v", err)
			}
		}
		return handler(ctx, req)
	}
}

// checkTHBConstraints recursively inspects protoreflect.Message for currency fields.
func checkTHBConstraints(m protoreflect.Message) error {
	var walkErr error
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if walkErr != nil {
			return false
		}

		if fd.IsList() {
			list := v.List()
			for i := 0; i < list.Len(); i++ {
				item := list.Get(i)
				if fd.Kind() == protoreflect.MessageKind {
					if err := checkTHBConstraints(item.Message()); err != nil {
						walkErr = err
						return false
					}
				}
			}
		} else if fd.IsMap() {
			// Maps generally aren't used for the complex nested DTOs here, but could be handled if needed
		} else if fd.Kind() == protoreflect.MessageKind {
			if err := checkTHBConstraints(v.Message()); err != nil {
				walkErr = err
				return false
			}
		} else if fd.Kind() == protoreflect.StringKind {
			// Check if field name implies currency
			name := strings.ToLower(string(fd.Name()))
			if strings.Contains(name, "cry") || strings.Contains(name, "currency") {
				val := v.String()
				// Only allow THB or empty (if optional)
				// We also check for 'Y' or 'N' because some fields might be flags like IsOthTxCry
				if val != "" && val != "THB" && val != "Y" && val != "N" && !strings.Contains(val, "THB") {
					// Some fields might be "THB,USD" (FundCrySet), we must reject if it contains anything other than THB
					parts := strings.Split(val, ",")
					for _, p := range parts {
						p = strings.TrimSpace(p)
						if p != "" && p != "THB" {
							// Exclude specific known flags that might coincidentally match "cry"
							if name == "iscry" || name == "is_oth_tx_cry" {
								continue
							}
							walkErr = fmt.Errorf("field %s contains unsupported currency: %s (only THB is supported)", fd.Name(), val)
							return false
						}
					}
				}
			}
		}
		return true
	})
	return walkErr
}
