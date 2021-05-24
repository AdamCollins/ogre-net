package types

type Constants struct {
	PingInterval int16 // amount of time between ping requests in ms
}

var Constant = Constants{
	PingInterval: 8000,
}
