// Package main provides ...
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gtldhawalgandhi/go-training/3.Intermediate/api"
	"github.com/gtldhawalgandhi/go-training/3.Intermediate/db"
	l "github.com/gtldhawalgandhi/go-training/3.Intermediate/logger"
	"github.com/gtldhawalgandhi/go-training/3.Intermediate/util"
	"github.com/jackc/pgx/v4"
)

func helo() {
	l.D("Helo debug")
	l.F("My Fatal example")
}

func simpleGETClient() {

	myClient := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := myClient.Get("http://127.0.0.1:5555/")
	if err != nil {
		// Try figure what error type it is and take decisions based on that
		// Dont throw fatal errors if not required
		l.F(err)
	}

	// var nB = 20 (read 20 bytes)
	// var nB = res.ContentLength
	// var p = make([]byte, nB)
	// n, err := res.Body.Read(p)

	// if err != nil {
	// 	l.F(err)
	// }
	// defer res.Body.Close()
	// fmt.Println("Bytes read", n)
	// fmt.Println("Bytes read data", string(p))

	fmt.Println(res.Header.Get("Content-Type"))

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		l.F(err)
	}
	defer res.Body.Close()
	fmt.Println(string(data))

	var mp = make(map[string]interface{})
	json.Unmarshal(data, &mp)
	fmt.Println(mp)
}

func simplePOSTRequest() {

	myClient := http.Client{
		Timeout: 5 * time.Second,
	}
	data := map[string]interface{}{
		"Hello": "Message from the Client",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		l.F(err)
	}

	payload := bytes.NewReader(jsonData)

	req, err := http.NewRequest("POST", "http://127.0.0.1:5555", payload)
	// jsonData2, err := json.Marshal(data)
	if err != nil {
		l.F(err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := myClient.Do(req)
	if err != nil {
		l.F(err)
	}

	readerData, err := ioutil.ReadAll(res.Body)

	if err != nil {
		l.F(err)
	}
	defer res.Body.Close()
	fmt.Println("Client >> Response", string(readerData))
}

func testServer(config util.Config) {
	var err error

	bgCtx := context.Background()

	conn, err := pgx.Connect(bgCtx, config.DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(bgCtx)

	store := db.NewPGStore(conn)

	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("Failed to create server", err)
	}

	err = server.Start(":5555")
	if err != nil {
		log.Fatal("Failed to start server", err)
	}
}

func testDB(config util.Config) {
	var err error

	conn, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	store := db.NewPGStore(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// ur, err := store.GetUserByEmail(context.Background(), "a@b.com")
	ur, err := store.UpdateUser(ctx, db.UserRequest{
		FirstName: "aa",
		LastName:  "bb",
		UserName:  "aabb",
		Email:     "aa@bb.com",
	})
	b, _ := json.Marshal(ur)
	fmt.Printf("%+v %v \n", ur, err)
	fmt.Printf("%v \n", string(b))

	fmt.Println(store.GetUsers(ctx))
}

// -trimpath will cut short file name everywhere in our code when displaying
// go build -trimpath -o app && ./app
// OR
// go run -trimpath .
func main() {
	// l.SetFileLogger("myLog.txt", l.TRACE)
	// defer l.CleanUp()
	// l.I("Entering main file")

	// l.T("My Trace data")
	// l.E("My error")
	// helo()

	// go func() {
	// 	time.Sleep(200 * time.Millisecond)
	// 	// simpleGETClient()
	// 	simplePOSTRequest()
	// }()
	// startServer()

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	fmt.Printf("%+v", config)

	// testDB(config)
	testServer(config)

}
