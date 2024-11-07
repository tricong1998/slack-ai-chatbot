package slack_handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
)

const (
	WorkingTime8001730 = "8:00-17:30"
	WorkingTime8301800 = "8:30-18:00"
	WorkingTime9001830 = "9:00-18:30"
)

func (s *SlackHandler) handleCreateLeaveRequestSubmission(payload slack.InteractionCallback) error {
	fmt.Println("handleCreateLeaveRequestSubmission----")
	startDate := payload.BlockActionState.Values["date_pickers"]["request_date_from_input"].SelectedDate
	endDate := payload.BlockActionState.Values["date_pickers"]["request_date_to_input"].SelectedDate
	hourFrom := payload.BlockActionState.Values["time_pickers"]["hour_from_input"].SelectedTime
	hourTo := payload.BlockActionState.Values["time_pickers"]["hour_to_input"].SelectedTime
	description := payload.BlockActionState.Values["description"]["description_input"].Value
	workingTime := payload.BlockActionState.Values["leave_type"]["working_time_input"].SelectedOption.Value
	leaveType := payload.BlockActionState.Values["leave_type"]["leave_type_input"].SelectedOption.Value
	workerEmail := payload.BlockActionState.Values["worker_email"]["email_input"].Value

	fmt.Println("leaveType: ", leaveType)
	fmt.Println("startDate: ", startDate)
	fmt.Println("endDate: ", endDate)
	fmt.Println("hourFrom: ", hourFrom)
	fmt.Println("hourTo: ", hourTo)
	fmt.Println("description: ", description)
	fmt.Println("workingTime: ", workingTime)
	fmt.Println("workerEmail: ", workerEmail)
	if startDate == "" || endDate == "" || hourFrom == "" || hourTo == "" || description == "" || workingTime == "" || leaveType == "" || workerEmail == "" {
		return s.slackService.SendMessage(context.Background(), &payload.Channel.ID, "All fields are required")
	}

	// Parse the dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return s.slackService.SendMessage(context.Background(), &payload.Channel.ID, "Invalid start date format")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return s.slackService.SendMessage(context.Background(), &payload.Channel.ID, "Invalid end date format")
	}

	hourFromCode := getHourFromCode(hourFrom)
	hourToCode := getHourFromCode(hourTo)
	calendarId, err := strconv.Atoi(workingTime)
	if err != nil {
		return s.slackService.SendMessage(context.Background(), &payload.Channel.ID, "Invalid working time")
	}
	holidayStatusId, err := strconv.Atoi(leaveType)
	if err != nil {
		return s.slackService.SendMessage(context.Background(), &payload.Channel.ID, "Invalid leave type")
	}
	fmt.Println("start: ", start)
	fmt.Println("end: ", end)
	fmt.Println("hourFromCode: ", hourFromCode)
	fmt.Println("hourToCode: ", hourToCode)
	fmt.Println("calendarId: ", calendarId)
	fmt.Println("holidayStatusId: ", holidayStatusId)
	fmt.Println("workerEmail: ", workerEmail)
	fmt.Println("start.Format: ", start.Format("02/01/2006"))
	fmt.Println("end.Format: ", end.Format("02/01/2006"))
	err = s.uiPathJobService.CreateLeaveRequestJob(dto.UIPathCreateLeaveRequestInput{
		RequestDateFrom: start.Format("02/01/2006"),
		RequestDateTo:   end.Format("02/01/2006"),
		HourFrom:        hourFromCode,
		HourTo:          hourToCode,
		Description:     description,
		CalendarId:      calendarId,
		WorkerEmail:     workerEmail,
		HolidayStatusId: holidayStatusId,
	}, payload.Channel.ID)
	return err
}

func (s *SlackHandler) handleLeaveRequestEvent(channelID string) error {
	return s.slackService.SendCreateLeaveRequestForm(context.Background(), channelID)
}

func getHourFromCode(hourFrom string) int {
	switch hourFrom {
	case "08:00":
		return 8
	case "08:30":
		return -9
	case "09:00":
		return 9
	case "09:30":
		return -10
	case "10:00":
		return 10
	case "10:30":
		return -11
	case "11:00":
		return 11
	case "11:30":
		return -12
	case "12:00":
		return 12
	case "13:30":
		return -14
	case "14:00":
		return 14
	case "14:30":
		return -15
	case "15:00":
		return 15
	case "15:30":
		return -16
	case "16:00":
		return 16
	case "16:30":
		return -17
	case "17:00":
		return 17
	case "17:30":
		return -18
	case "18:00":
		return 18
	case "18:30":
		return -19
	}

	return 0
}
