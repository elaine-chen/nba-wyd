package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	baseURL        = "http://data.nba.com/data"
	scoreboardURL  = "/10s/prod/v2/%s/scoreboard.json"
	genGamesURL    = "/10s/json/cms/noseason/scoreboard/%s/games.json"
	teamGamesURL   = "/10s/prod/v1/%d/teams/%s/schedule.json"
	boxscoreURL    = "/10s/prod/v1/%s/%s_mini_boxscore.json"
	standingsURL   = "/10s/prod/v1/%s/standings_conference.json"
	leagueTeamsURL = "/10s/prod/v1/%d/teams.json"
)

type Game struct {
	Date string `json:"gdate"`
}

const dateFormat = "20060102"

type NBAClient struct {
	client *http.Client
}

func NewNBAClient() *NBAClient {
	return &NBAClient{client: http.DefaultClient}
}

func (nc *NBAClient) getJSON(ctx context.Context, url string, data interface{}) error {
	log.Println(url)
	resp, err := nc.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 http.StatusCode %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return err
	}

	return nil
}

type LeagueTeamResponse struct {
	League struct {
		Standard []LeagueTeam `json:"standard"`
	} `json:"league"`
}

type LeagueTeam struct {
	TriCode  string `json:"tricode"`
	Nickname string `json:"nickname"`
	URLName  string `json:"urlName"`
	FullName string `json:"fullName"`
	TeamID   string `json:"teamId"`
}

func (nc *NBAClient) GetLeagueTeams(ctx context.Context, year int) ([]LeagueTeam, error) {
	var resp LeagueTeamResponse
	err := nc.getJSON(ctx, baseURL+fmt.Sprintf(leagueTeamsURL, year), &resp)
	return resp.League.Standard, err
}

type MiniBoxScoreResponse struct {
	BasicGameData MiniBoxScore `json:"basicGameData"`
}
type MiniBoxScore interface{}

func (nc *NBAClient) GetMiniBoxScore(ctx context.Context, date time.Time, gameID string) (MiniBoxScore, error) {
	var resp MiniBoxScoreResponse
	err := nc.getJSON(
		ctx,
		baseURL+fmt.Sprintf(boxscoreURL, date.Format(dateFormat), gameID),
		&resp,
	)
	return resp.BasicGameData, err
}

type TeamInfoResponse struct {
	League struct {
		Standard TeamInfo `json:"standard"`
	} `json:"league"`
}
type TeamInfo interface{}

func (nc *NBAClient) GetTeamInfo(ctx context.Context, year int, teamID string) (TeamInfo, error) {
	var resp TeamInfoResponse
	err := nc.getJSON(
		ctx,
		baseURL+fmt.Sprintf(teamGamesURL, year, teamID),
		&resp,
	)
	return resp.League.Standard, err
}
