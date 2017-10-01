package measure

type SizeInfo struct {
	Name  string
	MIMax int
	MJMax int
	MKMax int
}

var (
	Sizes = map[string]SizeInfo{
		"SSMALL": {"SSMALL", 33, 33, 65},
		"SMALL":  {"SMALL", 65, 65, 129},
		"MIDDLE": {"MIDDLE", 129, 129, 257},
		"LARGE":  {"LARGE", 257, 257, 513},
		"ELARGE": {"ELARGE", 513, 513, 1025},
	}
)

func FFlop(s string) float64 {
	return float64(Sizes[s].MKMax-3) * float64(Sizes[s].MJMax-3) * float64(Sizes[s].MIMax-3) * 34.0
}

func MFlops(nn int, cpu float64, size string) float64 {
	return FFlop(size) / cpu * 1.e-6 * float64(nn)
}
