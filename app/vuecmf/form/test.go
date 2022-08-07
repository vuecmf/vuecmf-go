package form

type Test struct {
	Name1 string `json:"name_1" form:"name_1" binding:"eq=10" ` // =
}

var ruleMaps = make(map[string]string)

func getRuleMaps() map[string]string {
	ruleMaps["="] = "eq"
	ruleMaps[">"] = "gt"
	ruleMaps["<"] = "lt"
	ruleMaps[">="] = "gte"
	ruleMaps["<="] = "lte"

	ruleMaps["!="] = "ne"

	return ruleMaps
}
