package main

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "flag"
    "log"
    "io"
    "os"
)

type Channel struct {
    GuildId string
    ChannelId string
}

type Options struct {
    Channel     Channel
    Token       string
    BufferSize  uint64
    Clear       bool
}

func DefaultOptions(token string, channel Channel) *Options {
    return &Options{
        Channel: channel,
        Token: token,
        BufferSize: 128,
        Clear: false,
    }
}

type DiscordLogger struct {
    Options *Options
    session *discordgo.Session
    Reader  io.Reader
}

func New(options *Options, reader io.Reader) (*DiscordLogger, error) {
    session, err := discordgo.New("Bot " + options.Token)
    if err != nil {
        return nil, err
    }

    return &DiscordLogger{
        Options: options,
        session: session,
        Reader: reader,
    }, nil
}

func (d *DiscordLogger) Log(msg string) error {
    _, err := d.session.ChannelMessageSend(d.Options.Channel.ChannelId, msg)
    return err
}

func (d *DiscordLogger) Close() {
    d.session.Close()
}

func ParseOptions() *Options {


    var channelId string
    flag.StringVar(&channelId, "c", "", "The id of the discord channel to write in")

    var clear bool
    flag.BoolVar(&clear, "clear", false, "Whether or not the channel should be cleared of previous messages.")

    var bufferSize uint64
    flag.Uint64Var(&bufferSize, "buffer", 0, "The max length of a log line")

    flag.Parse()

    token := os.Getenv("DISCORD_TOKEN")
    if len(token) == 0 {
        log.Fatal("missing required env var DISCORD_TOKEN")
    }

    if len(channelId) == 0 {
        log.Fatal("missing required param channel (-c)")
    }

    channel := Channel{"TODO: add guild id", channelId}

    options := DefaultOptions(token, channel)
    if bufferSize != 0 {
        options.BufferSize = bufferSize   
    }

    options.Clear = clear

    return options
}

func (d *DiscordLogger) Clear() {

    channel := d.Options.Channel.ChannelId

    for true {
        messages, err := d.session.ChannelMessages(channel, 100, "", "", "")
        if err != nil {
            log.Fatal(err)
            fmt.Fprintf(os.Stderr, "discord clear logs error: %s\n", err)
        }

        if len(messages) == 0 {
            break
        }
        ids := make([]string, len(messages))
        for i, msg := range messages {
            ids[i] = msg.ID
        }

        err = d.session.ChannelMessagesBulkDelete(channel, ids)
        if err != nil {
            // log.Fatal(err)
            for _, id := range ids {
                err = d.session.ChannelMessageDelete(channel, id)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "discord clear logs error: %s\n", err)
                    fmt.Fprintf(os.Stderr, "skipping message delete")
                    return
                }
            }
        }
    }
}

func main() {

    options := ParseOptions()

    bot, err := New(options, os.Stdin)
    if err != nil {
        log.Fatal(err)
    }

    if options.Clear {
        bot.Clear()
    }

    buffer := make([]byte, bot.Options.BufferSize)
    contentToRead := 0
    retries := 0
    for true {
        if contentToRead == 0 {
            contentToRead, err = bot.Reader.Read(buffer)
            if err != nil {
                break
            }
        }

        message := fmt.Sprintf("```\n%s\n```", string(buffer[:contentToRead]))
        err = bot.Log(message)
        if err != nil {
            retries++
            if retries >= 3 {
                log.Fatal("An error happened while writing logs")
            }
        } else {
            contentToRead = 0
            retries = 0
        }
    }

    bot.Close()
}
