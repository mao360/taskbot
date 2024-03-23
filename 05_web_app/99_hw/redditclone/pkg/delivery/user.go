package delivery

import (
	"encoding/json"
	"go.uber.org/zap"
	"html/template"
	"io/ioutil"
	"net/http"
	redditclone "reddit_clone"
	"reddit_clone/pkg/session"
)

type UserHandler struct {
	Tmpl      *template.Template
	Logger    *zap.SugaredLogger
	UsersRepo redditclone.UsersRepo
}

type Form struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Template err", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Index handler")
}

func (h *UserHandler) Registration(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "cant read body", http.StatusInternalServerError)
		return
	}
	f := new(Form)
	err = json.Unmarshal(body, f)
	if err != nil {
		http.Error(w, "cant unmarshall", http.StatusInternalServerError)
		return
	}

	user, err := h.UsersRepo.AppendUser(f.Password, f.Username) // что-то должно возвращать
	if err != nil || user == nil {
		http.Error(w, "cant append user", http.StatusInternalServerError)
		return
	}

	token, err := session.CreateSess(user)
	if err != nil && token == "" {
		http.Error(w, "empty token", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]interface{}{
		"token": token,
	})
	if err != nil {
		http.Error(w, "bad smth", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "bad end", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Can`t read body", http.StatusInternalServerError)
		return
	}

	var f Form
	err = json.Unmarshal(body, &f)
	if err != nil {
		http.Error(w, "err with unmarshall in login", http.StatusInternalServerError)
		return
	}

	user, err := h.UsersRepo.CheckUser(f.Password, f.Username)
	if err != nil {
		http.Error(w, "bad user", http.StatusInternalServerError)
	}

	tokenString, err := session.CreateSess(user)
	if err != nil {
		http.Error(w, "err in create sess", http.StatusInternalServerError)
	}

	resp, _ := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	w.Write(resp)
}
