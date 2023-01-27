package utils

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"strings"
)

type Responses struct {
	StatusCode int         `json:"status_code"`
	Method     string      `json:"method"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type ErrorResponse struct {
	StatusCode int         `json:"status_code"`
	Method     string      `json:"method"`
	Message    string      `json:"message"`
	Error      interface{} `json:"error"`
}

var mapGrpcHttpCodes = map[codes.Code]int{
	codes.OK:                 200,
	codes.Canceled:           499,
	codes.Unknown:            500,
	codes.InvalidArgument:    400,
	codes.DeadlineExceeded:   504,
	codes.NotFound:           111, //Important error code,This error msg will be shown to client UI
	codes.AlreadyExists:      409,
	codes.PermissionDenied:   403,
	codes.Unauthenticated:    401,
	codes.ResourceExhausted:  429,
	codes.Aborted:            409,
	codes.OutOfRange:         400,
	codes.Unimplemented:      501,
	codes.Internal:           500,
	codes.Unavailable:        503,
	codes.DataLoss:           500,
	codes.FailedPrecondition: 400,
}
func APIResponse(ctx *gin.Context, Message string, StatusCode codes.Code, Method string, Data interface{}) {
	var msg,errMsg string
	errMsg = Message
	code, ok := mapGrpcHttpCodes[StatusCode]
	if !ok {
		errResponse := ErrorResponse{
			StatusCode:  500,
			Method: Method,
			Message:    "Internal Server Error",
			Error:       Message,
		}
		ctx.JSON(500, errResponse) //Http.InternalServerError code
		defer ctx.AbortWithStatus(500)
	}
	if strings.Contains(Message,"EXTRA"){
		extraMsg := strings.Split(Message,"%!(EXTRA string=") //keyword used to filter is present is status Message
		if len(extraMsg) > 1{
			extraMsgFilter := strings.Split(extraMsg[1],")")//Removing unnecessary )
			msg = extraMsgFilter[0]
			errMsg = extraMsg[0]
		}else{
			extraMsg = strings.Split(Message,"(EXTRA string=")
			extraMsgFilter := strings.Split(extraMsg[1],")")//Removing unnecessary )
			msg = extraMsgFilter[0]
			errMsg = extraMsg[0]//This case is used  when status.Errof is written twice for the same error
		}
	}
	if errMsg == "" {
		errMsg = "Something went wrong!"
	}
	if msg == "" {
		msg = "Something went wrong!"
	}
	switch code {
	case 504:
		errResponse := ErrorResponse{
			StatusCode: code,
			Method: Method,
			Message:  "Request Timed out",
			Error:     Message,
		}
		ctx.JSON(code, errResponse)
		defer ctx.AbortWithStatus(code)
	case 200:
		jsonResponse := Responses{
			StatusCode: code,
			Method:     Method,
			Message:    Message,
			Data:       Data,
		}
		ctx.JSON(code, jsonResponse)
		defer ctx.AbortWithStatus(code)
	case 111:
		errResponse := ErrorResponse{
			StatusCode: code,
			Method: Method,
			Message:  msg,
			Error:      errMsg,
		}
		ctx.JSON(400, errResponse)
		defer ctx.AbortWithStatus(400)
	default:
		errResponse := ErrorResponse{
			StatusCode: code,
			Method: Method,
			Message:  msg,
			Error:      errMsg,
		}
		ctx.JSON(code, errResponse)
		defer ctx.AbortWithStatus(code)
	}
}
