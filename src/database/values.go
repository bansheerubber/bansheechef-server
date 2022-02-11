package database

func CreateArray(values ...interface{}) []interface{} {
	var result []interface{}
	result = append(result, values...)
	return result
}
