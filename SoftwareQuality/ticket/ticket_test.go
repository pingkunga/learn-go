package ticket

import "testing"

func TestPrice(t *testing.T) {

	t.Run("Should return 0 for age 0", func(t *testing.T) {
		//Arrange
		want := 0.0
		age := 0

		//Act
		got := Price(age)

		//Assert
		if got != want {
			t.Errorf("Price(%d) = %f; want %f", age, got, want)
		}

	})

	t.Run("Should return 0 for age 1", func(t *testing.T) {
		//Arrange
		want := 0.0
		age := 1

		//Act
		got := Price(age)

		//Assert
		if got != want {
			t.Errorf("Price(%d) = %f; want %f", age, got, want)
		}

	})

	t.Run("Should return 0 for age 5", func(t *testing.T) {
		//Arrange
		want := 0.0
		age := 5

		//Act
		got := Price(age)

		//Assert
		if got != want {
			t.Errorf("Price(%d) = %f; want %f", age, got, want)
		}

	})

	t.Run("Should return 0 for age 10", func(t *testing.T) {
		//Arrange
		want := 0.0
		age := 10

		//Act
		got := Price(age)

		//Assert
		if got != want {
			t.Errorf("Price(%d) = %f; want %f", age, got, want)
		}

	})

	t.Run("Should return 0 for age 11", func(t *testing.T) {
		//Arrange
		want := 0.5
		age := 11

		//Act
		got := Price(age)

		//Assert
		if got != want {
			t.Errorf("Price(%d) = %f; want %f", age, got, want)
		}

	})
}

func TestPriceWithTestable(t *testing.T) {
	testCases := []struct {
		name string
		age  int
		want float64
	}{
		{"T Should return 0 for age 0", 0, 0.0},
		{"T Should return 0 for age 1", 1, 0.0},
		{"T Should return 0 for age 5", 5, 0.0},
		{"T Should return 0 for age 10", 10, 0.0},
		{"T Should return 0 for age 11", 11, 0.5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Price(tc.age)
			if got != tc.want {
				t.Errorf("Price(%d) = %f; want %f", tc.age, got, tc.want)
			}
		})
	}
}
