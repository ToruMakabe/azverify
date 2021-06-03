package cmd

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDesiredResources(t *testing.T) {
	t.Run("get a single resource", func(t *testing.T) {
		f := "../testdata/desired/desired_1.json"
		r, err := getDesiredResources(f)
		require.Nil(t, err)

		t.Logf("resource: %v", r)
	})

	t.Run("get multiple resources", func(t *testing.T) {
		f := "../testdata/desired/desired_array.json"
		r, err := getDesiredResources(f)
		require.Nil(t, err)

		t.Logf("resource: %v", r)
	})

	t.Run("get invalid json (not array)", func(t *testing.T) {
		f := "../testdata/desired/invalid_not_array.json"
		_, err := getDesiredResources(f)
		require.NotNil(t, err)
		t.Log(err)
	})

	t.Run("get invalid json (format)", func(t *testing.T) {
		f := "../testdata/desired/invalid_format.json"
		_, err := getDesiredResources(f)
		require.NotNil(t, err)
		t.Log(err)
	})
}

func TestDiff(t *testing.T) {
	t.Run("match between actual and desired", func(t *testing.T) {
		af, err := os.Open("../testdata/actual/superset_match_1.json")
		require.Nil(t, err)
		defer af.Close()
		a, err := os.ReadFile(af.Name())
		require.Nil(t, err)

		df, err := os.Open("../testdata/desired/desired_1.json")
		require.Nil(t, err)
		defer df.Close()
		d, err := os.ReadFile(df.Name())
		require.Nil(t, err)

		var j []interface{}
		err = json.Unmarshal(d, &j)
		require.Nil(t, err)

		for _, v := range j {
			vv, err := json.Marshal(v)
			require.Nil(t, err)
			diff, diffStr := diff(a, vv)
			assert.Equal(t, diff, "SupersetMatch")
			t.Logf("details: %s", diffStr)
		}
	})

	t.Run("unmatch between actual and desired", func(t *testing.T) {
		af, err := os.Open("../testdata/actual/unmatch_1.json")
		require.Nil(t, err)
		defer af.Close()
		a, err := os.ReadFile(af.Name())
		require.Nil(t, err)

		df, err := os.Open("../testdata/desired/desired_1.json")
		require.Nil(t, err)
		defer df.Close()
		d, err := os.ReadFile(df.Name())
		require.Nil(t, err)

		var j []interface{}
		err = json.Unmarshal(d, &j)
		require.Nil(t, err)

		for _, v := range j {
			vv, err := json.Marshal(v)
			require.Nil(t, err)
			diff, diffStr := diff(a, vv)
			assert.Equal(t, diff, "NoMatch")
			t.Logf("details: %s", diffStr)
		}
	})

	t.Run("match between actual and desired (input multi resources)", func(t *testing.T) {
		af1, err := os.Open("../testdata/actual/superset_match_1.json")
		require.Nil(t, err)
		defer af1.Close()
		a1, err := os.ReadFile(af1.Name())
		require.Nil(t, err)
		var aa [][]byte
		aa = append(aa, a1)

		af2, err := os.Open("../testdata/actual/superset_match_2.json")
		require.Nil(t, err)
		defer af1.Close()
		a2, err := os.ReadFile(af2.Name())
		require.Nil(t, err)
		aa = append(aa, a2)

		df, err := os.Open("../testdata/desired/desired_array.json")
		require.Nil(t, err)
		defer df.Close()
		d, err := os.ReadFile(df.Name())
		require.Nil(t, err)

		var j []interface{}
		err = json.Unmarshal(d, &j)
		require.Nil(t, err)

		for k, v := range j {
			vv, err := json.Marshal(v)
			require.Nil(t, err)
			diff, diffStr := diff(aa[k], vv)
			assert.Equal(t, diff, "SupersetMatch")
			t.Logf("details: %s", diffStr)
		}
	})
}
