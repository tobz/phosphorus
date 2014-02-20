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
	tree := NewOctree(0.0, 0.0, 1024, 1024, 1024)
	assert.NotNil(t, tree)
}

func TestAddingObject(t *testing.T) {
	tree := NewOctree(0.0, 0.0, 1024.0, 1024.0, 1024.0)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{512.0, 512.0, 512.0}, 666}
	err := tree.AddObject(obj)
	assert.Nil(t, err)
}

func TestAddingOutOfBoundsObject(t *testing.T) {
	tree := NewOctree(0.0, 0.0, 1024.0, 1024.0, 1024.0)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{5012.0, 5102.0, 5120.0}, 666}
	err := tree.AddObject(obj)
	assert.NotNil(t, err)
}

func TestRetrievingInRangeObject(t *testing.T) {
	tree := NewOctree(0.0, 0.0, 1024.0, 1024.0, 1024.0)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{512.0, 512.0, 512.0}, 666}
	err := tree.AddObject(obj)
	assert.Nil(t, err)

	objs := tree.GetObjectsInRadius(utils.Point3D{500.0, 500.0, 500.0}, 100)
	assert.NotEmpty(t, objs)
	//assert.Equal(t, objs[0].ObjectID(), uint32(666))
}

func TestRetrievingOutOfRangeObject(t *testing.T) {
	tree := NewOctree(0.0, 0.0, 1024.0, 1024.0, 1024.0)
	assert.NotNil(t, tree)

	obj := &TestObject{utils.Point3D{768.0, 768.0, 768.0}, 666}
	err := tree.AddObject(obj)
	assert.Nil(t, err)

	objs := tree.GetObjectsInRadius(utils.Point3D{256.0, 256.0, 256.0}, 256)
	assert.Empty(t, objs)
}

func BenchmarkGetObjectsInRadiusNoObjectsSmall(b *testing.B) {
	tree := NewOctree(0.0, 0.0, 65536.0, 65536.0, 8192.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GetObjectsInRadius(utils.Point3D{24002.0, 13301.0, 875.0}, 300)
	}
}

func BenchmarkGetObjectsInRadiusNoObjectsBig(b *testing.B) {
	tree := NewOctree(0.0, 0.0, 65536.0, 65536.0, 8192.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GetObjectsInRadius(utils.Point3D{24002.0, 13301.0, 875.0}, 3600)
	}
}

func BenchmarkGetObjectsInRadiusSmall(b *testing.B) {
	tree := NewOctree(0.0, 0.0, 65536.0, 65536.0, 8192.0)

	// Seed our tree with 10000 objects.
	for j := 0; j < 10000; j++ {
		tree.AddObject(&TestObject{getRandomPoint(65536.0, 65536.0, 8192.0), uint32(j)})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GetObjectsInRadius(utils.Point3D{24002.0, 13301.0, 875.0}, 300)
	}
}

func BenchmarkGetObjectsInRadiusBig(b *testing.B) {
	tree := NewOctree(0.0, 0.0, 65536.0, 65536.0, 8192.0)

	// Seed our tree with 10000 objects.
	for j := 0; j < 10000; j++ {
		tree.AddObject(&TestObject{getRandomPoint(65536.0, 65536.0, 8192.0), uint32(j)})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GetObjectsInRadius(utils.Point3D{24002.0, 13301.0, 875.0}, 3600)
	}
}

func getRandomPoint(maxHeight, maxWidth, maxDepth float64) utils.Point3D {
	randHeight := rand.Float64() * maxHeight
	randWidth := rand.Float64() * maxWidth
	randDepth := rand.Float64() * maxDepth

	return utils.Point3D{randHeight, randWidth, randDepth}
}
