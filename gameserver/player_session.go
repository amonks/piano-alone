package main

import "monks.co/piano-alone/pianists"

type PlayerSession struct {
	Fingerprint string
	Name        string
}

func NewPlayerSession(fingerprint string) *PlayerSession {
	return &PlayerSession{
		Fingerprint: fingerprint,
		Name:        pianists.Hash(fingerprint),
	}
}
