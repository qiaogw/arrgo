package arrgo

import (
	"testing"
)

func TestArrf_AddC(t *testing.T) {
	arr := Arange(0, 10, 2)
	add := arr.AddC(2)
	if !add.Equal(Array([]float64{2, 4, 6, 8, 10})).AllTrues() {
		t.Error("Expected [2,4,6,8,10], got ", add)
	}
}

//测试nil
func TestArrf_AddC_ShapeERROR(t *testing.T) {
	var arr *Arrf = nil

	defer func() {
		var rec = recover()
		if rec != SHAPE_ERROR {
			t.Error("Expected Shape ERROR, got ", rec)
		}
	}()
	arr.AddC(10)
}

//测试空array
func TestArrf_AddC_ShapeERROR2(t *testing.T) {
	var arr *Arrf = Array([]float64{})

	defer func() {
		var rec = recover()
		if rec != SHAPE_ERROR {
			t.Error("Expected Shape ERROR, got ", rec)
		}
	}()
	arr.AddC(10)
}

func TestArrf_Add(t *testing.T) {
	var a = Array([]float64{1, 2, 3, 4, 5, 6}, 2, 3)
	var b = Array([]float64{6, 5, 4, 3, 2, 1}, 2, 3)
	var c = a.Add(b)
	if !c.Equal(Fill(7, 2, 3)).AllTrues() {
		t.Error("Expected [[7,7,7],[7,7,7]], got ", c)
	}
}

//func TestArrf_Add_NilException(t *testing.T) {
//    var a = Array([]float64{1,2,3,4,5,6}, 2, 3)
//
//    defer func(){
//       var rec = recover()
//        if rec != SHAPE_ERROR {
//            t.Error("Expected Shape ERROR, got ", rec)
//        }
//    }()
//    a.Add(nil)
//}

func TestArrf_Add_NDimException(t *testing.T) {
	var a = Array([]float64{1, 2, 3, 4, 5, 6})
	var b = Array([]float64{1, 2, 3}, 3, 1)
	defer func() {
		var rec = recover()
		if rec != SHAPE_ERROR {
			t.Error("Expected Shape ERROR, got ", rec)
		}
	}()
	a.Add(b)
}

func BenchmarkDotProd(b *testing.B) {
	a := Array([]float64{1, 2, 3})
	c := Array([]float64{4, 5, 6})
	for i := 0; i < 100; i++ {
		a.DotProd(c)
	}
}
