package arrgo

import "testing"

func TestArrf_ReShape(t *testing.T) {
	arr := Array([]float64{1, 2, 3, 4, 5, 6}, 2, 3)
	arr2 := arr.ReShape(3, 2)

	if !SameIntSlice(arr.Strides, []int{6, 2, 1}) {
		t.Error("Expected [6,2,1], got ", arr2.Strides)
	}
	if !SameIntSlice(arr.Shape, []int{3, 2}) {
		t.Error("Expected [3, 2], got ", arr.Shape)
	}
	if !SameIntSlice(arr2.Shape, []int{3, 2}) {
		t.Error("Expected [3, 2], got ", arr2.Shape)
	}
}

func TestArrf_ReShapeException(t *testing.T) {
	defer func() {
		r := recover()
		if r != SHAPE_ERROR {
			t.Error("Expected Shape error, got ", r)
		}
	}()

	Arange(4).ReShape(5)
}

func TestArrf_SameShapeTo(t *testing.T) {
	a := Arange(4).ReShape(2, 2)
	b := Array([]float64{3, 4, 5, 6}, 2, 2)
	if a.SameShapeTo(b) != true {
		t.Errorf("Expected true, got %t", a.SameShapeTo(b))
	}
}

func TestVstack(t *testing.T) {
	if Vstack() != nil {
		t.Errorf("Expected nil, got %s", Vstack())
	}

	a := Arange(3)
	stacked := Vstack(a)
	if !stacked.Equal(Arange(3)).AllTrues() {
		t.Errorf("Expected [0, 1, 2], got %s", stacked)
	}

	b := Array([]float64{3, 4, 5})
	stacked = Vstack(a, b)
	if !stacked.Equal(Array([]float64{0, 1, 2, 3, 4, 5}, 2, 3)).AllTrues() {
		t.Errorf("Expected [[0 1 2] [3 4 5]], got %s", stacked)
	}

	a = Arange(2)
	b = Arange(4).ReShape(2, 2)
	stacked = Vstack(a, b)
	if !stacked.Equal(Array([]float64{0, 1, 0, 1, 2, 3}, 3, 2)).AllTrues() {
		t.Errorf("Expected [[0,1], [0,1], [2, 3]], got %s", stacked)
	}
}

func TestVstackException(t *testing.T) {
	a := Arange(4).ReShape(1, 2, 2)
	defer func() {
		r := recover()
		if r != SHAPE_ERROR {
			t.Errorf("Expected Shape error, got %s", r)
		}
	}()

	Vstack(a)
}

func TestVstackException2(t *testing.T) {
	a := Arange(4)
	b := Arange(5)
	defer func() {
		r := recover()
		if r != SHAPE_ERROR {
			t.Error("Expected Shape error, got ", r)
		}
	}()

	Vstack(a, b)
}

func TestHstack(t *testing.T) {
	if Hstack() != nil {
		t.Error("Expected nil, got ", Hstack())
	}

	a := Arange(3)
	stacked := Hstack(a)
	if !stacked.Equal(Arange(3)).AllTrues() {
		t.Error("Expected [0, 1, 2], got ", stacked)
	}
	a = a.ReShape(3, 1)
	b := Array([]float64{3, 4, 5}).ReShape(3, 1)
	stacked = Hstack(a, b)
	if !stacked.Equal(Array([]float64{0, 3, 1, 4, 2, 5}, 3, 2)).AllTrues() {
		t.Error("Expected [[0 3] [1 4], [2 5]], got ", stacked)
	}

	a = Arange(2).ReShape(2, 1)
	b = Arange(4).ReShape(2, 2)
	stacked = Hstack(a, b)
	if !stacked.Equal(Array([]float64{0, 0, 1, 1, 2, 3}, 2, 3)).AllTrues() {
		t.Error("Expected [[0, 0, 1], [1, 2, 3]], got ", stacked)
	}
}

func TestHstackException(t *testing.T) {
	a := Arange(4).ReShape(1, 2, 2)
	defer func() {
		r := recover()
		if r != SHAPE_ERROR {
			t.Error("Expected Shape error, got ", r)
		}
	}()

	Hstack(a)
}

func TestHstackException2(t *testing.T) {
	a := Arange(4).ReShape(4, 1)
	b := Arange(5).ReShape(5, 1)
	defer func() {
		r := recover()
		if r != SHAPE_ERROR {
			t.Error("Expected Shape error, got ", r)
		}
	}()

	Hstack(a, b)
}

func TestConcat(t *testing.T) {
	if Concat(0) != nil {
		t.Error("Expected nil, got ", Concat(0))
	}
	concated := Concat(0, Arange(2))
	if !concated.Equal(Arange(2)).AllTrues() {
		t.Error("Expected [0, 1], got ", concated)
	}

	a := Arange(3)
	b := Arange(1, 4)

	concated = Concat(0, a, b)
	if !concated.Equal(Array([]float64{0, 1, 2, 1, 2, 3}, 2, 3)).AllTrues() {
		t.Error("Expected [[0,1,2], [1,2,3]], got ", concated)
	}

	a = Arange(3)
	b = Arange(1, 4)

	concated = Concat(1, a, b)
	t.Log(concated)
	if !concated.Equal(Array([]float64{0, 1, 2, 1, 2, 3}, 1, 6)).AllTrues() {
		t.Error("Expected [[0,1,2,1,2,3]], got ", concated)
	}

}

func TestConcatException(t *testing.T) {
	a := Arange(4)
	b := Arange(1, 4)

	defer func() {
		r := recover()
		if r != SHAPE_ERROR {
			t.Error("Expected Shape error, got ", r)
		}
	}()

	Concat(0, a, b)
}

func TestConcatException2(t *testing.T) {
	a := Arange(4)
	b := Arange(1, 4)

	defer func() {
		r := recover()
		if r != PARAMETER_ERROR {
			t.Error("Expected PARAMETER_ERROR, got ", r)
		}
	}()

	Concat(2, a, b)
}

func TestAtLeast2D(t *testing.T) {
	a := Arange(10)
	AtLeast2D(a)
	if !SameIntSlice(a.Shape, []int{1, 10}) {
		t.Error("Expected [1, 10], got ", a.Shape)
	}

	a.ReShape(1, 1, 10)
	AtLeast2D(a)
	if !SameIntSlice(a.Shape, []int{1, 1, 10}) {
		t.Error("Expected [1, 1, 10], got ", a.Shape)
	}

	if AtLeast2D(nil) != nil {
		t.Error("Expected nil, got ", AtLeast2D(nil))
	}
}

func TestAtLeast2D2(t *testing.T) {
	if AtLeast2D(nil) != nil {
		t.Error("Expected nil, got ", AtLeast2D(nil))
	}

	arr := Arange(3)
	AtLeast2D(arr)

	if !SameIntSlice(arr.Shape, []int{1, 3}) {
		t.Error("expected true, got false")
	}

	arr = Arange(3).ReShape(3, 1)
	AtLeast2D(arr)

	if !SameIntSlice(arr.Shape, []int{3, 1}) {
		t.Error("expected true, got false")
	}
}

func TestArrf_Flatten(t *testing.T) {
	arr := Arange(3).ReShape(3, 1)
	flattened := arr.Flatten()

	if !flattened.SameShapeTo(Arange(3)) {
		t.Error("expected [3], got ", flattened.Shape)
	}
}
