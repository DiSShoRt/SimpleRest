package taskstore

import (
	"fmt"
	"github.com/jackc/pgx"
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
	AddPostToDb(author string, text string, tags []string, due time.Time) int
	DeletePostFromDb(id int) error
	DeleteAllPostsFromDb() error
	GetAllPostsDb() []Posts
	GetPostByIdDb(id int) (Posts, error)
	GetPostsByTagDb(tag string) []Posts
	CreatePost(tx string, author string, tags []string, due time.Time) int
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

func (ps *PostStore) AddPostToDb(author, text string, tags []string, due time.Time) int {

	//create config string
	connStr := "user=anton password=123 dbname=postgres sslmode=disable"
	//create connect config
	conf, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		fmt.Errorf("connection string is bad %s", err)

	}
	//connect to db
	db, err := pgx.Connect(conf)
	if err != nil {
		fmt.Errorf("cant connect to db %s", err)
	}
	// close connection
	defer db.Close()

	var id int
	err = db.QueryRow("INSERT INTO posts  VALUES ( nextval('postsseq'), $1, $2, $3, $4) returning id", author, text, tags, due).Scan(&id)

	if err != nil {
		fmt.Errorf("cant connect to db %s", err)
	}
	return id
}

func (ps *PostStore) DeletePostFromDb(id int) error {
	//create config string
	connStr := "user=anton password=123 dbname=postgres sslmode=disable"
	//create connect config
	conf, err := pgx.ParseConnectionString(connStr)

	if err != nil {
		fmt.Errorf("connection string is bad %s", err)
	}
	//connect to db
	db, err := pgx.Connect(conf)
	if err != nil {
		fmt.Errorf("cant connect to db %s", err)
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM posts WHERE id = $1", id)

	if err != nil {
		fmt.Errorf("cant connect to db %s with id %d", err, id)
	}
	return nil
}

func (ps *PostStore) DeleteAllPostsFromDb() error {
	//create config string
	connStr := "user=anton password=123 dbname=postgres sslmode=disable"
	//create connect config
	conf, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		fmt.Errorf("connection string is bad %s", err)
	}
	//connect to db
	db, err := pgx.Connect(conf)
	if err != nil {
		fmt.Errorf("cant connect to db %s", err)
	}

	defer db.Close()
	//without answer
	_, err = db.Exec("DELETE FROM posts ")

	if err != nil {
		fmt.Errorf("cant connect to db %s ", err)
	}
	return nil
}

func (ps *PostStore) GetPostsByAuthorDb(author string) []Posts {
	// create slice
	posts := []Posts{}
	//create config string
	connStr := "user=anton password=123 dbname=postgres sslmode=disable"
	//create config
	conf, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		fmt.Errorf("connection string is bad %s", err)
	}
	//create connection
	db, err := pgx.Connect(conf)

	if err != nil {
		fmt.Errorf("cant connect to db %s", err)
	}

	all, err := db.Query("SELECT * FROM posts WHERE author = $1", author)

	defer func() {
		all.Close()
		db.Close()
	}()

	for all.Next() {
		p := Posts{}
		all.Scan(&p.ID, &p.Author, &p.Text, &p.Tags, &p.Due)
		if err != nil {
			fmt.Errorf("simthing go wrong")
		}

		posts = append(posts, p)
	}
	fmt.Println(posts)
	return posts
}

func (ps *PostStore) GetAllPostsDb() []Posts {
	//create slice
	posts := []Posts{}
	//create config string
	connStr := "user=anton password=123 dbname=postgres sslmode=disable "
	//create connect config
	conf, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		fmt.Errorf("connection string is bad %s", err)
	}
	//connect to db
	db, err := pgx.Connect(conf)
	if err != nil {
		fmt.Errorf("cant connect to db %s", err)
	}
	//get rows
	all, err := db.Query("SELECT * FROM posts")
	if err != nil {
		fmt.Errorf("cant connect to db %s", err)
	}

	defer all.Close()
	//walk to posts
	for all.Next() {
		//create new object type Posts
		p := Posts{}
		// scanning values
		err = all.Scan(&p.ID, &p.Author, &p.Text, &p.Tags, &p.Due)
		if err != nil {
			fmt.Errorf("simthing go wrong")
		}
		//fmt.Println(p)

		posts = append(posts, p)
	}

	if all.Err() != nil {
		fmt.Errorf("cant NULL")
	}

	return posts
}

func (ps *PostStore) GetPostByIdDb(id int) (Posts, error) {
	//create config string
	connStr := "user=anton password=123 dbname=postgres sslmode=disable"
	//create connect config
	conf, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		fmt.Errorf("connection string is bad %s", err)
	}
	//connect to db
	db, err := pgx.Connect(conf)
	if err != nil {
		fmt.Errorf("cant connect to db %s", err)
	}

	defer db.Close()
	//one row
	post := db.QueryRow("SELECT * FROM posts WHERE id = $1", id)
	if err != nil {
		fmt.Errorf("cant connect to db %s", err)
	}
	//defer post.Close()

	p := Posts{}

	err = post.Scan(&p.ID, &p.Author, &p.Text, &p.Tags, &p.Due)
	if err != nil {
		fmt.Errorf("simthing whet wrong")
	}
	fmt.Println(p)
	return p, nil
}

func (ps *PostStore) GetPostsByTagDb(tag string) []Posts {
	//create config string
	connStr := "user=anton password=123 dbname=postgres sslmode=disable"
	//create connect config
	conf, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		fmt.Errorf("connection string is bad %s", err)
	}
	//connect to db
	db, err := pgx.Connect(conf)
	if err != nil {
		fmt.Errorf("cant connect to db %s", err)
	}
	//close connection
	defer db.Close()

	all, err := db.Query("SELECT * FROM posts WHERE tags[0] = $1 OR tags[1] = $1 OR tags[2] = $1 OR tags[3] = $1 OR tags[4] = $1 OR tags[5] = $1   ", tag)
	if err != nil {
		fmt.Errorf("cant connect to db %s", err)
	}
	//close rows
	defer all.Close()

	products := []Posts{}

	for all.Next() {
		p := Posts{}
		err := all.Scan(&p.ID, &p.Author, &p.Text, &p.Tags, &p.Due)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}
	return products
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
