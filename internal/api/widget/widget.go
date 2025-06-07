package widget

import (
	"brosuite/internal/api"
	"brosuite/internal/user"
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

// Very scientific
const LeaveWidgetId string = "bhao3brtq4gjwqr5b3kdxeiz65"
const ApprovedWidgetId string = "jsxqrq6sx33fg2vfnpq4zdldn4"
const LinksWidgetId string = "giqlvoffyv6gypkrb32o7tj5et"
const InfoWidgetId string = "ucr4psiigncdzbtsnviy7akbtc"

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

type service struct {
	WidgetMap map[string]ContainerViewModel
	template  *template.Template
	Routes    []api.Route
}

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
		fmt.Printf("GET WIDGETS %s\n", id)
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
				if wdgt, ok := widgetMap[LeaveWidgetId]; ok {
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
					fmt.Printf("\n%+v\n", data)
					// err = template.ExecuteTemplate(w, "leave.html", data)
					err = template.ExecuteTemplate(w, wdgt.TemplateHtml, data)
				}

			case ApprovedWidgetId:
				if wdgt, ok := widgetMap[ApprovedWidgetId]; ok {
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

					err = template.ExecuteTemplate(w, wdgt.TemplateHtml, data)
				}

			case LinksWidgetId:
				if wdgt, ok := widgetMap[LinksWidgetId]; ok {
					err = template.ExecuteTemplate(w, wdgt.TemplateHtml, nil)
				}

			case InfoWidgetId:
				if wdgt, ok := widgetMap[InfoWidgetId]; ok {
					data := &infoWidget{
						LineManager:       user.LineManager,
						TechnicalDirector: user.TechnicalDirector,
						CurrentProject:    user.CurrentProject,
						CurrentProjectEnd: user.CurrentProjectEnd.Format("02 January 2006"),
					}
					err = template.ExecuteTemplate(w, wdgt.TemplateHtml, data)
				}

			default:
				err = fmt.Errorf("no template found for id '%s'", id)
			}
		}

	}
}
