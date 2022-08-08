package main

import (
	"log"
	"sort"
	"sync"
)

type player struct {
	subscribers []func()
	choice      uint8
	chosen      bool
	name        string
	model       *model
	sync.Mutex
}

type model struct {
	subscribers []func()
	players     map[string]*player
	disclosed   bool
	chosen      bool
	sync.Mutex
}

func (p *player) subscribe(f func()) {
	log.Printf("subscribe\n")
	p.Lock()
	defer p.Unlock()

	p.subscribers = append(p.subscribers, f)
}

func (p *player) notify() {
	log.Printf("notify\n")
	p.Lock()
	defer p.Unlock()
	subscribers := p.subscribers

	for _, fn := range subscribers {
		fn()
	}
}

func (m *model) getAverageChoice() float32 {
	log.Printf("getAverageChoice\n")
	m.Lock()
	defer m.Unlock()

	var sum float32
	var count float32

	count = 0
	sum = 0

	for _, p := range m.players {
		if p.chosen {
			count++
			sum += float32(p.choice)
		}
	}

	var average float32

	average = sum / count

	return average
}

func (m *model) clearChoices() {
	log.Printf("clearChoices\n")
	m.Lock()

	for _, p := range m.players {
		p.Lock()
		defer p.Unlock()

		p.chosen = false
		p.choice = 0
	}

	m.chosen = false
	m.Unlock()

	m.notify()
}

func (m *model) subscribe(f func()) {
	log.Printf("subscribe\n")
	m.Lock()
	defer m.Unlock()

	m.subscribers = append(m.subscribers, f)
}

func (m *model) hasChosen() bool {
	log.Printf("getChosen\n")
	m.Lock()
	defer m.Unlock()

	return m.chosen
}

func (m *model) getDisclosed() bool {
	log.Printf("getDisclosed\n")
	m.Lock()
	defer m.Unlock()

	return m.disclosed
}

func (m *model) setDisclose(r bool) {
	log.Printf("setDisclose")
	m.Lock()

	m.disclosed = r
	m.Unlock()

	m.notify()
}

func (m *model) toggleDisclose() {
	log.Printf("toggleDisclose")
	m.Lock()

	m.disclosed = !m.disclosed
	m.Unlock()

	m.notify()
}

func (m *model) notify() {
	log.Printf("notify")
	m.Lock()
	subscribers := m.subscribers
	m.Unlock()

	for _, fn := range subscribers {
		fn()
	}
}

func (m *model) getPlayers() []player {
	log.Printf("getPlayers")
	m.Lock()
	defer m.Unlock()

	players := make([]player, len(m.players))

	i := 0
	for _, p := range m.players {
		var newPlayer player

		p.Lock()
		newPlayer.choice = p.choice
		newPlayer.chosen = p.chosen
		newPlayer.name = p.name
		newPlayer.model = m
		p.Unlock()

		players[i] = newPlayer
		i++
	}

	sort.SliceStable(players, func(i, j int) bool {
		return players[i].name < players[j].name
	})
	return players
}

func (m *model) delPlayer(name string) {
	log.Printf("delPlayer")
	m.Lock()

	player, ok := m.players[name]
	if !ok {
		return
	}

	player.notify()

	delete(m.players, name)

	m.Unlock()

	m.notify()
}

func (m *model) addPlayer(p *player) bool {
	log.Printf("addPlayer")
	m.Lock()

	p.model = m

	_, ok := m.players[p.name]
	if ok {
		m.Unlock()
		return false
	}

	m.players[p.name] = p

	m.Unlock()

	m.notify()

	return true
}

func newModel() *model {
	var m model

	m.subscribers = make([]func(), 0)
	m.players = make(map[string]*player)

	return &m
}

func (p *player) hasChosen() bool {
	log.Printf("hasChosen")
	p.Lock()
	defer p.Unlock()

	return p.chosen
}

func (p *player) getChoice() uint8 {
	log.Printf("getChoice")
	p.Lock()
	defer p.Unlock()

	return p.choice
}

func (p *player) setChoice(choice uint8) {
	log.Printf("setChoice")
	p.Lock()
	defer p.Unlock()

	p.choice = choice
	p.chosen = true

	p.model.Lock()
	p.model.chosen = true
	p.model.Unlock()

	p.model.notify()
}

func (p *player) getName() string {
	log.Printf("getName")
	p.Lock()
	defer p.Unlock()

	name := p.name

	return name
}

func (p *player) setName(name string) {
	log.Printf("setName")
	p.Lock()
	defer p.Unlock()

	p.name = name

	if p.model != nil {
		p.model.notify()
	}
}
