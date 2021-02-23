package operator

// OP is the operator type
type OP string

var (
	AND OP = "AND"
	OR  OP = "OR"

	Eq    OP = "eq"
	NotEq OP = "not_eq"
	In    OP = "in"
	NotIn OP = "not_in"

	Contains    OP = "contains"
	NotContains OP = "not_contains"

	StartsWith    OP = "starts_with"
	NotStartsWith OP = "not_starts_with"

	EndsWith    OP = "ends_with"
	NotEndsWith OP = "not_ends_with"

	Lt  OP = "lt"
	Lte OP = "lte"
	Gt  OP = "gt"
	Gte OP = "gte"

	Any OP = "any"
)
