package services

type Luhn struct {
	key []byte
}

type LuhnInterface interface {
	IsValid(number int) bool
}

func NewLuhn(key string) CryptInterface {
	return Crypt{
		key: []byte(key),
	}
}

func (l Luhn) IsValid(number int) bool {
	return (number%10+l.checksum(number/10))%10 == 0
}

func (l Luhn) checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
