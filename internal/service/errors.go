package service

import "errors"

var (
	ErrCreateTender           = errors.New("cannot create tender (or newer version)")
	ErrUsername               = errors.New("user doesn't exists or username (id) is incorrect")
	ErrCheckResponsibility    = errors.New("cannot check responsibility")
	ErrForbidden              = errors.New("not ehough rights")
	ErrGetEmployeeByUsername  = errors.New("cannot get employee by username")
	ErrGetEmployeeById        = errors.New("cannot get employee by id")
	ErrGetTendersByUsername   = errors.New("cannot get tenders by username")
	ErrGetTenders             = errors.New("cannot get tenders")
	ErrNotFoundTender         = errors.New("tender not found (or exact tender version)")
	ErrGetTender              = errors.New("cannot get tender (or exact tender version)")
	ErrGetTenderLatestVersion = errors.New("cannot get latest version of tender")
	ErrCreateBid              = errors.New("cannot create tender (or newer version)")
	ErrNotFoundBid            = errors.New("bid not found (or exact bid version)")
	ErrGetBid                 = errors.New("cannot get bid (or exact bid version)")
)
