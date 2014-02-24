package world

import "testing"
import "math/rand"
import "github.com/tobz/phosphorus/utils"
import "github.com/stretchr/testify/assert"

type TestObject struct {
	pos utils.Point3D
	id  uint32
}

func (t *TestObject) Position() utils.Point3D {
	return t.pos
}

func (t *TestObject) ObjectID() uint32 {
	return t.id
}

func TestCreateOctree(t *testing.T) {
	tree := NewOctree(2048, 2048, 2048)
	assert.NotNil(t, tree)
}

func TestAddingObject(t *testing.T) {
	tree := NewOctree(2048, 2048, 2048)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{512, 512, 512}, 666}
	err := tree.AddObject(obj)
	assert.Nil(t, err)
}

func TestAddingOutOfBoundsObject(t *testing.T) {
	tree := NewOctree(2048, 2048, 2048)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{-12, -52, 51}, 666}
	err := tree.AddObject(obj)
	assert.NotNil(t, err)
}

func TestRemovingObject(t *testing.T) {
	tree := NewOctree(2048, 2048, 2048)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{512, 512, 512}, 666}
	err := tree.RemoveObject(obj)
	assert.Nil(t, err)
}

func TestRemovingOutOfBoundsObject(t *testing.T) {
	tree := NewOctree(2048, 2048, 2048)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{-12, -52, 51}, 666}
	err := tree.RemoveObject(obj)
	assert.NotNil(t, err)
}

func TestMovingObject(t *testing.T) {
	tree := NewOctree(2048, 2048, 2048)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{512, 512, 512}, 666}
	err := tree.MoveObject(obj)
	assert.Nil(t, err)
}

func TestMovingOutOfBoundsObject(t *testing.T) {
	tree := NewOctree(2048, 2048, 2048)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{-12, -52, 51}, 666}
	err := tree.MoveObject(obj)
	assert.NotNil(t, err)
}


func TestRetrievingInRangeObject(t *testing.T) {
	tree := NewOctree(2048, 2048, 2048)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{512, 512, 512}, 666}
	err := tree.AddObject(obj)
	assert.Nil(t, err)

	objs := tree.GetObjectsInRadius(utils.Point3D{500, 500, 500}, 256)
	assert.NotEmpty(t, objs)
	assert.Equal(t, objs[0].ObjectID(), uint32(666))
}


func TestRetrievingOutOfRangeObject(t *testing.T) {
	tree := NewOctree(2048, 2048, 2048)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{768, 768, 768}, 666}
	err := tree.AddObject(obj)
	assert.Nil(t, err)

	objs := tree.GetObjectsInRadius(utils.Point3D{256, 256, 256}, 256)
	assert.Empty(t, objs)
}

func BenchmarkGetObjectsInRadiusSmallNoObjects(b *testing.B) {
	tree := NewOctree(524288, 524288, 8192)

	// Pick a random point to search at.
	searchPoint := getRandomPoint(524288, 524288, 8192)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GetObjectsInRadius(searchPoint, 300)
	}
}

func BenchmarkGetObjectsInRadiusMediumNoObjects(b *testing.B) {
	tree := NewOctree(524288, 524288, 8192)

	// Pick a random point to search at.
	searchPoint := getRandomPoint(524288, 524288, 8192)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GetObjectsInRadius(searchPoint, 3600)
	}
}

func BenchmarkGetObjectsInRadiusLargeNoObjects(b *testing.B) {
	tree := NewOctree(524288, 524288, 8192)

	// Pick a random point to search at.
	searchPoint := getRandomPoint(524288, 524288, 8192)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GetObjectsInRadius(searchPoint, 65536)
	}
}

func BenchmarkGetObjectsInRadiusSmall100Objects(b *testing.B) {
	tree := NewOctree(524288, 524288, 8192)

	// Seed our tree with 100 objects.
	for j := 0; j < 100; j++ {
		tree.AddObject(&TestObject{getRandomPoint(524288, 524288, 8192), uint32(j)})
	}

	// Pick a random point to search at.
	searchPoint := getRandomPoint(524288, 524288, 8192)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GetObjectsInRadius(searchPoint, 300)
	}
}

func BenchmarkGetObjectsInRadiusMedium100Objects(b *testing.B) {
	tree := NewOctree(524288, 524288, 8192)

	// Seed our tree with 100 objects.
	for j := 0; j < 100; j++ {
		tree.AddObject(&TestObject{getRandomPoint(524288, 524288, 8192), uint32(j)})
	}

	// Pick a random point to search at.
	searchPoint := getRandomPoint(524288, 524288, 8192)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GetObjectsInRadius(searchPoint, 3600)
	}
}

func BenchmarkGetObjectsInRadiusLarge100Objects(b *testing.B) {
	tree := NewOctree(524288, 524288, 8192)

	// Seed our tree with 100 objects.
	for j := 0; j < 100; j++ {
		tree.AddObject(&TestObject{getRandomPoint(524288, 524288, 8192), uint32(j)})
	}

	// Pick a random point to search at.
	searchPoint := getRandomPoint(524288, 524288, 8192)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GetObjectsInRadius(searchPoint, 65536)
	}
}

func getRandomPoint(maxHeight, maxWidth, maxDepth int64) utils.Point3D {
	randHeight := rand.Int63n(maxHeight)
	randWidth := rand.Int63n(maxWidth)
	randDepth := rand.Int63n(maxDepth)

	return utils.Point3D{randHeight, randWidth, randDepth}
}
