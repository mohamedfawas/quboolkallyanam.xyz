package errutil

import "google.golang.org/grpc/status"

// IsGRPCError returns true if err is a gRPC Status error.
func IsGRPCError(err error) bool {
    if err == nil {
        return false
    }
    _, ok := status.FromError(err)
    return ok
}
