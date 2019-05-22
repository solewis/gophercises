package money

import "fmt"

type USD int64

func (m USD) String() string {
	x := float64(m)
	x /= 100
	return fmt.Sprintf("$%.2f", x)
}

func (m USD) Float64() float64 {
	x := float64(m)
	x = x / 100
	return x
}

func (m USD) Multiply(f float64) USD {
	x := (float64(m) * f) + 0.5
	return USD(x)
}

func ToUSD(f float64) USD {
	return USD((f * 100) + 0.5)
}
