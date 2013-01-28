package main

import (	
	"fmt"
	"html/template"
	"time"

	"github.com/badgerodon/web"
)

type (
	BlogController struct {}
	BlogModel struct {
		Posts []BlogPost
	}
	BlogPost struct {
		Title string
		Date time.Time
		Body template.HTML
		Author User
	}
	User struct {
		FirstName string
		LastName string
	}
)

func (this User) Name() string {
	return this.FirstName + " " + this.LastName
}

func (this *BlogController) Index(ctx web.Context) {
	ctx.Render(BlogModel{
		Posts: []BlogPost{
			BlogPost{
				Title: "Blog Post #1",
				Date: time.Now(),
				Body: template.HTML("An example blog post."),
				Author: User{
					FirstName: "Montana",
					LastName: "Banana",
				},
			},
			BlogPost{
				Title: "Blog Post #2",
				Date: time.Now(),
				Body: template.HTML("Another example blog post."),
				Author: User{
					FirstName: "Montana",
					LastName: "Banana",
				},
			},
		},
	})
}
func (this *BlogController) Show(ctx web.Context, id int) {
	ctx.Write(fmt.Sprint("ID: ", id))
}

func init() {
	web.Route("/blog", &BlogController{})
	web.Route("/", func(ctx web.Context) {
		ctx.Write("!")
	})
}