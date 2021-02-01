package idgenerator

import (
	"github.com/google/uuid"
)

type IDGenerator struct {
}

func (g *IDGenerator) GenerateID() (id string) {
	return uuid.New().String()
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}
