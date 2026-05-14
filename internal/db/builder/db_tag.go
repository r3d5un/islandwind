package builder

import "reflect"

func ColumnsFrom(v any) []string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	var cols []string
	for field := range t.Fields() {
		if tag, ok := field.Tag.Lookup("db"); ok {
			cols = append(cols, tag)
		}
	}

	return cols
}
