package log

import (
	"io"
	"os"
	"testing"

	api "github.com/getumen/proglog/api/v1"
	"github.com/golang/protobuf/proto"

	"github.com/stretchr/testify/require"
)

func TestSegment(t *testing.T) {
	dir, err := os.MkdirTemp(os.TempDir(), "segment-test")
	if err != nil {
		require.NoError(t, err)
	}
	defer os.RemoveAll(dir)

	want := &api.Record{
		Value: []byte("hello, world"),
	}

	c := Config{}
	c.Segment.MaxStoreBytes = 1024
	c.Segment.MaxIndexBytes = entWidth * 3

	s, err := newSegment(dir, 16, c)
	require.NoError(t, err)
	require.Equal(t, uint64(16), s.nextOffset)
	require.False(t, s.IsMaxed())

	for i := uint64(0); i < 3; i++ {
		off, err := s.Append(want)
		require.NoError(t, err)
		require.Equal(t, uint64(16+i), off)

		got, err := s.Read(off)
		require.NoError(t, err)
		require.Equal(t, want.Value, got.Value)
	}

	_, err = s.Append(want)
	require.Equal(t, io.EOF, err)

	require.True(t, s.IsMaxed())
	require.NoError(t, s.Close())

	p, err := proto.Marshal(want)
	require.NoError(t, err)

	c.Segment.MaxStoreBytes = uint64(len(p)+lenWidth) * 4
	c.Segment.MaxIndexBytes = 1024

	s, err = newSegment(dir, 16, c)
	require.NoError(t, err)

	require.True(t, s.IsMaxed())
	require.NoError(t, s.Remove())

	s, err = newSegment(dir, 16, c)
	require.NoError(t, err)

	require.False(t, s.IsMaxed())
	require.NoError(t, s.Close())
}