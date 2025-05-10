package handlers

type NilSetter[T any] interface {
	SetTo(T)
	SetToNull()
}

func PtrToApiNil[T any](ptr *T, setter NilSetter[T]) {
	if ptr == nil {
		setter.SetToNull()
		return
	}
	setter.SetTo(*ptr)
}

type Getter[T any] interface {
	Get() (T, bool)
}

func ApiValueToPtr[T any](getter Getter[T]) *T {
	if val, ok := getter.Get(); ok {
		return &val
	}
	return nil
}
