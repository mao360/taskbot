package redditclone

import "time"

type User struct {
	UserID   string `json:"id"`
	UserName string `json:"username"`
	Password string `json:"-"`
}

type Post struct {
	PostID           string    `json:"id"`
	Score            int       `json:"score"`
	Views            int       `json:"views"`
	Type             string    `json:"type"`
	Title            string    `json:"title"`
	Author           User      `json:"author"`
	Category         string    `json:"category"`
	Text             string    `json:"text"`
	Votes            []Vote    `json:"votes"`
	Comments         []Comment `json:"comments"` // CommentsID []uint64
	Created          time.Time `json:"created"`
	UpvotePercentage int       `json:"upvotePercentage"`
}

type Vote struct {
	UserID string `json:"user"`
	Grade  int    `json:"vote"`
}

type Comment struct {
	CommentID string    `json:"id"`
	Created   time.Time `json:"created"`
	Author    User      `json:"author"`
	Body      string    `json:"comment"`
}
