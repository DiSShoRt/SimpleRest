package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	poststore "SimpleRest/store"
)

type postStore struct {
	store *poststore.PostStore
}

func NewPostServer() *postStore {
	store := poststore.New()
	return &postStore{store: store}
}

func (ps *postStore) postHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/post/" {
		// Request is plain "/post/", without trailing ID.
		if req.Method == http.MethodPost {
			ps.createPostHandler(w, req)
		} else if req.Method == http.MethodGet {
			ps.getAllPostsHandler(w, req)
		} else if req.Method == http.MethodDelete {
			ps.deleteAllPostsHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /task/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		// Request has an ID, as in "/post/<id>".
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "expect /post/<id> in task handler", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodDelete {
			ps.deletePostHandler(w, req, int(id))
		} else if req.Method == http.MethodGet {
			ps.getPostHandler(w, req, int(id))
		} else {
			http.Error(w, fmt.Sprintf("expect method GET or DELETE at /task/<id>, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (ps *postStore) createPostHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task create at %s\n", req.URL.Path)

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.
	type RequestPost struct {
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

	id := ps.store.AddPostToDb(rt.Text, rt.Author, rt.Tags, rt.Due)
	fmt.Println(rt.Text, rt.Tags, rt.Due, ps.store)
	js, err := json.Marshal(ResponseId{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ps *postStore) getAllPostsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all tasks at %s\n", req.URL.Path)

	allTasks := ps.store.GetAllPostsDb()

	fmt.Println(allTasks)
	js, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ps *postStore) getPostHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handling get post at %s\n", req.URL.Path)

	task, err := ps.store.GetPostByIdDb(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	task.Tags = []string(task.Tags)

	js, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ps *postStore) deletePostHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handling delete post at %s\n", req.URL.Path)

	err := ps.store.DeletePostFromDb(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ps *postStore) deleteAllPostsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete all posts at %s\n", req.URL.Path)
	ps.store.DeleteAllPostsFromDb()
}

func (ps *postStore) tagHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling posts by tag at %s\n", req.URL.Path)

	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /tag/<tag>, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /tag/<tag> path", http.StatusBadRequest)
		return
	}
	tag := pathParts[1]

	tasks := ps.store.GetPostsByTagDb(tag)
	js, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ps *postStore) dueHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling posts by due at %s\n", req.URL.Path)

	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /due/<date>, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}
	if len(pathParts) != 4 {
		badRequestError()
		return
	}

	year, err := strconv.Atoi(pathParts[1])
	if err != nil {
		badRequestError()
		return
	}
	month, err := strconv.Atoi(pathParts[2])
	if err != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}
	day, err := strconv.Atoi(pathParts[3])
	if err != nil {
		badRequestError()
		return
	}

	tasks := ps.store.GetPostByDue(year, time.Month(month), day)
	js, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	mux := http.NewServeMux()
	server := NewPostServer()
	mux.HandleFunc("/post/", server.postHandler)
	mux.HandleFunc("/tag/", server.tagHandler)
	mux.HandleFunc("/due/", server.dueHandler)

	log.Fatal(http.ListenAndServe("localhost:"+"8080", mux))
}
