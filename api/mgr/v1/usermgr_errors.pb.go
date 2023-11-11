// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

// 用户已经存在
func IsUsersvrBadResp(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_USERSVR_BAD_RESP.String() && e.Code == 403
}

// 用户已经存在
func ErrorUsersvrBadResp(format string, args ...interface{}) *errors.Error {
	return errors.New(403, ErrorReason_USERSVR_BAD_RESP.String(), fmt.Sprintf(format, args...))
}

// 入参错误
func IsErrInvalidParam(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_ERR_INVALID_PARAM.String() && e.Code == 400
}

// 入参错误
func ErrorErrInvalidParam(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_ERR_INVALID_PARAM.String(), fmt.Sprintf(format, args...))
}
