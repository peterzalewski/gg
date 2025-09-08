package model

type Blob struct {
	*Object
}

func (b *Blob) ObjectType() string {
	return "blob"
}

func NewBlob(obj *Object) *Blob {
	return &Blob{Object: obj}
}
