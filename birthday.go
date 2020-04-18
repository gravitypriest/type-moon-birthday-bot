package main

import (
	"fmt"
	"strings"
	"os"
	"os/signal"
	"syscall"
	"time"
	"io/ioutil"
	"strconv"

	"gopkg.in/yaml.v2"
	"github.com/bwmarrin/discordgo"
)

var birthdays map[string]map[string][]string

func main() {

	var token string
	
	loadBirthdays()

	// initialize the discord object
	dg, err := discordgo.New("Bot " + token)
	dg.AddHandler(messageCreate)
	dg.AddHandler(guildCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Birthay bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!birthday") {
		// Find the channel that the message came from.
		_, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		_, err = s.ChannelMessageSend(m.ChannelID, getBirthdays());
		if err != nil {
			fmt.Println("something went wrong")
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "TYPE-MOON birthday bot is ready! Type !birthday to get birthdays.")
			return
		}
	}
}

func loadBirthdays() {
	// load YAML into birthdays map
	data, _ := ioutil.ReadFile("birthdays.yaml")
	_ = yaml.Unmarshal(data, &birthdays)
}

func getBirthdays() string {
	var bdayMessage string

	_, m, d := time.Now().Date()
	month := strconv.Itoa(int(m))
	day := strconv.Itoa(int(d))

	todaysBirthdays := birthdays[month][day]
	if len(todaysBirthdays) == 0 {
		bdayMessage = "No one is celebrating their birthday today.\n"
	} else{
		bdayMessage = ""
		for _, bday := range todaysBirthdays {
			bdayMessage = bdayMessage + fmt.Sprintf("Today is **%s**'s birthday!\n", bday)
		}
		bdayMessage = bdayMessage + "\nHappy birthday!\n"
	}

	thisMonthsBirthdays := "\nBirthdays this month:\n>>> "

	// thisMonthsBirthdays = thisMonthsBirthdays + monthBirthdays(month)

    bdayMessage = bdayMessage + thisMonthsBirthdays

    return bdayMessage
}

// func monthBirthdays(monthNum string) string {
//     for d := 1; d < len(birthdays[month]) + 1; d++ {
//     	day = strconv.Itoa(d)
//         if len(birthdays[month][day]) > 0 {
//         	for _, bday := range birthdays[month][day] {
//         		monthBirthdays = monthBirthdays + fmt.Sprintf("%s/%s - %s\n", month, day, bday)
//         	}
//         }
//     }
// }