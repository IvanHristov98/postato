package cluster

type FuzzyNum interface {
	MembershipDegree() float64
}

type gaussianFuzzyNum struct {
	center float64
	width  float64
}

func NewGaussianFuzzyNum(center float64, width float64) FuzzyNum {
	return &gaussianFuzzyNum{
		center: center,
		width:  width,
	}
}

// TODO: Implement.
func (gfn *gaussianFuzzyNum) MembershipDegree() float64 {
	return 0.0
}
