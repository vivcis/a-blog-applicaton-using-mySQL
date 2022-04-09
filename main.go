package main

import(
	"github.com/gorilla/mux"
	"log"
	"net/http"
)
func main(){
    r := mux.NewRouter()

	//route handlers/ endpoints
	r.HandleFunc("/", getBlogs).Methods("GET")
	r.HandleFunc("/add-post", addFile).Methods("GET")
	r.HandleFunc("/add-post", postPost).Methods("POST")
	r.HandleFunc("/view-post/{id}", viewPost).Methods("GET")
	r.HandleFunc("/edit-post/{id}", editPost).Methods("GET")
	r.HandleFunc("/update-post", updatePost).Methods("POST")
	r.HandleFunc("/delete-post/{id}", deletePost).Methods("GET")

	//handling the css file
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	//to run the server
	log.Println("Listening on :8000...")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		return
	}
}