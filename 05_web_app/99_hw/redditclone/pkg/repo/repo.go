package repo

import (
	redditclone "reddit_clone"
	"strconv"
	"sync"
	"time"
)

type PostsMemoryRepo struct {
	data []redditclone.Post
	mu   *sync.Mutex
}

type UsersMemoryRepo struct {
	data []redditclone.User
	mu   *sync.Mutex
}

func NewPostsMemRepo() *PostsMemoryRepo {
	var data []redditclone.Post
	mu := &sync.Mutex{}
	return &PostsMemoryRepo{
		data: data,
		mu:   mu,
	}
}

func NewUsersMemRepo() *UsersMemoryRepo {
	var data []redditclone.User
	mu := &sync.Mutex{}
	return &UsersMemoryRepo{
		data: data,
		mu:   mu,
	}
}

func (p *PostsMemoryRepo) AllPosts() ([]redditclone.Post, error) {
	return p.data, nil
}

func (p *PostsMemoryRepo) GetPostsByCategory(category string) ([]redditclone.Post, error) {
	var postsByCategory []redditclone.Post
	for _, val := range p.data {
		if val.Category == category {
			postsByCategory = append(postsByCategory, val)
		}
	}
	return postsByCategory, nil
}

func (p *PostsMemoryRepo) PostDetails(postID string) (*redditclone.Post, error) {

	for i := range p.data {
		if p.data[i].PostID == postID {
			return &p.data[i], nil
		}
	}
	return nil, nil
}

func (p *PostsMemoryRepo) GetPostsByUser(userID string) ([]redditclone.Post, error) {
	var postsByUser []redditclone.Post
	for _, val := range p.data {
		if val.Author.UserID == userID {
			postsByUser = append(postsByUser, val)
		}
	}
	return postsByUser, nil
}

func (p *PostsMemoryRepo) NewPost(post *redditclone.Post, user *redditclone.User) (*redditclone.Post, error) {
	// category, type, title, text
	post.PostID = strconv.Itoa(len(p.data))
	post.Score = 0
	post.Views = 0
	post.Author = *user
	post.Created = time.Now()
	post.UpvotePercentage = 0
	post.Comments = make([]redditclone.Comment, 0)
	post.Votes = make([]redditclone.Vote, 0)
	post.Votes = append(post.Votes, redditclone.Vote{
		UserID: user.UserID,
		Grade:  1,
	})
	p.mu.Lock()
	p.data = append(p.data, *post)
	p.mu.Unlock()
	return &p.data[len(p.data)-1], nil
}

func (p *PostsMemoryRepo) UpPost(postID string) error {
	for i := range p.data {
		if p.data[i].PostID == postID {
			p.data[i].Score++
			return nil
		}
	}
	return nil
}

func (p *PostsMemoryRepo) DownPost(postID string) error {
	for i := range p.data {
		if p.data[i].PostID == postID {
			p.data[i].Score--
			return nil
		}
	}
	return nil
}

func (p *PostsMemoryRepo) AddComment(comment *redditclone.Comment, postID string, user *redditclone.User) error {
	for i := range p.data {
		if p.data[i].PostID == postID {
			comment.CommentID = strconv.Itoa(len(p.data[i].Comments))
			comment.Author = *user // безпарольный юзер
			comment.Created = time.Now()
			p.data[i].Comments = append(p.data[i].Comments, *comment)
			return nil
		}
	}
	return nil
}

func (p *PostsMemoryRepo) DeleteComment(postID, commentID string) error {
	for i := range p.data {
		if p.data[i].PostID == postID {
			for j := range p.data[i].Comments {
				if p.data[i].Comments[j].CommentID == commentID {
					p.data[i].Comments = append(p.data[i].Comments[:j], p.data[i].Comments[j+1:]...) // удаление
					return nil
				}
			}
		}
	}
	return nil
}

func (p *PostsMemoryRepo) DeletePost(postID string) error {
	for i := range p.data {
		if p.data[i].PostID == postID {
			p.data = append(p.data[:i], p.data[i+1:]...) // удаление
			return nil
		}
	}
	return nil
}

func (u *UsersMemoryRepo) CheckUser(password, username string) (*redditclone.User, error) {
	for i := range u.data {
		if u.data[i].UserName == username && u.data[i].Password == password {
			return &u.data[i], nil
		}
	}
	return nil, nil
}

func (u *UsersMemoryRepo) AppendUser(password, username string) (*redditclone.User, error) {
	id := strconv.Itoa(len(u.data))
	user := redditclone.User{
		UserID:   id,
		UserName: username,
		Password: password,
	}
	u.data = append(u.data, user)
	return &u.data[len(u.data)-1], nil
}
