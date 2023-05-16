package validator

func IsValid(number int64) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int64) int64 {
	var luhn int64

	for i := 0; number > 0; i++ {
		cursor := number % 10
		if i%2 == 0 {
			cursor *= 2
			if cursor > 9 {
				cursor = cursor%10 + cursor / 10
			}
		}

		luhn += cursor
		number /= 10
	}

	return luhn % 10
}