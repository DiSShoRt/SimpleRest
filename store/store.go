package poststore

import (
	"fmt"
	"sync"
	"time"
)


type post struct {
	ID   int       `json:"id"`
	Name string    `json:"name"`
	Age  int       `json:"age"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`

}


type postStore struct {
	mux  sync.Mutex
	posts map[int]post
	nextID int
}

func New() *postStore {
	ts := &postStore{}
	ts.posts = make(map[int]post)
	ts.nextID = 0
	return ts
}

func (ts *postStore) Createpost(name string, age int, tx string, tags []string, due time.Time) int {
	ts.mux.Lock()

	post := post{
		ID: ts.nextID,
		Text: tx,
		Due: due,
		Age:age,
		Name:name,
	}

	post.Tags = make([]string, len(tags))
	copy(post.Tags, tags)
	

	ts.posts[ts.nextID] = post
	ts.nextID++
	return post.ID
}


func (ts *postStore) Getpost(id int) (post, error) {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	t, ok := ts.posts[id]
	if ok {
		return t, nil
	} else { 
		return post{}, fmt.Errorf("Please change input id = %d, post not found", id)
	}
}


func (ts *postStore) Deletepost(id int) error {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	if _, ok := ts.posts[id]; !ok {
		return fmt.Errorf("Please change input id = %d, post not found", id)

	} else { 

		delete(ts.posts, id)
		return nil
	}
}


func (ts *postStore) DeleteAllpost() error {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	ts.posts = make(map[int]post)
	return nil
}


func (ts *postStore) GetAllpost() []post {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	all := make([]post, 0, len(ts.posts))
	for _, post := range ts.posts {
		all = append(all, post)
	}
	return all
}


func (ts *postStore) GetpostByTags(tag string) []post {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	var posts []post

	for _,  post := range ts.posts {
		for _, posttag := range post.Tags {
			if posttag == tag {
				posts = append(posts, post)
				break
			}
		}
	}
	return posts
}


func (ts *postStore) GetpostByDue(year int, mn time.Month, day int) []post {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	var posts []post

	for _, post := range ts.posts {
		y, m, d := post.Due.Date()
		if y == year && m == mn && d == day {
			posts = append(posts, post)
		}
	}
	return posts
}