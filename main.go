//go:build js && wasm

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"syscall/js"
	"text/template"
)

func renderGoTemplate(_ js.Value, args []js.Value) any {
	log.Printf("renderGoTemplate: args=%v\n", args)
	if len(args) < 1 {
		return fmt.Errorf("error: missing template string")
	}
	if len(args) < 2 {
		return fmt.Errorf("error: missing template input values")
	}

	if args[0].Type() != js.TypeString {
		return fmt.Errorf("error: first argument must be a string, but got %[1]T (%[1]v)", args[0])
	}
	if args[1].Type() != js.TypeString {
		return fmt.Errorf("error: second argument must be JSON string, but got %[1]T (%[1]v)", args[1])
	}

	tpl := args[0].String()
	vals := map[string]any{}

	// Because toMap (below) doesn't work, instead JS must pass a JSON string
	// containing the template inputs, and we can parse that string into an
	// arbitrary map of values.
	err := json.Unmarshal([]byte(args[1].String()), &vals)
	if err != nil {
		return fmt.Errorf("error: second argument must be a JSON string: %w", err)
	}

	t, err := template.New("").Parse(tpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, vals); err != nil {
		return err
	}
	return buf.String()
}

// This doesn't work. I can't figure out a way to iterate over js object keys.
// You can Get and Set keys, and iterate arrays (which are considered objects)
// but seemingly can't iterate object keys.
func toMap(v js.Value) (map[string]any, error) {
	log.Printf("toMap: length=%d, %#v", v.Length(), v)
	result := map[string]any{}
	for i := 0; i < v.Length(); i++ {
		k := v.Index(i)
		log.Printf("toMap: k=%v", k)
		if k.Type() != js.TypeString {
			return nil, fmt.Errorf("non-string key when converting to Go map: %v", v)
		}
		key := v.String()
		val, err := toGo(v.Get(key))
		if err != nil {
			return nil, fmt.Errorf("unable to convert valid to Go type: %w", err)
		}
		result[key] = val
	}
	return result, nil
}

func toGo(v js.Value) (any, error) {
	switch v.Type() {
	case js.TypeUndefined, js.TypeNull:
		return nil, nil
	case js.TypeBoolean:
		return v.Bool(), nil
	case js.TypeNumber:
		return v.Float(), nil
	case js.TypeString:
		return v.String(), nil
	case js.TypeSymbol:
		// ???
		return nil, fmt.Errorf("cannot convert JS symbol to Go type: symbol=%v", v)
	case js.TypeObject:
		return toMap(v)
	case js.TypeFunction:
		return nil, fmt.Errorf("refusing to convert JS function to Go type: %v", v)
	}
	return nil, fmt.Errorf("unable to convert JS type to Go type: %v", v)
}

func main() {
	js.Global().Set("renderGoTemplate", js.FuncOf(renderGoTemplate))
	// Stay alive.
	<-make(chan struct{})
}
