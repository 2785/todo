package repositories

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlatFile(t *testing.T) {
	require := require.New(t)
	tmpDir, err := ioutil.TempDir(".", "flatfile")
	require.NoError(err)
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	r, err := NewFlatFile(tmpDir)
	require.NoError(err)

	ctx := context.Background()
	todo1 := Todo{
		Description: "test todo 1",
		Weight:      5,
		Done:        false,
	}

	added, err := r.AddTodo(ctx, todo1)
	require.NoError(err)
	require.NotEmpty(added.CreatedAt)
	require.NotEmpty(added.ID)

	all, err := r.GetAll(ctx)
	require.NoError(err)
	require.Equal(1, len(all))

	todo2 := todo1
	todo2.Weight = 10

	another, err := r.AddTodo(ctx, todo2)
	require.NoError(err)

	all, err = r.GetAll(ctx)
	require.NoError(err)
	require.Equal(2, len(all))

	// get by weight
	list, err := r.ListTodo(ctx, "weight", true, false)
	require.NoError(err)
	require.Equal(2, len(list))
	require.Equal(10, list[0].Weight)

	_, err = r.UpdateTodoStatus(ctx, another, true)
	require.NoError(err)

	list, err = r.ListTodo(ctx, "weight", true, false)
	require.NoError(err)
	require.Equal(1, len(list))
	require.Equal(5, list[0].Weight)
}
