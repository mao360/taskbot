package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"reddit_clone/pkg/delivery"
	"reddit_clone/pkg/middleware"
	"reddit_clone/pkg/repo"
)

func main() {

	templates := template.Must(template.ParseGlob("./static/html/*.html"))
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	usersRepo := repo.NewUsersMemRepo()
	postsRepo := repo.NewPostsMemRepo()

	userHandler := &delivery.UserHandler{
		Tmpl:      templates,
		Logger:    logger,
		UsersRepo: usersRepo,
	}
	postHandler := &delivery.PostHandler{
		Logger:    logger,
		PostsRepo: postsRepo,
	}

	r := mux.NewRouter()

	r.Handle("/static/css/main.74225161.chunk.css", http.FileServer(http.Dir("./")))
	r.Handle("/static/js/2.d59deea0.chunk.js", http.FileServer(http.Dir("./")))
	r.Handle("/static/js/main.32ebaf54.chunk.js", http.FileServer(http.Dir("./")))

	r.HandleFunc("/", userHandler.Index).Methods("GET")                     // обработка корневого запроса
	r.HandleFunc("/api/register", userHandler.Registration).Methods("POST") // регистрация
	r.HandleFunc("/api/login", userHandler.LogIn).Methods("POST")           // логин

	r.HandleFunc("/api/posts/", postHandler.AllPosts).Methods("GET")                          // список всех постов
	r.HandleFunc("/api/posts/{CATEGORY_NAME}", postHandler.GetPostsByCategory).Methods("GET") // список постов конкретной категории
	r.HandleFunc("/api/post/{POST_ID}", postHandler.PostDetails).Methods("GET")               // детали поста с комментами
	r.HandleFunc("/api/user/{USER_ID}", postHandler.GetPostsByUser).Methods("GET")            // получение всех постов конкртеного пользователя

	r.HandleFunc("/api/posts", postHandler.NewPost).Methods("POST")                               // добавление поста
	r.HandleFunc("/api/post/{POST_ID}", postHandler.AddComment).Methods("POST")                   // добавление коммента
	r.HandleFunc("/api/post/{POST_ID}/{COMMENT_ID}", postHandler.DeleteComment).Methods("DELETE") // удаление коммента
	r.HandleFunc("/api/post/{POST_ID}/upvote", postHandler.UpPost).Methods("GET")                 // рейтинг постп вверх
	r.HandleFunc("/api/post/{POST_ID}/downvote", postHandler.DownPost).Methods("GET")             // рейтинг поста вниз
	r.HandleFunc("/api/post/{POST_ID}", postHandler.DeletePost).Methods("DELETE")                 // удаление поста

	mux := middleware.Auth(r)
	mux = middleware.AccessLog(logger, mux)
	mux = middleware.Panic(mux)

	addr := ":8080"

	fmt.Println("server started")
	http.ListenAndServe(addr, mux)
}
