package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"encoding/json"
)

/**
 * Commands are basically functions that will take the command and return a response to it
 */
type Command func(*discordgo.Session, *discordgo.Message) string


var (
	BotID    string
	commands map[string]Command

	guilds *Guilds
)

const (
	TOKEN_FILE_NAME = "token.txt"
	RAIDS_FILE_NAME = "raids.txt"
)

/**
 * The "raids" map holds the sets of raids that should be happening.
 * The "commands" array maps string commands to a function in charge of returning the actual values
 */
func init() {
	commands = make(map[string]Command)
	
	commands["!draid"] = parseDefaultRaidCommand
	commands["!raid"] = parseRaidCommand
	commands["!help"] = showHelp
	commands["!mutha"] = lando
	commands["!towel"] = func(*discordgo.Session, *discordgo.Message) string { return "42" }
	commands["!shot_first"] = func(*discordgo.Session, *discordgo.Message) string { return "OF COURSE IT WAS HAN! This damn noobs..." }
	commands["!friend"] = func(*discordgo.Session, *discordgo.Message) string { return "My friend doesn't like you. I don't like you either." }

	guilds = &Guilds{Eu: &Guild{Location: euTime()}, Us: &Guild{Location: usTime()}}
}

/**
 * Register URL
 * https://discordapp.com/oauth2/authorize?&client_id=256141864548302849&scope=bot&permissions=0
 */
func main() {
	discord := startSession()

	loadSavedState()

	u, err := discord.User("@me")
	if err != nil {
		fmt.Println("error obtaining user details, ", err)
	}
	BotID = u.ID

	discord.AddHandler(onMessageCreated)

	err = discord.Open()
	defer discord.Close()

	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("Bot is now running.")

	// Simple way to keep program running until any kill signal is received (including ctrl + c).
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	err = saveState()
	if err == nil {
		fmt.Println("Data saved. Shutting down now. Cya")
	} else {
		fmt.Println("NOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO")
	}
}
func saveState() error {
	res, err := json.Marshal(guilds)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(RAIDS_FILE_NAME, res, 0600)
	return err
}

func loadSavedState() {
	raids, err := ioutil.ReadFile(RAIDS_FILE_NAME)
	if err != nil {
		fmt.Println("No old raids were found." + err.Error())
	}

	err = json.Unmarshal(raids, guilds)
	if err != nil {
		fmt.Println("Raids found but error while reading them: " + err.Error())
	}
	guilds.Eu.Location = euTime()
	guilds.Us.Location = usTime()
}

func startSession() *discordgo.Session {
	token, err := ioutil.ReadFile(TOKEN_FILE_NAME)
	if err != nil {
		panic(err)
	}

	//The readfile call will add a "\n" at the end of the string, so get rid of it.
	discord, err := discordgo.New("Bot " + string(token[0:len(token)-1]))
	if err != nil {
		panic(err)
	}
	return discord
}

func onMessageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	if botWasMentioned(m) {
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

func botWasMentioned(m *discordgo.MessageCreate) bool {
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

func parseDefaultRaidCommand(s *discordgo.Session, m *discordgo.Message) string {
	if !hasPermissionToSetRaids(s, m) {
		return "Nope"
	}
	return setDefaultRaid(m.Content)
}

func parseRaidCommand(s *discordgo.Session, m *discordgo.Message) string {
	return readRaidCommand(m.Content, hasPermissionToSetRaids(s, m))
}

func hasPermissionToSetRaids(s *discordgo.Session, m *discordgo.Message) bool {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	guildID := channel.GuildID
	member, err := s.GuildMember(guildID, m.Author.ID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	roles := member.Roles
		
	for _, role := range roles {
		r, _ := s.State.Role(guildID, role)
		if r.Name == "Inglorious Leaders" || r.Name == "Inglorious Officers" || r.Name == "US Officers" || r.Name == "EU Officers"{
			return true
		}
	}
	return false
}
