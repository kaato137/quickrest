package internal

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/robertkrimen/otto"
)

type RenderContext map[string]string

type Renderer struct {
	vm *otto.Otto
}

func NewRenderer() *Renderer {
	r := &Renderer{
		vm: otto.New(),
	}

	r.registerHelperFunctions()

	return r
}

func (r *Renderer) Render(template string, renderCtx RenderContext) ([]byte, error) {
	for k, v := range renderCtx {
		if err := r.vm.Set(k, v); err != nil {
			return nil, fmt.Errorf("set %q: %w", k, err)
		}
	}

	val, err := r.vm.Run(template)
	if err != nil {
		return nil, fmt.Errorf("js: %w", err)
	}

	return val.Object().MarshalJSON()
}

func (r *Renderer) registerHelperFunctions() {
	r.vm.Set("int", func(call otto.FunctionCall) otto.Value {
		intVal, err := call.Argument(0).ToInteger()
		if err != nil {
			return r.vm.MakeTypeError("failed to cast value to integer")
		}

		val, err := r.vm.ToValue(intVal)
		if err != nil {
			return r.vm.MakeTypeError("failed to make value of out int")
		}

		return val
	})

	r.vm.Set("uuid", func(call otto.FunctionCall) otto.Value {
		args := call.ArgumentList

		if len(args) == 0 {
			val, err := r.vm.ToValue(uuid.NewString())
			if err != nil {
				return r.vm.MakeTypeError(err.Error())
			}
			return val
		}

		arg, err := call.Argument(0).ToString()
		if err != nil {
			return r.vm.MakeTypeError(err.Error())
		}

		val, err := r.vm.ToValue(uuid.NewMD5(uuid.NameSpaceOID, []byte(arg)).String())
		if err != nil {
			return r.vm.MakeTypeError("failed to make value of out int")
		}

		return val
	})
}
