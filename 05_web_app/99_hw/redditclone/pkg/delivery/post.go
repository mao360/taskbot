package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	redditclone "reddit_clone"
	"reddit_clone/pkg/session"
)

type PostHandler struct {
	Logger    *zap.SugaredLogger
	PostsRepo redditclone.PostsRepo
}

// GET
func (h *PostHandler) AllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostsRepo.AllPosts()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "Can`t write response body", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("All posts showed")
}

// POST
func (h *PostHandler) NewPost(w http.ResponseWriter, r *http.Request) {
	// payload: {category: "videos", type: "text", title: "some_title", text: "some_text"}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Can`t read request body", http.StatusInternalServerError)
		return
	}

	post := new(redditclone.Post)
	err = json.Unmarshal(body, post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
		return
	}

	userMap, err := session.SessFromContext(r.Context())
	if err != nil {
		http.Error(w, "Bad session", http.StatusInternalServerError)
		return
	}

	user := redditclone.User{
		UserID:   userMap["userID"].(string),
		UserName: userMap["username"].(string),
		Password: userMap["password"].(string),
	}
	post, err = h.PostsRepo.NewPost(post, &user)
	if err != nil {
		http.Error(w, "Can`t create post", http.StatusInternalServerError)
	}

	data, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "Can`t convert to JSON", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Can`t write response body", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Post created, post ID: %v", post.PostID)
}

// GET
func (h *PostHandler) GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	posts, err := h.PostsRepo.GetPostsByCategory(vars["CATEGORY_NAME"])
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Can`t convert to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Can`t write response body", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Posts by category showed")
}

// GET
func (h *PostHandler) PostDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, ok := vars["POST_ID"]
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	post, err := h.PostsRepo.PostDetails(postID)
	if err != nil || post == nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "Can`t convert to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Can`t write response body", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Post details showed, postID: %v", post.PostID)
}

// POST
func (h *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["POST_ID"]

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Can`t read request body", http.StatusInternalServerError)
		return
	}

	comment := new(redditclone.Comment)
	err = json.Unmarshal(body, comment)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	userMap, err := session.SessFromContext(r.Context())
	if err != nil {
		http.Error(w, "Bad session", http.StatusInternalServerError)
		return
	}
	user := redditclone.User{
		UserID:   userMap["userID"].(string),
		UserName: userMap["username"].(string),
	}
	err = h.PostsRepo.AddComment(comment, postID, &user)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	post, err := h.PostsRepo.PostDetails(postID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "Can`t convert to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Can`t write response body", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Comment created, postID: %v, commentID: %v", post.PostID)
}

// DELETE
func (h *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, commentID := vars["POST_ID"], vars["COMMENT_ID"]

	err := h.PostsRepo.DeleteComment(postID, commentID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	post, err := h.PostsRepo.PostDetails(postID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "Can`t convert to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Can`t write response body", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Comment deleted, postID: %v, commentID: %v", post.PostID)
}

// GET
func (h *PostHandler) UpPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["POST_ID"]

	err := h.PostsRepo.UpPost(postID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	post, err := h.PostsRepo.PostDetails(postID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "Can`t convert to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Can`t write response body", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Post upped, postID: %v", post.PostID)
}

// GET
func (h *PostHandler) DownPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["POST_ID"]

	err := h.PostsRepo.DownPost(postID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	post, err := h.PostsRepo.PostDetails(postID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "Can`t convert to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Can`t write response body", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Post downed, postID: %v", post.PostID)
}

// DELETE
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["POST_ID"]

	err := h.PostsRepo.DeletePost(postID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(map[string]string{
		"message": "success",
	})
	if err != nil {
		http.Error(w, "Can`t convert to JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Can`t write response body", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Post deleted, postID: %v", postID)
}

// GET
func (h *PostHandler) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["USER_ID"]

	posts, err := h.PostsRepo.GetPostsByUser(userID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Can`t convert to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Can`t write response body", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Posts by user, userID: %v", userID)
}
