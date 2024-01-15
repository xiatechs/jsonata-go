package jsonata

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSorting(t *testing.T) {
	t.Run("normal nested structure", func(t *testing.T) {
		myMap := map[string]interface{}{
			"banana": 3,
			"apple":  1,
			"nested": map[string]interface{}{
				"orange": 2,
				"grape":  4,
			},
		}

		sortedMap := makeDeterministic(myMap, nil)

		jsonBytes, err := json.Marshal(sortedMap)
		assert.NoError(t, err)
		assert.Equal(t, "{\"apple\":1,\"banana\":3,\"nested\":{\"grape\":4,\"orange\":2}}", string(jsonBytes))
	})

	t.Run("circular dependency - caught and handled", func(t *testing.T) {
		myMap := map[string]interface{}{
			"banana": 3,
			"apple":  1,
			"nested": map[string]interface{}{
				"orange": 2,
				"grape":  4,
			},
		}

		// introduce circular reference
		myMap["nested"].(map[string]interface{})["circular"] = myMap

		sortedMap := makeDeterministic(myMap, nil)

		jsonBytes, err := json.Marshal(sortedMap)
		assert.NoError(t, err)
		assert.Equal(t, "{\"apple\":1,\"banana\":3,\"nested\":{\"grape\":4,\"orange\":2}}", string(jsonBytes))
	})

	t.Run("circular dependency - caught and handled", func(t *testing.T) {
		myMap := map[string]interface{}{
			"banana": 3,
			"apple":  1,
			"array": []interface{}{
				map[string]interface{}{"nested1": map[string]interface{}{
					"orange": 2,
					"grape":  4,
				}},	
				map[string]interface{}{"nested2": map[string]interface{}{
						"orange": 2,
						"grape":  4,
					}},

			},
		}

		// introduce circular reference
		myMap["array"].([]interface{})[0].(map[string]interface{})["circular"] = myMap

		sortedMap := makeDeterministic(myMap, nil)

		jsonBytes, err := json.Marshal(sortedMap)
		assert.NoError(t, err)
		assert.Equal(t, "{\"apple\":1,\"array\":[{\"nested1\":{\"grape\":4,\"orange\":2}},{\"nested2\":{\"grape\":4,\"orange\":2}}],\"banana\":3}",string(jsonBytes))
	})
}
