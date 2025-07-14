package scheduler

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

type (
	Schedule struct {
		Months  []int // Months of the year: `1 - 12`
		Weeks   []int // Weeks of the month: `1 - 5`
		Days    []int // Days of the week: `0 - 6` (Sunday first day of the week)
		Hours   []int // Hours of the day: `0 - 23`
		Minutes []int // Minutes of the hour: `0 - 59`
	}
)

var Errors = struct {
	InvalidMonths, InvalidWeeks, InvalidDays, InvalidHours, InvalidMinutes,
	ResolveMonth, ResolveWeek, ResolveDay, ResolveHour, ResolveMinute error
}{
	InvalidMonths:  errors.New("invalid value in months, valid values are [1-12]"),
	InvalidWeeks:   errors.New("invalid value in weeks, valid values are [1-5]"),
	InvalidDays:    errors.New("invalid value in days, valid values are [0-6]"),
	InvalidHours:   errors.New("invalid value in hours, valid values are [0-23]"),
	InvalidMinutes: errors.New("invalid value in minutes, valid values are [0-59]"),
	ResolveMonth:   errors.New("unable to resolve target month"),
	ResolveWeek:    errors.New("unable to resolve target week"),
	ResolveDay:     errors.New("unable to resolve target day"),
	ResolveHour:    errors.New("unable to resolve target hour"),
	ResolveMinute:  errors.New("unable to resolve target minute"),
}

func verifyScheduleData(schedule *Schedule) error {
	if i := slices.IndexFunc(schedule.Months, func(v int) bool { return v < 1 || v > 12 }); i != -1 {
		return Errors.InvalidMonths
	}
	if i := slices.IndexFunc(schedule.Weeks, func(v int) bool { return v < 1 || v > 5 }); i != -1 {
		return Errors.InvalidWeeks
	}
	if i := slices.IndexFunc(schedule.Days, func(v int) bool { return v < 0 || v > 6 }); i != -1 {
		return Errors.InvalidDays
	}
	if i := slices.IndexFunc(schedule.Hours, func(v int) bool { return v < 0 || v > 23 }); i != -1 {
		return Errors.InvalidHours
	}
	if i := slices.IndexFunc(schedule.Minutes, func(v int) bool { return v < 0 || v > 59 }); i != -1 {
		return Errors.InvalidMinutes
	}
	return nil
}

func setNextTimeByMonth(t *time.Time, months []int) error {
	currentMonth := int(t.Month())
	if slices.Contains(months, currentMonth) {
		return nil
	}

	offset := 0
	for range 10 {
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
	return Errors.ResolveMonth
}

func setNextTimeByWeek(t *time.Time, weeks []int) error {
	currentWeek := (int(t.Day()-1) / 7) + 1
	if slices.Contains(weeks, currentWeek) {
		return nil
	}

	offset := 0
	for range 2 {
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
	return Errors.ResolveWeek
}

func setNextTimeByDay(t *time.Time, days []int) error {
	currentDay := int(t.Weekday())
	if slices.Contains(days, currentDay) {
		return nil
	}

	offset := 0
	for range 2 {
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
	return Errors.ResolveDay
}

func setNextTimeByHour(t *time.Time, hours []int) error {
	currentHour := t.Hour()
	if slices.Contains(hours, currentHour) {
		return nil
	}

	offset := 0
	for range 2 {
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
	return Errors.ResolveHour
}

func setNextTimeByMinute(t *time.Time, minutes []int) error {
	currentMinute := t.Minute()
	if slices.Contains(minutes, currentMinute) {
		return nil
	}

	offset := 0
	for range 2 {
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
	return Errors.ResolveMinute
}

// Set time of t to next occurence in schedule
func SetNextTime(t *time.Time, schedule *Schedule) error {
	if err := verifyScheduleData(schedule); err != nil {
		return err
	}

	*t = t.Add(-time.Second*time.Duration(t.Second()) + -time.Nanosecond*time.Duration(t.Nanosecond()))

	if err := setNextTimeByMinute(t, schedule.Minutes); err != nil {
		return err
	}
	if err := setNextTimeByHour(t, schedule.Hours); err != nil {
		return err
	}
	if err := setNextTimeByDay(t, schedule.Days); err != nil {
		return err
	}
	if err := setNextTimeByWeek(t, schedule.Weeks); err != nil {
		return err
	}
	if err := setNextTimeByMonth(t, schedule.Months); err != nil {
		return err
	}
	return nil
}

// Sleep for a minimum of `dur`.
//
// If `msg` is not empty print `msg` + time remaining every `interval`.
func SleepFor(msg string, dur time.Duration, interval time.Duration) {
	endTime := time.Now().Add(dur)
	for time.Now().Before(endTime) {
		untilEndTime := time.Until(endTime)
		if msg != "" {
			fmt.Printf("\r\033[0J%v%v", msg, untilEndTime.Round(interval).String())
			if interval < untilEndTime {
				time.Sleep(interval)
				continue
			}
		}
		time.Sleep(untilEndTime)
	}

	if msg != "" {
		fmt.Print("\r\033[0J")
	}
}

// Short for `SleepFor(msg, time.Until(t), interval)`
func SleepUntil(msg string, t time.Time, interval time.Duration) {
	SleepFor(msg, time.Until(t), interval)
}
