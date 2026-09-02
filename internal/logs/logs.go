package logs

import (
	"factorio/internal/events"
	"fmt"
	"os"
	"sync"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

type luaWorker struct {
	state   *lua.LState
	version uint64
}

type LogParser struct {
	scriptPath string
	eventBus   *events.EventBus
	workerPool sync.Pool

	mutex   sync.RWMutex
	proto   *lua.FunctionProto
	version uint64
}

func compileLua(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	chunk, err := parse.Parse(file, filePath)
	if err != nil {
		return nil, err
	}

	return lua.Compile(chunk, filePath)
}

func (l *LogParser) Reload() error {
	newProto, err := compileLua(l.scriptPath)
	if err != nil {
		return fmt.Errorf("failed to compile lua script: %w", err)
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.proto = newProto
	l.version++

	return nil
}

func NewLogParser(scriptPath string, eventBus *events.EventBus) (*LogParser, error) {
	parser := &LogParser{
		scriptPath: scriptPath,
		eventBus:   eventBus,
		workerPool: sync.Pool{
			New: func() any {
				return &luaWorker{}
			},
		},
	}

	if err := parser.Reload(); err != nil {
		return nil, fmt.Errorf("failed initial lua compilation: %w", err)
	}

	return parser, nil
}

func (l *LogParser) ParseLine(line string) error {

	l.mutex.RLock()
	currentProto := l.proto
	currentVersion := l.version
	l.mutex.RUnlock()

	worker := l.workerPool.Get().(*luaWorker)
	defer l.workerPool.Put(worker)

	if worker.version != currentVersion || worker.state == nil {
		if worker.state != nil {
			worker.state.Close()
		}
		worker.state = lua.NewState()
		worker.state.Push(worker.state.NewFunctionFromProto(currentProto))
		if err := worker.state.PCall(0, 1, nil); err != nil {
			return fmt.Errorf("failed to initialize parser function: %w", err)
		}
		worker.state.SetGlobal("__parse_fn", worker.state.Get(-1))
		worker.state.Pop(1)
		worker.version = currentVersion
	}

	lstate := worker.state
	fn := lstate.GetGlobal("__parse_fn")
	lstate.Push(fn)
	lstate.Push(lua.LString(line))

	if err := lstate.PCall(1, 1, nil); err != nil {
		return err
	}

	ret := lstate.Get(-1)
	lstate.Pop(1)

	if ret == lua.LNil {
		return nil
	}

	tbl, ok := ret.(*lua.LTable)
	if !ok {
		return fmt.Errorf("expected table return, got %s", ret.Type().String())
	}

	typeVal := lstate.GetField(tbl, "type")
	if typeVal.Type() != lua.LTString {
		return fmt.Errorf("event table missing strict string field `type`")
	}

	dataVal := lstate.GetField(tbl, "data")

	l.eventBus.Publish(events.EventMessage{
		Type: typeVal.String(),
		Data: ConvertLuaValueToGo(dataVal),
	})

	return nil
}

func ConvertLuaValueToGo(val lua.LValue) any {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case *lua.LTable:
		maxN := v.MaxN()
		if maxN > 0 {
			arr := make([]any, 0, maxN)
			for i := 1; i <= maxN; i++ {
				arr = append(arr, ConvertLuaValueToGo(v.RawGetInt(i)))
			}
			return arr
		}

		message := make(map[string]any)
		v.ForEach(func(k, item lua.LValue) {
			message[k.String()] = ConvertLuaValueToGo(item)
		})
		return message
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return v.String()
	case lua.LBool:
		return bool(v)
	case *lua.LNilType:
		return nil
	default:
		return nil
	}
}
