package main

import (
	"fmt"
	"go-discord-bot/weather"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/enescakir/emoji"
	"github.com/joho/godotenv"
)

const prefix string = "!Botify"

func main() {

	godotenv.Load(".env")

	token := os.Getenv("DISCORD_TOKEN")
	sess, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(messageCreate)

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot is up and running. Press CTR-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	defer sess.Close()
}

func messageCreate(currentSess *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == currentSess.State.User.ID {
		return
	}

	if !strings.HasPrefix(message.Content, prefix) {
		return
	}

	reqCall := strings.Split(message.Content, " ")[1:]

	if reqCall[0] == "hello" {
		HelloWorldHandler(currentSess, message)
	}

	if reqCall[0] == "proverb" {
		ProverbsHandler(currentSess, message)

	}

	if reqCall[0] == "weather" {
		if len(reqCall) < 2 {
			currentSess.ChannelMessageSend(message.ChannelID, "Please include the location!")
		}
		countryCheck := strings.Split(message.Content, ",")

		if len(countryCheck) > 1 {
			CityAndCountryWeatherHandler(currentSess, message)
		} else {
			CityOnlyWeatherHandler(currentSess, message)
		}
	}
}

func HelloWorldHandler(currentSess *discordgo.Session, message *discordgo.MessageCreate) {
	response := "Hi " + message.Author.Username + "!"
	currentSess.ChannelMessageSend(message.ChannelID, response)
}

func ProverbsHandler(currentSess *discordgo.Session, message *discordgo.MessageCreate) {
	type BasicProverb struct {
		Language string
		Quote    string
	}

	proverbs := []BasicProverb{
		{Language: "African", Quote: "A king's child is a slave elsewhere."},
		{Language: "African", Quote: "What forgets is the ax, but the tree that has been axed will never forget."},
		{Language: "African", Quote: "It is no shame at all to work for money."},
		{Language: "African", Quote: "A loose tooth will not rest until it's pulled out."},
		{Language: "African", Quote: "He who digs too deep for a fish may come out with a snake."},
		{Language: "African", Quote: "The path is made by walking."},
		{Language: "Australian", Quote: "None are so deaf as those who would not hear."},
		{Language: "Australian", Quote: "A bad worker blames his tools."},
		{Language: "Australian", Quote: "In the planting season, visitors come singly, and in harvest time they come in crowds."},
		{Language: "Egyptian", Quote: "Do a good deed and throw it into the sea."},
		{Language: "Egyptian", Quote: "Time never gets tired of running."},
		{Language: "Bulgarian", Quote: "Tell me who your friends are, so I can tell you who you are."},
		{Language: "Bulgarian", Quote: "Measure thrice, cut once."},
		{Language: "English", Quote: "When the going gets tough, the tough get going."},
		{Language: "English", Quote: "The pen is mightier than the sword."},
		{Language: "English", Quote: "The squeaky wheel gets the grease."},
		{Language: "English", Quote: "No man is an island."},
		{Language: "English", Quote: "People who live in glass houses shouldn't throw stones."},
		{Language: "English", Quote: "Better late than never."},
		{Language: "English", Quote: "Two wrongs don't make a right."},
		{Language: "German", Quote: "He who rests grows rusty."},
		{Language: "German", Quote: "Starting is easy, persistence is an art."},
		{Language: "German", Quote: "The cheapest is always the most expensive."},
		{Language: "German", Quote: "Make haste with leisure."},
	}

	selection := rand.Intn(len(proverbs))

	origin := discordgo.MessageEmbedAuthor{
		Name: proverbs[selection].Language + " proverb",
	}

	embed := discordgo.MessageEmbed{
		Title:  proverbs[selection].Quote,
		Author: &origin,
	}

	currentSess.ChannelMessageSendEmbed(message.ChannelID, &embed)

}

type iconKey struct {
	openWeatherIcon string
	emoji           string
}

var iconSlice = []iconKey{
	{openWeatherIcon: "01d", emoji: string(emoji.Sun)},
	{openWeatherIcon: "01n", emoji: string(emoji.CrescentMoon)},
	{openWeatherIcon: "02d", emoji: string(emoji.SunBehindCloud)},
	{openWeatherIcon: "02n", emoji: string(emoji.SunBehindLargeCloud)},
	{openWeatherIcon: "03d", emoji: string(emoji.Cloud)},
	{openWeatherIcon: "03n", emoji: string(emoji.Cloud)},
	{openWeatherIcon: "04d", emoji: string(emoji.Cloud)},
	{openWeatherIcon: "04n", emoji: string(emoji.Cloud)},
	{openWeatherIcon: "09d", emoji: string(emoji.CloudWithRain)},
	{openWeatherIcon: "09n", emoji: string(emoji.CloudWithRain)},
	{openWeatherIcon: "10d", emoji: string(emoji.SunBehindRainCloud)},
	{openWeatherIcon: "10n", emoji: string(emoji.SunBehindRainCloud)},
	{openWeatherIcon: "11d", emoji: string(emoji.CloudWithLightningAndRain)},
	{openWeatherIcon: "11n", emoji: string(emoji.CloudWithLightningAndRain)},
	{openWeatherIcon: "13d", emoji: string(emoji.Snowman)},
	{openWeatherIcon: "13n", emoji: string(emoji.Snowman)},
	{openWeatherIcon: "50d", emoji: string(emoji.Fog)},
	{openWeatherIcon: "50n", emoji: string(emoji.Fog)},
}

func CityOnlyWeatherHandler(currentSess *discordgo.Session, message *discordgo.MessageCreate) {

	requestedLocationArray := strings.Split(message.Content, " ")[2:]
	requestedLocation := strings.Join(requestedLocationArray, " ")
	location, lat, lon, country := weather.GetCordsCityOnly(requestedLocation)
	ReturnData := weather.GetWeather(lat, lon)

	source := discordgo.MessageEmbedAuthor{
		Name: "Brought to you by the openweathermap API",
		URL:  "https://openweathermap.org/",
	}

	weatherEmoji := ""
	for _, value := range iconSlice {
		if value.openWeatherIcon == ReturnData.Icon {
			weatherEmoji = value.emoji
		}
	}

	embed := discordgo.MessageEmbed{
		Author: &source,
		Title:  "Today's weather in " + location + ", " + country,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  ReturnData.Primary + " " + weatherEmoji,
				Value: ReturnData.Description,
			},
			{
				Name: "Current Temperature: " + ReturnData.Temp + "C",
			},
			{
				Name: "With lows of " + ReturnData.MinTemp + "C" + " and highs of " + ReturnData.MaxTemp + "C",
			},
		},
	}

	currentSess.ChannelMessageSendEmbed(message.ChannelID, &embed)
}

func CityAndCountryWeatherHandler(currentSess *discordgo.Session, message *discordgo.MessageCreate) {
	countryCheck := strings.Split(message.Content, ",")
	countryStringDirty := strings.Join(countryCheck[1:], ",")
	countryString := strings.TrimSpace(countryStringDirty)
	countryStringCorrected := strings.Title(countryString)
	baseCode := weather.GetCountryCode(countryStringCorrected)

	countryRemovedString := countryCheck[0]
	requestedLocationArray := strings.Split(countryRemovedString, " ")[2:]
	requestedLocation := strings.Join(requestedLocationArray, " ") + "," + baseCode.ISO2
	location, lat, lon := weather.GetCordsCityAndCountry(requestedLocation)
	ReturnData := weather.GetWeather(lat, lon)

	source := discordgo.MessageEmbedAuthor{
		Name: "Brought to you by the openweathermap API",
		URL:  "https://openweathermap.org/",
	}

	weatherEmoji := ""
	for _, value := range iconSlice {
		if value.openWeatherIcon == ReturnData.Icon {
			weatherEmoji = value.emoji
		}
	}

	embed := discordgo.MessageEmbed{
		Author: &source,
		Title:  "Today's weather in " + location + ", " + baseCode.Name,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  ReturnData.Primary + " " + weatherEmoji,
				Value: ReturnData.Description,
			},
			{
				Name: "Current Temperature: " + ReturnData.Temp + "C",
			},
			{
				Name: "With lows of " + ReturnData.MinTemp + "C" + " and highs of " + ReturnData.MaxTemp + "C",
			},
		},
	}

	currentSess.ChannelMessageSendEmbed(message.ChannelID, &embed)
}
