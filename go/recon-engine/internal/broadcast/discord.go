package broadcast

import (
    "os"
    "fmt"
    "errors"
    "io"
    "github.com/bwmarrin/discordgo"
    storage "github.com/m1dugh/recon-engine/internal/database"
    "encoding/json"
    "sync"
)

type DiscordBot struct {
    session *discordgo.Session
    reader  io.Reader
    mut     *sync.Mutex
    logs    bool
    guildId string  
}

func (bot *DiscordBot) Guild() (*discordgo.Guild, error) {
    guild, err := bot.session.Guild(guildId)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("NewDiscord: Could not get guild: %s", err))
    }

    return guild, nil
}

const (
    programChannel = "1043318300375007343"
    logChannel = "1043318147706527824"
    guildId = "1043317896794874066"
)
const MAX_SIZE = 1800

func NewDiscord(token string) (*DiscordBot, error) {
    session, err := discordgo.New("Bot " + token)
    if err != nil {
        return nil, errors.New("could not create discord session")
    }

    return &DiscordBot{
        session: session,
        guildId: guildId,
    }, nil
}

func (bot *DiscordBot) Start() error {
    bot.clearLogs()
    err := bot.pingAll(programChannel)
    if err != nil {
        return err
    }
    err = bot.pingAll(logChannel)
    return err
}

func (bot *DiscordBot) sendDiffPart(chunk string, channel string) error {

    message := fmt.Sprintf("```\n%s\n```", chunk)
    _, err := bot.session.ChannelMessageSend(channel, message)
    if err != nil {
        return err
    }

    return nil
}

func (bot *DiscordBot) pingAll(channel string) error {
    guild, err := bot.Guild()
    if err != nil {
        return errors.New(fmt.Sprintf("DiscordBot.pingAll: error when fetching guild: %s", err))
    }
    var role *discordgo.Role
    for _, r := range guild.Roles {
        if r.Name == "hacker" {
            role = r
            break
        }
    }
    message := role.Mention()
    _, err = bot.session.ChannelMessageSend(channel, message)
    if err != nil {
        return err
    }

    return nil
}

func (bot *DiscordBot) SetReader(rd io.Reader) {
    bot.reader = rd
    bot.writeLogs(logChannel)
}

func (bot *DiscordBot) clearLogs() {
    channel := logChannel

    for true {
        messages, err := bot.session.ChannelMessages(channel, 100, "", "", "")
        if err != nil {
            // log.Fatal(err)
            fmt.Fprintf(os.Stderr, "discord clear logs error: %s\n", err)
        }

        if len(messages) == 0 {
            break
        }
        ids := make([]string, len(messages))
        for i, msg := range messages {
            ids[i] = msg.ID
        }

        err = bot.session.ChannelMessagesBulkDelete(channel, ids)
        if err != nil {
            // log.Fatal(err)
            for _, id := range ids {
                err = bot.session.ChannelMessageDelete(channel, id)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "discord clear logs error: %s\n", err)
                    fmt.Fprintf(os.Stderr, "skipping message delete")
                    return
                }
            }
        }
    }
}

func (bot *DiscordBot) writeLogs(channel string) {
    if bot.reader == nil || bot.logs {
        return
    }
    bot.logs = true

    go func() {
        buffer := make([]byte, 128)
        opened := true
        for opened {
            n, err := bot.reader.Read(buffer)
            if err != nil {
                opened = false
                break
            }
            message := fmt.Sprintf("```\n%s\n```", string(buffer[:n]))
            _, err = bot.session.ChannelMessageSend(channel, message)
            if err != nil {
                // log.Fatal(err)
                fmt.Fprintf(os.Stderr, "discord write logs error: %s\n", err)
            }
        }
    }()
}

func (bot *DiscordBot) SendDiffs(code string, diff storage.ProgramDiff) error {

    channel := programChannel

    b, err := json.MarshalIndent(diff, "", "\t")
    if err != nil {
        return errors.New("Could not marshal diffs")
    }
    _, err = bot.session.ChannelMessageSend(programChannel, code)
    if err != nil {
        return errors.New(fmt.Sprintf("Could not send message: %s", err))
    }

    msg := string(b)
    var count int = len(msg) / MAX_SIZE
    for i := 0; i < count; i++ {
        err = bot.sendDiffPart(msg[i * MAX_SIZE:(i + 1) * MAX_SIZE], channel)
        if err != nil {
            return errors.New(fmt.Sprintf("Could not send chunk %d: %s", i, err))
        }
    }
    if len(msg) % MAX_SIZE > 0 {
        err = bot.sendDiffPart(msg[count * MAX_SIZE:], channel)
        if err != nil {
            return errors.New(fmt.Sprintf("Could not send chunk %d: %s", count, err))
        }
    }

    return nil
}

func (bot *DiscordBot) Close() {
    bot.session.Close()
}
