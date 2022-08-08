package main

import (
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
	sync.Mutex
}

func (p *player) subscribe(f func()) {
	p.Lock()
	defer p.Unlock()

	p.subscribers = append(p.subscribers, f)
}

func (p *player) notify() {
	p.Lock()
	defer p.Unlock()
	subscribers := p.subscribers

	for _, fn := range subscribers {
		fn()
	}
}

func (m *model) clearChoices() {
	m.Lock()

	for _, p := range m.players {
		p.Lock()
		defer p.Unlock()

		p.chosen = false
		p.choice = 0
	}

	m.Unlock()

	m.notify()
}

func (m *model) subscribe(f func()) {
	m.Lock()
	defer m.Unlock()

	m.subscribers = append(m.subscribers, f)
}

func (m *model) getDisclosed() bool {
	m.Lock()
	defer m.Unlock()

	return m.disclosed
}

func (m *model) setDisclose(r bool) {
	m.Lock()

	m.disclosed = r
	m.Unlock()

	m.notify()
}

func (m *model) toggleDisclose() {
	m.Lock()

	m.disclosed = !m.disclosed
	m.Unlock()

	m.notify()
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

	sort.SliceStable(players, func(i, j int) bool {
		return players[i].name < players[j].name
	})
	return players
}

func (m *model) delPlayer(name string) {
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
