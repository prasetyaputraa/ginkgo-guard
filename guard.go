package ginkgo_guard

import (
	"context"

	"github.com/gin-gonic/gin"
)

var abortWithStatusCode uint = 1
var abortWithStatusJSONCode uint = 2
var abortWithStatusErrorCode uint = 3

type GuardContext struct {
	context.Context

	isAborted    bool
	abortCode    uint
	abortPayload interface{}
	statusCode   uint
}

func (c GuardContext) IsAborted() bool {
	return c.isAborted
}

func (c GuardContext) panicIfAlreadyAborted() {
	if c.isAborted {
		panic("Trying to abort an aborted guard context.")
	}
}

func (c GuardContext) GetStatusCode() uint {
	return c.statusCode
}

func (c GuardContext) GetAbortPayload() interface{} {
	return c.abortPayload
}

func (c GuardContext) AbortWithStatusJSON(status uint, jsonObject interface{}) {
	c.panicIfAlreadyAborted()

	c.isAborted = true
	c.abortCode = abortWithStatusJSONCode
	c.abortPayload = jsonObject
}

func (c GuardContext) AbortWithStatusError(status uint, err error) {
	c.panicIfAlreadyAborted()

	c.isAborted = true
	c.abortCode = abortWithStatusErrorCode
	c.abortPayload = err
}

func (c GuardContext) AbortWithStatus(status uint) {
	c.panicIfAlreadyAborted()

	c.isAborted = true
	c.abortCode = abortWithStatusCode
}

func (g GuardContext) AbortGinContext(c *gin.Context) {
	switch g.abortCode {
	case abortWithStatusCode:
		c.AbortWithStatus(int(g.statusCode))
	case abortWithStatusErrorCode:
		c.AbortWithError(int(g.statusCode), g.abortPayload.(error))
	case abortWithStatusJSONCode:
		c.AbortWithStatusJSON(int(g.statusCode), g.abortPayload)
	default:
		c.Abort()
	}
}
