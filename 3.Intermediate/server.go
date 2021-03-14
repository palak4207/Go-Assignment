package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func dump(v ...interface{}) {
	// spew.Config.MaxDepth = 3
	spew.Config.Indent = "  "
	spew.Config.DisableMethods = true
	spew.Config.SortKeys = true
	spew.Dump(v)
}

func middleware1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Entered 1st Middleware")
		ctx := context.WithValue(r.Context(), ctxKey, "myMiddlewareData")
		next.ServeHTTP(w, r.WithContext(ctx))
		fmt.Println("Exited 1st Middleware")
	})
}

func middleware2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Entered 2nd Middleware")
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Since(t1)
		fmt.Println("Request Time >> ", t2)
		fmt.Println("Exited 2nd Middleware")
	})
}

func ejectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Entered eject Middleware")
		if strings.Contains(r.RequestURI, "eject") {
			fmt.Fprint(w, "This request has been ejected, will not execute any further")
		} else {
			next.ServeHTTP(w, r)
		}
		fmt.Println("Exited eject Middleware")
	})
}

// Context is always request scoped. Dont not use context outside http request
type contextKey string

const ctxKey contextKey = "someUniqueKey"

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func errResponse(w http.ResponseWriter, msg ...string) {
	var response string = "Something did not work nicely"
	if len(msg) > 0 {
		response = msg[0]
	}

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, response)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Upload Handler Called:")

	fmt.Printf("Upload::Request %+v \n", r)
	dump(r)

	// 5 * 2^20 = 5 MB
	r.ParseMultipartForm(5 << 20)
	file, fHeader, err := r.FormFile("myFile")
	if perr, ok := err.(*http.ProtocolError); ok {

		var b []byte
		if b, err = ioutil.ReadAll(r.Body); err != nil {
			errResponse(w, "Upload type not supported")

			return
		}

		fmt.Println(">>>>>>>>>", len(b), "<<<<<<<<<<<<<<<<<")

	} else {
		fmt.Printf("%v", perr)
		errResponse(w, "myFile not found")
	}
	if file != nil {
		defer file.Close()
		f, err := os.OpenFile(filepath.Join(".", fHeader.Filename), os.O_WRONLY|os.O_CREATE, 0660)
		if err != nil {
			// Dont send specific err info back to the client, because it can reveal internal working of the server opening room for mallicious attack

			// I am putting here for demo purpose only
			errResponse(w, "Error while opening file")
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}

	response := map[string]interface{}{
		"Upload": "Success",
	}
	b, err := json.Marshal(response)
	if err != nil {
		errResponse(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprint(w, string(b))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	// Do somthing here

	fmt.Fprint(w, "Login was called")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Root Handler Init")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println("Root Handler Ctx Value >> ", r.Context().Value(ctxKey))

	fmt.Println("Server::Body >> ", string(b))
	fmt.Println("Server::RequestURI >> ", r.RequestURI)
	fmt.Println("Server::URL.Path >> ", r.URL.Path)
	fmt.Println("Server::Cookies >> ", r.Cookies())
	fmt.Println("Server::Method >> ", r.Method)

	// localhost:5555?hello=world&something=here
	r.ParseForm()
	fmt.Println("Server::Form >> ", r.Form)
	fmt.Println("Server::Request::Header >> ", r.Header.Get("Content-Type"))

	var reqMap map[string]interface{}
	json.Unmarshal(b, &reqMap)
	fmt.Println("Server::Request::JSON.Unmarshalled >> ", reqMap)

	var mp = map[string]interface{}{
		"RootHandler": "From server",
		"Route":       r.URL.Path,
	}

	for k, v := range reqMap {
		mp[k] = v
	}

	b, err = json.Marshal(mp)
	if err != nil {
		fmt.Println(err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Server", "My Custom Server")

	if strings.Contains(r.RequestURI, "bad") {
		http.NotFound(w, r)

		return
	} else if strings.Contains(r.RequestURI, "private") {
		redirect(w, r)

		return
	}

	fmt.Fprint(w, string(b))
}

func server1() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/", rootHandler)
}

type myServer struct {
}

func (p *myServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rootHandler(w, r)
}

// In second method we must provide a type which will implement http.Handler
// That is to say, a type with ServeHTTP method on it
func server2() http.Handler {
	http.HandleFunc("/", rootHandler)

	return &myServer{}
}

func server3() http.Handler {
	mux := http.NewServeMux()
	rootMuxHandler := http.HandlerFunc(rootHandler)
	loginMuxHandler := http.HandlerFunc(loginHandler)
	mux.Handle("/", ejectMiddleware(middleware2(middleware1(rootMuxHandler))))
	mux.Handle("/login", ejectMiddleware(middleware2(middleware1(loginMuxHandler))))

	return mux
}

func startServer() {
	var srv http.Handler
	// server1()
	// srv = server2()
	srv = server3()
	err := http.ListenAndServe(":5555", srv)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
