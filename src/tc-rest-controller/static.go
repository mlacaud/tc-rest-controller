package main

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Var and Struct
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var serverPort string

var uploadIface string

var downloadIface string

type netem struct {
	Delay string
	Loss  string
}