package invite

import (
	"dj-api/app/models"
)

type Tree struct {
	models.ModelId

	Parent int
	Level  int
	Left   int
	Right  int
	Root   int
}

func (Tree) TableName() string {
	return models.TableInviteTree
}
