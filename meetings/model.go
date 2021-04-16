package meetings

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type ZoomToken struct {
	Token string `json:"token"`
}
type Meeting struct {
	UUID              string    `json:"uuid"`
	ID                int       `json:"id"`
	HostID            int       `json:"host_id"`
	Type              int       `json:"type"`
	Topic             string    `json:"topic"`
	UserName          string    `json:"user_name"`
	UserEmail         string    `json:"user_email"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	Duration          int       `json:"duration"`
	TotalMinutes      int       `json:"total_minutes"`
	ParticipantsCount int       `json:"participants_count"`
}

type ResultMeetings struct {
	Meetings      []Meeting `json:"meetings"`
	NextPageToken string    `json:"next_page_token"`
	PageSize      int       `json:"page_size"`
	TotalRecords  int       `json:"total_records"`
}

type MeetingInstance struct {
	MeetingId int       `json:"meeting_id"`
	StartTime time.Time `json:"start_time"`
	UUID      string    `json:"uuid"`
}

func (u *MeetingInstance) ModifyId(id int) {

	u.MeetingId = id
}

type ResultInstance struct {
	Meetings []MeetingInstance `json:"meetings"`
}

type Participant struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	UserEmail string `json:"user_email"`
}

type Group struct {
	ZoomMeetingId string `json:"zoom_meeting_id"`
}

type GroupResponse struct {
	Results []Group `json:"results"`
}

type ResultParticipants struct {
	NextPageToken string        `json:"next_page_token"`
	PageSize      int           `json:"page_size"`
	PageCount     int           `json:"page_count"`
	Participants  []Participant `json:"participants"`
}
type MeetingAndParticipants struct {
	InstanceId   string        `json:"instance_id" pg:",pk"`
	Participants []Participant `json:"participants"`
	StartTime    time.Time     `json:"start_time"`
}

func InsertMeeting(db *pg.DB, meeting *Meeting) error {
	_, err := db.Model(meeting).Insert()

	if err != nil {
		return err
	}

	return nil
}

func InsertAllInstances(db *pg.DB, instance []MeetingInstance) error {
	_, err := db.Model(instance).Insert()

	if err != nil {
		return err
	}

	return nil
}

func InsertInstance(db *pg.DB, instance MeetingInstance) error {
	_, err := db.Model(instance).Insert()

	if err != nil {
		return err
	}

	return nil
}

func InstertParticipant(db *pg.DB, model *[]MeetingAndParticipants) error {
	_, err := db.Model(model).Insert()

	if err != nil {
		return err
	}

	return nil
}

func GetParticipant(db *pg.DB, model *MeetingAndParticipants) error {
	err := db.Model(model).Where("meeting_and_participants.instance_id = ?", model.InstanceId).Select()
	if err != nil {
		return err
	}

	return nil
}

func CreateSchema(db *pg.DB) error {
	model := (*MeetingAndParticipants)(nil)

	q := db.Model(model).Table()
	exists, err := q.Exists()
	if err != nil {
		fmt.Println(err)
	}

	if !exists {

		err := db.Model(model).CreateTable(&orm.CreateTableOptions{})

		if err != nil {
			return err
		}
	}
	return nil
}
