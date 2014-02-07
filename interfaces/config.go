package interfaces

type Config interface {
	GetAsInteger(string) (int64, error)
	GetAsFloat(string) (float64, error)
	GetAsString(string) (string, error)

	GetAsManyIntegers(string) ([]int64, error)
	GetAsManyFloats(string) ([]float64, error)
	GetAsManyStrings(string) ([]string, error)
}
