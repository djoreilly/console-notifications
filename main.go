package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/term"
)

type Notification struct {
	// org.freedesktop.Notifications.Notify message fields
	// https://specifications.freedesktop.org/notification-spec/latest/ar01s09.html
	// Only need the string fields.
	appName string
	summary string
	body    string
}

func (n Notification) String() string {
	return n.appName + "\n\n" + n.summary + "\n\n" + n.body + "\n\n"
}

func (n Notification) isEmpty() bool {
	return n.appName == "" && n.summary == "" && n.body == ""
}

func getNotification(msg *dbus.Message) Notification {
	var note Notification
	if len(msg.Body) < 5 {
		return note
	}
	if appName, ok := msg.Body[0].(string); ok {
		note.appName = appName
	}
	if summary, ok := msg.Body[3].(string); ok {
		note.summary = summary
	}
	if body, ok := msg.Body[4].(string); ok {
		note.body = body
	}
	return note
}

func printNotifications() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	rules := []string{
		"type='method_call',member='Notify',path='/org/freedesktop/Notifications',interface='org.freedesktop.Notifications',destination='org.freedesktop.Notifications'",
	}
	var flag uint = 0
	call := conn.BusObject().Call("org.freedesktop.DBus.Monitoring.BecomeMonitor", 0, rules, flag)
	if call.Err != nil {
		fmt.Fprintln(os.Stderr, "Failed to become monitor:", call.Err)
		os.Exit(1)
	}

	msgCh := make(chan *dbus.Message, 100)
	conn.Eavesdrop(msgCh)

	for msg := range msgCh {
		note := getNotification(msg)
		if note.isEmpty() {
			continue
		}
		output := time.Now().Format(time.DateTime) + " - " + note.String()

		cols := screenWidth()
		fmt.Print(wordwrap.WrapString(output, uint(cols)))
		fmt.Println(strings.Repeat("-", cols))
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running clear: %s", err)
	}
}

func screenWidth() int {
	width, _, err := term.GetSize(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in getColumns: %s", err)
		return 0
	}
	return width
}

func main() {
	go printNotifications()
	for {
		clearScreen()
		fmt.Println("Monitoring notifications. Enter to clear screen. Ctrl-c to quit.")
		fmt.Println(strings.Repeat("-", screenWidth()))
		fmt.Scanln() // block for Enter
	}
}
