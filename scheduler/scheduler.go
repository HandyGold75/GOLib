package scheduler

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"golang.org/x/term"
)

type (
	Scedule struct {
		// Months of the year: 1 - 12
		Months []int

		// Weeks of the month: 1 - 5
		Weeks []int

		// Days of the week: 0 - 6 (Sunday first day of the week)
		Days []int

		// Hours of the day: 0 - 23
		Hours []int

		// Minutes of the hour: 0 - 59
		Minutes []int
	}
)

func verifySceduleData(scedule *Scedule) error {
	for _, v := range scedule.Months {
		if v < 1 || v > 12 {
			return fmt.Errorf(fmt.Sprintf("Value \"%v\" for month is not valid!\nValid values: <1-12>", v))
		}
	}

	for _, v := range scedule.Weeks {
		if v < 1 || v > 5 {
			return fmt.Errorf(fmt.Sprintf("Value \"%v\" for week is not valid!\nValid values: <1-5>", v))
		}
	}

	for _, v := range scedule.Days {
		if v < 0 || v > 6 {
			return fmt.Errorf(fmt.Sprintf("Value \"%v\" for day is not valid!\nValid values: <0-6>", v))
		}
	}

	for _, v := range scedule.Hours {
		if v < 0 || v > 23 {
			return fmt.Errorf(fmt.Sprintf("Value \"%v\" for hours is not valid!\nValid values: <0-23>", v))
		}
	}

	for _, v := range scedule.Minutes {
		if v < 0 || v > 59 {
			return fmt.Errorf(fmt.Sprintf("Value \"%v\" for minutes is not valid!\nValid values: <0-59>", v))
		}
	}

	return nil
}

func setNextTimeByMonth(currentTime *time.Time, months []int) {
	currentMonth := int(currentTime.Month())

	if slices.Contains(months, currentMonth) {
		return
	}

	offset := 0
	for i := 0; i < 10; i++ {
		for _, possibleTargetMonth := range months {
			if possibleTargetMonth < currentMonth {
				continue
			}

			*currentTime = currentTime.AddDate(0, (possibleTargetMonth+offset)-currentMonth, 0)

			return
		}

		offset += 13 - currentMonth
		currentMonth = 1
	}

	panic("unable to resolve target month!")
}

func setNextTimeByWeek(currentTime *time.Time, weeks []int) {
	currentWeek := (int(currentTime.Day()-1) / 7) + 1

	if slices.Contains(weeks, currentWeek) {
		return
	}

	offset := 0
	for i := 0; i < 2; i++ {
		for _, possibleTargetWeek := range weeks {
			if possibleTargetWeek < currentWeek {
				continue
			}

			*currentTime = currentTime.AddDate(0, 0, ((possibleTargetWeek+offset)-currentWeek)*7)

			return
		}

		offset = 6 - currentWeek
		currentWeek = 1
	}

	panic("unable to resolve target week!")
}

func setNextTimeByDay(currentTime *time.Time, days []int) {
	currentDay := int(currentTime.Weekday())

	if slices.Contains(days, currentDay) {
		return
	}

	offset := 0
	for i := 0; i < 2; i++ {
		for _, possibleTargetDay := range days {
			if possibleTargetDay < currentDay {
				continue
			}

			*currentTime = currentTime.AddDate(0, 0, (possibleTargetDay+offset)-currentDay)

			return
		}

		offset = 8 - currentDay
		currentDay = 1
	}

	panic("unable to resolve target day!")
}

func setNextTimeByHour(currentTime *time.Time, hours []int) {
	currentHour := currentTime.Hour()

	if slices.Contains(hours, currentHour) {
		return
	}

	offset := 0
	for i := 0; i < 2; i++ {
		for _, possibleTargetHour := range hours {
			if possibleTargetHour < currentHour {
				continue
			}

			*currentTime = currentTime.Add(time.Hour * time.Duration((possibleTargetHour+offset)-currentHour))

			return
		}

		offset = 24 - currentHour
		currentHour = 0
	}

	panic("unable to resolve target hour!")
}

func setNextTimeByMinute(currentTime *time.Time, minutes []int) {
	currentMinute := currentTime.Minute()

	if slices.Contains(minutes, currentMinute) {
		return
	}

	offset := 0
	for i := 0; i < 2; i++ {
		for _, possibleTargetMinute := range minutes {
			if possibleTargetMinute < currentMinute {
				continue
			}

			*currentTime = currentTime.Add(time.Minute * time.Duration((possibleTargetMinute+offset)-currentMinute))

			return
		}

		offset = 60 - currentMinute
		currentMinute = 0
	}

	panic("unable to resolve target minute!")
}

func SetNextTime(currentTime *time.Time, scedule *Scedule) error {
	if err := verifySceduleData(scedule); err != nil {
		return err
	}

	*currentTime = currentTime.Add(-time.Second*time.Duration(currentTime.Second()) + -time.Nanosecond*time.Duration(currentTime.Nanosecond()))

	setNextTimeByMinute(currentTime, scedule.Minutes)
	setNextTimeByHour(currentTime, scedule.Hours)
	setNextTimeByDay(currentTime, scedule.Days)
	setNextTimeByWeek(currentTime, scedule.Weeks)
	setNextTimeByMonth(currentTime, scedule.Months)

	return nil
}

func SleepFor(msg string, timeDur time.Duration, updateInterval time.Duration) {
	endTime := time.Now().Add(timeDur)
	for {
		if endTime.Before(time.Now()) {
			break
		}

		untilEndTime := time.Until(endTime)

		if msg != "" {
			width, _, _ := term.GetSize(0)
			fmt.Printf("\r%"+strconv.Itoa(width)+"v\r%v%v", "", msg, untilEndTime.Round(updateInterval).String())
		}

		if updateInterval < untilEndTime {
			time.Sleep(updateInterval)
			continue
		}
		time.Sleep(untilEndTime)
	}

	if msg != "" {
		width, _, _ := term.GetSize(0)
		fmt.Printf("\r%"+strconv.Itoa(width)+"v\r", "")
	}
}
