package model

type Tag struct {
	*Object
}

func NewTag(obj *Object) *Tag {
	return &Tag{Object: obj}
}
