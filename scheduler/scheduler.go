package scheduler

import (
	"errors"
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

	errScheduler struct{ ErrInvalidMonths, ErrInvalidWeeks, ErrInvalidDays, ErrInvalidHours, ErrInvalidMinutes, ErrResolveMonth, ErrResolveWeek, ErrResolveDay, ErrResolveHour, ErrResolveMinute error }
)

var (
	// All errors that can be raised by scheduler
	ErrScheduler = errScheduler{
		ErrInvalidMonths:  errors.New("invalid value in months, valid values are [1-12]"),
		ErrInvalidWeeks:   errors.New("invalid value in weeks, valid values are [1-5]"),
		ErrInvalidDays:    errors.New("invalid value in days, valid values are [0-6]"),
		ErrInvalidHours:   errors.New("invalid value in hours, valid values are [0-23]"),
		ErrInvalidMinutes: errors.New("invalid value in minutes, valid values are [0-59]"),
		ErrResolveMonth:   errors.New("unable to resolve target month"),
		ErrResolveWeek:    errors.New("unable to resolve target week"),
		ErrResolveDay:     errors.New("unable to resolve target day"),
		ErrResolveHour:    errors.New("unable to resolve target hour"),
		ErrResolveMinute:  errors.New("unable to resolve target minute"),
	}
)

func verifySceduleData(scedule *Scedule) error {
	if i := slices.IndexFunc(scedule.Months, func(v int) bool { return v < 1 || v > 12 }); i != -1 {
		return ErrScheduler.ErrInvalidMonths
	}
	if i := slices.IndexFunc(scedule.Weeks, func(v int) bool { return v < 1 || v > 5 }); i != -1 {
		return ErrScheduler.ErrInvalidWeeks
	}
	if i := slices.IndexFunc(scedule.Days, func(v int) bool { return v < 0 || v > 6 }); i != -1 {
		return ErrScheduler.ErrInvalidDays
	}
	if i := slices.IndexFunc(scedule.Hours, func(v int) bool { return v < 0 || v > 23 }); i != -1 {
		return ErrScheduler.ErrInvalidHours
	}
	if i := slices.IndexFunc(scedule.Minutes, func(v int) bool { return v < 0 || v > 59 }); i != -1 {
		return ErrScheduler.ErrInvalidMinutes
	}
	return nil
}

func setNextTimeByMonth(t *time.Time, months []int) error {
	currentMonth := int(t.Month())
	if slices.Contains(months, currentMonth) {
		return nil
	}

	offset := 0
	for i := 0; i < 10; i++ {
		for _, possibleTargetMonth := range months {
			if possibleTargetMonth < currentMonth {
				continue
			}
			*t = t.AddDate(0, (possibleTargetMonth+offset)-currentMonth, 0)
			return nil
		}

		offset += 13 - currentMonth
		currentMonth = 1
	}

	return ErrScheduler.ErrResolveMonth
}

func setNextTimeByWeek(t *time.Time, weeks []int) error {
	currentWeek := (int(t.Day()-1) / 7) + 1
	if slices.Contains(weeks, currentWeek) {
		return nil
	}

	offset := 0
	for i := 0; i < 2; i++ {
		for _, possibleTargetWeek := range weeks {
			if possibleTargetWeek < currentWeek {
				continue
			}
			*t = t.AddDate(0, 0, ((possibleTargetWeek+offset)-currentWeek)*7)
			return nil
		}

		offset = 6 - currentWeek
		currentWeek = 1
	}

	return ErrScheduler.ErrResolveWeek
}

func setNextTimeByDay(t *time.Time, days []int) error {
	currentDay := int(t.Weekday())
	if slices.Contains(days, currentDay) {
		return nil
	}

	offset := 0
	for i := 0; i < 2; i++ {
		for _, possibleTargetDay := range days {
			if possibleTargetDay < currentDay {
				continue
			}
			*t = t.AddDate(0, 0, (possibleTargetDay+offset)-currentDay)
			return nil
		}

		offset = 8 - currentDay
		currentDay = 1
	}

	return ErrScheduler.ErrResolveDay
}

func setNextTimeByHour(t *time.Time, hours []int) error {
	currentHour := t.Hour()
	if slices.Contains(hours, currentHour) {
		return nil
	}

	offset := 0
	for i := 0; i < 2; i++ {
		for _, possibleTargetHour := range hours {
			if possibleTargetHour < currentHour {
				continue
			}
			*t = t.Add(time.Hour * time.Duration((possibleTargetHour+offset)-currentHour))
			return nil
		}

		offset = 24 - currentHour
		currentHour = 0
	}

	return ErrScheduler.ErrResolveHour
}

func setNextTimeByMinute(t *time.Time, minutes []int) error {
	currentMinute := t.Minute()
	if slices.Contains(minutes, currentMinute) {
		return nil
	}

	offset := 0
	for i := 0; i < 2; i++ {
		for _, possibleTargetMinute := range minutes {
			if possibleTargetMinute < currentMinute {
				continue
			}
			*t = t.Add(time.Minute * time.Duration((possibleTargetMinute+offset)-currentMinute))
			return nil
		}

		offset = 60 - currentMinute
		currentMinute = 0
	}

	return ErrScheduler.ErrResolveMinute
}

// Set time of t to next occurence in schedule
func SetNextTime(t *time.Time, scedule *Scedule) error {
	if err := verifySceduleData(scedule); err != nil {
		return err
	}

	*t = t.Add(-time.Second*time.Duration(t.Second()) + -time.Nanosecond*time.Duration(t.Nanosecond()))

	if err := setNextTimeByMinute(t, scedule.Minutes); err != nil {
		return err
	}
	if err := setNextTimeByHour(t, scedule.Hours); err != nil {
		return err
	}
	if err := setNextTimeByDay(t, scedule.Days); err != nil {
		return err
	}
	if err := setNextTimeByWeek(t, scedule.Weeks); err != nil {
		return err
	}
	if err := setNextTimeByMonth(t, scedule.Months); err != nil {
		return err
	}

	return nil
}

// Sleep for dur and print time remaining every interval
func SleepFor(msg string, dur time.Duration, interval time.Duration) {
	endTime := time.Now().Add(dur)
	for {
		if endTime.Before(time.Now()) {
			break
		}
		untilEndTime := time.Until(endTime)
		if msg != "" {
			width, _, _ := term.GetSize(0)
			fmt.Printf("\r%"+strconv.Itoa(width)+"v\r%v%v", "", msg, untilEndTime.Round(interval).String())
		}
		if interval < untilEndTime {
			time.Sleep(interval)
			continue
		}
		time.Sleep(untilEndTime)
	}

	if msg != "" {
		width, _, _ := term.GetSize(0)
		fmt.Printf("\r%"+strconv.Itoa(width)+"v\r", "")
	}
}
