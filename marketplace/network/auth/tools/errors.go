package tools

import "errors"

var ErrLogin = errors.New("Incorrect password or login")
var ErrLoginFormat = errors.New("Incorrect login format")
var ErrPasswordFormat = errors.New("Incorrect password format")
var ErrLoginIsTaken = errors.New("Login is taken")
var ErrServer = errors.New("Server error")
