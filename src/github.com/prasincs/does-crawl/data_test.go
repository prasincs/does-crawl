package main

import (
	"github.com/stretchr/testify/assert"
	//"strings"
	"testing"
)

func TestDBInitialization(t *testing.T) {
	InitializeDB()
	assert.NotNil(t, db)
}

func TestDBInsertion(t *testing.T) {
	id, err := db.Add(&Url{Id: 1, Link: "google.com", Parent: ""})
	assert.Nil(t, err)
	assert.Equal(t, 1, id, "Id does not look right")
}

func TestDBCount(t *testing.T) {
	id, err := db.Add(&Url{Id: 1, Link: "google.com", Parent: ""})
	assert.Nil(t, err)
	assert.Equal(t, 1, id, "Id does not look right")
	count := db.GetCount()
	assert.Equal(t, 1, count, "Not the right count")
}
