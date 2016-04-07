package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UserRestObject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Respond to URLs of the form /v1/...
func APIHandler(response http.ResponseWriter, request *http.Request) {
	log.Print("APIhandler")
	//Connect to database
	db, e := sql.Open("mysql", "root:@tcp(localhost:3306)/test")

	defer db.Close()

	if e != nil {
		fmt.Print(e)
	}

	//set mime type to JSON
	response.Header().Set("Content-type", "application/json")

	err := request.ParseForm()
	if err != nil {
		http.Error(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	// var jsonResult []string
	
    switch request.Method {
	case "GET":
		userId := request.FormValue("userId")
        
		if userId == "" {
	        var users []UserRestObject
			log.Println("userId", userId)
			log.Println("GET ALL")
			rows, err := db.Query("select id, name, age from users limit 100")
			if err != nil {
				fmt.Print(err)
			}

			defer rows.Close()
			for rows.Next() {
				var name string
				var age int
				var id int
				err = rows.Scan(&id, &name, &age)
                if err != nil {
                    fmt.Println(err)
                    return
                }
				user := UserRestObject{Id: id, Name: name, Age: age}
			    users = append(users, user)
			}
            jsonResult, err := json.Marshal(users)
            if err != nil {
                fmt.Println(err)
                return
            }
            response.WriteHeader(http.StatusOK)
            fmt.Fprintf(response, "%v", string(jsonResult))
            return
		} else {
			log.Println("userId", userId)
			log.Println("GET search")
			st, err := db.Prepare("select id, name, age from users WHERE id = ? ")
			defer st.Close()
			if err != nil {
				log.Fatal("error preparing statement:", err)
			}

			var name string
			var age int
			var id int

			err = st.QueryRow(userId).Scan(&id, &name, &age)

			if err != nil {
				log.Fatal(err)
			}
			user := &UserRestObject{Id: id, Name: name, Age: age}
			jsonResult, err := json.Marshal(user)
			if err != nil {
				fmt.Println(err)
				return
			}
            fmt.Fprintf(response, "%v", string(jsonResult))            
            return
		}

	case "POST":
		log.Println("POST")
		var userRestObject UserRestObject
		userRestObject = jsonToUserObject(response, request)
		log.Println("name:" + userRestObject.Name + " age: " + strconv.Itoa(userRestObject.Age))
		st, err := db.Prepare("INSERT INTO users(name, age) VALUES(?, ?)")
		if err != nil {
			fmt.Print(err)
		}
		result, err := st.Exec(userRestObject.Name, userRestObject.Age)
		if err != nil {
			fmt.Print(err)
		}
        
        jsonResult, err := json.Marshal(result)
        if err != nil {
            fmt.Println(err)
            return
        }

        fmt.Fprintf(response, "%v", string(jsonResult))
        return

	case "PUT":
		id := strings.Replace(request.URL.Path, "/v1/", "", -1)
		log.Print("DELETE  id:" + id + "")
		var userRestObject UserRestObject
		userRestObject = jsonToUserObject(response, request)

		name := userRestObject.Name
		age := userRestObject.Age
		log.Println("PUT name:" + name + "-id:" + id + " age:" + strconv.Itoa(age) + "")

		st, err := db.Prepare("UPDATE users SET name=?, age=? WHERE id=?")
		if err != nil {
			fmt.Print(err)
		}
		result, err := st.Exec(name, age, id)
		if err != nil {
			fmt.Print(err)
		}
        
        jsonResult, err := json.Marshal(result)
        if err != nil {
            fmt.Println(err)
            return
        }
        
        fmt.Fprintf(response, "%v", string(jsonResult))
        return

	case "DELETE":
		id := strings.Replace(request.URL.Path, "/v1/", "", -1)
		log.Print("DELETE  id:" + id + "")
		st, err := db.Prepare("DELETE FROM users WHERE id=?")
		if err != nil {
			fmt.Print(err)
		}
		result, err := st.Exec(id)
		if err != nil {
			fmt.Print(err)
		}

		jsonResult, err := json.Marshal(result)
        if err != nil {
            fmt.Println(err)
            return
        }
        
        fmt.Fprintf(response, "%v", string(jsonResult))
        return
	    default:
	}
}

func jsonToUserObject(response http.ResponseWriter, request *http.Request) UserRestObject {
	var userRestObject UserRestObject
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := request.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &userRestObject); err != nil {
		response.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(response).Encode(err); err != nil {
			panic(err)
		}
	}
	return userRestObject
}

func main() {
	port := 4000

	var err string
	portstring := strconv.Itoa(port)

	mux := http.NewServeMux()

	mux.Handle("/v1/", http.HandlerFunc(APIHandler))
	//	mux.Handle("/", http.HandlerFunc(Handler))

	log.Print("Listening on port " + portstring + " ... ")
	errs := http.ListenAndServe(":"+portstring, mux)
	if errs != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
