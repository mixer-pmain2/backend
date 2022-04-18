package consts

const (
	REAS_DEAD = "S008"
)

const (
	VISIT_WORK_WITH_DOCUMENTS = 131072
)

var ArrErrors = map[int]string{
	0: "Успешно",

	20: "Ошибка подключения к базе",
	21: "Ошибка транзации",
	22: "Ошибка сохранения",

	100: "Не смог получить данные учета пациента",
	101: "Ошибка записи посещения умершему пациенту",
	102: "Не смог проверить посещения в выбранный день",
	103: "Посещение уже было в этот день",

	200: "Ошибка записи посещения",
	201: "Ошибка записи СРЦ",
}
