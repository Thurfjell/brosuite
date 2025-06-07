package user

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type ApprovedLeave struct {
	From time.Time
	To   time.Time
}

type ReadData struct {
	Name               string
	LineManager        string
	TechnicalDirector  string
	Company            string
	CurrentProject     string
	CurrentProjectEnd  time.Time
	Role               string
	AnnualLeaveTotal   uint16
	AnnualLeaveUsed    uint16
	TrainingLeaveTotal uint16
	TrainingLeaveUsed  uint16
	SauceLeaveTotal    uint16
	SauceLeaveUsed     uint16
	RightWidgetIDS     []string
	LeftWidgetIDS      []string
	MainWidgetIDS      []string
	ApprovedLeaves     []ApprovedLeave
}

type RandomUserBuster struct {
	user *ReadData
}

func (b *RandomUserBuster) GetUser() (user *ReadData, err error) {
	user = b.user
	return
}

func randz(from time.Time) time.Time {
	diff := time.Now().Unix() - from.Unix()
	if diff < 0 {
		diff = -diff
	}
	return from.Add(time.Duration(rand.Int63n(diff+1)) * time.Second)
}

func (b *RandomUserBuster) RandomizeUser() {
	b.user = RandomUser()
}

func RandomUser() (user *ReadData) {
	annualLeaveTotal := uint16(20 + rand.Intn(15))
	leaveLeft := annualLeaveTotal
	leavesC := 1 + rand.Intn(5)
	leaveStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	leaves := make([]ApprovedLeave, 0, leavesC)

	for range leavesC {

		if leaveLeft <= 0 {
			break
		}

		from := randz(leaveStart)

		to := from
		added := 1 + rand.Intn(int(leaveLeft))

		to = to.Add(time.Duration(added*24) * time.Hour)
		leaveStart = to
		leaveLeft = leaveLeft - uint16(added)

		leaves = append(leaves, ApprovedLeave{
			From: from,
			To:   to,
		})
	}

	user = &ReadData{
		Name:               gofakeit.Name(),
		Company:            gofakeit.Company(),
		Role:               gofakeit.JobTitle(),
		LineManager:        gofakeit.Name(),
		TechnicalDirector:  gofakeit.Name(),
		CurrentProjectEnd:  time.Now().Add(time.Duration(24*rand.Intn(365)) * time.Hour),
		CurrentProject:     fmt.Sprintf("%s %s", gofakeit.AppName(), gofakeit.BS()),
		AnnualLeaveTotal:   annualLeaveTotal,
		AnnualLeaveUsed:    annualLeaveTotal - leaveLeft,
		TrainingLeaveTotal: 5,
		TrainingLeaveUsed:  0,
		SauceLeaveTotal:    3,
		SauceLeaveUsed:     1,
		LeftWidgetIDS:      []string{"bhao3brtq4gjwqr5b3kdxeiz65", "jsxqrq6sx33fg2vfnpq4zdldn4"},
		MainWidgetIDS:      []string{"giqlvoffyv6gypkrb32o7tj5et", "ucr4psiigncdzbtsnviy7akbtc"},
		RightWidgetIDS:     []string{},
		ApprovedLeaves:     leaves,
	}
	return
}

func New() *RandomUserBuster {
	buster := &RandomUserBuster{}
	buster.user = RandomUser()
	return buster
}
