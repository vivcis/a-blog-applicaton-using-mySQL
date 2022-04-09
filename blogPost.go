package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Blog struct {
	ID          string
	Title       string
	Ingredients string
	Content     string
	Time        string
	AuthorName  string
}


var blog []Blog

func init() {
	//mock data - implement DB
	// blog = append(blog, Blog{
	// 	ID:          "1",
	// 	Title:       "Pilaf Rice",
	// 	Ingredients: "3/4 cup unsalted raw cashews\n1 tablespoon extra-virgin olive oil\n3 finely chopped, plus 2 whole small garlic cloves\n3/4 pound (about 18 spears), ends trimmed, cut into 2-inch pieces",
	// 	Content: "How To Make a Simple Rice Pilaf:" + " " +
	// 		"Place the rice in a strainer and rinse it thoroughly under cool water. The water running through the rice will look milky at first, " +
	// 		"but will then become clearer and only lightly clouded. It's fine if there's still some haze in the water. There is no need to dry the rice before cooking; " +
	// 		"a bit of moisture on the rice is fine. " +
	// 		"Set the strainer of rice aside while you cook the onion.",
	// 	Time:       time.Now().Format(time.RFC822),
	// 	AuthorName: "Sidney Sheldon",
	// })
	// blog = append(blog, Blog{
	// 	ID:          uuid.NewString(),
	// 	Title:       "Concoction Rice",
	// 	Ingredients: "3/4 cup unsalted raw cashews\n1 tablespoon extra-virgin olive oil\n3 finely chopped, plus 2 whole small garlic cloves\n3/4 pound (about 18 spears), ends trimmed, cut into 2-inch pieces",
	// 	Content:     "In a small saucepan, cover the cashews with water to a depth of about 2 inches. Bring to a boil over medium-high heat. Remove from the heat and let soak for at least 10 minutes. Alternatively, soak the cashews in room temperature water overnight in the refrigerator.",
	// 	Time:        time.Now().Format(time.RFC822),
	// 	AuthorName:  "Cecilia",
	// })
	// blog = append(blog, Blog{
	// 	ID:          uuid.NewString(),
	// 	Title:       "Gizzdodo",
	// 	Ingredients: "3 tablespoons fresh lemon juice (from about 1 large lemon), divided\n3/4 cup unsweetened non-dairy milk\n1 tablespoon white miso\n3/4 teaspoon onion powder\n1/2 (14-ounce) can quartered artichoke hearts, drained, and coarsely chopped",
	// 	Content: "In a large skillet over medium heat, warm the oil. Cook the chopped garlic, stirring constantly, for about 30 seconds, until fragrant. Add the asparagus and cook, stirring constantly, for about 2 minute more, until bright green; season with salt and pepper." +
	// 		"Transfer the drained cashews to a blender. Add the milk, miso, onion powder, 2 whole garlic cloves, the remaining 2 tablespoons of the lemon juice, and ¼ cup water. Blend on high speed, adding more water if needed, until completely smooth. ",
	// 	Time:       time.Now().Format(time.RFC822),
	// 	AuthorName: "Lovey",
	// })
}


//-------------------------GET/VIEW A POST-----------------------------------------------
func getBlogs(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")

	if err != nil {
		fmt.Println(err)
		// handle the error properly
		return
	}
	//err = t.Execute(w, blog)
	//if err != nil{
	//	log.Println(err)
	//	return
	//}
	// DB.Query("SELECT * FROM BLOG")
	rows, err := DB.Query("SELECT * FROM BLOG")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	// define the variable
	for rows.Next() {
		var b Blog
		err = rows.Scan(&b.ID, &b.Title, &b.Time, &b.AuthorName, &b.Ingredients, &b.Content)
		if err != nil {
			log.Println(err)
			return
		}
		blog = append(blog, b)
	}

	err = t.ExecuteTemplate(w, "index.html", blog)
	if err != nil {
		log.Println(err)
		return
	}
}

func viewPost(w http.ResponseWriter, r *http.Request) {
	blogInstance := Blog{}
	vars := mux.Vars(r)
	id := vars["id"]
	for _, v := range blog {
		if id == v.ID {
			blogInstance = v
		}
	}
	// TODO: SELECT * FROM blog where id = ?
	t, err := template.ParseFiles("templates/view.html")

	if err != nil {
		fmt.Println(err)
	}
	
	_ = t.Execute(w, blogInstance)
}


//-------------------------CREATE A POST-----------------------------------------------
func postPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.PostFormValue("title")
	author := r.PostFormValue("author")
	ingredients := r.PostFormValue("ingredients")
	content := r.PostFormValue("content")

	blogPost := Blog{
		ID:          uuid.NewString(),
		Title:       title,
		Ingredients: ingredients,
		Content:     content,
		Time:        time.Now().Format(time.RFC822),
		AuthorName:  author,
	}

	
	//blog = append(blog, blogPost)
	q, err := DB.Prepare("INSERT INTO BLOG (id, title, ingredients, content, time, author_name) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("%s", err)
		// handle errors properly
		return
	}
	_, err = q.Exec(blogPost.ID, blogPost.Title, blogPost.Ingredients, blogPost.Content, time.Now().Format(time.RFC822), blogPost.AuthorName)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func addFile(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/add.html")

	if err != nil {
		fmt.Println(err)
	}
	_ = t.Execute(w, nil)
}


//-------------------------EDIT A POST HANDLER-----------------------------------------------
func editPost(w http.ResponseWriter, r *http.Request) {
	editInstance := Blog{}
	vars := mux.Vars(r)
	id := vars["id"]
	for _, v := range blog {
		if id == v.ID {
			editInstance = v
		}
	}
	t, err := template.ParseFiles("templates/edit.html")

	if err != nil {
		fmt.Println(err)
	}
	_ = t.Execute(w, editInstance)
}

//-------------------------UPDATE A POST-----------------------------------------------
func updatePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostFormValue("id")
	title := r.PostFormValue("title")
	author := r.PostFormValue("author")
	ingredients := r.PostFormValue("ingredients")
	content := r.PostFormValue("content")

	//fmt.Println(ingredients, content)

	postUpdate := Blog{
		ID:          id,
		Title:       title,
		Content:     content,
		Ingredients: ingredients,
		Time:        time.Now().Format(time.RFC822),
		AuthorName:  author,
	}
	// for i, v := range blog {
	// 	if id == v.ID {
	// 		blog[i] = postUpdate
	// 	}
	// }

	for i, v := range blog {
		if id == v.ID {
			blog = append(blog[:i], blog[i+1:]...)
		}
	}
	q, err := DB.Prepare("UPDATE BLOG SET title=?, ingredients=?, content=?, author_name=? WHERE id=?")
	if err != nil {
		log.Printf("%s", err)
		return
	}
	_, err = q.Exec(postUpdate.Title, postUpdate.Ingredients, postUpdate.Content, postUpdate.AuthorName, id)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//----------------------------------DELETE A POST---------------------------------------------
func deletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for i, v := range blog {
		if id == v.ID {
			blog = append(blog[:i], blog[i+1:]...)
		}
	}

	q, err := DB.Prepare("DELETE FROM BLOG WHERE id=?")
	if err != nil {
		log.Printf("%s", err)
		return
	}
	_, err = q.Exec(id)
	if err != nil {
		log.Printf("%s", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
