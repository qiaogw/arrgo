package arrgo

func (a *Arrb) LogicalAnd(b *Arrb) *Arrb {
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = v && b.Data[i]
	}
	return t
}

func (a *Arrb) LogicalOr(b *Arrb) *Arrb {
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = v || b.Data[i]
	}
	return t
}

func (a *Arrb) LogicalNot() *Arrb {
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = !v
	}
	return t
}

func LogicalAnd(a, b *Arrb) *Arrb {
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = v && b.Data[i]
	}
	return t
}

func LogicalOr(a, b *Arrb) *Arrb {
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = v || b.Data[i]
	}
	return t
}

func LogicalNot(a *Arrb) *Arrb {
	var t = EmptyB(a.Shape...)
	for i, v := range a.Data {
		t.Data[i] = !v
	}
	return t
}
