package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Hello struct{}

type RestObject struct {
	Id   int
	Name string
	Age  int
}

func (h Hello) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/html")
	//	response.Header().Set("Content-type", "application/json")
	fmt.Fprint(response, "<h1>GORILLA</h1>Hello! what is your name")
	fmt.Fprint(response, "what is your name!")
}

func Handler(response http.ResponseWriter, request *http.Request) {
	log.Print("handler")
	response.Header().Set("Content-type", "text/html")
	webpage, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(response, fmt.Sprintf("home.html file error %v", err), 500)
	}
	fmt.Fprint(response, string(webpage))
}

func AccountAPIHandler(response http.ResponseWriter, request *http.Request) {
	log.Print("AccountAPIHandler")
	response.Header().Set("Content-type", "text/html")
	webpage, err := ioutil.ReadFile("account.html")
	if err != nil {
		http.Error(response, fmt.Sprintf("home.html file error %v", err), 500)
	}
	fmt.Fprint(response, string(webpage))
}

// Respond to URLs of the form /generic/...
func APIHandler(response http.ResponseWriter, request *http.Request) {
	log.Print("APIhandler")
	//Connect to database
	db, e := sql.Open("mysql", "root:@tcp(localhost:3306)/test")
	if e != nil {
		fmt.Print(e)
	}

	//set mime type to JSON
	response.Header().Set("Content-type", "application/json")

	//	err := request.ParseMultipartForm(20000)
	err := request.ParseForm()
	if err != nil {
		http.Error(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	//can't define dynamic slice in golang
	var result = make([]string, 1000)

	switch request.Method {
	case "GET":

		log.Print("GET")
		log.Print((request.Body))
		//		log.Print(response)
		//		st, err := db.Prepare("select * from users limit 10")
		//		if err != nil {
		//			fmt.Print(err)
		//		}
		//		rows, err := st.Query()
		rows, err := db.Query("select id, name from users")
		if err != nil {
			fmt.Print(err)
		}
		i := 0
		for rows.Next() {
			var name string
			var id int
			err = rows.Scan(&id, &name)
			log.Print(rows.Columns())
			log.Print(id)
			log.Print(name)
			panda := &RestObject{Id: id, Name: name}
			b, err := json.Marshal(panda)
			if err != nil {
				fmt.Println(err)
				return
			}
			result[i] = fmt.Sprintf("%s", string(b))
			i++
		}
		result = result[:i]

	case "POST":
		log.Println(request.Body)
		log.Println("POST")
		name := request.FormValue("name")
		age := request.FormValue("age")
		postName := request.PostFormValue("name")

		ageNumber := 44
		log.Println("name:" + name + " postName: " + postName)
		log.Println(age)
		log.Println(request)
		st, err := db.Prepare("INSERT INTO users(name, age) VALUES(?, ?)")
		//		st, err := db.Prepare("INSERT INTO users(name) VALUES(?)")
		if err != nil {
			fmt.Print(err)
		}
		res, err := st.Exec(name, ageNumber)
		//		res, err := st.Exec(name)
		if err != nil {
			fmt.Print(err)
		}

		if res != nil {
			result[0] = "true"
		}
		result = result[:1]

	case "PUT":
		log.Print("PUT")
		log.Print((request.Body))
		log.Print((request.Form))
		log.Print((request.PostForm))

		name := request.PostFormValue("name")
		id := request.PostFormValue("id")
		log.Print("PUT name:" + name + "-id:" + id + "|")

		st, err := db.Prepare("UPDATE users SET name=? WHERE id=?")
		if err != nil {
			fmt.Print(err)
		}
		res, err := st.Exec(name, id)
		if err != nil {
			fmt.Print(err)
		}

		if res != nil {
			result[0] = "true"
		}
		result = result[:1]
	case "DELETE":
		id := strings.Replace(request.URL.Path, "/api/", "", -1)
		st, err := db.Prepare("DELETE FROM users WHERE id=?")
		if err != nil {
			fmt.Print(err)
		}
		res, err := st.Exec(id)
		if err != nil {
			fmt.Print(err)
		}

		if res != nil {
			result[0] = "true"
		}
		result = result[:1]

	default:
	}

	json, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send the text diagnostics to the client.
	fmt.Fprintf(response, "%v", string(json))
	//fmt.Fprintf(response, " request.URL.Path   '%v'\n", request.Method)
	db.Close()
}

func main() {
	port := 4000

	var err string
	portstring := strconv.Itoa(port)

	mux := http.NewServeMux()
	mux.Handle("/v1/accountholder/", http.HandlerFunc(AccountAPIHandler))

	mux.Handle("/v1/", http.HandlerFunc(APIHandler))
	mux.Handle("/", http.HandlerFunc(Handler))

	log.Print("Listening on port " + portstring + " ... ")
	errs := http.ListenAndServe(":"+portstring, mux)
	if errs != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
