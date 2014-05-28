package interfaces

type ExecutorResult interface {
	String() string
	Integer() int64
	Float() float64
	Boolean() bool
}

type ScriptExecutor interface {
	Execute(engineName, functionName string, arguments ...interface{}) (ExecutorResult, error)
}
