package routes

import(
	"net/http"
	"CaptureTheSoul/controllers"
)

func SetupRoutes(){
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/challenges", controllers.ListChallenges)
	http.HandleFunc("/submit", controllers.SubmitFlag)
	http.HandleFunc("/scoreboard", controllers.Scoreboard)
	http.HandleFunc("/defend", controllers.DefendStatus)
}