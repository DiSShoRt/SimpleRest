package main

import (
	poststore "SimpleRest/store"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"time"
)

type postStore struct {
	store *poststore.PostStore
}

func NewPostServer() *postStore {
	store := poststore.New()
	return &postStore{store: store}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ps *postStore) createPostHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task create at %s\n", req.URL.Path)

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.
	type RequestPost struct {
		ID     int       `json:"id"`
		Text   string    `json:"text"`
		Author string    `json:"author"`
		Tags   []string  `json:"tags"`
		Due    time.Time `json:"due"`
	}

	type ResponseId struct {
		ID int `json:"id"`
	}

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestPost
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ps.store.AddPostToDb(rt.Author, rt.Text, rt.Tags, rt.Due)
	rt.ID = id
	fmt.Println(rt.Text, rt.Tags, rt.Due, ps.store)
	renderJSON(w, rt)
}

func (ps *postStore) getPostsByAuthor(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all tasks at %s\n", req.URL.Path)

	author := mux.Vars(req)["author"]

	fmt.Println(author)
	allPosts := ps.store.GetPostsByAuthorDb(author)

	js, err := json.Marshal(allPosts)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ps *postStore) getAllPostsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all tasks at %s\n", req.URL.Path)

	allTasks := ps.store.GetAllPostsDb()

	fmt.Println(allTasks)
	renderJSON(w, allTasks)
}

func (ps *postStore) getPostHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get post at %s\n", req.URL.Path)

	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	task, err := ps.store.GetPostByIdDb(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	task.Tags = []string(task.Tags)

	renderJSON(w, task)
}

func (ps *postStore) deletePostHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete post at %s\n", req.URL.Path)
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	err := ps.store.DeletePostFromDb(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (ps *postStore) deleteAllPostsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete all posts at %s\n", req.URL.Path)
	ps.store.DeleteAllPostsFromDb()
}

func (ps *postStore) tagHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling posts by tag at %s\n", req.URL.Path)

	tag := mux.Vars(req)["tag"]

	tasks := ps.store.GetPostsByTagDb(tag)
	renderJSON(w, tasks)
}

func (ps *postStore) dueHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling posts by due at %s\n", req.URL.Path)

	vars := mux.Vars(req)
	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}

	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])
	if month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}
	day, _ := strconv.Atoi(vars["day"])
	tasks := ps.store.GetPostByDue(year, time.Month(month), day)
	renderJSON(w, tasks)
}

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := NewPostServer()

	router.HandleFunc("/post/", server.createPostHandler).Methods("POST")
	router.HandleFunc("/post/", server.getAllPostsHandler).Methods("GET")
	router.HandleFunc("/post/", server.deleteAllPostsHandler).Methods("DELETE")
	router.HandleFunc("/post/{id:[0-9]+}/", server.getPostHandler).Methods("GET")
	router.HandleFunc("/post/{id:[0-9]+}/", server.deletePostHandler).Methods("DELETE")
	router.HandleFunc("/tag/{tag}/", server.tagHandler).Methods("GET")
	router.HandleFunc("/author/{author}/", server.getPostsByAuthor).Methods("GET")
	router.HandleFunc("/due/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}/", server.dueHandler).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:"+"8080", router))
}
