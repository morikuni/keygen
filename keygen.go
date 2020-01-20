package keygen

import (
	"encoding/base64"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/minio/highwayhash"
)

var hashKey = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

func joinKeys(keys []string) []byte {
	var size int
	for _, k := range keys {
		size += len(k)
	}

	key := make([]byte, 0, size)
	for _, k := range keys {
		key = append(key, k...)
	}

	return key
}

func hashUint64(keys []string) uint64 {
	return highwayhash.Sum64(joinKeys(keys), hashKey)
}

func hash(keys []string) [32]byte {
	res := highwayhash.Sum(joinKeys(keys), hashKey)
	return res
}

type Generator struct {
	Reporter func(err error)
}

func New() *Generator {
	return &Generator{Reporter: func(err error) {
		panic(err)
	}}
}

var global = New()

func (g *Generator) Int(keys ...string) *int {
	v := int(*g.Int64(keys...))
	return &v
}

func (g *Generator) Int64(keys ...string) *int64 {
	v := int64(hashUint64(keys))
	return &v
}

func (g *Generator) Uint(keys ...string) *uint {
	v := uint(*g.Uint64(keys...))
	return &v
}

func (g *Generator) Uint64(keys ...string) *uint64 {
	v := hashUint64(keys)
	return &v
}

func (g *Generator) Length(keys ...string) *int {
	v := int(*g.Uint(keys...)%5 + 1)
	return &v
}

func (g *Generator) String(keys ...string) *string {
	h := hash(keys)
	// Take the first 240 bit to encode base64 without "=" suffix.
	// The hash function returns 256 bit (32 byte), but base64 is based on 6 bit.
	v := base64.URLEncoding.EncodeToString(h[:30])
	return &v
}

func (g *Generator) Bool(keys ...string) *bool {
	v := hashUint64(keys)%2 == 1
	return &v
}

func (g *Generator) Float64(keys ...string) *float64 {
	v := math.Float64frombits(hashUint64(keys))
	return &v
}

func (g *Generator) Float32(keys ...string) *float32 {
	v := float32(*g.Float64(keys...))
	return &v
}

// Time returns random time between 2009-11-10 23:00:00 and 2030-01-01 00:00:00.
func (g *Generator) Time(keys ...string) *time.Time {
	t := time.Unix(
		int64(*Uint64(keys...)%uint64(176545*time.Hour/time.Second)+1257894000),
		int64(*Uint64(keys...)%uint64(time.Second)),
	)
	return &t
}

func (g *Generator) URL(keys ...string) *url.URL {
	u, err := url.Parse(fmt.Sprintf("%s://%s/%s",
		[]string{"https", "http"}[*Uint(keys...)%2],
		*String(append(keys, "host")...)+[]string{".com", ".org", ".net"}[*Uint(keys...)%3],
		*String(append(keys, "path")...),
	))
	if err != nil {
		g.Reporter(err)
		return nil
	}
	return u
}

func (g *Generator) Any(dst interface{}, keys ...string) interface{} {
	rv := reflect.New(reflect.TypeOf(dst).Elem())

	if rv.Kind() != reflect.Ptr {
		g.Reporter(fmt.Errorf("dst should be pointer but %T", dst))
		return nil
	}

	g.gen(rv.Elem(), keys...)

	return rv.Interface()
}

func (g *Generator) gen(rv reflect.Value, keys ...string) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv.SetInt(*g.Int64(keys...))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv.SetUint(*g.Uint64(keys...))
	case reflect.Float32, reflect.Float64:
		rv.SetFloat(*g.Float64(keys...))
	case reflect.String:
		rv.SetString(*g.String(keys...))
	case reflect.Bool:
		rv.SetBool(*g.Bool(keys...))
	case reflect.Array:
		for i, l := 0, rv.Len(); i < l; i++ {
			g.gen(rv.Index(i), append(keys, strconv.Itoa(i))...)
		}
	case reflect.Slice:
		l := *g.Length(append(keys, "len")...)
		rv.Set(reflect.MakeSlice(rv.Type(), l, l))
		for i, l := 0, rv.Len(); i < l; i++ {
			g.gen(rv.Index(i), append(keys, strconv.Itoa(i))...)
		}
	case reflect.Map:
		l := *g.Length(append(keys, "len")...)
		t := rv.Type()
		m := reflect.MakeMapWithSize(t, 100)
		rv.Set(m)
		keyT := t.Key()
		valT := t.Elem()
		for i, l := 0, l; i < l; i++ {
			key := reflect.New(keyT).Elem()
			val := reflect.New(valT).Elem()
			g.gen(key, append(keys, strconv.Itoa(i), "key")...)
			g.gen(val, append(keys, strconv.Itoa(i), "value")...)
			rv.SetMapIndex(key, val)
		}
	case reflect.Struct:
		rt := rv.Type()
		for i, l := 0, rv.NumField(); i < l; i++ {
			f := rt.Field(i)
			g.gen(rv.Field(i), append(keys, strcase.ToSnake(f.Name))...)
		}
	case reflect.Ptr:
		rv.Set(reflect.New(rv.Type().Elem()))
		g.gen(rv.Elem(), keys...)
	default:
		g.Reporter(fmt.Errorf("not supported kind: %v", rv.Kind().String()))
	}
}

func Int(keys ...string) *int {
	return global.Int(keys...)
}

func Int64(keys ...string) *int64 {
	return global.Int64(keys...)
}

func Uint(keys ...string) *uint {
	return global.Uint(keys...)
}

func Uint64(keys ...string) *uint64 {
	return global.Uint64(keys...)
}

func Length(keys ...string) *int {
	return global.Length(keys...)
}

func String(keys ...string) *string {
	return global.String(keys...)
}

func Bool(keys ...string) *bool {
	return global.Bool(keys...)
}

func Float64(keys ...string) *float64 {
	return global.Float64(keys...)
}

func Float32(keys ...string) *float32 {
	return global.Float32(keys...)
}

// Time returns random time between 2009-11-10 23:00:00 and 2030-01-01 00:00:00.
func Time(keys ...string) *time.Time {
	return global.Time(keys...)
}

func URL(keys ...string) *url.URL {
	return global.URL(keys...)
}

func Any(dst interface{}, keys ...string) interface{} {
	return global.Any(dst, keys...)
}
