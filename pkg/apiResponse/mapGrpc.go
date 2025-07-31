package apiresponse

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// map gRPC codes to string codes and use original message
func mapGRPC(st *status.Status) (string, string) {
	switch st.Code() {
	case codes.OK:
		return "SUCCESS", st.Message()
	case codes.Canceled:
		return "REQUEST_CANCELED", st.Message()
	case codes.Unknown:
		return "UNKNOWN_ERROR", st.Message()
	case codes.InvalidArgument:
		return "INVALID_INPUT", st.Message()
	case codes.DeadlineExceeded:
		return "TIMEOUT", st.Message()
	case codes.NotFound:
		return "NOT_FOUND", st.Message()
	case codes.AlreadyExists:
		return "ALREADY_EXISTS", st.Message()
	case codes.PermissionDenied:
		return "FORBIDDEN", st.Message()
	case codes.ResourceExhausted:
		return "RATE_LIMITED", st.Message()
	case codes.FailedPrecondition:
		return "PRECONDITION_FAILED", st.Message()
	case codes.Aborted:
		return "CONFLICT", st.Message()
	case codes.OutOfRange:
		return "OUT_OF_RANGE", st.Message()
	case codes.Unimplemented:
		return "NOT_IMPLEMENTED", st.Message()
	case codes.Internal:
		return "INTERNAL_SERVER_ERROR", st.Message()
	case codes.Unavailable:
		return "SERVICE_UNAVAILABLE", st.Message()
	case codes.DataLoss:
		return "DATA_CORRUPTION", st.Message()
	case codes.Unauthenticated:
		return "UNAUTHORIZED", st.Message()
	default:
		return "INTERNAL_SERVER_ERROR", st.Message()
	}
}
