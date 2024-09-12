package util

const (
	Unauthorized    = "Пользователь не существует или некорректен."
	Forbidden       = "Недостаточно прав для выполнения действия."
	NotFound        = "Тендер или предложение не найдено."
	VersionNotFound = "Версия не найдена."
)

type MalformedRequestError struct {
	Status int
	Msg    string
}

func (mr *MalformedRequestError) Error() string {
	return mr.Msg
}

type MyResponseError struct {
	Status int
	Msg    string
}

func (er MyResponseError) Error() string {
	return er.Msg
}