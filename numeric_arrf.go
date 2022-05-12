package arrgo

import (
	"fmt"
	"math"
	"strings"
)

type Arrf struct {
	Shape   []int
	Strides []int
	Data    []float64
}

//通过[]float64，形状来创建多维数组。
//输入参数1：data []float64，以·C· 顺序存储，作为多维数组的输入数据，内部复制一份新的internalData，不改变data。
//输入参数2：Shape ...int，指定多维数组的形状，多维，类似numpy中的Shape。
//	如果某一个（仅支持一个维度）维度为负数，则根据len(Data)推断该维度的大小。
//情况1：如果不指定Shape，而且data为nil，则创建一个空的*Arrf。
//情况2：如果不指定Shape，而且data不为nil，则创建一个len(Data)大小的一维*Arrf。
//情况3：如果指定Shape，而且data不为nil，则根据data大小创建多维数组，如果len(Data)不等于Shape，或者len(Data)不能整除Shape，抛出异常。
//情况4：如果指定Shape，而且data为nil，则创建Shape大小的全为0.0的多维数组。
func Array(Data []float64, Shape ...int) *Arrf {
	if len(Shape) == 0 && Data == nil {
		return &Arrf{
			Shape:   []int{0},
			Strides: []int{0, 1},
			Data:    []float64{},
		}
	}

	if len(Shape) == 0 && Data != nil {
		internalData := make([]float64, len(Data)) //复制data，不影响输入的值。
		copy(internalData, Data)
		return &Arrf{
			Shape:   []int{len(Data)},
			Strides: []int{len(Data), 1},
			Data:    internalData,
		}
	}

	if Data == nil {
		for _, v := range Shape {
			if v <= 0 {
				fmt.Println("Shape should be positive when Data is nill")
				panic(SHAPE_ERROR)
			}
		}
		length := ProductIntSlice(Shape)
		internalShape := make([]int, len(Shape))
		copy(internalShape, Shape)
		Strides := make([]int, len(Shape)+1)
		Strides[len(Shape)] = 1
		for i := len(Shape) - 1; i >= 0; i-- {
			Strides[i] = Strides[i+1] * internalShape[i]
		}

		return &Arrf{
			Shape:   internalShape,
			Strides: Strides,
			Data:    make([]float64, length),
		}
	}

	var dataLength = len(Data)
	negativeIndex := -1
	internalShape := make([]int, len(Shape))
	copy(internalShape, Shape)
	for k, v := range Shape {
		if v < 0 {
			if negativeIndex < 0 {
				negativeIndex = k
				internalShape[k] = 1
			} else {
				fmt.Println("Shape can only have one negative demention.")
				panic(SHAPE_ERROR)
			}
		}
	}
	ShapeLength := ProductIntSlice(internalShape)

	if dataLength < ShapeLength {
		fmt.Println("Data length is shorter than Shape length.")
		panic(SHAPE_ERROR)
	}
	if (dataLength % ShapeLength) != 0 {
		fmt.Println("Data length cannot divided by Shape length")
		panic(SHAPE_ERROR)
	}

	if negativeIndex >= 0 {
		internalShape[negativeIndex] = dataLength / ShapeLength
	}

	Strides := make([]int, len(internalShape)+1)
	Strides[len(internalShape)] = 1
	for i := len(internalShape) - 1; i >= 0; i-- {
		Strides[i] = Strides[i+1] * internalShape[i]
	}

	internalData := make([]float64, len(Data))
	copy(internalData, Data)

	return &Arrf{
		Shape:   internalShape,
		Strides: Strides,
		Data:    internalData,
	}
}

// 通过指定起始、终止和步进量来创建一维Array。
// 输入参数： vals，可以有三种情况，详见下面描述。
// 情况1：Arange(stop): 以0开始的序列，创建Array [0, 0+(-)1, ..., stop)，不包括stop，stop符号决定升降序。
// 情况2：Arange(start, stop):创建Array [start, start +(-)1, ..., stop)，如果start小于start则递增，否则递减。
// 情况3：Arange(start, stop, step):创建Array [start, start + step, ..., stop)，step符号决定升降序。
// 输入参数多于三个的都会被忽略。
// 输出序列为“整型数”序列。
func Arange(vals ...int) *Arrf {
	var start, stop, step int = 0, 0, 1

	switch len(vals) {
	case 0:
		fmt.Println("range function should have range")
		panic(PARAMETER_ERROR)
	case 1:
		if vals[0] <= 0 {
			step = -1
			stop = vals[0] + 1
		} else {
			stop = vals[0] - 1
		}
	case 2:
		if vals[1] < vals[0] {
			step = -1
			stop = vals[1] + 1
		} else {
			stop = vals[1] - 1
		}
		start = vals[0]
	default:
		if vals[1] < vals[0] {
			if vals[2] >= 0 {
				fmt.Println("increment should be negative.")
				panic(PARAMETER_ERROR)
			}
			stop = vals[1] + 1
		} else {
			if vals[2] <= 0 {
				fmt.Println("increment should be positive.")
				panic(PARAMETER_ERROR)
			}
			stop = vals[1] - 1
		}
		start, step = vals[0], vals[2]
	}

	a := Array(nil, int(math.Abs(float64((stop-start)/step)))+1)
	for i, v := 0, start; i < len(a.Data); i, v = i+1, v+step {
		a.Data[i] = float64(v)
	}
	return a
}

//判断Arrf是否为空数组。
//如果内部的data长度为0或者为nil，返回true，否则位false。
func (a *Arrf) IsEmpty() bool {
	return len(a.Data) == 0 || a.Data == nil
}

//创建Shape形状的多维数组，全部填充为fillvalue。
//必须指定Shape，否则抛出异常。
func Fill(fillValue float64, Shape ...int) *Arrf {
	if len(Shape) == 0 {
		fmt.Println("Shape is empty!")
		panic(SHAPE_ERROR)
	}
	arr := Array(nil, Shape...)
	for i := range arr.Data {
		arr.Data[i] = fillValue
	}

	return arr
}

//根据Shape创建全为1.0的多维数组。
func Ones(Shape ...int) *Arrf {
	return Fill(1, Shape...)
}

//根据输入的多维数组的形状创建全1的多维数组。
func OnesLike(a *Arrf) *Arrf {
	return Ones(a.Shape...)
}

//根据Shape创建全为0的多维数组。
func Zeros(Shape ...int) *Arrf {
	return Fill(0, Shape...)
}

//根据输入的多维数组的形状创建全0的多维数组。
func ZerosLike(a *Arrf) *Arrf {
	return Zeros(a.Shape...)
}

// String Satisfies the Stringer interface for fmt package
func (a *Arrf) String() (s string) {
	switch {
	case a == nil:
		return "<nil>"
	case a.Data == nil || a.Shape == nil || a.Strides == nil:
		return "<nil>"
	case a.Strides[0] == 0:
		return "[]"
	case len(a.Shape) == 1:
		return fmt.Sprint(a.Data)
		//strs := make([]string, len(a.Data))
		//for i := range a.Data {
		//	strs[i] = string(strconv.FormatFloat(a.Data[i], 'f', -1, 64))
		//
		//}
		//return strings.Join(strs, ", ")
	}

	stride := a.Shape[len(a.Shape)-1]

	for i, k := 0, 0; i+stride <= len(a.Data); i, k = i+stride, k+1 {

		t := ""
		for j, v := range a.Strides {
			if i%v == 0 && j < len(a.Strides)-2 {
				t += "["
			}
		}

		s += strings.Repeat(" ", len(a.Shape)-len(t)-1) + t
		s += fmt.Sprint(a.Data[i : i+stride])

		t = ""
		for j, v := range a.Strides {
			if (i+stride)%v == 0 && j < len(a.Strides)-2 {
				t += "]"
			}
		}

		s += t + strings.Repeat(" ", len(a.Shape)-len(t)-1)
		if i+stride != len(a.Data) {
			s += "\n"
			if len(t) > 0 {
				s += "\n"
			}
		}
	}
	return
}

//获取index指定位置的元素。
//index必须在Shape规定的范围内，否则会抛出异常。
//index的长度必须小于等于维度的个数，否则会抛出异常。
//如果index的个数小于维度个数，则会取后面的第一个值。
func (a *Arrf) At(index ...int) float64 {
	idx := a.valIndex(index...)
	return a.Data[idx]
}

//详见At函数。
func (a *Arrf) Get(index ...int) float64 {
	return a.At(index...)
}

//At函数的内部实现，返回index指定的元素在切片中的位置，如果有错误，则返回error。
func (a *Arrf) valIndex(index ...int) int {
	idx := 0
	if len(index) > len(a.Shape) {
		fmt.Println("index len should not longer than Shape.")
		panic(INDEX_ERROR)
	}
	for i, v := range index {
		if v >= a.Shape[i] || v < 0 {
			fmt.Println("index value out of range.")
			panic(INDEX_ERROR)
		}
		idx += v * a.Strides[i+1]
	}
	return idx
}

//获取多维数组元素的个数。
func (a *Arrf) Length() int {
	return len(a.Data)
}

//创建一个n X n 的2维单位矩阵(数组)。
func Eye(n int) *Arrf {
	arr := Zeros(n, n)
	for i := 0; i < n; i++ {
		arr.Set(1, i, i)
	}
	return arr
}

//Eye的另一种称呼，详见Eye函数。
func Identity(n int) *Arrf {
	return Eye(n)
}

//指定位置的元素被新值替换。
//如果index的超出范围则会抛出异常。
//返回当前数组的指引，方便后续的连续操作。
func (a *Arrf) Set(value float64, index ...int) *Arrf {
	idx := a.valIndex(index...)

	a.Data[idx] = value
	return a
}

//返回多维数组的内部数组元素。
//对返回值的操作会影响多维数组，一定谨慎操作。
func (a *Arrf) Values() []float64 {
	return a.Data
}

//根据[start, stop]指定的区间，创建包含num个元素的一维数组。
func Linspace(start, stop float64, num int) *Arrf {
	var Data = make([]float64, num)
	var startF, stopF = start, stop
	if startF <= stopF {
		var step = (stopF - startF) / (float64(num - 1.0))
		for i := range Data {
			Data[i] = startF + float64(i)*step
		}
		return Array(Data, num)
	} else {
		var step = (startF - stopF) / (float64(num - 1.0))
		for i := range Data {
			Data[i] = startF - float64(i)*step
		}
		return Array(Data, num)
	}
}

//复制一个形状一样，但是数据被深度复制的多维数组。
func (a *Arrf) Copy() *Arrf {
	b := ZerosLike(a)
	copy(b.Data, a.Data)
	return b
}

//返回多维数组的维度数目。
func (a *Arrf) Ndims() int {
	return len(a.Shape)
}

//Returns ta view of the array with axes transposed.
//根据指定的轴顺序，生成一个新的调整后的多维数组。
//如果是1维数组，则没有任何变化。
//如果是2维数组，则行列交换。
//如果是n维数组，则根据指定的顺序调整，生成新的多维数组。
//输入参数1：如果不指定输入参数，则轴顺序全部反序；如果指定参数则个数必须和轴个数相同，否则抛出异常。
//fixme 这里的实现效率不高，后面有时间需要提升一下。
func (a *Arrf) Transpose(axes ...int) *Arrf {
	var n = a.Ndims()
	var permutation []int
	var nShape []int

	switch len(axes) {
	case 0:
		permutation = make([]int, n)
		nShape = make([]int, n)
		for i := range permutation {
			permutation[i] = n - i
		}
		for i := 0; i < n; i++ {
			permutation[i] = n - 1 - i
			nShape[i] = a.Shape[permutation[i]]
		}

	case n:
		permutation = axes
		nShape = make([]int, n)
		for i := range nShape {
			nShape[i] = a.Shape[permutation[i]]
		}

	default:
		fmt.Println("axis number wrong.")
		panic(DIMENTION_ERROR)
	}

	var totalIndexSize = 1
	for i := range a.Shape {
		totalIndexSize *= a.Shape[i]
	}

	var indexsSrc = make([][]int, totalIndexSize)
	var indexsDst = make([][]int, totalIndexSize)

	var b = Zeros(nShape...)
	var index = make([]int, n)
	for i := 0; i < totalIndexSize; i++ {
		tindexSrc := make([]int, n)
		copy(tindexSrc, index)
		indexsSrc[i] = tindexSrc
		var tindexDst = make([]int, n)
		for j := range tindexDst {
			tindexDst[j] = index[permutation[j]]
		}
		indexsDst[i] = tindexDst

		var j = n - 1
		index[j]++
		for {
			if j > 0 && index[j] >= a.Shape[j] {
				index[j-1]++
				index[j] = 0
				j--
			} else {
				break
			}
		}
	}
	for i := range indexsSrc {
		b.Set(a.Get(indexsSrc[i]...), indexsDst[i]...)
	}
	return b
}
