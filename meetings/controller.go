package meetings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

func GetAllMeetings(db *pg.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var meetings []Meeting

		err := db.Model(&meetings).Select()

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "someError2",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"meetings": meetings,
		})
	}

	return gin.HandlerFunc(fn)
}

func GetMeeting(db *pg.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var meeting Meeting
		meetingId := c.Param("id")
		token := os.Getenv("ZOOM_TOKEN")
		baseUrl := os.Getenv("BASE_URL")
		url := baseUrl + "meetings/" + meetingId

		client := &http.Client{}
		r, _ := http.NewRequest(http.MethodGet, url, nil)

		r.Header.Add("Authorization", "Bearer "+token)

		resp, err := client.Do(r)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		json.Unmarshal(data, &meeting)
		c.JSON(http.StatusOK, meeting)

	}

	return gin.HandlerFunc(fn)
}

func SetToken(c *gin.Context) {
	var token ZoomToken
	err := c.ShouldBind(&token)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	oldToken := os.Getenv("ZOOM_TOKEN")

	err = os.Setenv("ZOOM_TOKEN", token.Token)

	if err != nil {
		fmt.Println(err)
	}

	newToken := os.Getenv("ZOOM_TOKEN")
	if oldToken == newToken {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token is identical",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func FetchPastMeetingsFromZoom(db *pg.DB) gin.HandlerFunc {

	fn := func(c *gin.Context) {
		token := os.Getenv("ZOOM_TOKEN")
		accountId := c.Param("accountId")
		url := "https://api.zoom.us/v2/users/" + accountId + "/meetings"
		client := &http.Client{}

		r, _ := http.NewRequest(http.MethodGet, url, nil)

		r.Header.Add("Authorization", "Bearer "+token)

		resp, err := client.Do(r)

		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)

		var result ResultMeetings
		json.Unmarshal(data, &result)
		if err != nil {
			fmt.Printf("%s", err)
		}

		c.JSON(http.StatusOK, result)

	}

	return gin.HandlerFunc(fn)
}

func FetchAllMeetingInstances(db *pg.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		meetingId := c.Param("id")
		result := RequestInstances(meetingId)
		c.JSON(http.StatusOK, result)
	}

	return gin.HandlerFunc(fn)
}

func FetchAllParticipants(db *pg.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		meetingId := c.Param("meetingId")
		res := RequestParticipants(meetingId)
		c.JSON(http.StatusOK, res)
	}

	return gin.HandlerFunc(fn)
}

func FetchParticipant(db *pg.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		participantId := c.Param("participantId")
		model := &MeetingAndParticipants{InstanceId: participantId}
		err := GetParticipant(db, model)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "404",
				"message": "Not found",
			})
			return
		}

		c.JSON(http.StatusOK, model)
	}

	return gin.HandlerFunc(fn)
}

func RequestParticipants(meetingId string) []MeetingAndParticipants {
	token := os.Getenv("ZOOM_TOKEN")
	baseUrl := os.Getenv("BASE_URL")
	instances := RequestInstances(meetingId)
	client := &http.Client{}
	result := make([]MeetingAndParticipants, len(instances.Meetings))
	for i, meeting := range instances.Meetings {

		url := baseUrl + "past_meetings/" + meeting.UUID + "/participants?page_size=120"

		r, _ := http.NewRequest(http.MethodGet, url, nil)
		r.Header.Add("Authorization", "Bearer "+token)
		resp, err := client.Do(r)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		var res ResultParticipants
		data, _ := io.ReadAll(resp.Body)
		json.Unmarshal(data, &res)
		meet := MeetingAndParticipants{InstanceId: meeting.UUID, Participants: res.Participants, StartTime: meeting.StartTime}
		result[i] = meet
	}
	return result
}

func RequestInstances(meetingId string) ResultInstance {
	token := os.Getenv("ZOOM_TOKEN")
	baseUrl := os.Getenv("BASE_URL")
	url := baseUrl + "past_meetings/" + meetingId + "/instances"

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodGet, url, nil)

	r.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	var result ResultInstance
	json.Unmarshal(data, &result)
	if err != nil {
		fmt.Printf("%s", err)
	}

	for i := range result.Meetings {
		id, err := strconv.Atoi(meetingId)

		if err != nil {
			fmt.Println(err)

		}
		result.Meetings[i].ModifyId(id)

	}
	return result
}

func GetAllClassMeetings() []Group {
	token := os.Getenv("CLASS_CENTER_TOKEN")
	endpoint := "https://class-center.herokuapp.com/api/v1/groups"
	client := &http.Client{}

	r, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		panic(err)
	}

	r.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(r)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var result GroupResponse

	json.Unmarshal(data, &result)
	return result.Results
}

func CronJob(db *pg.DB) {
	meetingIds := GetAllClassMeetings()

	participants := make([][]MeetingAndParticipants, 20)

	for _, v := range meetingIds {
		p := RequestParticipants(v.ZoomMeetingId)
		participants = append(participants, p)
	}

	for _, i := range participants {
		len := len(i)
		if len > 0 {
			err := InstertParticipant(db, &i)
			if err != nil {
				pgErr, ok := err.(pg.Error)
				if ok && pgErr.IntegrityViolation() {
					fmt.Println("Participant already exists:", err)
				} else if pgErr.Field('S') == "PANIC" {
					panic(err)
				}
			}
		}

	}
}
