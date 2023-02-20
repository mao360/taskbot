package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Room struct {
	Name      string
	Ways      []string
	Inventory map[string]bool
	Actions   map[string]func() string
	IsOpen    bool
}

type World struct {
	Rooms []Room
	P     Player
}

type Player struct {
	CurRoom   Room
	Inventory map[string]bool
}

func (w *World) GetRoom(s string) Room {
	for i := range w.Rooms {
		if w.Rooms[i].Name == s {
			return w.Rooms[i]
		}
	}
	return Room{}
}

func (w *World) LookAround() string { // осмотреться
	return w.P.CurRoom.Actions["осмотреться"]()
}

func (w *World) GoTo(s string) string {
	for i := range w.P.CurRoom.Ways {
		if w.P.CurRoom.Ways[i] == s && w.GetRoom(s).IsOpen == true {
			w.P.CurRoom = w.GetRoom(s)
			return w.P.CurRoom.Actions["идти"]()
		} else if w.P.CurRoom.Ways[i] == s && w.GetRoom(s).IsOpen == false {
			return "дверь закрыта"
		}
	}
	return fmt.Sprintf("нет пути в %s", s)
}

func (w *World) PutOnItem(s string) string {
	w.P.Inventory[s] = true
	w.P.CurRoom.Inventory[s] = false
	return fmt.Sprintf("вы надели: %s", s)
}

func (w *World) TakeItem(s string) string { // взять
	if w.P.CurRoom.Inventory[s] == true && w.P.Inventory["рюкзак"] == true {
		w.P.CurRoom.Inventory[s] = false
		w.P.Inventory[s] = true
		return fmt.Sprintf("предмет добавлен в инвентарь: %s", s)
	} else if w.P.CurRoom.Inventory[s] == false && w.P.Inventory["рюкзак"] == true {
		return "нет такого"
	}
	return "некуда класть"
}

func (w *World) UseItem(s1, s2 string) string { // применить

	if w.P.Inventory[s1] == false {
		return fmt.Sprintf("нет предмета в инвентаре - %s", s1)
	} else if w.P.Inventory[s1] == true && s2 == "дверь" && w.P.CurRoom.Name == "коридор" {
		for i := range w.Rooms {
			if w.Rooms[i].Name == "улица" {
				w.Rooms[i].IsOpen = true
			}
		}
		return "дверь открыта"
	}
	return "не к чему применить"
}

var w = &World{
	Rooms: []Room{},
}

func main() {
	initGame()

	var command string
	var err error
	for {
		command, err = bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		command = strings.TrimSuffix(command, "\n")
		fmt.Println(handleCommand(command))
	}
}

func initGame() {
	var (
		corridor Room
		kitchen  Room
		bedroom  Room
		street   Room
		player   Player
	)

	w.Rooms = []Room{}
	w.P = Player{}

	bedroom.Name = "комната"
	bedroom.IsOpen = true
	bedroom.Inventory = map[string]bool{
		"ключи":     true,
		"конспекты": true,
		"рюкзак":    true,
	}

	bedroom.Actions = map[string]func() string{
		"осмотреться": func() string {
			if w.P.Inventory["ключи"] == false && w.P.Inventory["конспекты"] == false && w.P.Inventory["рюкзак"] == false {
				return "на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор"
			} else if w.P.Inventory["ключи"] == false && w.P.Inventory["конспекты"] == false && w.P.Inventory["рюкзак"] == true {
				return "на столе: ключи, конспекты. можно пройти - коридор"
			} else if w.P.Inventory["ключи"] == false && w.P.Inventory["конспекты"] == true {
				return ""
			} else if w.P.Inventory["ключи"] == true && w.P.Inventory["конспекты"] == false {
				return "на столе: конспекты. можно пройти - коридор"
			}
			return "пустая комната. можно пройти - коридор"
		},
		"идти": func() string {
			return "ты в своей комнате. можно пройти - коридор"
		},
	}

	kitchen.Name = "кухня"
	kitchen.IsOpen = true
	kitchen.Inventory = map[string]bool{}

	kitchen.Actions = map[string]func() string{
		"осмотреться": func() string {
			if w.P.Inventory["рюкзак"] == true {
				return "ты находишься на кухне, на столе: чай, надо идти в универ. можно пройти - коридор"
			}
			return "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"
		},
		"идти": func() string {
			return "кухня, ничего интересного. можно пройти - коридор"
		},
	}

	street.Name = "улица"
	street.IsOpen = false
	street.Inventory = map[string]bool{}

	street.Actions = map[string]func() string{
		"осмотреться": func() string {
			return ""
		},
		"идти": func() string {
			return "на улице весна. можно пройти - домой"
		},
	}

	corridor.Name = "коридор"
	corridor.IsOpen = true
	corridor.Inventory = map[string]bool{}
	corridor.Actions = map[string]func() string{
		"идти": func() string {
			return "ничего интересного. можно пройти - кухня, комната, улица"
		},
	}

	player.Inventory = map[string]bool{
		"ключи":     false,
		"конспекты": false,
		"рюкзак":    false,
	}

	corridor.Ways = append(corridor.Ways, "улица", "кухня", "комната")
	street.Ways = append(street.Ways, "коридор")
	kitchen.Ways = append(kitchen.Ways, "коридор")
	bedroom.Ways = append(bedroom.Ways, "коридор")

	player.CurRoom = kitchen
	w.Rooms = append(w.Rooms, bedroom, kitchen, corridor, street)
	w.P = player
}

func handleCommand(command string) string {

	parsedCommand := strings.Split(command, " ")
	str := ""

	switch parsedCommand[0] {
	case "осмотреться":
		str = w.LookAround()
	case "идти":
		str = w.GoTo(parsedCommand[1])
	case "применить":
		str = w.UseItem(parsedCommand[1], parsedCommand[2])
	case "взять":
		str = w.TakeItem(parsedCommand[1])
	case "надеть":
		str = w.PutOnItem(parsedCommand[1])
	default:
		str = "неизвестная команда"
	}

	return str
}
