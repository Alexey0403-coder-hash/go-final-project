package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const format = "20060102"

// NextDate реализует расчет следующей даты повторения задачи
func NextDate(now time.Time, start string, repeat string) (string, error) {
	if start == "" {
		return "", errors.New("ошибка: строка первоначальной даты задачи пустая")
	}

	startTime, err := time.Parse(format, start)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга первоначальной даты задачи: %w", err)
	}

	if repeat == "" {
		return "", errors.New("ошибка: строка повторения пустая")
	}

	repeatSlice := strings.Split(repeat, " ")

	switch repeatSlice[0] {
	case "d":
		if len(repeatSlice) != 2 {
			return "", errors.New("ошибка: некорректный формат повторения в днях")
		}
		dayInt, err := strconv.Atoi(repeatSlice[1])
		if err != nil {
			return "", fmt.Errorf("ошибка преобразования типов дней: %w", err)
		}

		if dayInt <= 0 {
			return "", errors.New("ошибка: количество дней повторения не положительное")
		}

		if dayInt > 400 {
			return "", errors.New("ошибка: недопустимое максимальное количество дней для повторения")
		}

		for {
			startTime = startTime.AddDate(0, 0, dayInt)
			if afterNow(startTime, now) {
				repeatDateString := startTime.Format(format)
				return repeatDateString, nil
			}
		}
	case "y":
		if len(repeatSlice) != 1 {
			return "", errors.New("ошибка: некорректный формат повторения в годах")
		}
		for {
			startTime = startTime.AddDate(1, 0, 0)
			if afterNow(startTime, now) {
				repeatDateString := startTime.Format(format)
				return repeatDateString, nil
			}
		}
	default:
		return "", errors.New("ошибка: недопустимый символ для повторения")
	}
}

// afterNow реализует сравнение сегодняшней даты и следующей даты выполнения
func afterNow(date time.Time, now time.Time) bool {
	now = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	date = time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	return date.After(now)
}

// nextDayHandler выступает в роли обработчика и релизует на основании полученных параметров запроса ответ в виде следующей даты типа string
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	var timeNow time.Time

	now := r.URL.Query().Get("now")
	if now == "" {
		timeNow = time.Now()
	} else {
		var err error
		timeNow, err = time.Parse(format, now)
		if err != nil {
			http.Error(w, "ошибка парсинга начальной даты", http.StatusBadRequest)
			return
		}
	}

	date := r.URL.Query().Get("date")

	repeat := r.URL.Query().Get("repeat")

	repeatString, err := NextDate(timeNow, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, repeatString)
}
