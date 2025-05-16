package constants

const (
	ErrServerError = "Ошибка сервера"
	ErrLoadEnv     = "Не удалось загрузить env"

	ErrUserNotFound = "Пользователь не найден"
	ErrInvalidData  = "Данные введены не верно"
	ErrGetData      = "Ошибка в получении данных"

	ErrStatusNotFound     = "Статус не найден"
	ErrGroupNotFound      = "Группа не найдена"
	ErrEvaluationNotFound = "Оценка не найдена"
	ErrStudentNotFound    = "Студент не найден"
	ErrObjectNotFound     = "Предмет не найден"
	ErrObjectNameNotFound = "Имя предмета не найден"

	ErrDeleteEvaluation = "Ошибка удалении оценки"
	ErrSendEvaluation   = "Ошибка в добалвении оценки"

	ErrEntrySystem     = "Пожалуйста, войдите в систему"
	ErrSessionNotFound = "Сессия не найдена"
)

const (
	SuccConnectMongo = "Подключение к MONGODB - успешно"
	SuccConnectRedis = "Подключение к REDIS - успешно"

	SuccDeleteEvaluation = "Оценка была удалена"
	SuccSendEvaluation   = "Оценка отправлена"
)
