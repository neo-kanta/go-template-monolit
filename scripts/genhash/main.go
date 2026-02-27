package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	users := []struct{ user, pass string }{
		{"admin", "admin123"},
		{"user", "user123"},
		{"manager1", "manager123"},
		{"manager2", "manager234"},
		{"maker1", "maker123"},
		{"maker2", "maker234"},
		{"checker1", "checker123"},
		{"checker2", "checker234"},
		{"viewer1", "viewer123"},
		{"viewer2", "viewer234"},
		{"auditor1", "auditor123"},
		{"auditor2", "auditor234"},
	}
	for _, u := range users {
		h, _ := bcrypt.GenerateFromPassword([]byte(u.pass), bcrypt.DefaultCost)
		fmt.Printf("-- %s / %s\n", u.user, u.pass)
		fmt.Printf("INSERT INTO sys_users (username, password, role, created_by, updated_by) VALUES ('%s', '%s', ", u.user, string(h))
		switch u.user {
		case "admin":
			fmt.Print("'admin'")
		case "user":
			fmt.Print("'user'")
		case "manager1", "manager2":
			fmt.Print("'manager'")
		case "maker1", "maker2":
			fmt.Print("'maker'")
		case "checker1", "checker2":
			fmt.Print("'checker'")
		case "viewer1", "viewer2":
			fmt.Print("'viewer'")
		case "auditor1", "auditor2":
			fmt.Print("'auditor'")
		}
		fmt.Println(", 'system', 'system');")
	}
}
