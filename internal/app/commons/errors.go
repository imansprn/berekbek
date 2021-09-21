package commons

import (
	"errors"
	"fmt"
	"net/http"

	phttp "github.com/valbury-repos/gotik/http"
	"github.com/valbury-repos/gotik/structs"
	"github.com/gobliggg/berekbek/config"
)

var cfg = config.Config()

func ErrServiceCode(errCode string) string {
	return fmt.Sprintf("%s%s", cfg.GetString("app.code"), errCode)
}

// InjectErrors injecting all error response to the handler context
func InjectErrors(handlerCtx *phttp.HandlerContext) {
	handlerCtx.AddError(ErrDBConn, &ErrDBConnResp)
	handlerCtx.AddError(ErrCacheConn, &ErrCacheConnResp)
	handlerCtx.AddError(ErrBodyRequestInvalid, &ErrBodyRequestInvalidResp)
	handlerCtx.AddError(ErrDataNotFound, &ErrDataNotFoundResp)
	// etc...
}

func getErrorResponse(httpStatus int, message string, errorCode string) structs.ErrorResponse {
	return structs.ErrorResponse{
		Meta: structs.Meta{
			Status: httpStatus,
			Message: message,
			Errors: []structs.Errors{
				{
					ID: cfg.GetString(fmt.Sprintf("%s%s", "response_code.ID.", errorCode)),
					EN: cfg.GetString(fmt.Sprintf("%s%s", "response_code.EN.", errorCode)),
				},
			},
		},
	}
}

var (
	// ErrDBConn error type for Error DB Connection
	ErrDBConn = errors.New("ErrDBConn")
	// ErrDBConnResp ErrDBConn's response
	ErrDBConnResp = getErrorResponse(http.StatusInternalServerError, "", ErrServiceCode("1001"))
)

var (
	// ErrCacheConn error type for Error Cache Connection
	ErrCacheConn = errors.New("ErrCacheConn")

	// ErrCacheConnResp ErrCacheConn's response
	ErrCacheConnResp = getErrorResponse(http.StatusInternalServerError, "", ErrServiceCode("1002"))
)

var (
	// ErrBodyRequestInvalid error type for Error Invalid Body Request
	ErrBodyRequestInvalid = errors.New("ErrBodyRequestInvalid")

	// ErrBodyRequestInvalidResp ErrBodyRequestInvalid's response
	ErrBodyRequestInvalidResp = getErrorResponse(http.StatusBadRequest, "", ErrServiceCode("1003"))
)

var (
	// ErrDataNotFound error type for Data Not Found Error
	ErrDataNotFound = errors.New("ErrDataNotFound")
	// ErrUnknownResp ErrUnknown's response
	ErrDataNotFoundResp = getErrorResponse(http.StatusNotFound, "", ErrServiceCode("1004"))
)