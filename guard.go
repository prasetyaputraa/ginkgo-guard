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

	GinContext *gin.Context
}

type GuardOperand func(interface{})

func (c *GuardContext) Copy() GuardContext {
	return GuardContext{
		isAborted:    c.isAborted,
		abortCode:    c.abortCode,
		abortPayload: c.abortPayload,
		statusCode:   c.statusCode,
		GinContext:   c.GinContext,
	}
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

func (c GuardContext) GetAbortCode() uint {
	return c.abortCode
}

func (c GuardContext) GetAbortPayload() interface{} {
	return c.abortPayload
}

func (c *GuardContext) CopyAbortionStateOf(carbon GuardContext) {
	if !carbon.IsAborted() {
		panic("carbon code is not aborted")
	}

	if c.GinContext != carbon.GinContext {
		panic("carbon is not derived from c")
	}

	c.abortCode = carbon.GetAbortCode()
	c.statusCode = carbon.GetStatusCode()
	c.abortPayload = carbon.GetAbortPayload()
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

func PGuardAnd(midw1 GuardOperand, midw2 GuardOperand) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func PGuardOr(midw1 GuardOperand, midw2 GuardOperand) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GuardAnd(midw1 GuardOperand, midw2 GuardOperand) GuardOperand {
	return func(c interface{}) {
		guardContext, isGuardContext := c.(*GuardContext)
		ginContext, isGinContext := c.(*gin.Context)

		if !isGinContext && !isGuardContext {
			panic("expecting gin.Context or ginkgo_guard.GuardContext")
		}

		if isGuardContext {
			if guardContext.IsAborted() {
				return
			} else {
				carbon1 := guardContext.Copy()
				carbon2 := guardContext.Copy()

				midw1(&carbon1)
				midw2(&carbon2)

				if carbon1.IsAborted() {
					c.(*GuardContext).CopyAbortionStateOf(carbon1)
					return
				}

				if carbon2.IsAborted() {
					c.(*GuardContext).CopyAbortionStateOf(carbon2)
					return
				}
			}
		} else if isGinContext {
			if ginContext.IsAborted() {
				return
			} else {
				carbon1 := GuardContext{
					GinContext: ginContext,
				}
				carbon2 := GuardContext{
					GinContext: ginContext,
				}

				midw1(&carbon1)
				midw2(&carbon2)

				if carbon1.IsAborted() {
					c.(*GuardContext).CopyAbortionStateOf(carbon1)
					return
				}

				if carbon2.IsAborted() {
					c.(*GuardContext).CopyAbortionStateOf(carbon2)
					return
				}
			}
		}
	}
}

func GuardOr(midw1 GuardOperand, midw2 GuardOperand) GuardOperand {
	return func(c interface{}) {
		guardContext, isGuardContext := c.(*GuardContext)
		ginContext, isGinContext := c.(*gin.Context)

		if !isGinContext && !isGuardContext {
			panic("expecting gin.Context or ginkgo_guard.GuardContext")
		}

		if isGuardContext {
			if guardContext.IsAborted() {
				return
			} else {
				carbon1 := guardContext.Copy()
				carbon2 := guardContext.Copy()

				midw1(&carbon1)
				midw2(&carbon2)

				if carbon1.IsAborted() && carbon2.IsAborted() {
					c.(*GuardContext).CopyAbortionStateOf(carbon1)
					return
				}
			}
		} else if isGinContext {
			if ginContext.IsAborted() {
				return
			} else {
				carbon1 := GuardContext{
					GinContext: ginContext,
				}
				carbon2 := GuardContext{
					GinContext: ginContext,
				}

				midw1(&carbon1)
				midw2(&carbon2)

				if carbon1.IsAborted() && carbon2.IsAborted() {
					c.(*GuardContext).CopyAbortionStateOf(carbon1)
					return
				}
			}
		}
	}
}
