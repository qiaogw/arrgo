package arrgo

import (
	"sort"

	asm "github.com/qiaogw/arrgo/internal"
	//"github.com/ledao/arrgo/internal"
)

func (a *Arrf) Sum(axis ...int) *Arrf {
	if len(axis) == 0 || len(axis) >= a.Ndims() {
		tot := float64(0)
		for _, v := range a.Data {
			tot += v
		}
		return Fill(tot, 1)
	}

	//对axis进行排序，按照从大到小的顺序进行规约
	sort.IntSlice(axis).Sort()
	//规约后的数组的形状
	restAxis := make([]int, len(a.Shape)-len(axis))
	//对a进行复制，所有的操作都作用于临时变量ta中，最后将ta返回
	ta := a.Copy()

axisR:
	for i, t := 0, 0; i < len(ta.Shape); i++ {
		for _, w := range axis {
			if i == w {
				continue axisR
			}
		}
		restAxis[t] = ta.Shape[i]
		t++
	}

	//数组的元素的个数保存到ln中
	ln := ta.Strides[0]
	//对每个指定的轴，顺寻进行规约
	for k := 0; k < len(axis); k++ {
		//如果轴大小为1，则不需要任何操作
		if ta.Shape[axis[k]] == 1 {
			continue
		}
		//获取当前轴的大小v，当前轴的跨度wd，以及下一个轴的跨度st
		v, wd, st := ta.Shape[axis[k]], ta.Strides[axis[k]], ta.Strides[axis[k]+1]
		//如果下一个轴st的跨度为1，则说明当前轴为最后一个轴，只需要每wd个跨度进行一个规约即可
		if st == 1 {
			//每wd个数据进行一次规约，结果依次放到开始的位置
			asm.Hadd(uint64(wd), ta.Data)
			ln /= v
			ta.Data = ta.Data[:ln]
			continue
		}
		//如果不是最后一个轴，则在该轴上进行规约
		for w := 0; w < ln; w += wd {
			t := ta.Data[w/wd*st : (w/wd+1)*st]
			copy(t, ta.Data[w:w+st])
			for i := 1; i*st+1 < wd; i++ {
				asm.Vadd(t, ta.Data[w+(i)*st:w+(i+1)*st])
			}
		}
		ln /= v
		ta.Data = ta.Data[:ln]
	}
	ta.Shape = restAxis

	tmp := 1
	for i := len(restAxis); i > 0; i-- {
		ta.Strides[i] = tmp
		tmp *= restAxis[i-1]
	}
	ta.Strides[0] = tmp
	ta.Data = ta.Data[:tmp]
	ta.Strides = ta.Strides[:len(restAxis)+1]
	return ta
}

func Sum(a *Arrf, axis ...int) *Arrf {
	return a.Sum(axis...)
}

func (a *Arrf) Mean(axis ...int) *Arrf {
	if len(axis) == 0 || len(axis) >= a.Ndims() {
		tot := float64(0)
		for _, v := range a.Data {
			tot += v
		}
		return Fill(tot/float64(a.Strides[0]), 1)
	}

	sort.IntSlice(axis).Sort()
	selectShape := make([]int, len(axis))
	for i := range selectShape {
		selectShape[i] = a.Shape[axis[i]]
	}
	N := ProductIntSlice(selectShape)

	ta := a.Sum(axis...)

	return ta.DivC(float64(N))
}

func Mean(a *Arrf, axis ...int) *Arrf {
	return a.Mean(axis...)
}

func (a *Arrf) Var(axis ...int) *Arrf {
	a2 := a.Mul(a).Sum(axis...)
	m := a.Mean(axis...)
	var N int
	if len(axis) == 0 || len(axis) >= a.Ndims() {
		N = ProductIntSlice(a.Shape)
	} else {
		selectShape := make([]int, len(axis))
		for i := range selectShape {
			selectShape[i] = a.Shape[axis[i]]
		}
		N = ProductIntSlice(selectShape)
	}

	m2 := m.Mul(m).MulC(float64(N))
	a_m_2 := a.Sum(axis...).Mul(m).MulC(2)
	return a2.Sub(a_m_2).Add(m2).DivC(float64(N))
}

func Var(a *Arrf, axis ...int) *Arrf {
	return a.Var(axis...)
}

func (a *Arrf) Std(axis ...int) *Arrf {
	return Sqrt(a.Var(axis...))
}

func Std(a *Arrf, axis ...int) *Arrf {
	return a.Std(axis...)
}

func (a *Arrf) Min(axis ...int) *Arrf {
	if len(axis) == 0 || len(axis) >= a.Ndims() {
		minValue := a.Data[0]
		for _, v := range a.Data {
			if minValue > v {
				minValue = v
			}
		}
		return Fill(minValue, 1)
	}

	sort.IntSlice(axis).Sort()
	restAxis := make([]int, len(a.Shape)-len(axis))
	ta := a.Copy()
axisR:
	for i, t := 0, 0; i < len(ta.Shape); i++ {
		for _, w := range axis {
			if i == w {
				continue axisR
			}
		}
		restAxis[t] = ta.Shape[i]
		t++
	}

	//数组的元素的个数保存到ln中
	ln := ta.Strides[0]
	//对每个指定的轴，顺寻进行规约
	for k := 0; k < len(axis); k++ {
		//如果轴大小为1，则不需要任何操作
		if ta.Shape[axis[k]] == 1 {
			continue
		}
		//获取当前轴的大小v，当前轴的跨度wd，以及下一个轴的跨度st
		v, wd, st := ta.Shape[axis[k]], ta.Strides[axis[k]], ta.Strides[axis[k]+1]
		//如果下一个轴st的跨度为1，则说明当前轴为最后一个轴，只需要每wd个跨度进行一个规约即可
		if st == 1 {
			//每wd个数据进行一次规约，结果依次放到开始的位置
			Hmin(wd, ta.Data)
			ln /= v
			ta.Data = ta.Data[:ln]
			continue
		}
		//如果不是最后一个轴，则在该轴上进行规约
		for w := 0; w < ln; w += wd {
			t := ta.Data[w/wd*st : (w/wd+1)*st]
			copy(t, ta.Data[w:w+st])
			for i := 1; i*st+1 < wd; i++ {
				Vmin(t, ta.Data[w+(i)*st:w+(i+1)*st])
			}
		}
		ln /= v
		ta.Data = ta.Data[:ln]
	}

	ta.Shape = restAxis

	tmp := 1
	for i := len(restAxis); i > 0; i-- {
		ta.Strides[i] = tmp
		tmp *= restAxis[i-1]
	}
	ta.Strides[0] = tmp
	ta.Strides = ta.Strides[:len(restAxis)+1]
	return ta
}

func Min(a *Arrf, axis ...int) *Arrf {
	return a.Min(axis...)
}

func (a *Arrf) Max(axis ...int) *Arrf {
	if len(axis) == 0 || len(axis) >= a.Ndims() {
		maxValue := a.Data[0]
		for _, v := range a.Data {
			if maxValue < v {
				maxValue = v
			}
		}
		return Fill(maxValue, 1)
	}

	sort.IntSlice(axis).Sort()
	restAxis := make([]int, len(a.Shape)-len(axis))
	ta := a.Copy()
axisR:
	for i, t := 0, 0; i < len(ta.Shape); i++ {
		for _, w := range axis {
			if i == w {
				continue axisR
			}
		}
		restAxis[t] = ta.Shape[i]
		t++
	}

	//数组的元素的个数保存到ln中
	ln := ta.Strides[0]
	//对每个指定的轴，顺寻进行规约
	for k := 0; k < len(axis); k++ {
		//如果轴大小为1，则不需要任何操作
		if ta.Shape[axis[k]] == 1 {
			continue
		}
		//获取当前轴的大小v，当前轴的跨度wd，以及下一个轴的跨度st
		v, wd, st := ta.Shape[axis[k]], ta.Strides[axis[k]], ta.Strides[axis[k]+1]
		//如果下一个轴st的跨度为1，则说明当前轴为最后一个轴，只需要每wd个跨度进行一个规约即可
		if st == 1 {
			//每wd个数据进行一次规约，结果依次放到开始的位置
			Hmax(wd, ta.Data)
			ln /= v
			ta.Data = ta.Data[:ln]
			continue
		}
		//如果不是最后一个轴，则在该轴上进行规约
		for w := 0; w < ln; w += wd {
			t := ta.Data[w/wd*st : (w/wd+1)*st]
			copy(t, ta.Data[w:w+st])
			for i := 1; i*st+1 < wd; i++ {
				Vmax(t, ta.Data[w+(i)*st:w+(i+1)*st])
			}
		}
		ln /= v
		ta.Data = ta.Data[:ln]
	}

	ta.Shape = restAxis

	tmp := 1
	for i := len(restAxis); i > 0; i-- {
		ta.Strides[i] = tmp
		tmp *= restAxis[i-1]
	}
	ta.Strides[0] = tmp
	ta.Strides = ta.Strides[:len(restAxis)+1]
	return ta
}

func Max(a *Arrf, axis ...int) *Arrf {
	return a.Max(axis...)
}

func (a *Arrf) ArgMax(axis int) *Arrf {
	if axis < 0 {
		axis = axis + len(a.Shape)
	}
	restAxis := make([]int, len(a.Shape)-1)
	ta := a.Copy()
	for i, t := 0, 0; i < len(ta.Shape); i++ {
		if i == axis {
			continue
		}
		restAxis[t] = ta.Shape[i]
		t++
	}

	//数组的元素的个数保存到ln中
	ln := ta.Strides[0]

	//获取当前轴的大小v，当前轴的跨度wd，以及下一个轴的跨度st
	v, wd, st := ta.Shape[axis], ta.Strides[axis], ta.Strides[axis+1]
	//如果下一个轴st的跨度为1，则说明当前轴为最后一个轴，只需要每wd个跨度进行一个规约即可
	if st == 1 {
		//每wd个数据进行一次规约，结果依次放到开始的位置
		Hargmax(wd, ta.Data)
		ln /= v
		ta.Data = ta.Data[:ln]
	} else {
		//如果不是最后一个轴，则在该轴上进行规约
		td := make([]float64, 0, ln/wd)
		for w := 0; w < ln; w += wd {
			Vargmax(st, ta.Data[w:w+wd])
			td = append(td, ta.Data[w : w+wd][:st]...)
		}
		ln /= v
		ta.Data = td
	}

	ta.Shape = restAxis

	tmp := 1
	for i := len(restAxis); i > 0; i-- {
		ta.Strides[i] = tmp
		tmp *= restAxis[i-1]
	}
	ta.Strides[0] = tmp
	ta.Strides = ta.Strides[:len(restAxis)+1]
	return ta
}

func ArgMax(a *Arrf, axis int) *Arrf {
	return a.ArgMax(axis)
}

//fixme has bug
func (a *Arrf) ArgMin(axis int) *Arrf {
	if axis < 0 {
		axis = axis + len(a.Shape)
	}
	restAxis := make([]int, len(a.Shape)-1)
	ta := a.Copy()
	for i, t := 0, 0; i < len(ta.Shape); i++ {
		if i == axis {
			continue
		}
		restAxis[t] = ta.Shape[i]
		t++
	}

	//数组的元素的个数保存到ln中
	ln := ta.Strides[0]

	//获取当前轴的大小v，当前轴的跨度wd，以及下一个轴的跨度st
	v, wd, st := ta.Shape[axis], ta.Strides[axis], ta.Strides[axis+1]
	//如果下一个轴st的跨度为1，则说明当前轴为最后一个轴，只需要每wd个跨度进行一个规约即可
	if st == 1 {
		//每wd个数据进行一次规约，结果依次放到开始的位置
		Hargmin(wd, ta.Data)
		ln /= v
		ta.Data = ta.Data[:ln]
	} else {
		//如果不是最后一个轴，则在该轴上进行规约
		td := make([]float64, 0, ln/wd)
		for w := 0; w < ln; w += wd {
			Vargmin(st, ta.Data[w:w+wd])
			td = append(td, ta.Data[w : w+wd][:st]...)
		}
		ln /= v
		ta.Data = td
	}

	ta.Shape = restAxis

	tmp := 1
	for i := len(restAxis); i > 0; i-- {
		ta.Strides[i] = tmp
		tmp *= restAxis[i-1]
	}
	ta.Strides[0] = tmp
	ta.Strides = ta.Strides[:len(restAxis)+1]
	return ta
}

func ArgMin(a *Arrf, axis int) *Arrf {
	return a.ArgMin(axis)
}
