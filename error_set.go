package arrgo

import "fmt"

var (
	INDEX_ERROR       error = fmt.Errorf("INDEX ERROR")
	SHAPE_ERROR       error = fmt.Errorf("Shape ERROR")
	DIMENTION_ERROR   error = fmt.Errorf("DIMENTION ERROR")
	TYPE_ERROR        error = fmt.Errorf("DATA TYPE ERROR")
	EMPTY_ARRAY_ERROR error = fmt.Errorf("EMPTY ARRAY ERROR")
	PARAMETER_ERROR   error = fmt.Errorf("PARAMETER ERROR")

	UNIMPLEMENT_ERROR error = fmt.Errorf("UNIMPLEMENT ERROR")
)
