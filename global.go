package keygen

import (
	"net/url"
	"time"
)

func New() *Generator {
	return &Generator{Reporter: func(err error) {
		panic(err)
	}}
}

var global = New()

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
