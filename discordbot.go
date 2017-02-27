package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

/**
 * Commands are basically functions that will take the command and return a response to it
 */
type Command func(*discordgo.Session, *discordgo.Message) string

var (
	BotID    string
	commands map[string]Command
)

/**
 * The "raids" map holds the sets of raids that should be happening.
 * The "commands" array maps string commands to a function in charge of returning the actual values
 */
func init() {
	raids = make(map[string]map[string]Raid)

	commands = make(map[string]Command)

	commands["!raid"] = parseRaidCommand
	commands["!help"] = showHelp
	commands["!mutha"] = lando
	commands["!towel"] = func(*discordgo.Session, *discordgo.Message) string { return "42" }
}

/**
 * Register URL
 * https://discordapp.com/oauth2/authorize?&client_id=256141864548302849&scope=bot&permissions=0
 */
func main() {
	token, err := ioutil.ReadFile("token.txt")
	if err != nil {
		panic(err)
	}

	//The readfile call will add a "\n" at the end of the string, so get rid of it.
	discord, err := discordgo.New("Bot " + string(token[0:len(token)-1]))
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	u, err := discord.User("@me")
	if err != nil {
		fmt.Println("error obtaining user details, ", err)
	}
	BotID = u.ID

	discord.AddHandler(messageCreated)

	err = discord.Open()
	defer discord.Close()

	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("Bot is now running.")
	fmt.Println("Press CTRL-C to exit.")

	// Simple way to keep program running until CTRL-C is pressed.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	fmt.Println("Shutting down now. Cya")
	<-sc
}

func messageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	u, _ := s.User("@me")
	if botWasMentioned(m, u) {
		sendMessage(s, m.ChannelID, "I'm not the droid you're looking for. Type \"!help\" for the available commands... and stop tagging me.")
		return
	}

	chunks := strings.Split(m.Content, " ")
	if val, ok := commands[chunks[0]]; ok {
		channel, _ := s.Channel(m.ChannelID)
		fmt.Printf("received command %s on %s by %s\n", m.Content, channel.Name, m.Author.Username)
		sendMessage(s, m.ChannelID, val(s, m.Message))
	}
}

func botWasMentioned(m *discordgo.MessageCreate, u *discordgo.User) bool {
	for _, user := range m.Mentions {
		if user.ID == BotID {
			return true
		}
	}
	return false
}

func sendMessage(s *discordgo.Session, channel string, text string) {
	_, _ = s.ChannelMessageSend(channel, text)
}

func showHelp(s *discordgo.Session, m *discordgo.Message) string {
	return "Type \"!raid\" to get the next timings for raids.\n" + "Type \"!raid eu\" to get the next timings for raids in the eu guild.\n" + "Type \"!raid us\" to get the next timings for raids in the us guild.\n" + "Type \"!mutha fuk'n LANDO\" to learn the TRUTH.\n"
}

func lando(s *discordgo.Session, m *discordgo.Message) string {
	chunks := strings.Split(m.Content, " ")
	if len(chunks) == 3 && chunks[1] == "fuk'n" && chunks[2] == "LANDO" {
		return "Yeah! We started that shit :D"
	} else {
		return "Learn your shit, it's !mutha fuk'n LANDO"
	}
}
