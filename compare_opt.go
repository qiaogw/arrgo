package arrgo

import "fmt"

func (a *Arrf) Greater(b *Arrf) *Arrb {
	if len(a.Data) == 0 || len(b.Data) == 0 {
		panic(EMPTY_ARRAY_ERROR)
	}
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = v > b.Data[i]
	}
	return t
}

func (a *Arrf) GreaterEqual(b *Arrf) *Arrb {
	if len(a.Data) == 0 || len(b.Data) == 0 {
		panic(EMPTY_ARRAY_ERROR)
	}
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = v >= b.Data[i]
	}
	return t
}

func (a *Arrf) Less(b *Arrf) *Arrb {
	if len(a.Data) == 0 || len(b.Data) == 0 {
		panic(EMPTY_ARRAY_ERROR)
	}
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = v < b.Data[i]
	}
	return t
}

func (a *Arrf) LessEqual(b *Arrf) *Arrb {
	if len(a.Data) == 0 || len(b.Data) == 0 {
		panic(EMPTY_ARRAY_ERROR)
	}
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = v <= b.Data[i]
	}
	return t
}

//判断两个Array相对位置的元素是否相同，返回Arrb。
//如果两个Array任一为空，或者形状不同，则抛出异常。
func (a *Arrf) Equal(b *Arrf) *Arrb {
	if len(a.Data) == 0 || len(b.Data) == 0 {
		fmt.Println("empty array.")
		panic(EMPTY_ARRAY_ERROR)
	}
	if !SameIntSlice(a.Shape, b.Shape) {
		fmt.Println("Shape not same.")
		panic(SHAPE_ERROR)
	}
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = v == b.Data[i]
	}
	return t
}

func (a *Arrf) NotEqual(b *Arrf) *Arrb {
	if len(a.Data) == 0 || len(b.Data) == 0 {
		panic(EMPTY_ARRAY_ERROR)
	}
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = v != b.Data[i]
	}
	return t
}
func Greater(a, b *Arrf) *Arrb {
	return a.Greater(b)
}

func GreaterEqual(a, b *Arrf) *Arrb {
	return a.GreaterEqual(b)
}

func Less(a, b *Arrf) *Arrb {
	return a.Less(b)
}

func LessEqual(a, b *Arrf) *Arrb {
	return a.LessEqual(b)
}

func Equal(a, b *Arrf) *Arrb {
	return a.Equal(b)
}

func NotEqual(a, b *Arrf) *Arrb {
	return a.NotEqual(b)
}

func (a *Arrf) Sort(axis ...int) *Arrf {
	ax := -1
	if len(axis) == 0 {
		ax = a.Ndims() - 1
	} else {
		ax = axis[0]
	}

	axisShape, axisSt, axis1St := a.Shape[ax], a.Strides[ax], a.Strides[ax+1]
	if axis1St == 1 {
		Hsort(axisSt, a.Data)
	} else {
		Vsort(axis1St, a.Data[0:axisShape*axis1St])
	}

	return a
}

func Sort(a *Arrf, axis ...int) *Arrf {
	return a.Copy().Sort(axis...)
}

func (a *Arrf) Size() int {
	return ProductIntSlice(a.Shape)
}
