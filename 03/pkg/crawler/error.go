package crawler

import "errors"

var ErrRequest = errors.New("request error")
var ErrInternal = errors.New("internal error")
var ErrResponse = errors.New("response error")
var ErrInvalidURL = errors.New("invalid URL")
var ErrParse = errors.New("parse error")
