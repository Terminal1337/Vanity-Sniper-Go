package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"
)

func ReadConfig() (vanity interface{}, token interface{}, serverid interface{}, proxy interface{}) {
	jsonfile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Read Config....")
	defer jsonfile.Close()
	byteresp, _ := ioutil.ReadAll(jsonfile)
	var res map[string]interface{}
	json.Unmarshal([]byte(byteresp), &res)
	return res["vanity"], res["token"], res["serverid"], res["proxy"]
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}
func convert(t interface{}) {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			fmt.Println(s.Index(i))
		}
	}
}
func proxies() []string {
	file, err := os.Open("proxies.txt")
	if err != nil {
		log.Fatalln("Error: " + err.Error())
	}
	defer file.Close()
	var newlines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		newlines = append(newlines, scanner.Text())
	}
	return newlines

}

func sniper() {

	vanity, token, serverid, proxy := ReadConfig()
	jsonStr := map[string]string{"code": token.(string)}
	jsonvalue, _ := json.Marshal(jsonStr)
	str, _ := vanity.(string)
	fmt.Println(vanity)
	proxyUrl, _ := url.Parse("http://" + proxy.(string))
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}, Timeout: time.Second * 15}
	for true {
		resp, err := client.Get("https://discord.com/api/v9/invites/" + str)
		if err != nil {
			fmt.Println("Error: ", err.Error())

		}
		if resp.StatusCode == 404 {
			req, err := http.NewRequest("PATCH", "https://discord.com/api/v9/guilds/"+serverid.(string)+"/"+str, bytes.NewBuffer(jsonvalue))
			if err != nil {
				fmt.Println("Error : ", err.Error())
			}
			req.Header.Set("authorization", token.(string))
			req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
			if req.Response.StatusCode == 200 {
				fmt.Printf("[SUCCESS] Successfully Sniped Vanity %s \n", token.(string))

			}

		} else {
			fmt.Println("[FAILED] Attempt To Claim Failed")
		}
	}
}

func main() {
	fmt.Println("[INPUT] Threads : ")
	var threads int
	fmt.Scan(&threads)
	a := make(chan int)
	for i := 0; i < threads; i++ {

		go func() {

			for {
				sniper()
			}
		}()
	}
	<-a
}
