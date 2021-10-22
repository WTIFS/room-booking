package booking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
)


var machines = []string{"10.30.102.147", "10.30.102.148", "10.30.102.149"}

type BucketInfo struct {
	Num  int
	Lock *sync.Mutex
}

func Book() {

	getRoomInfo("1")
	roomNum := BucketInfo{
		Num: 0,
	}
	bookNum := BucketInfo{
		Num: 0,
	}
	fmt.Println(roomNum, bookNum)

}

func getRoomInfo(roomId string) string {
	values := map[string]string{"user_token": "6001ff8445378791b8d8f1524660c983f0dc70c8", "meeting_room_id": roomId}
	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Println(err)
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/api/info", machines[rand.Intn(3)]), "application/json", bytes.NewBuffer(jsonData))
	//resp, err := http.Post("http://127.0.0.1:8082", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))
	return string(content)
}

func book(roomId string) int {
	values := map[string]string{"user_token": "6001ff8445378791b8d8f1524660c983f0dc70c8", "meeting_room_id": roomId}
	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Println(err)
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/api/book", machines[rand.Intn(3)]), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return 400
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	return code
}
