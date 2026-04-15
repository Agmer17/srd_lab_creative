package model

import (
	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/google/uuid"
)

type Category struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"desc"`
}

func MapToCategoryModel(gr sqlcgen.Category) Category {
	return Category{
		Id:          gr.ID,
		Name:        gr.Name,
		Slug:        gr.Slug,
		Description: gr.Description,
	}
}

func MapListToCategoryModel(ls []sqlcgen.Category) []Category {

	tempList := make([]Category, len(ls))

	for i, v := range ls {
		tempList[i] = MapToCategoryModel(v)
	}

	return tempList
}
