package dashboard

import (
	"brosuite/internal/api"
	"brosuite/internal/api/widget"
	"brosuite/internal/user"
	"embed"
	"html/template"
	"net/http"
)

type UserService interface {
	GetUser() (user *user.ReadData, err error)
	RandomizeUser()
}

//go:embed template/*
var html embed.FS

type Dashboard struct {
	template *template.Template
	Routes   []api.Route
}

type WidgetService interface {
	Derp() map[string]widget.ContainerViewModel
}

func New(ws WidgetService, us UserService) (dashboard *Dashboard, err error) {
	template, err := template.ParseFS(html, "template/*.html")

	if err != nil {
		return
	}

	dashboard = &Dashboard{
		template: template,
		Routes: []api.Route{
			{
				Path:        "/",
				HandlerFunc: getDashboard(template, ws, us),
			},
		},
	}
	return
}

type layout struct {
	UserName    string
	CompanyName string
	CompanyRole string
	Left        []widget.ContainerViewModel
	Main        []widget.ContainerViewModel
	Right       []widget.ContainerViewModel
}

func getDashboard(template *template.Template, ws WidgetService, us UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, _ := us.GetUser()
		derp := ws.Derp()
		l := &layout{
			UserName:    user.Name,
			CompanyName: user.Company,
			CompanyRole: user.Role,
			Left:        make([]widget.ContainerViewModel, 0, len(user.LeftWidgetIDS)),
			Main:        make([]widget.ContainerViewModel, 0, len(user.MainWidgetIDS)),
			Right:       make([]widget.ContainerViewModel, 0, len(user.RightWidgetIDS)),
		}

		for _, x := range user.LeftWidgetIDS {
			l.Left = append(l.Left, derp[x])
		}

		for _, x := range user.RightWidgetIDS {
			l.Right = append(l.Right, derp[x])
		}
		for _, x := range user.MainWidgetIDS {
			l.Main = append(l.Main, derp[x])
		}

		err := template.ExecuteTemplate(w, "index.html", l)

		if err != nil {
			w.Write([]byte(":("))
		}
		us.RandomizeUser()
	}
}
