package main

import (
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
)


type User struct {
	Id   string
	Name string
}

type UserResource struct {
	users map[string]User
}

func (u UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/users").Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{user-id}").To(u.findUser))
	ws.Route(ws.POST("/").To(u.updateUser))

	container.Add(ws)
}

func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	usr, ok := u.users[id]
	if !ok {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "user cloud not be found")
	} else {
		response.WriteEntity(usr)
	}
}

func (u UserResource) updateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.Id] = *usr
		response.WriteEntity(usr)
	} else {
		response.AddHeader("Context-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	delete(u.users, id)
}

func main() {
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	u := UserResource{map[string]User{}}
	u.Register(wsContainer)

	log.Printf("start listening on localhost:8080")
	server := &http.Server{
		Addr:    ":8080",
		Handler: wsContainer,
	}
	log.Fatal(server.ListenAndServe())
}
