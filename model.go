package main

import (
	"sync"
)

type player struct {
	choice uint8
	chosen bool
	name   string
	model  *model
	sync.Mutex
}

type model struct {
	subscribers []func()
	players     map[string]*player
	sync.Mutex
}

func (m *model) subscribe(f func()) {
	m.Lock()
	defer m.Unlock()

	m.subscribers = append(m.subscribers, f)
}

func (m *model) notify() {
	m.Lock()
	defer m.Unlock()
	subscribers := m.subscribers

	for _, fn := range subscribers {
		fn()
	}
}

func (m *model) getPlayers() []player {
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

	return players
}

func (m *model) delPlayer(name string) {
	m.Lock()

	delete(m.players, name)

	m.Unlock()

	m.notify()
}

func (m *model) addPlayer(p *player) {
	m.Lock()

	p.model = m
	m.players[p.name] = p

	m.Unlock()

	m.notify()
}

func newModel() *model {
	var m model

	m.subscribers = make([]func(), 0)
	m.players = make(map[string]*player)

	return &m
}

func (p *player) hasChosen() bool {
	p.Lock()
	defer p.Unlock()

	return p.chosen
}

func (p *player) getChoice() uint8 {
	p.Lock()
	defer p.Unlock()

	return p.choice
}

func (p *player) setChoice(choice uint8) {
	p.Lock()
	defer p.Unlock()

	p.choice = choice
	p.chosen = true

	p.model.notify()
}

func (p *player) getName() string {
	p.Lock()
	defer p.Unlock()

	name := p.name

	return name
}

func (p *player) setName(name string) {
	p.Lock()
	defer p.Unlock()

	p.name = name

	if p.model != nil {
		p.model.notify()
	}
}
