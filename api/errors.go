package api

import (
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

type HttpError struct {
	Code int
	Err  error
	Type gin.ErrorType
}

func (e HttpError) Error() string {
	return e.Err.Error()
}

func (e HttpError) Abort(c *gin.Context) {
	c.AbortWithError(e.Code, e.Err).SetType(e.Type)
}

func NewHttpError(code int, err error, t gin.ErrorType) *HttpError {
	return &HttpError{Code: code, Err: err, Type: t}
}

func NewNotFoundError() *HttpError {
	return &HttpError{Code: http.StatusNotFound, Type: gin.ErrorTypePublic}
}

func NewBadRequestError(err error) *HttpError {
	return NewHttpError(http.StatusBadRequest, err, gin.ErrorTypePublic)
}

func NewForbiddenError() *HttpError {
	return &HttpError{Code: http.StatusForbidden, Type: gin.ErrorTypePublic}
}

func NewInternalError(err error) *HttpError {
	return NewHttpError(http.StatusInternalServerError, err, gin.ErrorTypePrivate)
}

type FileNotFound struct {
	Sha1 string
}

func (x FileNotFound) Error() string {
	return fmt.Sprintf("File not found, sha1 = %s", x.Sha1)
}

type UpChainOperationNotFound struct {
	FileID int64
	opType string
}

func (x UpChainOperationNotFound) Error() string {
	return fmt.Sprintf("Up chain operation %s not found for file_id %d", x.opType, x.FileID)
}

type CollectionNotFound struct {
	CaptureID interface{}
}

func (x CollectionNotFound) Error() string {
	return fmt.Sprintf("Collection not found, CaptureID = %s", x.CaptureID)
}
