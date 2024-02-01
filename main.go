package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type CityInfo struct {
	Temperature string  `json:"temperature"`
	FeelsLike   string  `json:"feels_like"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
}

type Reminder struct {
	User    string
	Message string
	Timer   time.Duration
}

func readTokenFromFile(filename string) (string, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(fileContent), nil
}

func main() {
	botToken, err := readTokenFromFile("token.txt")
	if err != nil {
		log.Fatal("error reading the file with bot's token")
		return
	}
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection to Discord:", err)
		return
	}

	fmt.Println("Bot is now running. Press Ctrl+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "!start" {
		go func() {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hello, %s!", m.Author.Username))
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I am a bot. To see what i can do, type '!help'"))
		}()
		return
	}

	if m.Content == "!help" {
		go func() {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'!weather [name of your city]' to see the current weather in your city"))
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'!define [word]' to look up a word's meaning in english"))
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'!reminder [hh:mm:ss] [message] to set a timer with a message. I will remind you with the message at the specified moment"))
		}()
		return
	}

	if strings.HasPrefix(m.Content, "!weather") {
		go func() {
			lng, lat := fetchCoordinates(m.Content[9:])
			if lng == "" || lat == "" {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I can't find a city with that name. Please make sure you typed the name correctly"))
				return
			}
			city, err := fetchWeather(lng, lat)
			if err != nil {
				log.Fatal("error fetching weather data, please try again")
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Temperature: "+city.Temperature+" celcuis")
			s.ChannelMessageSend(m.ChannelID, "Feels like: "+city.FeelsLike+" celcius")
			s.ChannelMessageSend(m.ChannelID, "Description: "+city.Description)
			s.ChannelMessageSend(m.ChannelID, "Windspeed: "+strconv.FormatFloat(city.WindSpeed, 'f', 2, 64)+" m/s")
			s.ChannelMessageSend(m.ChannelID, "Humidity: "+strconv.Itoa(city.Humidity))
		}()
		return
	}

	if strings.HasPrefix(m.Content, "!reminder") {
		go func() {
			if len(m.Content) < 19 {
				s.ChannelMessageSend(m.ChannelID, "wrong format. Do it in the following way: !reminder hh:mm:ss [message]")
				return
			}
			hour, _ := strconv.Atoi(m.Content[10:12])
			minute, _ := strconv.Atoi(m.Content[13:15])
			second, _ := strconv.Atoi(m.Content[16:18])
			message := m.Content[19:]
			untilAction := time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute + time.Duration(second)*time.Second
			time.AfterFunc(untilAction, func() {
				s.ChannelMessageSend(m.ChannelID, "timer invoked with your message: "+message)
			})
			s.ChannelMessageSend(m.ChannelID, "the reminder is set!")
			fmt.Print(untilAction)
		}()
		return
	}
	if strings.HasPrefix(m.Content, "!define ") {
		go func() {
			content := strings.Split(m.Content, " ")
			if len(content) > 2 {
				s.ChannelMessageSend(m.ChannelID, "You can translate only one word at a time")
				return
			}
			s.ChannelMessageSend(m.ChannelID, defineWord(content[1]))
		}()
		return
	}
	s.ChannelMessageSend(m.ChannelID, "There is no such command. Type '!help' to see available commands")
}

func fetchCoordinates(city string) (lng string, lat string) {
	apiKey := "2ydJz6WD+XG2joDa4LdGhA==mm8kWNd3mVgg5xKs"
	url := "https://api.api-ninjas.com/v1/city?name=" + city
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("X-API-KEY", apiKey)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Error making request: ", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if len(body) < 10 {
		return lng, lat
	}
	if err != nil {
		log.Fatal("error reading body: ", err)
	}
	cityInfo := string(body)
	var arr []string = strings.Split(cityInfo, ",")
	lat = arr[1][13:] // в json на вторых и третьих индексах находятся долгота и широта с которыми я иду в другой API для просмотра погоды
	lng = arr[2][14:]
	return lng, lat
}

func fetchWeather(lng, lat string) (CityInfo, error) {
	apiKey := "5358b6e76276fad90bfc3250e7d81827"
	url := "https://api.openweathermap.org/data/2.5/weather?" + "lat=" + lat + "&lon=" + lng + "&appid=" + apiKey
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal("error creating a request ", err)
	}
	response, err := client.Do(request)

	if err != nil {
		log.Fatal("error making the request ", err)
	}
	defer response.Body.Close()

	var weatherData map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&weatherData)
	if err != nil {
		return CityInfo{}, err
	}

	mainData := weatherData["main"].(map[string]interface{})
	windData := weatherData["wind"].(map[string]interface{})
	weatherArray := weatherData["weather"].([]interface{})
	weatherInfo := weatherArray[0].(map[string]interface{})

	cityInfo := CityInfo{
		Temperature: fmt.Sprintf("%.2f", mainData["temp"].(float64)-273),
		FeelsLike:   fmt.Sprintf("%.2f", mainData["feels_like"].(float64)-273),
		Description: weatherInfo["description"].(string),
		Humidity:    int(mainData["humidity"].(float64)),
		WindSpeed:   windData["speed"].(float64),
	}
	return cityInfo, nil
}

func defineWord(word string) string {
	apiKey := "2ydJz6WD+XG2joDa4LdGhA==mm8kWNd3mVgg5xKs"
	url := "https://api.api-ninjas.com/v1/dictionary?word=" + word

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("X-API-KEY", apiKey)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Error making request: ", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	fmt.Print(string(body))
	valid := strings.Split(string(body), "valid")[1]
	fmt.Println(valid)
	if strings.Contains(valid, "false") {
		return "This is not a valid english word"
	}
	fmt.Print(string(body))
	translation := strings.Split(string(body), "1.")[1]
	translation = strings.Split(translation, ".")[0]
	return translation
}
