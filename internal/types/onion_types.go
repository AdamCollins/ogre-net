package types

type OnionMessage struct {
	NextHop   NodeAddress   // address of nextNode
	NextLayer *OnionMessage // NextLayer down in package
	Payload   string        // data payload. only to be used on bottom layer
}
