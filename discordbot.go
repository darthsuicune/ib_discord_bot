package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

/**
 * Commands are basically functions that will take the command and return a response to it
 */
type Command func(*discordgo.Session, *discordgo.Message) string

/**
 * The Raid structure holds all the required information for the raid itself
 */
type Raid struct {
	Type   string
	Timing time.Time
	Guild  string
}

/**
 * Post the Raid information in an elegant way
 */
func (r Raid) String() string {
	//this will show something like 1h1m12.123651235612s, so better trim the ending@
	waitTime := r.Timing.Sub(time.Now()).String()
	minuteMark := strings.LastIndex(waitTime, "m") + 1
	return "Next **" + r.Type + "** raid in **" + strings.ToUpper(r.Guild) + "** guild in **" + waitTime[0:minuteMark] + "**\n"
}

var (
	BotID    string
	raids    map[string]map[string]Raid
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

func parseRaidCommand(s *discordgo.Session, m *discordgo.Message) string {
	chunks := strings.Split(m.Content, " ")

	var buffer bytes.Buffer
	if len(chunks) == 1 {
		for _, guildRaids := range raids {
			for _, items := range guildRaids {
				buffer.Write([]byte(items.String()))
			}
		}
	} else if chunks[1] == "eu" || chunks[1] == "us" {
		for _, guildRaids := range raids[chunks[1]] {
			buffer.Write([]byte(guildRaids.String()))
		}
	} else if chunks[1] == "set" {
		channel, _ := s.Channel(m.ChannelID)
		member, _ := s.GuildMember(channel.GuildID, m.Author.ID)
		if canSetRaids(s, channel.GuildID, member.Roles) {
			if len(chunks) >= 5 {
				err := setRaid(chunks[2:])
				if err != nil {
					buffer.Write([]byte("Error when creating the raid " + err.Error()))
				} else {
					buffer.Write([]byte("I've altered the deal. New details:\n"))
					buffer.Write([]byte(raids[chunks[2]][chunks[3]].String()))
				}
			} else {
				buffer.Write([]byte("You are missing like half the stuff, dude. Do it like this: !raid set us Rancor 1970-01-01 00:00 TZ"))
			}
		} else {
			buffer.Write([]byte("You can't do that"))
		}
	} else if chunks[1] == "delete" {
		channel, _ := s.Channel(m.ChannelID)
		member, _ := s.GuildMember(channel.GuildID, m.Author.ID)
		if canSetRaids(s, channel.GuildID, member.Roles) {
			if len(chunks) < 4 || len(chunks) > 5 {
				buffer.Write([]byte("Wrong parameters, use !raid delete <us/eu> <raid_name>"))
			} else {
				delete(raids[chunks[2]], chunks[3])
				buffer.Write([]byte("I've altered the deal again. Pray that I don't alter it any further."))
			}
		}
	}
	return buffer.String()
}

func setRaid(chunks []string) error {
	if chunks[0] != "eu" && chunks[0] != "us" {
		return errors.New("Guild must be eu or us")
	}
	guild := strings.TrimSpace(chunks[0])
	if _, ok := raids[guild]; !ok {
		raids[guild] = make(map[string]Raid)
	}
	raidType := strings.TrimSpace(strings.ToLower(chunks[1]))
	when := strings.TrimSpace(strings.Join(chunks[2:], " "))
	
	var location *time.Location
	if(len(chunks) > 4) {
		location, _ = time.LoadLocation(chunks[4])
	} else {
		location, _ = time.LoadLocation("GMT")
	}
	raidTime, err := time.ParseInLocation("2006-01-02 15:04 MST", when, location)
	//Only set the new raid if there wasn't any error while parsing the date.
	if err == nil {
		raids[guild][raidType] = Raid{Guild: guild, Type: raidType, Timing: raidTime}
	}
	return err

}

func canSetRaids(s *discordgo.Session, guildID string, roles []string) bool {
	for _, role := range roles {
		r, _ := s.State.Role(guildID, role)
		if r.Name == "Inglorious Leaders" || r.Name == "Inglorious Officers" {
			return true
		}
	}
	return false
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
