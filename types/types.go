package types

import lua "github.com/yuin/gopher-lua"

type String struct {
	Key, Value string
}

func NewString(key, value string) *String {
	return &String{
		Key:   key,
		Value: value,
	}
}

func (s *String) GetKey() string {
	return s.Key
}

func (s *String) GetValue() lua.LValue {
	return lua.LString(s.Value)
}

type Number struct {
	Key   string
	Value float64
}

func NewNumber(key string, value float64) *Number {
	return &Number{
		Key:   key,
		Value: value,
	}
}

func (n *Number) GetKey() string {
	return n.Key
}

func (n *Number) GetValue() lua.LValue {
	return lua.LNumber(n.Value)
}

type Bool struct {
	Key   string
	Value bool
}

func NewBool(key string, value bool) *Bool {
	return &Bool{
		Key:   key,
		Value: value,
	}
}

func (b *Bool) GetKey() string {
	return b.Key
}

func (b *Bool) GetValue() lua.LValue {
	return lua.LBool(b.Value)
}
