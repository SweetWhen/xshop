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
func IsUserExisted(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_USER_EXISTED.String() && e.Code == 400
}

// 用户已经存在
func ErrorUserExisted(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_USER_EXISTED.String(), fmt.Sprintf(format, args...))
}

// 用户找不到
func IsUserNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_USER_NOT_FOUND.String() && e.Code == 400
}

// 用户找不到
func ErrorUserNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_USER_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

// 入参不对
func IsInvaildParam(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_INVAILD_PARAM.String() && e.Code == 400
}

// 入参不对
func ErrorInvaildParam(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_INVAILD_PARAM.String(), fmt.Sprintf(format, args...))
}

// 密码错误
func IsWrongPasswd(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_WRONG_PASSWD.String() && e.Code == 400
}

// 密码错误
func ErrorWrongPasswd(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_WRONG_PASSWD.String(), fmt.Sprintf(format, args...))
}
