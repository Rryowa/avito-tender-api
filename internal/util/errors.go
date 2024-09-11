package util

const (
	Unathorized     = "Пользователь не существует или некорректен."
	Forbidden       = "Недостаточно прав для выполнения действия."
	NotFound        = "Тендер или предложение не найдено."
	VersionNotFound = "Версия не найдена."
)

type MalformedRequest struct {
	Status int
	Msg    string
}

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}

type MyErrorResponse struct {
	Status int
	Msg    string
}

func (er MyErrorResponse) Error() string {
	return er.Msg
}