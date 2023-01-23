package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"github.com/m1dugh/program-browser/pkg/browser"
	"github.com/m1dugh/program-browser/pkg/types"
	"github.com/m1dugh/program-browser/pkg/utils"
)

func CreateFolder(root, platform, program string) (string, error) {
    url := path.Join(root, platform, program)
    
    err := os.MkdirAll(url, 0777)
    if err != nil {
        if os.IsExist(err) {
            return url, nil
        } else {
            return "", err
        }
    }

    return url, nil

}


func CreateProgramFiles(root string, program *types.Program) error {
    
    scope := program.GetScope(types.API, types.Website)
    urls, domains := scope.ExtractInfo()

    filePath := path.Join(root, "urls.txt")
    file, err := os.Create(filePath)
    if err != nil {
        return errors.New(fmt.Sprintf("Error while creating %s", filePath))
    }
    for _, url := range urls.ToArray() {
        file.WriteString(fmt.Sprintln(url))
    }
    err = file.Close()
    if err != nil {
        return errors.New(fmt.Sprintf("Error while closing %s", filePath))
    }

    filePath = path.Join(root, "domains.txt")
    file, err = os.Create(filePath)
    if err != nil {
        return errors.New(fmt.Sprintf("Error while creating %s", filePath))
    }

    for _, domain := range domains.ToArray() {
        file.WriteString(fmt.Sprintln(domain))
    }
    err = file.Close()
    if err != nil {
        return errors.New(fmt.Sprintf("Error while closing %s", filePath))
    }

    filePath = path.Join(root, "program.json")
    file, err = os.Create(filePath)
    if err != nil {
        return errors.New(fmt.Sprintf("Error while creating %s", filePath))
    }
    body, err := json.MarshalIndent(program, "", "\t")
    if err != nil {
        return errors.New(fmt.Sprintf("Error while marshaling program: %s", err))
    }

    file.Write(body)
    err = file.Close()
    if err != nil {
        return errors.New(fmt.Sprintf("Error while closing %s", filePath))
    }

    filePath = path.Join(root, "scope.json")
    file, err = os.Create(filePath)
    if err != nil {
        return errors.New(fmt.Sprintf("Error while creating %s", filePath))
    }

    body, err = scope.ToBurpScope()
    file.Write(body)
    err = file.Close()
    if err != nil {
        return errors.New(fmt.Sprintf("Error while closing %s", filePath))
    }

    return nil
}

func main() {

    var outputFolder string
    flag.StringVar(&outputFolder, "o", "", "The path to output at")

    var settingsFile string
    flag.StringVar(&settingsFile, "settings", "", "The path to the settings file")
    
    flag.Parse()

    if len(outputFolder) == 0 {
        log.Fatal("missing -o flag")
    }

    var options *browser.Options

    if len(settingsFile) == 0 {
        options = browser.DefaultOptions()
    } else {
        body, err := os.ReadFile(settingsFile)
        if err != nil {
            log.Fatal(err)
        }
        options, err = utils.DeserializeOptions(body)
        if err != nil {
            log.Fatal(err)
        }
    }


    browser := browser.New(options)

    results, err := browser.GetPrograms()
    if err != nil {
        log.Fatal(err)
    }

    for _, program := range results {

        path, err := CreateFolder(outputFolder, program.Platform, program.Name)
        if err != nil {
            log.Fatal(fmt.Sprintf("Could not create folder for %s: %s", program.Code(),  err))
        }

        err = CreateProgramFiles(path, program)
        if err != nil {
            log.Fatal(err)
        }
    }
}
