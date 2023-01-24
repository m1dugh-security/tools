package broadcast

import (
    "testing"
    "os"
    "fmt"
    storage "github.com/m1dugh-security/tools/go/recon-engine/internal/database"
)

func TestSendMessage(t *testing.T) {
    diff := storage.ProgramDiff{
        ProgId: -1,
    }
    token := os.Getenv("DISCORD_TOKEN")
    bot, err := NewDiscord(token)
    if err != nil {
        t.Errorf(fmt.Sprintf("%s", err))
    }

    err = bot.SendDiffs("test-program", diff)
    if err != nil {
        t.Errorf(fmt.Sprintf("%s", err))
    }

    bot.Close()
}

func TestPing(t *testing.T) {
    token := os.Getenv("DISCORD_TOKEN")
    if len(token) == 0 {
        t.Errorf("Missing DISCORD_TOKEN")
    }

    bot, err := NewDiscord(token)
    if err != nil {
        t.Errorf(fmt.Sprintf("%s", err))
    }

    err = bot.pingAll(programChannel)
    if err != nil {
        t.Errorf(fmt.Sprintf("%s", err))
    }

    bot.Close()
}
