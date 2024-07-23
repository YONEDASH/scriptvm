package vm

type TypeId uint8

const (
	Invalid TypeId = iota
	Int
	Float
	Bool
	Function
	Array
)

func TypeOf(v any) TypeId {
	switch v.(type) {
	case int:
		return Int
	case float64:
		return Float
	case Func:
		return Function
	case []any:
		return Array
	default:
		return Invalid
	}
}

type Type struct {
	Id TypeId
}

type Func struct {
	//Params []string
	Address int
}
