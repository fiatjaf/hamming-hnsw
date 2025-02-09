package hnsw

import (
	"cmp"
	crand "crypto/rand"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_maxLevel(t *testing.T) {
	var m int

	m = maxLevel(0.5, 10)
	require.Equal(t, 4, m)

	m = maxLevel(0.5, 1000)
	require.Equal(t, 11, m)
}

func newTestGraph[K cmp.Ordered]() *Graph[K] {
	return &Graph[K]{
		M:        6,
		Ml:       0.5,
		EfSearch: 20,
		Rng:      rand.New(rand.NewSource(0)),
	}
}

func TestGraph_AddSearch(t *testing.T) {
	t.Parallel()

	g := newTestGraph[int]()

	for i := 0; i < 128; i++ {
		g.Add(
			Node[int]{
				Key:   i,
				Value: BinaryString{byte(i * 212), byte(i * i * i * 177), byte(i * i * 37), byte((1 + i) * 134), byte(i), byte(256 - ((i * 48) % 256))},
			},
		)
	}

	al := Analyzer[int]{Graph: g}

	// Layers should be approximately log2(128) = 7
	// Look for an approximate doubling of the number of nodes in each layer.
	require.Equal(t, []int{
		128,
		67,
		28,
		12,
		6,
		2,
		1,
		1,
	}, al.Topography())

	nearest := g.Search(
		BinaryString{16, 72, 244, 8, 53, 18},
		4,
	)

	require.Len(t, nearest, 4)
	require.Equal(t, nearest[0].Key, 49)
	require.Equal(t, nearest[1].Key, 33)
	require.Equal(t, nearest[2].Key, 17)
	require.Equal(t, nearest[3].Key, 113)
}

func TestGraph_AddDelete(t *testing.T) {
	t.Parallel()

	g := newTestGraph[int]()
	for i := 0; i < 128; i++ {
		g.Add(Node[int]{
			Key:   i,
			Value: BinaryString{byte(i)},
		})
	}

	require.Equal(t, 128, g.Len())
	an := Analyzer[int]{Graph: g}

	preDeleteConnectivity := an.Connectivity()

	// Delete every even node.
	for i := 0; i < 128; i += 2 {
		ok := g.Delete(i)
		require.True(t, ok)
	}

	require.Equal(t, 64, g.Len())

	postDeleteConnectivity := an.Connectivity()

	// Connectivity should be the same for the lowest layer.
	require.Equal(
		t, preDeleteConnectivity[0],
		postDeleteConnectivity[0],
	)

	t.Run("DeleteNotFound", func(t *testing.T) {
		ok := g.Delete(-1)
		require.False(t, ok)
	})
}

func Benchmark_HSNW(b *testing.B) {
	b.ReportAllocs()

	sizes := []int{100, 1000, 10000}

	// Use this to ensure that complexity is O(log n) where n = h.Len().
	for _, size := range sizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			g := Graph[int]{}
			g.Ml = 0.5
			for i := 0; i < size; i++ {
				g.Add(Node[int]{
					Key:   i,
					Value: BinaryString{byte(i)},
				})
			}
			b.ResetTimer()

			b.Run("Search", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					g.Search(
						BinaryString{byte(i % size)},
						4,
					)
				}
			})
		})
	}
}

func randBytes(n int) BinaryString {
	x := make(BinaryString, n)
	crand.Read(x)
	return x
}

func Benchmark_HNSW_1536(b *testing.B) {
	b.ReportAllocs()

	g := newTestGraph[int]()
	const size = 1000
	points := make([]Node[int], size)
	for i := 0; i < size; i++ {
		points[i] = Node[int]{
			Key:   i,
			Value: randBytes(1536),
		}
		g.Add(points[i])
	}
	b.ResetTimer()

	b.Run("Search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g.Search(
				points[i%size].Value,
				4,
			)
		}
	})
}

func TestGraph_DefaultCosine(t *testing.T) {
	g := NewGraph[int]()
	g.Add(
		Node[int]{Key: 1, Value: BinaryString{1, 1}},
		Node[int]{Key: 2, Value: BinaryString{0, 1}},
		Node[int]{Key: 3, Value: BinaryString{1, 0}},
	)

	neighbors := g.Search(
		BinaryString{1, 1},
		1,
	)

	require.Equal(
		t,
		[]Node[int]{
			{1, BinaryString{1, 1}},
		},
		neighbors,
	)
}
