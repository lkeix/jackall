package usecase

import "github.com/lkeix/jackall/testdata/application/entity"

type HogeJoinner interface {
	Join(*entity.Hoge, *entity.Hoge) *entity.Hoge
}

type HogeNameJoinner struct{}

func NewHogeNameJoineer() HogeJoinner {
	return &HogeNameJoinner{}
}

func (h *HogeNameJoinner) Join(h1, h2 *entity.Hoge) *entity.Hoge {
	return entity.NewHoge(h1.Name + h2.Name)
}
