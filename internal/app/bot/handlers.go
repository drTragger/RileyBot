package riley

import (
	"encoding/json"
	"fmt"
	"github.com/drTragger/rileyBot/internal/app/models"
	"github.com/yanzay/tbot/v2"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var (
	picks   = []string{"rock", "paper", "scissors"}
	options = map[string]string{
		"paper":    "rock",
		"scissors": "paper",
		"rock":     "scissors",
	}
	translations = map[string]string{"scissors": "ножницы✂", "rock": "камень🗿", "paper": "бумагу\U0001F9FB"}
)

func handleChatActionError(err error) {
	if err != nil {
		log.Println("Error chat action: ", err.Error())
	}
}

func handleMessageError(message *tbot.Message, err error) {
	if err != nil {
		log.Printf("Message: %s\nError: %s", message.Text, err.Error())
	}
}

func (b *Bot) StartHandler(m *tbot.Message) {
	handleChatActionError(b.client.SendChatAction(m.Chat.ID, tbot.ActionTyping))
	time.Sleep(1 * time.Second)
	var stdMessage = "Привет, я бот Райли🖖\n\n/play\tКамень-Ножницы-Бумага\n\n/weather\tУзнать прогноз погоды"
	b.LogHandler(m)
	var msg string
	userId, err := strconv.Atoi(m.Chat.ID)
	if err != nil {
		b.logger.Info("Failed to convert user ID ", err.Error())
	}
	userExists, err := b.storage.User().UserExists(m.From.Username)
	if err != nil {
		b.logger.Info("Failed to find user: ", err.Error())
	}
	user, ok, err := b.storage.User().FindByTelegramUsernameWithGreetings(m.From.Username)
	if err != nil {
		b.logger.Info("Failed to find user: ", err.Error())
	}
	if ok {
		msg = user.Greeting.Message
	} else if !ok && !userExists {
		err = b.storage.User().Create(&models.User{Username: m.From.Username, TelegramId: &userId})
		if err != nil {
			b.logger.Info("Failed to create new user: ", err.Error())
		}
		msg = stdMessage
	} else {
		msg = stdMessage
	}

	handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
}

func (b *Bot) PlayHandler(m *tbot.Message) {
	b.LogHandler(m)
	handleChatActionError(b.client.SendChatAction(m.Chat.ID, tbot.ActionTyping))
	time.Sleep(500 * time.Millisecond)
	buttons := makeButtons()
	handleMessageError(b.client.SendMessage(m.Chat.ID, "Твой ход:", tbot.OptInlineKeyboardMarkup(buttons)))
}

func (b *Bot) CallbackHandler(cq *tbot.CallbackQuery) {
	b.LogCallbackHandler(cq)
	handleChatActionError(b.client.SendChatAction(cq.Message.Chat.ID, tbot.ActionTyping))
	time.Sleep(500 * time.Millisecond)
	humanMove := cq.Data
	msg := playGame(humanMove)
	handleChatActionError(b.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID))
	handleMessageError(b.client.SendMessage(cq.Message.Chat.ID, msg))
}

type weather struct {
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

func (b *Bot) cityRequestHandler(m *tbot.Message) {
	b.LogHandler(m)
	handleChatActionError(b.client.SendChatAction(m.Chat.ID, tbot.ActionTyping))
	time.Sleep(500 * time.Millisecond)

	var user *models.User
	user, ok, err := b.storage.User().FindByTelegramUsername(m.From.Username)
	if err != nil {
		b.logger.Info("Error during fetching user data: ", err.Error())
	}
	if !ok {
		userId, err := strconv.Atoi(m.Chat.ID)
		if err != nil {
			b.logger.Info("Failed to convert user ID ", err.Error())
		}

		user = &models.User{Username: m.From.Username, TelegramId: &userId}
		err = b.storage.User().Create(user)
		if err != nil {
			b.logger.Info("Failed to create new user: ", err.Error())
		}
	}
	err = b.storage.Dialog().Create(&models.Dialog{Name: "weather", UserId: user.ID, Status: true})
	if err != nil {
		b.logger.Info("Failed to create new dialog: ", err.Error())
	}
	handleMessageError(b.client.SendMessage(m.Chat.ID, "Напишите мне название города, в котором хотите узнать погоду"))
}

func (b *Bot) weatherHandler(m *tbot.Message) {
	b.LogHandler(m)
	handleChatActionError(b.client.SendChatAction(m.Chat.ID, tbot.ActionTyping))

	user, ok, err := b.storage.User().FindByTelegramUsername(m.From.Username)
	if err != nil {
		b.logger.Info("Error during fetching user data: ", err.Error())
		return
	}

	var msg string
	if !ok {
		b.logger.Info("User and dialog not found")
		msg = "Пожалуйста, запустите меня, выполнив команду /start"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}

	dialog, ok, err := b.storage.Dialog().FindLatestUserDialog(user.ID)
	if err != nil {
		b.logger.Error("Error during fetching dialog data: ", err.Error())
		return
	}

	if !ok || dialog.Status != true {
		b.logger.Info("No active dialog status")
		msg = "Прошу прощения, я пока не умею распознавать такие сообщения. Попробуйте:\n\n/play - Поиграть в Камень-Ножницы-Бумага\n\n/weather - Узнать, какая сейчас погода"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}

	url := "https://community-open-weather-map.p.rapidapi.com/find?q=" + strings.ReplaceAll(strings.TrimSpace(m.Text), " ", "+") + "&lang=ru&mode=null&lon=0&type=link%2C%20accurate&lat=0&units=metric"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-host", "community-open-weather-map.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", b.config.WeatherKey)

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			b.logger.Warning("Failed to close HTTP connection")
		}
	}(res.Body)

	body, _ := ioutil.ReadAll(res.Body)

	w := weather{}

	if err := json.Unmarshal(body, &w); err != nil {
		b.logger.Errorf("Error during unmarshalling weather JSON: %s\nResponse: %s", err.Error(), string(body))
		msg = "Извините, временно туплю.\nНе могу обработать данные о погоде.\nПожалуйста, попробуйте позже.\nА пока можете поиграть в Камень-Ножницы-Бумага /play"
	} else {
		if w.Count < 1 {
			handleMessageError(b.client.SendMessage(m.Chat.ID, "Хмм...🤔\nЧто-то я не слышал о таком городе.\nПопробуйте другой."))
			return
		}
		b.logger.Info(w)
		for i := 0; i < w.Count; i++ {
			var coordinates string
			if w.Count > 2 {
				coordinates = fmt.Sprintf("Координаты: %f, %f\n", w.List[i].Coord.Lat, w.List[i].Coord.Lon)
			}
			var weatherDescription string
			for k, val := range w.List[i].Weather {
				var icon string
				switch val.Icon {
				case "01d":
					fallthrough
				case "01n":
					icon = "☀"
				case "02d":
					fallthrough
				case "02n":
					icon = "🌤"
				case "03d":
					fallthrough
				case "03n":
					icon = "⛅"
				case "04d":
					fallthrough
				case "04n":
					icon = "☁"
				case "09d":
					fallthrough
				case "09n":
					icon = "🌧"
				case "10d":
					fallthrough
				case "10n":
					icon = "🌦"
				case "11d":
					fallthrough
				case "11n":
					icon = "⛈"
				case "13d":
					fallthrough
				case "13n":
					icon = "❄"
				case "50d":
					fallthrough
				case "50n":
					icon = "🌫"
				}
				weatherDescription += strings.Title(val.Description) + " " + icon
				if k < len(w.List[i].Weather)-1 {
					weatherDescription += "\n"
				}
			}
			clouds := w.List[i].Clouds.All
			weatherData := map[string]interface{}{
				"temp":        int(math.Round(w.List[i].Main.Temp)),
				"feelsLike":   int(math.Round(w.List[i].Main.FeelsLike)),
				"humidity":    w.List[i].Main.Humidity,
				"windSpeed":   int(math.Round(w.List[i].Wind.Speed)),
				"city":        w.List[i].Name,
				"country":     w.List[i].Sys.Country,
				"description": weatherDescription,
			}

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

			msg += fmt.Sprintf(""+
				"Город/Страна: %s [%s]\n%s"+
				"%s\n\n"+
				"Температура🌡: %d°C\n"+
				"Ощущается как🌡: %d°C\n\n"+
				"Влажность💧: %d%%\n"+
				"Скорость ветра💨: %d м/с\n"+
				"Облачность: %d%% %s\n"+
				"", weatherData["city"], weatherData["country"], coordinates, weatherData["description"], weatherData["temp"], weatherData["feelsLike"], weatherData["humidity"], weatherData["windSpeed"], clouds, emoji)
			if i < w.Count-1 {
				msg += "\n"
				if len(coordinates) == 0 {
					for j := 0; j < utf8.RuneCountInString(w.List[i].Sys.Country)+utf8.RuneCountInString(w.List[i].Name)+16; j++ {
						if j < 32 {
							msg += "="
						} else {
							break
						}
					}
				} else {
					for j := 0; j < utf8.RuneCountInString(coordinates)-2; j++ {
						if j < 32 {
							msg += "="
						} else {
							break
						}
					}
				}
				msg += "\n\n"
			}
		}
	}

	if err := b.storage.Dialog().UpdateStatus(dialog.ID); err != nil {
		b.logger.Error("Error during updating dialog status: ", err.Error())
	}

	handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
}

func makeButtons() *tbot.InlineKeyboardMarkup {
	btnRock := tbot.InlineKeyboardButton{
		Text:         "Камень",
		CallbackData: "rock",
	}
	btnPaper := tbot.InlineKeyboardButton{
		Text:         "Бумага",
		CallbackData: "paper",
	}
	btnScissors := tbot.InlineKeyboardButton{
		Text:         "Ножницы",
		CallbackData: "scissors",
	}

	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			{btnRock, btnScissors, btnPaper},
		},
	}
}

func playGame(humanMove string) (msg string) {
	var result string
	botMove := picks[rand.Intn(len(picks))]
	switch humanMove {
	case botMove:
		result = "Ничья"
	case options[botMove]:
		result = "Ты проиграл"
	default:
		result = "Ты выиграл"
	}
	msg = fmt.Sprintf("%s!\nТы выбрал %s\nЯ выбрал %s", result, translations[humanMove], translations[botMove])
	return
}
