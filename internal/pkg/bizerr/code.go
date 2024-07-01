package bizerr

import (
	"fmt"
	"sync"
)

type CodeMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// implement the error interface
func (e CodeMsg) Error() string {
	return e.Msg
}

// Wrap the error code and message
func (e CodeMsg) Wrap(err error) CodeMsg {
	e.Msg = err.Error()
	return e
}

// WithMsg returns an error with the supplied message.
func (e CodeMsg) WithMsg(msg string) CodeMsg {
	e.Msg = msg
	return e
}

// WithMsgf returns an error with the supplied message.
func (e CodeMsg) WithMsgf(format string, args ...interface{}) CodeMsg {
	e.Msg = fmt.Sprintf(format, args...)
	return e
}

// codeMsgMap stores the error code and message.
var codeMsgMap = make(map[int]string)
var mutex sync.Mutex

// New returns an error with the supplied message.
func New(code int, msg string) CodeMsg {
	// Store the error code and message, and determine if it is a duplicate definition
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := codeMsgMap[code]; ok {
		panic("error code already exists: " + msg)
	}
	codeMsgMap[code] = msg
	return CodeMsg{Code: code, Msg: msg}
}
