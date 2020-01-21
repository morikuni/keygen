package keygen

import (
	"errors"
	"fmt"
	"reflect"
)

type CustomGenerator func(g *Generator, args []string, keys []string) (interface{}, error)

func (g *Generator) RegisterGenerator(name string, gen CustomGenerator) {
	if _, ok := g.CustomGenerator[name]; ok {
		g.Reporter(fmt.Errorf("name %q is already registered", name))
		return
	}

	g.CustomGenerator[name] = gen
}

func (g *Generator) genCustom(rv reflect.Value, rt reflect.Type, args []string, keys []string) {
	if len(args) == 0 {
		g.Reporter(errors.New("name of generator is empty"))
		return
	}

	name, args := args[0], args[1:]

	gen, ok := g.CustomGenerator[name]
	if !ok {
		g.Reporter(fmt.Errorf("generator %q is not found", name))
		return

	}

	v, err := gen(g, args, keys)
	if err != nil {
		g.Reporter(fmt.Errorf("generator %q: %v", name, err))
		return
	}

	if rt.Kind() == reflect.String {
		rv.SetString(fmt.Sprint(v))
		return
	}

	vv := reflect.ValueOf(v)
	vt := vv.Type()

	if !vt.ConvertibleTo(rt) {
		g.Reporter(fmt.Errorf("generator %q: type mismatch %s and %s", name, vt.String(), rt.String()))
		return
	}

	rv.Set(vv.Convert(rt))
	return
}
