package world

import "fmt"
import "sync"
import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/utils"

type Octree struct {
	bl       utils.Point3D
	tr       utils.Point3D
	mlock    *sync.RWMutex
	children [8]*Octree
	objects  map[uint32]interfaces.WorldObject
}

func NewOctree(height, width, depth int64) *Octree {
	bl := utils.Point3D{0, 0, 0}
	tr := utils.Point3D{height, width, depth}

	tree := &Octree{bl: bl, tr: tr, mlock: &sync.RWMutex{}, objects: make(map[uint32]interfaces.WorldObject)}
	tree.Subdivide(2)

	return tree
}

func (o *Octree) Subdivide(divideDepth int) {
	// Lock ourselves so we can subdivide.
	o.mlock.Lock()
	defer o.mlock.Unlock()

	o.subdivideImpl(divideDepth)
}

func (o *Octree) subdivideImpl(divideDepth int) {
	// Do the subdivision.
	o.children[0] = &Octree{bl: utils.Point3D{o.bl.X, o.bl.Y, o.bl.Z}, tr: utils.Point3D{o.tr.X / 2, o.tr.Y / 2, o.tr.Z / 2}}
	o.children[1] = &Octree{bl: utils.Point3D{o.tr.X / 2, o.bl.Y, o.bl.Z}, tr: utils.Point3D{o.tr.X, o.tr.Y / 2, o.tr.Z / 2}}
	o.children[2] = &Octree{bl: utils.Point3D{o.bl.X, o.tr.Y / 2, o.bl.Z}, tr: utils.Point3D{o.tr.X / 2, o.tr.Y, o.tr.Z / 2}}
	o.children[3] = &Octree{bl: utils.Point3D{o.tr.X / 2, o.tr.Y / 2, o.bl.Z}, tr: utils.Point3D{o.tr.X, o.tr.Y, o.tr.Z / 2}}
	o.children[4] = &Octree{bl: utils.Point3D{o.bl.X, o.bl.Y, o.tr.Z / 2}, tr: utils.Point3D{o.tr.X / 2, o.tr.Y / 2, o.tr.Z}}
	o.children[5] = &Octree{bl: utils.Point3D{o.tr.X / 2, o.bl.Y, o.tr.Z / 2}, tr: utils.Point3D{o.tr.X, o.tr.Y / 2, o.tr.Z}}
	o.children[6] = &Octree{bl: utils.Point3D{o.bl.X, o.tr.Y / 2, o.tr.Z / 2}, tr: utils.Point3D{o.tr.X / 2, o.tr.Y, o.tr.Z}}
	o.children[7] = &Octree{bl: utils.Point3D{o.tr.X / 2, o.tr.Y / 2, o.tr.Z / 2}, tr: utils.Point3D{o.tr.X, o.tr.Y, o.tr.Z}}

	// See if we need to keep dividing.
	for _, child := range o.children {
		if divideDepth > 1 {
			child.subdivideImpl(divideDepth - 1)
		} else {
			// We're not dividing, so make sure these children can hold objects.
			child.objects = make(map[uint32]interfaces.WorldObject)
		}
	}

	// Now go through any children we have and move them into their new homes.
	for _, obj := range o.objects {
		o.addObjectImpl(obj)
	}

	// Clear ourselves out.
	for k := range o.objects {
		delete(o.objects, k)
	}

	o.objects = nil
}

func (o *Octree) contains(obj interfaces.WorldObject) bool {
	p := obj.Position()
	return (o.bl.X <= p.X && p.X <= o.tr.X) && (o.bl.Y <= p.Y && p.Y <= o.tr.Y) && (o.bl.Z <= p.Z && p.Z <= o.tr.Z)
}

func (o *Octree) containsRadius(p utils.Point3D, radius int64) bool {
	distSquared := square(radius)

	if p.X < o.bl.X {
		distSquared -= square(p.X - o.bl.X)
	}

	if p.Y < o.bl.Y {
		distSquared -= square(p.Y - o.bl.Y)
	}

	if p.Z < o.bl.Z {
		distSquared -= square(p.Z - o.bl.Z)
	}

	return distSquared > 0
}

func (o *Octree) AddObject(obj interfaces.WorldObject) error {
	if !o.contains(obj) {
		return fmt.Errorf("object position is not within this tree")
	}

	// Lock everything down while we add the object.
	o.mlock.Lock()
	defer o.mlock.Unlock()

	// Call our unsafe implementation without locking.
	o.addObjectImpl(obj)

	return nil
}

func (o *Octree) addObjectImpl(obj interfaces.WorldObject) {
	// See if this fits in any of our children, if we have any.
	for _, child := range o.children {
		if child != nil && child.contains(obj) {
			child.addObjectImpl(obj)
			return
		}
	}

	if o.objects != nil {
		o.objects[obj.ObjectID()] = obj
		return
	}
}

func (o *Octree) RemoveObject(obj interfaces.WorldObject) error {
	if !o.contains(obj) {
		return fmt.Errorf("object position is not within this tree")
	}

	// Lock everything down while we remove the object.
	o.mlock.Lock()
	defer o.mlock.Unlock()

	// Call our unsafe implementation without locking.
	o.removeObjectImpl(obj)

	return nil
}

func (o *Octree) removeObjectImpl(obj interfaces.WorldObject) {
	// See if this object is in any of our children nodes, if we have any.
	for _, child := range o.children {
		if child != nil && child.contains(obj) {
			child.removeObjectImpl(obj)
			return
		}
	}

	delete(o.objects, obj.ObjectID())
}

func (o *Octree) MoveObject(obj interfaces.WorldObject) error {
	if !o.contains(obj) {
		return fmt.Errorf("object position is not within this tree")
	}

	// We have to lock early to ensure we keep the tree consistent and don't have
	// search results where the object is missing as we remove and add it back.
	o.mlock.Lock()
	defer o.mlock.Unlock()

	o.removeObjectImpl(obj)
	o.addObjectImpl(obj)

	return nil
}

func (o *Octree) GetObjectsInRadius(p utils.Point3D, radius int64) []interfaces.WorldObject {
	o.mlock.RLock()
	defer o.mlock.RUnlock()

	// Create our holder for objects we find.
	objects := make([]interfaces.WorldObject, 0, 16)
	objects = o.getObjectsInRadiusImpl(p, radius, objects)

	return objects
}

func (o *Octree) getObjectsInRadiusImpl(p utils.Point3D, radius int64, objects []interfaces.WorldObject) []interfaces.WorldObject {
	// Find out which of our children intersect with this sphere and rope them into the check.
	for _, child := range o.children {
		if child != nil && child.containsRadius(p, radius) {
			objects = child.getObjectsInRadiusImpl(p, radius, objects)
		}
	}

	// If we're a leaf node, we need to check our own children.
	if o.objects != nil {
		radiusSquared := square(radius)
		for _, obj := range o.objects {
			p0 := obj.Position()
			if (square(p.X-p0.X) + square(p.Y-p0.Y) + square(p.Z-p0.Z)) <= radiusSquared {
				objects = append(objects, obj)
			}
		}
	}

	return objects
}

func square(i int64) int64 {
	return i * i
}
