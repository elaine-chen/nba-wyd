package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNBAClient(t *testing.T) {
	ctx := context.Background()
	nc := NewNBAClient()
	lastYear := time.Now().AddDate(-1, 0, 0)
	teams, err := nc.GetLeagueTeams(ctx, lastYear.Year())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(nc.GetTeamInfo(ctx, lastYear.Year(), teams[0].TeamID))
}
