package scripting

import "fmt"
import "strings"
import "io/ioutil"
import "path/filepath"
import "github.com/tobz/phosphorus/log"
import "github.com/tobz/phosphorus/interfaces"
import "github.com/idada/v8.go"

type ScriptExecutor struct {
	scriptEngines map[string]*ScriptEngine
}

func NewScriptExecutor(config interfaces.Config) (*ScriptExecutor, error) {
	executor := &ScriptExecutor{scriptEngines: make(map[string]*ScriptEngine)}

	absPath, err := filepath.Abs("scripts")
	if err != nil {
		return nil, fmt.Errorf("Failed getting absolute path to scripts directory: %s", err)
	}

	entries, err := ioutil.ReadDir(absPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read all entries in the scripts directory: %s", err)
	}

	engineBuilder := func(engineName, scriptPath string) error {
		entries, err := ioutil.ReadDir(scriptPath)
		if err != nil {
			return fmt.Errorf("Failed to read all entries in the %s directory: %s", engineName, err)
		}

		engine := NewScriptEngine(engineName)

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".js") {
				fullScriptPath := filepath.Join(scriptPath, entry.Name())

				err := engine.addScript(fullScriptPath)
				if err != nil {
					return fmt.Errorf("Failed to load %s into '%s' engine: %s", fullScriptPath, engineName, err)
				}
			}
		}

		err = executor.addEngine(engine)
		if err != nil {
			return fmt.Errorf("Failed to add engine to executor: %s", err)
		}

		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			log.Server.Info("scripting", "Creating new engine with scripts from '%s'...", entry.Name())

			err := engineBuilder(entry.Name(), filepath.Join(absPath, entry.Name()))
			if err != nil {
				return nil, err
			}
		}
	}

	return executor, nil
}

func (se *ScriptExecutor) addEngine(engine *ScriptEngine) error {
	if _, ok := se.scriptEngines[engine.Name]; ok {
		return fmt.Errorf("Script engine '%s' already exists!", engine.Name)
	}

	se.scriptEngines[engine.Name] = engine

	return nil
}

func (se *ScriptExecutor) bindToEngines(targetName string, targetValue interface{}) error {
	for _, engine := range se.scriptEngines {
		err := engine.addBinding(targetName, targetValue)
		if err != nil {
			return err
		}
	}

	return nil
}

func (se *ScriptExecutor) Execute(engineName, functionName string, arguments ...interface{}) (interfaces.ExecutorResult, error) {
	engine, ok := se.scriptEngines[engineName]
	if !ok {
		return nil, fmt.Errorf("Script engine '%s' does not exist or isn't registered with this executor!", engineName)
	}

	var err error
	var result interfaces.ExecutorResult

	engine.runInContext(func(scope v8.ContextScope) {
		wrappedCallee := scope.Eval(fmt.Sprintf("(function() { return %s })()", functionName))
		if !wrappedCallee.IsFunction() {
			err = fmt.Errorf("Failed to wrap function '%s' during Execute call! Result: %#v", functionName, wrappedCallee)
			return
		}

		wrappedArguments := engine.GetWrappedArguments(arguments...)
		res := wrappedCallee.ToFunction().Call(wrappedArguments...)

		result = NewExecutorResult(res)
	})

	return result, err
}
