package form

import "time"

type Test struct {
	Name1 string `json:"name_1" form:"name_1" binding:"eq=10" `
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn,bookabledate" time_format:"2006-01-02"`
}

var ruleMaps = make(map[string]string)

//https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme
func getRuleMaps() map[string]string {
	ruleMaps["="] = "eq"
	ruleMaps[">"] = "gt"
	ruleMaps["<"] = "lt"
	ruleMaps[">="] = "gte"
	ruleMaps["<="] = "lte"
    ruleMaps["alpha"] = "alpha"
	ruleMaps["alphaNum"] = "alphanum"
	ruleMaps["boolean"] = "boolean"
	ruleMaps["lower"] = "lowercase"
	ruleMaps["upper"] = "uppercase"
	ruleMaps["integer"] = "numeric"
	ruleMaps["number"] = "number"
	ruleMaps["date"] = "datetime"
	ruleMaps["email"] = "email"
	ruleMaps["file"] = "file"
	ruleMaps["ip"] = "ip"
	ruleMaps["macAddr"] = "mac"
	ruleMaps["length"] = "len"
	ruleMaps["max"] = "max"
	ruleMaps["min"] = "min"
	ruleMaps["require"] = "required"
	ruleMaps["requireIf"] = "required_if"
	ruleMaps["requireWith"] = "required_with"
	ruleMaps["requireWithout"] = "required_without"
	ruleMaps["unique"] = "unique"
	ruleMaps["url"] = "url"
	ruleMaps["zip"] = "postcode_iso3166_alpha2"
	ruleMaps["regex"] = "regex"

	return ruleMaps
}
