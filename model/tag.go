package model

type Tag struct {
	*Object
}

func (t *Tag) ObjectType() string {
	return "tag"
}

func NewTag(obj *Object) *Tag {
	return &Tag{Object: obj}
}
