package post_logic

type Post struct {
	ID            int
	Title         string
	Content       string
	Username      string
	Image         string
	Gif           string
	IsLike        bool
	IsDislike     bool
	LikesCount    int
	DislikesCount int
	ReportID      int
}
type Like struct {
	PostID    int
	ID        int
	Username  string
	PostTitle string
}

func ReversePosts(posts []Post) []Post {
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}
	return posts
}
func GetUsername(post Post) string {
	return post.Username
}
