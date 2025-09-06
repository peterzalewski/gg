package model

type Blob struct {
	*Object
}

func NewBlob(obj *Object) *Blob {
	return &Blob{Object: obj}
}
