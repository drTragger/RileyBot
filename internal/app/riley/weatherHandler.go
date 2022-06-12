package riley

import (
	"encoding/json"
	"fmt"
	"github.com/drTragger/RileyBot/internal/app/models"
	"github.com/yanzay/tbot/v2"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type openWeather struct {
	Message string `json:"message"`
	Cod     string `json:"cod"`
	Count   int    `json:"count"`
	List    []struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		} `json:"main"`
		Dt   int `json:"dt"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		} `json:"wind"`
		Sys struct {
			Country string `json:"country"`
		} `json:"sys"`
		Rain   interface{} `json:"rain"`
		Snow   interface{} `json:"snow"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Weather []struct {
			Id          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"list"`
}

type weather struct {
	temp        int
	feelsLike   int
	humidity    int
	windSpeed   int
	clouds      int
	city        string
	country     string
	description string
	coordinates string
}

func (b *Bot) cityRequestHandler(m *tbot.Message) {
	handleChatActionError(b.client.SendChatAction(m.Chat.ID, tbot.ActionTyping))
	time.Sleep(500 * time.Millisecond)

	var msg string
	var user *models.User
	user, ok, err := b.storage.User().FindByTelegramUsername(m.From.Username)
	if err != nil {
		b.logger.Info("Error during fetching user data: ", err.Error())
		msg = "Вибачте, тимчасово туплю.\nБудь ласка, спробуйте пізніше.\nА поки можете пограти у Камінь-Ножиці-Папір /play"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}
	if !ok {
		userId, err := strconv.Atoi(m.Chat.ID)
		if err != nil {
			b.logger.Info("Failed to convert user ID ", err.Error())
			msg = "Вибачте, тимчасово туплю.\nБудь ласка, спробуйте пізніше.\nА поки можете пограти у Камінь-Нножиці-Папір /play"
			handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
			return
		}

		user = &models.User{Username: m.From.Username, TelegramId: &userId}
		err = b.storage.User().Create(user)
		if err != nil {
			b.logger.Info("Failed to create new user: ", err.Error())
			msg = "Вибачте, тимчасово туплю.\nБудь ласка, спробуйте пізніше.\nА поки можете пограти у Камінь-Ножиці-Папір /play"
			handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
			return
		}
	}
	err = b.storage.Dialog().Create(&models.Dialog{Name: "weather", UserId: user.ID, Status: true})
	if err != nil {
		b.logger.Info("Failed to create new dialog: ", err.Error())
		msg = "Вибачте, тимчасово туплю.\nБудь ласка, спробуйте пізніше.\nА поки можете пограти у Камінь-Ножиці-Папір /play"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}

	msg = "Напишіть мені назву міста, у якому хочете дізнатись погоду"

	b.LogHandler(m, msg)
	handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
}

func (b *Bot) weatherHandler(m *tbot.Message) {
	handleChatActionError(b.client.SendChatAction(m.Chat.ID, tbot.ActionTyping))

	user, ok, err := b.storage.User().FindByTelegramUsername(m.From.Username)
	var msg string
	if err != nil {
		b.logger.Info("Error during fetching user data: ", err.Error())
		msg = "Вибачте, тимчасово туплю.\nБудь ласка, спробуйте пізніше.\nА поки можете пограти у Камінь-Ножиці-Папір /play"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}

	if !ok {
		b.logger.Info("User and dialog not found")
		msg = "Будь ласка, запустіть мене, виконавши команду /start"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}

	dialog, ok, err := b.storage.Dialog().FindLatestUserDialog(user.ID, "weather")
	if err != nil {
		b.logger.Error("Error during fetching dialog data: ", err.Error())
		msg = "Вибачте, тимчасово туплю.\nБудь ласка, спробуйте пізніше.\nА поки можете пограти у Камінь-Ножиці-Папір /play"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}

	if !ok || dialog.Status != true {
		b.logger.Info("No active dialog status")
		msg = "Перепрошую, я поки не вмію розпізнавати такі повідомлення. Спробуйте:\n\n/play - Пограти у Камінь-Ножиці-Папір\n\n/weather - Дізнатись, яка зараз погода"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}

	ow, err, response := getWeatherData(b.config.WeatherKey, m.Text)
	fmt.Printf("%+v\n", ow)

	if err != nil {
		b.logger.Errorf("Error during unmarshalling weather JSON: %s\nResponse: %s", err.Error(), response)
		msg = "Вибачте, тимчасово туплю.\nНе можу обробити дані про погоду.\nБудь ласка, спробуйте пізніше.\nА поки можете пограти у Камінь-Ножиці-Папір /play"
	} else {
		if ow.Count < 1 {
			handleMessageError(b.client.SendMessage(m.Chat.ID, "Хмм...🤔\nНе чув про таке місто.\nСпробуйте інше."))
			return
		}
		b.logger.Info(ow)
		for i := 0; i < ow.Count; i++ {
			if ow.List[i].Sys.Country == "RU" {
				msg = "Фу, я не буду лізти у той смітник щоб дізнатись погоду"
				break
			}
			var weatherDescription string
			for k, val := range ow.List[i].Weather {
				weatherDescription += getWeatherDescription(val.Description, val.Icon)
				if k < len(ow.List[i].Weather)-1 {
					weatherDescription += "\n"
				}
			}

			weatherData := weather{
				temp:        int(math.Round(ow.List[i].Main.Temp)),
				feelsLike:   int(math.Round(ow.List[i].Main.FeelsLike)),
				humidity:    ow.List[i].Main.Humidity,
				windSpeed:   int(math.Round(ow.List[i].Wind.Speed)),
				city:        ow.List[i].Name,
				country:     ow.List[i].Sys.Country,
				description: weatherDescription,
				clouds:      ow.List[i].Clouds.All,
				coordinates: "",
			}

			if ow.Count > 2 {
				weatherData.coordinates = fmt.Sprintf("Координаты: %f, %f\n", ow.List[i].Coord.Lat, ow.List[i].Coord.Lon)
			}

			msg += fmt.Sprintf(""+
				"Місто/Країна: %s [%s]\n%s"+
				"%s\n\n"+
				"Температура🌡: %d°C\n"+
				"Відчувається як🌡: %d°C\n\n"+
				"Вологість💧: %d%%\n"+
				"Швидкість вітру💨: %d м/с\n"+
				"Хмарність: %d%% %s\n",
				weatherData.city, weatherData.country, weatherData.coordinates, weatherData.description, weatherData.temp,
				weatherData.feelsLike, weatherData.humidity, weatherData.windSpeed, weatherData.clouds, getCloudsEmoji(weatherData.clouds),
			)
			if i < ow.Count-1 {
				msg += getCitiesDelimiter(&weatherData)
			}
		}
	}

	if err := b.storage.Dialog().UpdateStatus(dialog.ID); err != nil {
		b.logger.Error("Error during updating dialog status: ", err.Error())
	}

	b.LogHandler(m, msg)
	handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
}

func getWeatherData(weatherKey string, requestCity string) (*openWeather, error, string) {
	url := "https://community-open-weather-map.p.rapidapi.com/find?q=" + strings.ReplaceAll(strings.TrimSpace(requestCity), " ", "+") + "&lang=ua&mode=null&lon=0&type=link%2C%20accurate&lat=0&units=metric"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-host", "community-open-weather-map.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", weatherKey)

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println("Failed to close HTTP connection")
		}
	}(res.Body)

	body, _ := ioutil.ReadAll(res.Body)

	w := openWeather{}

	if err := json.Unmarshal(body, &w); err != nil {
		return nil, err, string(body)
	}
	return &w, nil, ""
}

func getWeatherDescription(description string, icon string) string {
	var emoji string
	switch icon {
	case "01d":
		fallthrough
	case "01n":
		emoji = "☀"
	case "02d":
		fallthrough
	case "02n":
		emoji = "🌤"
	case "03d":
		fallthrough
	case "03n":
		emoji = "⛅"
	case "04d":
		fallthrough
	case "04n":
		emoji = "☁"
	case "09d":
		fallthrough
	case "09n":
		emoji = "🌧"
	case "10d":
		fallthrough
	case "10n":
		emoji = "🌦"
	case "11d":
		fallthrough
	case "11n":
		emoji = "⛈"
	case "13d":
		fallthrough
	case "13n":
		emoji = "❄"
	case "50d":
		fallthrough
	case "50n":
		emoji = "🌫"
	}

	return strings.Title(description) + " " + emoji
}

func getCloudsEmoji(clouds int) string {
	var emoji string
	if clouds >= 0 && clouds < 26 {
		emoji = "☀"
	} else if clouds > 1 && clouds < 25 {
		emoji = "🌤"
	} else if clouds > 25 && clouds < 51 {
		emoji = "⛅"
	} else if clouds > 50 && clouds < 76 {
		emoji = "🌥"
	} else {
		emoji = "☁"
	}

	return emoji
}

func getCitiesDelimiter(w *weather) string {
	delimiter := "\n"
	var lineLen int
	if len(w.coordinates) == 0 {
		lineLen = utf8.RuneCountInString(w.country) + utf8.RuneCountInString(w.city) + 16
	} else {
		lineLen = utf8.RuneCountInString(w.coordinates) - 2
	}
	delimiter += fmt.Sprintf("%s", strings.Repeat("=", lineLen))
	delimiter += "\n\n"

	return delimiter
}
