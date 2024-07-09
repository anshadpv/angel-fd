package errors

import (
	"net/http"

	"github.com/angel-one/goerr"
)

// Error Messages
const (
	AuthInvalidSigningMethod = "unexpected signing method"
)

// controller errros
var (
	HeaderAuthMissingInvalid = goerr.New(nil, http.StatusBadRequest, "header authorization is invalid/missing")
	ErrInvalidRequest        = goerr.New(nil, http.StatusBadRequest, "one or more required parameter missing")
	NotAuthorized            = goerr.New(nil, http.StatusUnauthorized, "user not authorized")
	ErrClientData            = goerr.New(nil, http.StatusBadRequest, "error validating client")
	ErrClientValidation      = goerr.New(nil, http.StatusBadRequest, "Client does not belong to the partner")
)
