package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"gomap/internal/db"
	"gomap/internal/gps_utils"
	"log"
	"os"
	"regexp"
	"strconv"
)

func extractCoordinates(pointStr string) (float64, float64, error) {
	re := regexp.MustCompile(`POINT\(([-+]?[0-9]*\.?[0-9]+) ([-+]?[0-9]*\.?[0-9]+)\)`)
	matches := re.FindStringSubmatch(pointStr)
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("invalid POINT format")
	}

	lng, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, 0, err
	}

	lat, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return 0, 0, err
	}

	return lat, lng, nil
}

func Bot() {
	// .env init
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Bot in gorutine
	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		for update := range updates {
			if update.Message == nil {
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Location != nil {
				latitude := update.Message.Location.Latitude
				longitude := update.Message.Location.Longitude

				var myCoords gps_utils.GpsCoordinates
				myCoords.Lat = strconv.FormatFloat(latitude, 'f', -1, 64)
				myCoords.Lng = strconv.FormatFloat(longitude, 'f', -1, 64)

				places, err := db.GetNearPlaces(myCoords)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error fetching nearby places.")
					bot.Send(msg)
					continue
				}

				if len(places) == 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No nearby places found.")
					bot.Send(msg)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "\nБлижайшие интересные места: \n")
				bot.Send(msg)

				for _, place := range places {
					responseText := fmt.Sprintf("%s\n%d метров от вас.\n\n%s", place.Name, int(place.Distance), place.Desc)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
					bot.Send(msg)

					lat, lng, err := extractCoordinates(place.Geom)
					if err != nil {
						log.Printf("Error extracting coordinates: %v", err)
						continue
					}

					locationMsg := tgbotapi.NewLocation(update.Message.Chat.ID, lat, lng)
					bot.Send(locationMsg)

					// Добавление отступа между сообщениями
					separatorMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "\n\n ")
					bot.Send(separatorMsg)
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please send a location.")
				bot.Send(msg)
			}
		}
	}()
}
