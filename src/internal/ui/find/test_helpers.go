package find

func setupTestModel() Model {
	return GenerateModel(50, 80) //nolint:mnd // test dimensions
}

func setupTestModelWithResults(resultCount int) Model {
	m := setupTestModel()
	m.results = make([]FindResult, resultCount)
	for i := range resultCount {
		m.results[i] = FindResult{
			Path:  "/test/path" + string(rune('0'+i)),
			IsDir: i%2 == 0,
		}
	}
	return m
}
