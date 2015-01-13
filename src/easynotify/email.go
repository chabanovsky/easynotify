package main

import (
	"encoding/base64"
	"exp/utf8string"
	"log"
	"net/smtp"
)

type Message struct {
	Subject       string
	Message       string
	Receiver      string
	ReceiverEmail string
}

type SenderInfo struct {
	Name     string
	Email    string
	Password string
	Port     string
	Server   string
}

func ByteToBase64(value []byte) string {
	return base64.StdEncoding.EncodeToString(value)
}
func Base64ToByte(value string) []byte {
	bytes, _ := base64.StdEncoding.DecodeString(value)
	return bytes
}
func makeMessageHeader(senderEmail, receiverEmail, subj, from string) string {
	from = "From: =?UTF-8?B?" + ByteToBase64([]byte(from)) + "?= <" + senderEmail + ">\r\n"
	to := "To: " + receiverEmail + "\r\n"
	subject := "Subject: =?UTF-8?B?" + ByteToBase64([]byte(subj)) + "?=\r\n"
	mime := "MIME-Version: 1.0\r\n"
	contentType := "Content-Type: text/html; charset=utf-8\r\n"
	contentTransferEncoding := "Content-Transfer-Encoding: 8bit\r\n"

	return from + to + subject + mime + contentTransferEncoding + contentType + "\r\n"
}

func sendMessage(message *Message, senderInfo *SenderInfo) {

	if message == nil || message.ReceiverEmail == "" || senderInfo == nil {
		log.Fatal("SendMessageHelper: invalid input data")
	}

	sender := senderInfo.Email
	password := senderInfo.Password
	port := senderInfo.Port
	server := senderInfo.Server

	receiver := message.ReceiverEmail
	subject := message.Subject

	var msg = message.Message
	header := makeMessageHeader(senderInfo.Email, receiver, subject, senderInfo.Name)

	auth := smtp.PlainAuth(
		"",
		sender,
		password,
		server,
	)
	err := smtp.SendMail(
		server+":"+port,
		auth,
		sender,
		[]string{receiver},
		[]byte((utf8string.NewString(header + msg)).String()),
	)
	if err != nil {
		log.Fatalf("sendMessageHelper error: %s", err)
	}
}
