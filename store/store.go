package taskstore

import (
	"fmt"
	"sync"
	"time"
)

type Posts struct {
	ID     int       `json:"id"`
	Author string    `json:"author"'`
	Text   string    `json:"text"`
	Tags   []string  `json:"tags"`
	Due    time.Time `json:"due"`
}

type PostStore struct {
	mux    sync.Mutex
	Post   map[int]Posts
	nextID int
}

type PostStoreManager interface {
	CreatePost(tx string, tags []string, due time.Time) int
	GetPost(id int) (Posts, error)
	DeletePost(id int) error
	DeleteAllPost() error
	GetAllPost() []Posts
	GetPostByTags(tag string) []Posts
	GetPostByDue(year int, mn time.Month, day int) []Posts
}

func New() *PostStore {
	ts := &PostStore{}
	ts.Post = make(map[int]Posts)
	ts.nextID = 0
	return ts
}

func (p *PostStore) CreatePost(tx string, author string, tags []string, due time.Time) int {
	p.mux.Lock()
	defer p.mux.Unlock()

	post := Posts{
		ID:     p.nextID,
		Author: author,
		Text:   tx,
		Due:    due,
	}

	post.Tags = make([]string, len(tags))
	copy(post.Tags, tags)

	p.Post[p.nextID] = post
	p.nextID++
	return post.ID
}

func (p *PostStore) GetPost(id int) (Posts, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	t, ok := p.Post[id]
	if ok {
		return t, nil
	} else {
		return Posts{}, fmt.Errorf("Please change input id = %d, task not found", id)
	}
}

func (p *PostStore) DeletePost(id int) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if _, ok := p.Post[id]; !ok {
		return fmt.Errorf("Please change input id = %d, task not found", id)

	} else {

		delete(p.Post, id)
		return nil
	}
}

func (p *PostStore) DeleteAllPost() error {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.Post = make(map[int]Posts)
	return nil
}

func (p *PostStore) GetAllPost() []Posts {
	p.mux.Lock()
	defer p.mux.Unlock()

	all := make([]Posts, 0, len(p.Post))
	for _, task := range p.Post {
		all = append(all, task)
	}
	return all
}

func (p *PostStore) GetPostByTags(tag string) []Posts {
	p.mux.Lock()
	defer p.mux.Unlock()

	var posts []Posts

	for _, task := range p.Post {
		for _, tasktag := range task.Tags {
			if tasktag == tag {
				posts = append(posts, task)
				break
			}
		}
	}
	return posts
}

func (p *PostStore) GetPostByDue(year int, mn time.Month, day int) []Posts {
	p.mux.Lock()
	defer p.mux.Unlock()

	var posts []Posts

	for _, post := range p.Post {
		y, m, d := post.Due.Date()
		if y == year && m == mn && d == day {
			posts = append(posts, post)
		}
	}
	return posts
}
