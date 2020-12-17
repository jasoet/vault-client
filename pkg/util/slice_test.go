package util_test

import (
	"fmt"
	. "github.com/jasoet/vault-client/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToArrStr(t *testing.T) {
	t.Run("should convert all string elements", func(t *testing.T) {
		input := []interface{}{
			"One",
			"Two",
			"Three",
			"Four",
		}

		expectedOutput := []string{
			"One",
			"Two",
			"Three",
			"Four",
		}

		output := ToArrStr(input)
		assert.NotNil(t, output)

		for i, o := range output {
			assert.Equal(t, expectedOutput[i], o)
		}
	})

	t.Run("should convert  mixed typed elements, will result string representation of the elements", func(t *testing.T) {
		input := []interface{}{
			8,
			4.5,
			true,
			"Four",
			[]int{5, 3, 2},
		}

		expectedOutput := []string{
			"8",
			"4.5",
			"true",
			"Four",
			"[5 3 2]",
		}

		output := ToArrStr(input)
		assert.NotNil(t, output)

		for i, o := range output {
			assert.Equal(t, expectedOutput[i], o)
		}
	})

}

func TestToArrStrPrefixPath(t *testing.T) {
	prefix := "test-prefix"

	t.Run("should convert all string elements with prefix", func(t *testing.T) {
		input := []interface{}{
			"One",
			"Two",
			"Three",
			"Four",
		}

		expectedOutput := []string{
			fmt.Sprintf("%v/%v", prefix, "One"),
			fmt.Sprintf("%v/%v", prefix, "Two"),
			fmt.Sprintf("%v/%v", prefix, "Three"),
			fmt.Sprintf("%v/%v", prefix, "Four"),
		}

		output := ToArrStrPrefixPath(input, prefix)
		assert.NotNil(t, output)

		for i, o := range output {
			assert.Equal(t, expectedOutput[i], o)
		}
	})

	t.Run("should convert mixed typed elements, will result string representation with prefix", func(t *testing.T) {
		input := []interface{}{
			8,
			4.5,
			true,
			"Four",
			[]int{5, 3, 2},
		}

		expectedOutput := []string{
			fmt.Sprintf("%v/%v", prefix, "8"),
			fmt.Sprintf("%v/%v", prefix, "4.5"),
			fmt.Sprintf("%v/%v", prefix, "true"),
			fmt.Sprintf("%v/%v", prefix, "Four"),
			fmt.Sprintf("%v/%v", prefix, "[5 3 2]"),
		}

		output := ToArrStrPrefixPath(input, prefix)
		assert.NotNil(t, output)

		for i, o := range output {
			assert.Equal(t, expectedOutput[i], o)
		}
	})

}

func TestToArrStrPrefix(t *testing.T) {
	prefix := "test-prefix"

	t.Run("should convert all string elements with prefix", func(t *testing.T) {
		input := []interface{}{
			"One",
			"Two",
			"Three",
			"Four",
		}

		expectedOutput := []string{
			fmt.Sprintf("%v%v", prefix, "One"),
			fmt.Sprintf("%v%v", prefix, "Two"),
			fmt.Sprintf("%v%v", prefix, "Three"),
			fmt.Sprintf("%v%v", prefix, "Four"),
		}

		output := ToArrStrPrefix(input, prefix)
		assert.NotNil(t, output)

		for i, o := range output {
			assert.Equal(t, expectedOutput[i], o)
		}
	})

	t.Run("should convert mixed typed elements, will result string representation with prefix", func(t *testing.T) {
		input := []interface{}{
			8,
			4.5,
			true,
			"Four",
			[]int{5, 3, 2},
		}

		expectedOutput := []string{
			fmt.Sprintf("%v%v", prefix, "8"),
			fmt.Sprintf("%v%v", prefix, "4.5"),
			fmt.Sprintf("%v%v", prefix, "true"),
			fmt.Sprintf("%v%v", prefix, "Four"),
			fmt.Sprintf("%v%v", prefix, "[5 3 2]"),
		}

		output := ToArrStrPrefix(input, prefix)
		assert.NotNil(t, output)

		for i, o := range output {
			assert.Equal(t, expectedOutput[i], o)
		}
	})

}
