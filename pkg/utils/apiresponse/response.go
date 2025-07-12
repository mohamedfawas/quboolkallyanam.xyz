package apiresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// Fail inspects err; if it's a gRPC error, maps its code→HTTP, otherwise 500.
// Writes JSON: { status, message, error }
func Fail(c *gin.Context, err error) {
	httpCode := http.StatusInternalServerError
	msg := "internal error"

	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.InvalidArgument:
			httpCode = http.StatusBadRequest
			msg = st.Message()
		case codes.Unauthenticated:
			httpCode = http.StatusUnauthorized
			msg = st.Message()
		case codes.PermissionDenied:
			httpCode = http.StatusForbidden
			msg = st.Message()
		case codes.NotFound:
			httpCode = http.StatusNotFound
			msg = st.Message()
		case codes.AlreadyExists:
			httpCode = http.StatusConflict
			msg = st.Message()
		case codes.Unavailable:
			httpCode = http.StatusServiceUnavailable
			msg = st.Message()
		default:
			httpCode = http.StatusInternalServerError
			msg = st.Message()
		}
	} else {
		// non-grpc errors: you can add application-specific checks here
		msg = err.Error()
	}

	c.JSON(httpCode, Response{
		Status:  httpCode,
		Message: msg,
		Error:   msg,
	})
	c.Abort()
}
