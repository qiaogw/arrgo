package arrgo

import "fmt"

//改变原始多维数组的形状，并返回改变后的多维数组的指引引用。
//不会创建新的数据副本。
//如果新的Shape的大小和原来多维数组的大小不同，则抛出异常。
func (a *Arrf) ReShape(Shape ...int) *Arrf {
	if a.Length() != ProductIntSlice(Shape) {
		fmt.Println("new Shape length does not equal to original array length.")
		panic(SHAPE_ERROR)
	}

	internalShape := make([]int, len(Shape))
	copy(internalShape, Shape)
	a.Shape = internalShape

	a.Strides = make([]int, len(a.Shape)+1)
	a.Strides[len(a.Shape)] = 1
	for i := len(a.Shape) - 1; i >= 0; i-- {
		a.Strides[i] = a.Strides[i+1] * a.Shape[i]
	}

	return a
}

//两个多维数组形状相同，则返回true， 否则返回false。
func (a *Arrf) SameShapeTo(b *Arrf) bool {
	return SameIntSlice(a.Shape, b.Shape)
}

//将多个两维数组在垂直方向上组合起来，形成新的多维数组。
//不影响原多维数组。
func Vstack(arrs ...*Arrf) *Arrf {
	for i := range arrs {
		if arrs[i].Ndims() > 2 {
			fmt.Println("in Vstack function, array dimension cannot bigger than 2.")
			panic(SHAPE_ERROR)
		}
	}
	if len(arrs) == 0 {
		return nil
	}
	if len(arrs) == 1 {
		return arrs[0].Copy()
	}

	return Concat(0, arrs...)
	//
	//var vlenSum int = 0
	//
	//var hlen int
	//if arrs[0].Ndims() == 1 {
	//	hlen = arrs[0].Shape[0]
	//	vlenSum += 1
	//} else {
	//	hlen = arrs[0].Shape[1]
	//	vlenSum += arrs[0].Shape[0]
	//}
	//for i := 1; i < len(arrs); i++ {
	//	var nextHen int
	//	if arrs[i].Ndims() == 1 {
	//		nextHen = arrs[i].Shape[0]
	//		vlenSum += 1
	//	} else {
	//		nextHen = arrs[i].Shape[1]
	//		vlenSum += arrs[i].Shape[0]
	//	}
	//	if hlen != nextHen {
	//		panic(SHAPE_ERROR)
	//	}
	//}
	//
	//Data := make([]float64, vlenSum*hlen)
	//var offset = 0
	//for i := range arrs {
	//	copy(Data[offset:], arrs[i].Data)
	//	offset += len(arrs[i].Data)
	//}
	//
	//return Array(Data, vlenSum, hlen)
}

//将多个两维数组在水平方向上组合起来，形成新的多维数组。
//不影响原多维数组。
func Hstack(arrs ...*Arrf) *Arrf {
	for i := range arrs {
		if arrs[i].Ndims() > 2 {
			panic(SHAPE_ERROR)
		}
	}
	if len(arrs) == 0 {
		return nil
	}
	if len(arrs) == 1 {
		return arrs[0].Copy()
	}

	return Concat(1, arrs...)

	//var hlenSum int = 0
	//var hBlockLens = make([]int, len(arrs))
	//var vlen int
	//if arrs[0].Ndims() == 1 {
	//	vlen = 1
	//	hlenSum += arrs[0].Shape[0]
	//	hBlockLens[0] = arrs[0].Shape[0]
	//} else {
	//	vlen = arrs[0].Shape[0]
	//	hlenSum += arrs[0].Shape[1]
	//	hBlockLens[0] = arrs[0].Shape[1]
	//}
	//for i := 1; i < len(arrs); i++ {
	//	var nextVlen int
	//	if arrs[i].Ndims() == 1 {
	//		nextVlen = 1
	//		hlenSum += arrs[i].Shape[0]
	//		hBlockLens[i] = arrs[i].Shape[0]
	//	} else {
	//		nextVlen = arrs[i].Shape[0]
	//		hlenSum += arrs[i].Shape[1]
	//		hBlockLens[i] = arrs[i].Shape[1]
	//	}
	//	if vlen != nextVlen {
	//		panic(SHAPE_ERROR)
	//	}
	//}
	//
	//Data := make([]float64, hlenSum*vlen)
	//for i := 0; i < vlen; i++ {
	//	var curPos = 0
	//	for j := 0; j < len(arrs); j++ {
	//		copy(Data[curPos+i*hlenSum:curPos+i*hlenSum+hBlockLens[j]], arrs[j].Data[i*hBlockLens[j]:(i+1)*hBlockLens[j]])
	//		curPos += hBlockLens[j]
	//	}
	//}
	//
	//return Array(Data, vlen, hlenSum)
}

//将多个多维数组在指定的轴上组合起来。
//一维数组默认扩充为2维，参考AtLeast2D函数。
func Concat(axis int, arrs ...*Arrf) *Arrf {
	if len(arrs) == 0 {
		return nil
	}
	if len(arrs) == 1 {
		return arrs[0].Copy()
	}

	for i := range arrs {
		AtLeast2D(arrs[i])
	}

	if axis >= arrs[0].Ndims() {
		fmt.Println("axis is bigger than dimensions num.")
		panic(PARAMETER_ERROR)
	}

	var newShape = make([]int, arrs[0].Ndims())
	for index, firstL := range arrs[0].Shape {
		if index == axis {
			newShape[index] += firstL
			for j := 1; j < len(arrs); j++ {
				newShape[index] += arrs[j].Shape[index]
			}
		} else {
			newShape[index] = firstL
			for j := 1; j < len(arrs); j++ {
				if firstL != arrs[j].Shape[index] {
					panic(SHAPE_ERROR)
				}
			}
		}
	}

	var times = 0
	if axis == 0 {
		times = 1
	} else {
		times = ProductIntSlice(arrs[0].Shape[0:axis])
	}

	var Data = make([]float64, ProductIntSlice(newShape))

	var curPos = 0
	for i := 0; i < times; i++ {
		for j := 0; j < len(arrs); j++ {
			var l = ProductIntSlice(arrs[j].Shape[axis:])
			copy(Data[curPos:curPos+l], arrs[j].Data[i*l:(i+1)*l])
			curPos += l
		}
	}

	return Array(Data, newShape...)
}

//将一维数组扩充为二维
func AtLeast2D(a *Arrf) *Arrf {
	if a == nil {
		return nil
	} else if a.Ndims() >= 2 {
		return a
	} else {
		newShpae := make([]int, 2)
		newShpae[0] = 1
		newShpae[1] = a.Shape[0]
		a.Shape = newShpae
		return a
	}
}

//将数组内部的元素铺平返回，创建新的数据副本。
func (a *Arrf) Flatten() *Arrf {
	ra := make([]float64, len(a.Data))
	copy(ra, a.Data)
	return Array(ra, len(a.Data))
}
