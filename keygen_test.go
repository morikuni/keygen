package keygen

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func NotEqual(tb testing.TB, want, got interface{}) {
	tb.Helper()

	if reflect.DeepEqual(want, got) {
		tb.Errorf("equal\n\twant: %T(%v)\n\tgot: %T(%v)", want, want, got, got)
	}
}

func Equal(tb testing.TB, want, got interface{}) {
	tb.Helper()

	if !reflect.DeepEqual(want, got) {
		tb.Errorf("not equal\n\twant: %T(%v)\n\tgot: %T(%v)", want, want, got, got)
	}
}

func TestEquality(t *testing.T) {
	cases := map[string]struct {
		gen func(keys ...string) interface{}
	}{
		"Int":     {gen: func(keys ...string) interface{} { return Int(keys...) }},
		"Int64":   {gen: func(keys ...string) interface{} { return Int64(keys...) }},
		"Uint":    {gen: func(keys ...string) interface{} { return Uint(keys...) }},
		"Uint64":  {gen: func(keys ...string) interface{} { return Uint64(keys...) }},
		"Float64": {gen: func(keys ...string) interface{} { return Float64(keys...) }},
		"Float32": {gen: func(keys ...string) interface{} { return Float32(keys...) }},
		"String":  {gen: func(keys ...string) interface{} { return String(keys...) }},
		"Bool":    {gen: func(keys ...string) interface{} { return Bool(keys...) }},
		"Time":    {gen: func(keys ...string) interface{} { return Time(keys...) }},
		"URL":     {gen: func(keys ...string) interface{} { return URL(keys...) }},
	}

	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			i1 := tc.gen("data", "1")
			i2 := tc.gen("data", "2")
			i3 := tc.gen("data1")

			NotEqual(t, i1, i2)
			Equal(t, i1, i3)
		})
	}
}

func TestRange(t *testing.T) {
	cases := map[string]struct {
		gen      func(keys ...string) interface{}
		validate func(v interface{}) bool
	}{
		"Time": {
			gen: func(keys ...string) interface{} { return Time(keys...) },
			validate: func(v interface{}) bool {
				t := v.(*time.Time)
				return t.After(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)) &&
					t.Before(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))
			},
		},
		"Length": {
			gen: func(keys ...string) interface{} { return Length(keys...) },
			validate: func(v interface{}) bool {
				i := *v.(*int)
				return 1 <= i && i <= 5
			},
		},
	}

	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			for i := 0; i < 10000; i++ {
				v := tc.gen("data", strconv.Itoa(i))
				if !tc.validate(v) {
					t.Fatalf("validation failed: %v", v)
				}
			}
		})
	}
}

type Embed struct {
	BoolP  *bool
	String string
}
type Object struct {
	Int     int
	Uint    uint32
	Float   float64
	String  string
	Array   [5]int64
	Slice   []bool
	Struct  Embed
	StructP *Embed
}

func TestAny(t *testing.T) {
	obj := *Any((*Object)(nil), "fill").(*Object)

	Equal(t, *Int("fill", "int"), obj.Int)
	Equal(t, uint32(*Uint("fill", "uint")), obj.Uint)
	Equal(t, *Float64("fill", "float"), obj.Float)
	Equal(t, *String("fill", "string"), obj.String)
	Equal(t,
		[5]int64{
			*Int64("fill", "array", "0"),
			*Int64("fill", "array", "1"),
			*Int64("fill", "array", "2"),
			*Int64("fill", "array", "3"),
			*Int64("fill", "array", "4"),
		},
		obj.Array,
	)
	Equal(t, *Length("fill", "slice", "len"), len(obj.Slice)) // this value == 4
	Equal(t,
		[]bool{
			*Bool("fill", "slice", "0"),
			*Bool("fill", "slice", "1"),
			*Bool("fill", "slice", "2"),
			*Bool("fill", "slice", "3"),
		},
		obj.Slice,
	)
	Equal(t, *Bool("fill", "struct", "bool_p"), *obj.Struct.BoolP)
	Equal(t, *String("fill", "struct", "string"), obj.Struct.String)
	Equal(t, *Bool("fill", "struct_p", "bool_p"), *obj.StructP.BoolP)
	Equal(t, *String("fill", "struct_p", "string"), obj.StructP.String)
}

func BenchmarkAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Any((*Object)(nil), "fill").(*Object)
	}
}
