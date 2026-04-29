package main

import (
	"log"
	"net/http"
	"social-network-backend/pkg/db"
	"social-network-backend/pkg/handlers"
	"github.com/gorilla/sessions"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// Разрешаем localhost и любые локальные адреса
		if origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// Для других случаев (из браузера или других хостов)
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	database, err := db.InitDB("./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	store := sessions.NewCookieStore([]byte("your-secret-key-change-in-production"))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/register", handlers.RegisterHandler(database))
	mux.HandleFunc("/api/login", handlers.LoginHandler(database, store))
	mux.HandleFunc("/api/logout", handlers.LogoutHandler(store))
	mux.HandleFunc("/api/profile", handlers.GetProfileHandler(database, store))
	mux.HandleFunc("/api/profile/update", handlers.UpdateProfileHandler(database, store))
	mux.HandleFunc("/api/posts", handlers.GetPostsHandler(database, store))
	mux.HandleFunc("/api/posts/create", handlers.CreatePostHandler(database, store))
	mux.HandleFunc("/api/feed", handlers.GetFeedHandler(database, store))
	mux.HandleFunc("/api/comments", handlers.GetCommentsHandler(database, store))
	mux.HandleFunc("/api/comments/create", handlers.CreateCommentHandler(database, store))
	mux.HandleFunc("/api/followers", handlers.GetFollowersHandler(database, store))
	mux.HandleFunc("/api/following", handlers.GetFollowingHandler(database, store))
	mux.HandleFunc("/api/follow", handlers.FollowHandler(database, store))
	mux.HandleFunc("/api/unfollow", handlers.UnfollowHandler(database, store))
	mux.HandleFunc("/api/follow/accept", handlers.AcceptFollowHandler(database, store))
	mux.HandleFunc("/api/groups", handlers.GetGroupsHandler(database, store))
	mux.HandleFunc("/api/groups/create", handlers.CreateGroupHandler(database, store))
	mux.HandleFunc("/api/groups/join", handlers.JoinGroupHandler(database, store))
	mux.HandleFunc("/api/groups/accept", handlers.AcceptGroupJoinHandler(database, store))
	mux.HandleFunc("/api/group-posts", handlers.GetGroupPostsHandler(database, store))
	mux.HandleFunc("/api/group-posts/create", handlers.CreateGroupPostHandler(database, store))
	mux.HandleFunc("/api/events", handlers.GetGroupEventsHandler(database, store))
	mux.HandleFunc("/api/events/create", handlers.CreateEventHandler(database, store))
	mux.HandleFunc("/api/events/respond", handlers.RespondEventHandler(database, store))
	mux.HandleFunc("/ws/chat", handlers.ChatHandler(database, store))
	mux.HandleFunc("/api/messages", handlers.GetMessagesHandler(database, store))
	mux.HandleFunc("/api/notifications", handlers.GetNotificationsHandler(database, store))
	mux.HandleFunc("/api/notifications/read", handlers.MarkNotificationReadHandler(database, store))

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", enableCORS(mux))
}