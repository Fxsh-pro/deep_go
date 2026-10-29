package main

type Option func(*GamePerson)

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	houseFlag uint8 = 1 << iota
	gunFlag
	familyFlag
)

const (
	personTypeShift        = 3
	personTypeMask  uint8  = 0b11 << personTypeShift
	statMask        uint16 = 0xF
	resourceMask    uint32 = (1 << 10) - 1
)

type GamePerson struct {
	x, y, z   int32    // 12B
	gold      uint32   // 4B
	stats     uint16   // strength, experience, level, respect: 2B
	resources [3]byte  // mana and health: 3B
	flags     uint8    // house, gun, family, player type (2 bit): 1B
	name      [42]byte // 42B
}

func NewGamePerson(options ...Option) GamePerson {
	var person GamePerson
	for _, option := range options {
		option(&person)
	}
	return person
}

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		clear(person.name[:])
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x, person.y, person.z = int32(x), int32(y), int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.setResource(mana, 0)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.setResource(health, 10)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.setStat(respect, 12)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.setStat(strength, 0)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.setStat(experience, 4)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.setStat(level, 8)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= houseFlag
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= gunFlag
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= familyFlag
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags = person.flags&^personTypeMask | (uint8(personType)<<personTypeShift)&personTypeMask
	}
}

func (p *GamePerson) Name() string {
	length := 0
	for length < len(p.name) && p.name[length] != 0 {
		length++
	}
	return string(p.name[:length])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.resourceBits() & resourceMask)
}

func (p *GamePerson) Health() int {
	return int((p.resourceBits() >> 10) & resourceMask)
}

func (p *GamePerson) Respect() int {
	return int((p.stats >> 12) & statMask)
}

func (p *GamePerson) Strength() int {
	return int(p.stats & statMask)
}

func (p *GamePerson) Experience() int {
	return int((p.stats >> 4) & statMask)
}

func (p *GamePerson) Level() int {
	return int((p.stats >> 8) & statMask)
}

func (p *GamePerson) HasHouse() bool {
	return p.flags&houseFlag != 0
}

func (p *GamePerson) HasGun() bool {
	return p.flags&gunFlag != 0
}

func (p *GamePerson) HasFamilty() bool {
	return p.flags&familyFlag != 0
}

func (p *GamePerson) Type() int {
	return int((p.flags & personTypeMask) >> personTypeShift)
}

func (p *GamePerson) setStat(value int, shift uint) {
	mask := statMask << shift
	p.stats = p.stats&^mask | (uint16(value)&statMask)<<shift
}

func (p *GamePerson) resourceBits() uint32 {
	return uint32(p.resources[0]) | uint32(p.resources[1])<<8 | uint32(p.resources[2])<<16
}

func (p *GamePerson) setResource(value int, shift uint) {
	bits := p.resourceBits()
	mask := resourceMask << shift
	bits = bits&^mask | (uint32(value)&resourceMask)<<shift
	p.resources = [3]byte{byte(bits), byte(bits >> 8), byte(bits >> 16)}
}
