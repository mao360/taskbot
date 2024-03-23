package redditclone

type PostsRepo interface {
	AllPosts() ([]Post, error)                                    // список всех постов
	GetPostsByCategory(category string) ([]Post, error)           // список постов конкретной категории
	PostDetails(postID string) (*Post, error)                     // детали поста
	GetPostsByUser(userID string) ([]Post, error)                 // получение всех постов конкртеного пользователя
	NewPost(post *Post, user *User) (*Post, error)                // добавление поста
	UpPost(postID string) error                                   // рейтинг постп вверх
	DownPost(postID string) error                                 // рейтинг поста вниз
	AddComment(comment *Comment, postID string, user *User) error // добавление коммента
	DeleteComment(postID, commentID string) error                 // удаление коммента
	DeletePost(postID string) error                               // удаление поста
}

type UsersRepo interface {
	CheckUser(password, username string) (*User, error)
	AppendUser(password, username string) (*User, error)
}
