package scripting

import "github.com/idada/v8.go"

type ExecutorResult struct {
	str     string
	integer int64
	float   float64
	boolean bool
}

func (er *ExecutorResult) String() string {
	return er.str
}

func (er *ExecutorResult) Integer() int64 {
	return er.integer
}

func (er *ExecutorResult) Float() float64 {
	return er.float
}

func (er *ExecutorResult) Boolean() bool {
	return er.boolean
}

func NewExecutorResult(v *v8.Value) *ExecutorResult {
	res := &ExecutorResult{}

	switch {
	case v.IsString():
		res.str = v.ToString()
	case v.IsInt32(), v.IsUint32():
		res.integer = v.ToInteger()
	case v.IsBoolean():
		res.boolean = v.ToBoolean()
	case v.IsNumber():
		res.float = v.ToNumber()
	}

	return res
}
