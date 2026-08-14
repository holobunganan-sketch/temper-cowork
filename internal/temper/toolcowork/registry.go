package toolcowork

import (
	"encoding/json"

	"reasonix/internal/tool"
)

// registered 是 process-global 的 temper 工具集,由 init() 填充。
var registered = map[string]tool.Tool{}

// Register 注册一个 temper 工具(init 中使用)。
func Register(t tool.Tool) {
	name := t.Name()
	if _, dup := registered[name]; dup {
		panic("toolcowork: duplicate tool " + name)
	}
	registered[name] = t
}

// All 返回全部注册的 temper 工具(按名排序)。
func All() []tool.Tool {
	out := make([]tool.Tool, 0, len(registered))
	for _, t := range registered {
		out = append(out, t)
	}
	return out
}

// RegisterInto 把 temper 工具注册进 Reasonix 运行时 Registry。
func RegisterInto(r *tool.Registry) {
	for _, t := range All() {
		r.Add(t)
	}
}

// mustSchema 把 map 编码为 JSON Schema;编码失败时 panic(编译期错误)。
func mustSchema(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic("toolcowork: bad schema: " + err.Error())
	}
	return b
}
