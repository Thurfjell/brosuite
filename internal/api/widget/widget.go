package widget

import (
	"brosuite/internal/api"
	"brosuite/internal/user"
	"embed"
	"fmt"
	"html/template"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// Very scientific
const LeaveWidgetId string = "bhao3brtq4gjwqr5b3kdxeiz65"
const ApprovedWidgetId string = "jsxqrq6sx33fg2vfnpq4zdldn4"
const LinksWidgetId string = "giqlvoffyv6gypkrb32o7tj5et"
const InfoWidgetId string = "ucr4psiigncdzbtsnviy7akbtc"
const MyTimeWidgetId string = "ahuecy3xft2u6pzcry27ycn5o4"
const ReminderWidgetId string = "yulgmb5mfvzfov6wqq3qiwwclj"
const SocialMediaWidgetId string = "t3vt6s3kdwbboye3mvo6fh7dun"

type ContainerViewModel struct {
	ID           string
	Class        string
	TemplateHtml string
}

type UserService interface {
	GetUser() (user *user.ReadData, err error)
}

type leaveWidget struct {
	AnnualLeaveTotal   uint16
	AnnualLeaveLeft    uint16
	TrainingLeaveTotal uint16
	TrainingLeaveUsed  uint16
	SauceLeaveTotal    uint16
	SauceLeaveUsed     uint16
	TrainingPercent    uint16
	SaucePercent       uint16
	AnnualPercent      uint16
}

type fromTo struct {
	FromDay  string
	ToDay    string
	FromDate string // eg 15 Jul
	ToDate   string
	Total    uint8
}

type approvedWidget struct {
	TotalDays uint16
	Approves  []fromTo
}

type infoWidget struct {
	LineManager       string
	TechnicalDirector string
	CurrentProject    string
	CurrentProjectEnd string
}

type myTimeWidget struct {
	Times []myTimeItem
	Total uint16
}

type myTimeItem struct {
	Week         string
	Date         string
	CustomerName string
	Day          string
	Note         string
}

type socialMediaWidget struct {
	Posts []socialMediaPost
}

func timeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return fmt.Sprintf("%ds ago", int(duration.Seconds()))
	case duration < time.Hour:
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(duration.Hours()/24))
	}
}

func randomUserPortrait() string {
	id := 1 + rand.IntN(100)
	g := "women"
	if rand.IntN(2) == 0 {
		g = "men"
	}
	return fmt.Sprintf("https://randomuser.me/api/portraits/%s/%d.jpg", g, id)
}

func (w *socialMediaWidget) Random() (post socialMediaPost) {
	post = socialMediaPost{
		Hashtag: fmt.Sprintf("#%s", gofakeit.BuzzWord()),
		Tag:     fmt.Sprintf("@%s", gofakeit.CelebrityActor()),
		Content: gofakeit.HipsterSentence(5),
		TimeAgo: timeAgo(time.Now().Add(time.Duration(-rand.IntN(50)) * time.Hour)),
		Emoji:   gofakeit.Emoji(),
		ImgSrc:  randomUserPortrait(),
	}
	return
}

type socialMediaPost struct {
	Tag     string
	Content string
	Hashtag string
	TimeAgo string
	Emoji   string
	ImgSrc  string
}

type service struct {
	WidgetMap map[string]ContainerViewModel
	template  *template.Template
	Routes    []api.Route
}

// Imagine like GetWidget(id string) (widget *ContainerViewModel, err error) lols :)
func (s *service) Derp() map[string]ContainerViewModel {
	return s.WidgetMap
}

//go:embed template/*
var html embed.FS

func New(userService UserService) (widget *service, err error) {
	template, err := template.ParseFS(html, "template/*.html")
	if err != nil {
		return
	}

	widgetMap := make(map[string]ContainerViewModel)
	widgetMap[LeaveWidgetId] = ContainerViewModel{ID: LeaveWidgetId, Class: "bg-slate-800 rounded-md p-4 shadow-sm border border-slate-700 text-white", TemplateHtml: "leave.html"}
	widgetMap[ApprovedWidgetId] = ContainerViewModel{ID: ApprovedWidgetId, Class: "bg-slate-800 rounded-md p-4 shadow-sm border border-slate-700 text-white", TemplateHtml: "approved.html"}
	widgetMap[LinksWidgetId] = ContainerViewModel{ID: LinksWidgetId, Class: "bg-slate-800 rounded-md p-4 shadow-sm border border-slate-700 text-white", TemplateHtml: "links.html"}
	widgetMap[InfoWidgetId] = ContainerViewModel{ID: InfoWidgetId, Class: "bg-slate-800 rounded-md p-4 shadow-sm border border-slate-700 text-white", TemplateHtml: "userInfo.html"}
	widgetMap[MyTimeWidgetId] = ContainerViewModel{ID: MyTimeWidgetId, Class: "bg-slate-800 rounded-md p-4 shadow-sm border border-slate-700 text-white", TemplateHtml: "myTime.html"}
	widgetMap[ReminderWidgetId] = ContainerViewModel{ID: ReminderWidgetId, Class: "bg-gradient-to-r from-purple-700 via-pink-600 to-red-500 rounded-lg p-5 shadow-lg text-white", TemplateHtml: "reminder.html"}
	widgetMap[SocialMediaWidgetId] = ContainerViewModel{ID: SocialMediaWidgetId, Class: "bg-gradient-to-br from-blue-500 via-cyan-400 to-indigo-600 rounded-lg p-5 shadow-lg text-white", TemplateHtml: "socialMedia.html"}

	widget = &service{
		WidgetMap: widgetMap,
		template:  template,
		Routes: []api.Route{
			{
				Path:        "/widgets/{id}",
				HandlerFunc: getWidget(userService, widgetMap, template),
			},
		},
	}
	return
}

func getWidget(userService UserService, widgetMap map[string]ContainerViewModel, template *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var err error
		defer func() {
			if err != nil {
				w.Write([]byte(":("))
			}
		}()

		user, err := userService.GetUser()

		if err != nil {
			return
		}

		if widget, ok := widgetMap[id]; ok {
			switch widget.ID {
			case LeaveWidgetId:
				left := user.AnnualLeaveTotal - user.AnnualLeaveUsed
				data := &leaveWidget{
					AnnualLeaveTotal:   user.AnnualLeaveTotal,
					AnnualLeaveLeft:    left,
					AnnualPercent:      (left * 100) / user.AnnualLeaveTotal,
					TrainingLeaveTotal: user.TrainingLeaveTotal,
					TrainingLeaveUsed:  user.TrainingLeaveUsed,
					TrainingPercent:    (user.TrainingLeaveUsed * 100) / user.TrainingLeaveTotal,
					SauceLeaveTotal:    user.SauceLeaveTotal,
					SauceLeaveUsed:     user.SauceLeaveUsed,
					SaucePercent:       (user.SauceLeaveUsed * 100) / user.SauceLeaveTotal,
				}

				err = template.ExecuteTemplate(w, widget.TemplateHtml, data)

			case ApprovedWidgetId:
				data := &approvedWidget{
					TotalDays: 0,
					Approves:  make([]fromTo, 0, len(user.ApprovedLeaves)),
				}
				for _, x := range user.ApprovedLeaves {
					hours := x.To.Sub(x.From).Hours()
					xtotal := uint8(hours / 24)
					data.TotalDays = data.TotalDays + uint16(xtotal)
					ft := fromTo{
						FromDay:  x.From.Format("Mon"),
						ToDay:    x.To.Format("Mon"),
						FromDate: x.From.Format("02 Jan"),
						ToDate:   x.To.Format("02 Jan"),
						Total:    xtotal,
					}
					data.Approves = append(data.Approves, ft)
				}

				err = template.ExecuteTemplate(w, widget.TemplateHtml, data)

			case LinksWidgetId:
				err = template.ExecuteTemplate(w, widget.TemplateHtml, nil)

			case InfoWidgetId:
				data := &infoWidget{
					LineManager:       user.LineManager,
					TechnicalDirector: user.TechnicalDirector,
					CurrentProject:    user.CurrentProject,
					CurrentProjectEnd: user.CurrentProjectEnd.Format("02 January 2006"),
				}
				err = template.ExecuteTemplate(w, widget.TemplateHtml, data)

			case MyTimeWidgetId:

				data := &myTimeWidget{
					Times: make([]myTimeItem, 0, len(user.TimeData)),
					Total: uint16(user.TimeDataTotal),
				}
				for _, d := range user.TimeData {
					_, week := d.Date.ISOWeek()
					data.Times = append(data.Times, myTimeItem{
						Date:         d.Date.Format("2006-01-02"),
						Day:          d.Date.Format("Mon"),
						Note:         d.Note,
						CustomerName: d.Customer,
						Week:         strconv.Itoa(week),
					})
				}
				err = template.ExecuteTemplate(w, widget.TemplateHtml, data)

			case ReminderWidgetId:
				err = template.ExecuteTemplate(w, widget.TemplateHtml, nil)

			case SocialMediaWidgetId:
				data := socialMediaWidget{
					Posts: make([]socialMediaPost, 0, 6),
				}

				for range 6 {
					data.Posts = append(data.Posts, data.Random())
				}
				err = template.ExecuteTemplate(w, widget.TemplateHtml, data)

			default:
				err = fmt.Errorf("no template found for id '%s'", id)
			}
		}

	}
}
