// SPDX-License-Identifier: Apache-2.0

package asset_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/perun-network/nerd-op/asset"
)

func TestFileStorage(t *testing.T) {
	t.Run("non-exist", func(t *testing.T) {
		s, err := asset.NewFileStorage("foo")
		require.Error(t, err)
		require.Nil(t, s)
	})

	t.Run("no-dir", func(t *testing.T) {
		s, err := asset.NewFileStorage("testdata/nodir")
		require.Error(t, err)
		require.Nil(t, s)
	})

	t.Run("no-ext", func(t *testing.T) {
		s, err := asset.NewFileStorage("testdata/noext")
		require.NoError(t, err)
		require.NotNil(t, s)
		_, err = s.Get(big.NewInt(123))
		require.Error(t, err)

		for _, id := range []string{
			"1",
			"42",
			"18446744073709551616",
		} {
			t.Log("id:", id)
			data, err := s.Get(bigIntStr(id))
			require.NoError(t, err)
			require.Equal(t, append([]byte(id), '\n'), data)
		}
	})

	t.Run("with-ext", func(t *testing.T) {
		s, err := asset.NewFileStorage("testdata/withext")
		s.SetExtension("a")
		require.NoError(t, err)
		require.NotNil(t, s)
		_, err = s.Get(big.NewInt(123))
		require.Error(t, err)

		for _, id := range []string{
			"1",
			"5",
			"256",
		} {
			t.Log("id:", id)
			idn := bigIntStr(id)
			data, err := s.Get(idn)
			require.NoError(t, err)
			require.Len(t, data, int(idn.Int64()+1))
		}
	})
}

func bigIntStr(s string) *big.Int {
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("invalid number string " + s)
	}
	return x
}
