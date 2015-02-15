package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

const (
	password    = "password"
	email       = "email"
	port        = "port"
	name        = "name"
	server      = "server"
	UserDataNum = 2
)

type InitData struct {
	DataFilePath       string
	TmplFilePath       string
	SenderInfoFilePath string
	MessageSubject     string
}

type UserData struct {
	UserName string
	Email    string
}

type Notifier struct {
	InitData     InitData
	SenderInfo   SenderInfo
	RawUserData  []string
	UserData     []UserData
	Set          *template.Template
	RenderedData []string
}

func (self *Notifier) Init() {
	if _, err := os.Stat(self.InitData.DataFilePath); os.IsNotExist(err) {
		log.Fatalf("Data file does not exsist: %s", self.InitData.DataFilePath)
	}
	self.readData()
	self.readSenderInfo()
	self.parse()
	self.Set = template.New("gloabal")
}

func (self *Notifier) ProcessTemplate() {
	if _, err := os.Stat(self.InitData.TmplFilePath); os.IsNotExist(err) {
		log.Fatalf("Tmpl file does not exsist: %s", self.InitData.DataFilePath)
	}
	if _, err := self.Set.ParseFiles(self.InitData.TmplFilePath); err != nil {
		log.Fatalf("set.ParseFiles error: %s", err.Error())
	}
	parts := strings.Split(self.InitData.TmplFilePath, "/")
	tmplName := self.InitData.TmplFilePath
	if len(parts) > 0 {
		tmplName = parts[len(parts)-1]
	}

	var data []string
	for index := 0; index < len(self.UserData); index++ {
		ctx := make(Context)
		ctx["UserData"] = self.UserData[index]
		buf := new(bytes.Buffer)
		RenderTemplate(self.Set, buf, tmplName, ctx)
		item := buf.String()
		data = append(data, item)
	}
	self.RenderedData = data
}

func (self *Notifier) Send() {
	for index := 0; index < len(self.RenderedData); index++ {
		message := self.RenderedData[index]

		sendMessage(&Message{
			Subject:       self.InitData.MessageSubject,
			Message:       message,
			ReceiverEmail: self.UserData[index].Email,
			Receiver:      self.UserData[index].UserName,
		}, &self.SenderInfo)
	}
}

func (self *Notifier) readSenderInfo() {
	if _, err := os.Stat(self.InitData.SenderInfoFilePath); os.IsNotExist(err) {
		log.Fatalf("Sender info file does not exsist: %s", self.InitData.DataFilePath)
	}
	data, err := readLines(self.InitData.SenderInfoFilePath)
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}
	if len(data) == 0 {
		log.Fatalf("Sender info file is empty")
	}
	info := make(map[string]string)
	for index := 0; index < len(data); index++ {
		items := strings.Split(data[index], "=")
		if len(items) != 2 {
			log.Fatalf("Cannot read sender info file. Error in line (%i): %s", index+1, data[index])
		}
		info[items[0]] = items[1]
	}
	senderInfo := SenderInfo{
		Email:    info[email],
		Password: info[password],
		Name:     info[name],
		Port:     info[port],
		Server:   info[server],
	}
	self.SenderInfo = senderInfo
}

func (self *Notifier) readData() {
	data, err := readLines(self.InitData.DataFilePath)
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}
	if len(data) == 0 {
		log.Fatalf("Data file is empty")
	}
	self.RawUserData = data
}

func (self *Notifier) parse() {
	var data []UserData
	for index := 0; index < len(self.RawUserData); index++ {
		items := strings.Split(self.RawUserData[index], ";")
		dataSet := make(map[string]string)
		for subIndex := 0; subIndex < len(items); subIndex++ {
			if len(strings.Trim(items[subIndex], " ")) == 0 {
				continue
			}
			keyValue := strings.Split(items[subIndex], "=")
			if len(keyValue) != 2 {
				log.Fatalf("Cannot read data file. Error in line (%i): %s", index+1, self.RawUserData[index])
				continue
			}
			dataSet[keyValue[0]] = keyValue[1]
		}
		data = append(data, UserData{Email: strings.Trim(dataSet[email], " "), UserName: dataSet[name]})
		fmt.Println(email + ": " + dataSet[email] + ", " + name + ": " + dataSet[name])
	}
	self.UserData = data
}

func CreteNotifier(initData InitData) *Notifier {
	notifier := Notifier{InitData: initData}
	return &notifier
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		lines = append(lines, text)
		fmt.Println("Read: " + text)
	}
	return lines, scanner.Err()
}
