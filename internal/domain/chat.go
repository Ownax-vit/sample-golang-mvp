package domain

import "time"


type Chat struct{
	ID 		  int 			`json:"Id"       example:"125216"`
	Title 	  string		`json:"Title"    example:"Тестовый чат"`
	CreatedAt time.Time		`json:"source"`
	Messages []Message		`json:"messages"`
}
