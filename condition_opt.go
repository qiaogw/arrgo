package arrgo

func Where(cond *Arrb, tv, fv interface{}) *Arrf {
	t := Zeros(cond.Shape...)
	for i, v := range cond.Data {
		if v {
			switch tv.(type) {
			case float64:
				t.Data[i] = tv.(float64)
			case float32:
				t.Data[i] = float64(tv.(float32))
			case int:
				t.Data[i] = float64(tv.(int))
			case *Arrf:
				t.Data[i] = tv.(*Arrf).Data[i]
			default:
				panic(TYPE_ERROR)
			}
		} else {
			switch fv.(type) {
			case float64:
				t.Data[i] = fv.(float64)
			case float32:
				t.Data[i] = float64(fv.(float32))
			case int:
				t.Data[i] = float64(fv.(int))
			case *Arrf:
				t.Data[i] = fv.(*Arrf).Data[i]
			default:
				panic(TYPE_ERROR)
			}
		}
	}
	return t
}
