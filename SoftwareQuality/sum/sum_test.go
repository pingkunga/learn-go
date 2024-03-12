package sum

import (
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("should return 3 when 1 and 2", func(t *testing.T) {
		// Arrange
		want := 3
		// Act
		got := Sum(1, 2)
		// Assert
		if got != want {
			t.Errorf("sum(1, 2) = %d; want %d", got, want)
		}
	})

	//Next Case
	t.Run("should return -1 when 1 and -2", func(t *testing.T) {
		// Arrange
		want := -1
		// Act
		got := Sum(1, -2)
		// Assert
		if got != want {
			t.Errorf("sum(1, -2) = %d; want %d", got, want)
		}
	})
}

func TestSumVariadic(t *testing.T) {
	t.Run("should return 3 when 1, 2", func(t *testing.T) {
		// Arrange
		want := 10
		// Act
		got := SumVariadic(1, 2, 3, 4)
		// Assert
		if got != want {
			t.Errorf("sum(1, 2, 3, 4) = %d; want %d", got, want)
		}
	})

	//Next Case
	t.Run("should return 0 when 0", func(t *testing.T) {
		// Arrange
		want := 0
		// Act
		got := SumVariadic(0)
		// Assert
		if got != want {
			t.Errorf("sum(0) = %d; want %d", got, want)
		}
	})

	//Next Case
	t.Run("should return 0 when empty", func(t *testing.T) {
		// Arrange
		want := 0
		// Act
		got := SumVariadic()
		// Assert
		if got != want {
			t.Errorf("sum() = %d; want %d", got, want)
		}
	})
}
