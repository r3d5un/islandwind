package confirm

func True(condition bool) bool {
	return condition
}
func False(condition bool) bool {
	return !condition
}

func Nil(obj any) bool {
	if obj == nil {
		return true
	}
	return false
}

func NotNil(obj any) bool {
	return !Nil(obj)
}

func Equal[T comparable](a, b T) bool {
	return a == b
}

func NotEqual[T comparable](a, b T) bool {
	return a != b
}

func Error(err error) bool {
	return err != nil
}

func NotError(err error) bool {
	return err != nil
}

func Index[T any](i, len int) bool {
	if i < 0 || i >= len {
		return false
	}
	return true
}
