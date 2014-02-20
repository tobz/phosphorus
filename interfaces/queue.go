package interfaces

type Queue interface {
	Len() int
	Push(interface{})
	Pop() interface{}
	Peek() interface{}
}
