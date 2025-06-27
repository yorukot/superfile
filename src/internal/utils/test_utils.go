package utils

type TestTOMLType struct {
	SampleBool  bool     `toml:"sample_bool"`
	SampleInt   int      `toml:"sample_int"`
	SampleStr   string   `toml:"sample_str"`
	SampleSlice []string `toml:"sample_slice"`
}

type TestTOMLMissingIgnorerType struct {
	SampleBool    bool     `toml:"sample_bool"`
	SampleInt     int      `toml:"sample_int"`
	SampleStr     string   `toml:"sample_str"`
	SampleSlice   []string `toml:"sample_slice"`
	IgnoreMissing bool     `toml:"ignore_missing"`
}

func (t TestTOMLMissingIgnorerType) GetIgnoreMissingFields() bool {
	return t.IgnoreMissing
}

func (t TestTOMLMissingIgnorerType) WithIgnoreMissing(val bool) TestTOMLMissingIgnorerType {
	t.IgnoreMissing = val
	return t
}
