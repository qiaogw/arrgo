package arrgo

import (
	"fmt"
	"strings"
)

type Arrb struct {
	Shape   []int
	Strides []int
	Data    []bool
}

//通过[]bool，形状来创建多维数组。
//输入参数1：Data []bool，以·C· 顺序存储，作为多维数组的输入数据，内部复制一份新的internalData，不改变data。
//输入参数2：Shape ...int，指定多维数组的形状，多维，类似numpy中的Shape。
//	如果某一个（仅支持一个维度）维度为负数，则根据len(Data)推断该维度的大小。
//情况1：如果不指定Shape，而且data为nil，则创建一个空的*Arrb。
//情况2：如果不指定Shape，而且data不为nil，则创建一个len(Data)大小的一维*Arrb。
//情况3：如果指定Shape，而且data不为nil，则根据data大小创建多维数组，如果len(Data)不等于Shape，或者len(Data)不能整除Shape，抛出异常。
//情况4：如果指定Shape，而且data为nil，则创建Shape大小的全为false的多维数组。
func ArrayB(Data []bool, Shape ...int) *Arrb {
	if len(Shape) == 0 && Data == nil {
		return &Arrb{
			Shape:   []int{0},
			Strides: []int{0, 1},
			Data:    []bool{},
		}
	}

	if len(Shape) == 0 && Data != nil {
		internalData := make([]bool, len(Data)) //复制data，不影响输入的值。
		copy(internalData, Data)
		return &Arrb{
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

		return &Arrb{
			Shape:   internalShape,
			Strides: Strides,
			Data:    make([]bool, length),
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

	internalData := make([]bool, len(Data))
	copy(internalData, Data)

	return &Arrb{
		Shape:   internalShape,
		Strides: Strides,
		Data:    internalData,
	}
}

//创建Shape形状的多维布尔数组，全部填充为fillvalue。
//必须指定Shape，否则抛出异常。
func FillB(fullValue bool, Shape ...int) *Arrb {
	if len(Shape) == 0 {
		fmt.Println("Shape is empty!")
		panic(SHAPE_ERROR)
	}
	arr := ArrayB(nil, Shape...)
	for i := range arr.Data {
		arr.Data[i] = fullValue
	}

	return arr
}

//创建全为false，形状位Shape的多维布尔数组
func EmptyB(Shape ...int) (a *Arrb) {
	a = FillB(false, Shape...)
	return
}

func (a *Arrb) String() (s string) {
	switch {
	case a == nil:
		return "<nil>"
	case a.Shape == nil || a.Strides == nil || a.Data == nil:
		return "<nil>"
	case a.Strides[0] == 0:
		return "[]"
	}

	stride := a.Strides[len(a.Strides)-2]
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

//如果多维布尔数组元素都为真，返回true，否则返回false。
func (ab *Arrb) AllTrues() bool {
	if len(ab.Data) == 0 {
		return false
	}
	for _, v := range ab.Data {
		if v == false {
			return false
		}
	}
	return true
}

//如果多维布尔数组元素都为假，返回false，否则返回true。
func (ab *Arrb) AnyTrue() bool {
	if len(ab.Data) == 0 {
		return false
	}
	for _, v := range ab.Data {
		if v == true {
			return true
		}
	}
	return false
}

//返回多维数组中真值的个数。
func (a *Arrb) Sum() int {
	sum := 0
	for _, v := range a.Data {
		if v {
			sum++
		}
	}
	return sum
}
