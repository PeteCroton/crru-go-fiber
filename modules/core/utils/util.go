package module

func Sum(a, b int) int {
	return a + b
}

func SumLoop(a, b int) int {
	i := 0
	for i := 0; i < 10; i++ {
		a += b
		i++
	}
	return a + i
}

func TestPointer(a *int) int {
	return *a
}

func sum_array(a ...int) int {
	s := 0
	for _, v := range a {
		s += v
	}
	return s
}
