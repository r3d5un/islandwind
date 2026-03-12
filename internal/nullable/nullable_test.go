package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/r3d5un/islandwind/internal/nullable"
	"github.com/stretchr/testify/assert"
)

func TestNewNullableNull(t *testing.T) {
	null := nullable.NewNullableNull[string]()
	assert.True(t, null.IsNull())
	assert.True(t, null.IsSpecified())
}

func TestNewNullableValue(t *testing.T) {
	null := nullable.NewNullableValue[string](t.Name())
	assert.False(t, null.IsNull())
	str, err := null.Get()
	assert.NoError(t, err)
	assert.Equal(t, t.Name(), str)
	assert.True(t, null.IsSpecified())
}

func TestNewNullableUnspecified(t *testing.T) {
	null := nullable.NewNullableUnspecified[string]()
	assert.False(t, null.IsNull())
	assert.False(t, null.IsSpecified())
	val, err := null.Get()
	assert.Error(t, err)
	assert.Empty(t, val)
}

type Example struct {
	Field nullable.Nullable[string] `json:"field,omitempty"`
}

func TestNullableUnmarshalling(t *testing.T) {
	t.Run("Null", func(t *testing.T) {
		example := Example{}
		err := json.Unmarshal([]byte(`{"field": null}`), &example)
		assert.NoError(t, err)
		assert.True(t, example.Field.IsSpecified())
		assert.True(t, example.Field.IsNull())
		_, err = example.Field.Get()
		assert.ErrorIs(t, err, nullable.ErrNullValue)
	})

	t.Run("Unspecified", func(t *testing.T) {
		example := Example{}
		err := json.Unmarshal([]byte(`{}`), &example)
		assert.NoError(t, err)
		assert.False(t, example.Field.IsSpecified())
		assert.False(t, example.Field.IsNull())
		_, err = example.Field.Get()
		assert.ErrorIs(t, err, nullable.ErrNotSpecified)
	})

	t.Run("Set", func(t *testing.T) {
		example := Example{}
		err := json.Unmarshal([]byte(`{"field": "value"}`), &example)
		assert.NoError(t, err)
		assert.True(t, example.Field.IsSpecified())
		assert.False(t, example.Field.IsNull())
		val, err := example.Field.Get()
		assert.NoError(t, err)
		assert.Equal(t, "value", val)
	})
}

func TestNullableMarshalling(t *testing.T) {
	t.Run("Null", func(t *testing.T) {
		example := Example{}
		example.Field = nullable.NewNullableNull[string]()
		data, err := json.Marshal(example)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"field": null}`, string(data))
	})

	t.Run("Unspecified", func(t *testing.T) {
		example := Example{}
		example.Field = nullable.NewNullableUnspecified[string]()
		data, err := json.Marshal(example)
		assert.NoError(t, err)
		assert.JSONEq(t, `{}`, string(data))
	})

	t.Run("Set", func(t *testing.T) {
		example := Example{}
		example.Field = nullable.NewNullableValue("value")
		data, err := json.Marshal(example)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"field": "value"}`, string(data))
	})
}
