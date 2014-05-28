package scripting

import "io/ioutil"
import "reflect"
import "path/filepath"
import "github.com/idada/v8.go"

type ScriptEngine struct {
	Name          string
	engine        *v8.Engine
	globalObject  *v8.ObjectTemplate
	globalContext *v8.Context
}

func NewScriptEngine(name string) *ScriptEngine {
	engine := v8.NewEngine()
	globalObject := engine.NewObjectTemplate()
	globalContext := engine.NewContext(globalObject)

	return &ScriptEngine{
		Name:          name,
		engine:        engine,
		globalObject:  globalObject,
		globalContext: globalContext,
	}
}

func (se *ScriptEngine) addScript(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	contents, err := ioutil.ReadFile(abs)
	if err != nil {
		return err
	}

	se.globalContext.Scope(func(scope v8.ContextScope) {
		scope.Eval(string(contents))
	})

	return nil
}

func (se *ScriptEngine) addBinding(targetName string, targetValue interface{}) error {
	return se.globalObject.Bind(targetName, targetValue)
}

func (se *ScriptEngine) GetWrappedArguments(args ...interface{}) []*v8.Value {
	wrappedArgs := []*v8.Value{}

	for _, arg := range args {
		wrappedArgs = append(wrappedArgs, se.engine.GoValueToJsValue(reflect.ValueOf(arg)))
	}

	return wrappedArgs
}

func (se *ScriptEngine) runInContext(f func(v8.ContextScope)) error {
	var err error

	se.globalContext.Scope(func(scope v8.ContextScope) {
		err = scope.TryCatch(func() {
			f(scope)
		})
	})

	return err
}
