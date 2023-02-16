package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/lkeix/jackall/testdata/application/entity"
	"github.com/lkeix/jackall/testdata/application/usecase"
)

func main() {
	hogesan := entity.NewHoge("hoge")
	fmt.Println(hogesan)

	r := echo.New()

	joinner := usecase.NewHogeNameJoineer()

	r.GET("/join", func(c echo.Context) error {
		h1 := entity.NewHoge("hoge")
		h2 := entity.NewHoge("hogehoge")
		joinner.Join(h1, h2)
		return nil
	})

	r.Start(":8080")
}
