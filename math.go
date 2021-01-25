package util

type min struct{}

var Min min

func (*min) Int(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (*min) Uint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func (*min) Int32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func (*min) Ints(i ...int) int {
	if len(i) == 0 {
		panic("at least one element")
	}
	ret := i[0]
	for _, v := range i {
		if v < ret {
			ret = v
		}
	}
	return ret
}

func (*min) Uint32s(i ...uint32) uint32 {
	if len(i) == 0 {
		panic("at least one element")
	}
	ret := i[0]
	for _, v := range i {
		if v < ret {
			ret = v
		}
	}
	return ret
}

func (*min) Int32s(i ...int32) int32 {
	if len(i) == 0 {
		panic("at least one element")
	}
	ret := i[0]
	for _, v := range i {
		if v < ret {
			ret = v
		}
	}
	return ret
}

type max struct{}

var Max max

func (*max) Int(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (*max) Uint32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func (*max) Int32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func (*max) Ints(i ...int) int {
	if len(i) == 0 {
		panic("at least one element")
	}
	ret := i[0]
	for _, v := range i {
		if v > ret {
			ret = v
		}
	}
	return ret
}

func (*max) Uint32s(i ...uint32) uint32 {
	if len(i) == 0 {
		panic("at least one element")
	}
	ret := i[0]
	for _, v := range i {
		if v > ret {
			ret = v
		}
	}
	return ret
}

func (*max) Int32s(i ...int32) int32 {
	if len(i) == 0 {
		panic("at least one element")
	}
	ret := i[0]
	for _, v := range i {
		if v > ret {
			ret = v
		}
	}
	return ret
}
