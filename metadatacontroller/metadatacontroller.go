package metadatacontroller

import (
	"github.com/clawio/entities"
)

// MetaDataController is an interface to perform metadata operations.
type MetaDataController interface {
	Init(user entities.User) error
	ExamineObject(user entities.User, pathSpec string) (entities.ObjectInfo, error)
	ListTree(user entities.User, pathSpec string) ([]entities.ObjectInfo, error)
}
