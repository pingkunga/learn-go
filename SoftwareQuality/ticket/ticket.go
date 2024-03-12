package ticket

func Price(age int) float64 {
	if age <= 10 {
		return 0
	} else if age < 18 || age > 60 {
		return 0.5
	} else {
		return 1.5
	}
}
