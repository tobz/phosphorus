package interfaces

import "github.com/tobz/phosphorus/utils"

type WorldObject interface {
    Position() utils.Point3D
    ObjectID() uint32
}
