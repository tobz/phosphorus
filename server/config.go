package phosphorus

type ServerConfig struct {
}

func (sc *ServerConfig) GetAsInteger(string) (int64, error) {
	return 0, nil
}

func (sc *ServerConfig) GetAsFloat(string) (float64, error) {
	return 0.0, nil
}

func (sc *ServerConfig) GetAsString(string) (string, error) {
	return "", nil
}

func (sc *ServerConfig) GetAsManyIntegers(string) ([]int64, error) {
	return []int64{}, nil
}

func (sc *ServerConfig) GetAsManyFloats(string) ([]float64, error) {
	return []float64{}, nil
}

func (sc *ServerConfig) GetAsManyStrings(string) ([]string, error) {
	return []string{}, nil
}
