package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	riotapi "github.com/renja-g/riotgo/schemas/RiotAPI"
)

type Region string

const (
	Americas Region = "americas"
	Europe   Region = "europe"
	Asia     Region = "asia"
	Sea      Region = "sea"
	Esports  Region = "esports"

	Br1  Region = "br1"
	Eun1 Region = "eun1"
	Euw1 Region = "euw1"
	Jp1  Region = "jp1"
	Kr   Region = "kr"
	La1  Region = "la1"
	La2  Region = "la2"
	Me1  Region = "me1"
	Na1  Region = "na1"
	Oc1  Region = "oc1"
	Tr1  Region = "tr1"
	Ru   Region = "ru"
	Ph2  Region = "ph2"
	Sg2  Region = "sg2"
	Th2  Region = "th2"
	Tw2  Region = "tw2"
	Vn2  Region = "vn2"
)

type RiotAPIClient struct {
	*BaseClient
}

func NewRiotClient(apiKey string, opts ...Option) *RiotAPIClient {
	base := NewBaseClient(
		"https://%s.api.riotgames.com",
		append(opts,
			WithDefaultHeader("X-Riot-Token", apiKey),
		)...,
	)
	return &RiotAPIClient{BaseClient: base}
}

func (rc *RiotAPIClient) urlFor(r Region, path string) string {
	return fmt.Sprintf(rc.BaseURL, r) + path
}

// ----- Endpoints -----

// Account-V1

func (rc *RiotAPIClient) GetAccountV1ByPUUID(
	ctx context.Context,
	r Region,
	puuid string,
) (*riotapi.AccountV1Account, error) {
	return invokeJSON[riotapi.AccountV1Account](
		rc,
		ctx,
		r,
		http.MethodGet,
		"/riot/account/v1/accounts/by-puuid/{puuid}",
		map[string]string{"puuid": puuid},
		nil,
		nil,
		nil,
	)
}

func (rc *RiotAPIClient) GetAccountV1ByRiotID(
	ctx context.Context,
	r Region,
	gameName string,
	tagLine string,
) (*riotapi.AccountV1Account, error) {
	return invokeJSON[riotapi.AccountV1Account](
		rc,
		ctx,
		r,
		http.MethodGet,
		"/riot/account/v1/accounts/by-riot-id/{game_name}/{tag_line}",
		map[string]string{
			"game_name": gameName,
			"tag_line":  tagLine,
		},
		nil,
		nil,
		nil,
	)
}

func (rc *RiotAPIClient) GetAccountV1Me(
	ctx context.Context,
	r Region,
	authorization string,
) (*riotapi.AccountV1Account, error) {
	return invokeJSON[riotapi.AccountV1Account](
		rc,
		ctx,
		r,
		http.MethodGet,
		"/riot/account/v1/accounts/me",
		nil,
		nil,
		map[string]string{"Authorization": authorization},
		nil,
	)
}

func (rc *RiotAPIClient) GetAccountV1ActiveShardByPUUID(
	ctx context.Context,
	r Region,
	puuid string,
	game string,
) (*riotapi.AccountV1ActiveShard, error) {
	return invokeJSON[riotapi.AccountV1ActiveShard](
		rc,
		ctx,
		r,
		http.MethodGet,
		"/riot/account/v1/active-shards/by-game/{game}/by-puuid/{puuid}",
		map[string]string{"puuid": puuid, "game": game},
		nil,
		nil,
		nil,
	)
}

func expandPath(template string, params map[string]string) string {
	if len(params) == 0 {
		return template
	}
	for k, v := range params {
		template = strings.ReplaceAll(template, fmt.Sprintf("{%s}", k), v)
	}
	return template
}

func invokeJSON[T any](
	rc *RiotAPIClient,
	ctx context.Context,
	region Region,
	method, pathTemplate string,
	pathParams map[string]string,
	queries url.Values,
	headers map[string]string,
	body io.Reader,
) (*T, error) {

	// 1) Expand template variables in path.
	resolvedPath := expandPath(pathTemplate, pathParams)

	// 2) Build full URL with region placeholder.
	fullURL := rc.urlFor(region, resolvedPath)

	// 3) Dispatch HTTP request.
	resp, err := rc.Invoke(ctx, method, fullURL, body, headers, queries)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 4) Basic status-code validation. Upstream errors are surfaced to the caller
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	// 5) Decode JSON payload directly into caller-provided generic type.
	var v T
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}
