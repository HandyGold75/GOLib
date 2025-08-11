package yts

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type (
	searchFilter string
	sortOrder    string

	path []any

	thumbnail struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	}
	channel struct {
		ID         string      `json:"id"`
		Title      string      `json:"title"`
		Thumbnails []thumbnail `json:"thumbnails"`
		URL        string      `json:"url"`
	}
	videoItem struct {
		ID            string      `json:"id"`
		Title         string      `json:"title"`
		PublishedTime string      `json:"publishedTime"`
		Duration      int         `json:"duration"`
		ViewCount     int         `json:"viewCount"`
		Thumbnails    []thumbnail `json:"thumbnails"`
		RichThumbnail thumbnail   `json:"richThumbnail"`
		Description   string      `json:"description"`
		Channel       channel     `json:"channel"`
		URL           string      `json:"url"`
	}
	channelItem struct {
		ID          string      `json:"id"`
		Title       string      `json:"title"`
		Thumbnails  []thumbnail `json:"thumbnails"`
		VideoCount  int         `json:"videoCount"`
		Description string      `json:"description"`
		Subscribers string      `json:"subscribers"`
		URL         string      `json:"url"`
	}
	playlistItem struct {
		ID         string      `json:"id"`
		Title      string      `json:"title"`
		VideoCount int         `json:"videoCount"`
		Channel    channel     `json:"channel"`
		Thumbnails []thumbnail `json:"thumbnails"`
		URL        string      `json:"url"`
	}
	shelfItem struct {
		Title string       `json:"title"`
		Items []*videoItem `json:"items"`
	}
	searchResult struct {
		EstimatedResults int             `json:"estimatedResults"`
		Videos           []*videoItem    `json:"videos"`
		Channels         []*channelItem  `json:"channels"`
		Playlists        []*playlistItem `json:"playlists"`
		Shelves          []*shelfItem    `json:"shelves"`
		Suggestions      []string        `json:"suggestions"`
	}

	SearchClient struct {
		Query string

		Language string
		Region   string

		SearchFilter searchFilter
		SortOrder    sortOrder

		// CustomParams allows you to copy and paste params from YouTube.
		CustomParams string

		HTTPClient *http.Client

		continuationKey string
		newPage         bool
	}
)

const (
	searchKey = "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"

	FilterAll      searchFilter = "%253D"
	FilterVideo    searchFilter = "SAhAB"
	FilterChannel  searchFilter = "SAhAC"
	FilterPlaylist searchFilter = "SAhAD"

	OrderRelevance  sortOrder = "CAA"
	OrderUploadDate sortOrder = "CAI"
	OrderViewCount  sortOrder = "CAM"
	OrderRating     sortOrder = "CAE"
)

// Create a new SearchClient with default parameters.
func NewSearch(query string) *SearchClient {
	return &SearchClient{
		Query:        query,
		Language:     "en",
		Region:       "US",
		SearchFilter: FilterAll,
		SortOrder:    OrderRelevance,
		HTTPClient:   &http.Client{},
	}
}

// Create a new SearchClient for video search.
func NewSearchVideo(query string) *SearchClient {
	return &SearchClient{
		Query:        query,
		Language:     "en",
		Region:       "US",
		SearchFilter: FilterVideo,
		SortOrder:    OrderRelevance,
		HTTPClient:   &http.Client{},
	}
}

// Create a new SearchClient for channel search.
func NewSearchChannel(query string) *SearchClient {
	return &SearchClient{
		Query:        query,
		Language:     "en",
		Region:       "US",
		SearchFilter: FilterChannel,
		SortOrder:    OrderRelevance,
		HTTPClient:   &http.Client{},
	}
}

// Create a new SearchClient for playlist search.
func NewSearchPlaylist(query string) *SearchClient {
	return &SearchClient{
		Query:        query,
		Language:     "en",
		Region:       "US",
		SearchFilter: FilterPlaylist,
		SortOrder:    OrderRelevance,
		HTTPClient:   &http.Client{},
	}
}

// NextExists returns whether the Next call will return new content.
func (search *SearchClient) NextExists() bool {
	if !search.newPage {
		return true
	}
	return search.continuationKey != ""
}

// Next returns content from the next page.
func (search *SearchClient) Next() (*searchResult, error) {
	if !search.NextExists() {
		return nil, errors.New("page does not exist")
	}
	response, err := search.makeReq()
	if err != nil {
		return nil, err
	}
	responseSource, continuationKey, estimatedResults := parseSource(response, search.newPage)
	result := parseComponents(responseSource)
	result.EstimatedResults = estimatedResults
	result.Suggestions = parseSuggestions(response)

	search.continuationKey = continuationKey
	search.newPage = true
	return result, nil
}

func (search *SearchClient) makeReq() (map[string]any, error) {
	payload := map[string]any{
		"query": search.Query,
		"context": map[string]map[string]any{
			"client": {
				"clientName":       "WEB",
				"clientVersion":    "2.20210224.06.00",
				"newVisitorCookie": true,
			},
			"user": {
				"lockedSafetyMode": false,
			},
		},
	}
	if search.CustomParams == "" {
		payload["params"] = string(search.SortOrder) + string(search.SearchFilter)
	} else {
		payload["params"] = search.CustomParams
	}
	clientData := payload["context"].(map[string]map[string]any)["client"]
	clientData["hl"] = search.Language
	clientData["gl"] = search.Region
	if search.continuationKey != "" {
		payload["continuation"] = search.continuationKey
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	url := "https://www.youtube.com/youtubei/v1/search?key=" + searchKey
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	request.Header = map[string][]string{
		"Content-Type": {"application/json; charset=utf-8"},
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"},
	}

	response, err := search.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result map[string]any
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func durationToInt(duration any) int {
	items := strings.Split(duration.(string), ":")
	result := 0
	for i := range len(items) {
		durationInt, _ := strconv.Atoi(items[i])
		result += durationInt * int(math.Pow(60, float64(len(items)-i-1)))
	}
	return result
}

func viewCountToInt(viewCount any) int {
	str := strings.Split(viewCount.(string), " ")[0]
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, "\u00A0", "")
	result, _ := strconv.Atoi(str)
	return result
}

func descriptionToStr(description []any) string {
	result := make([]string, len(description))
	for i := range len(description) {
		result[i] = description[i].(map[string]any)["text"].(string)
	}
	return strings.Join(result, "")
}

func getValue(source any, path path) any {
	value := source
	for _, element := range path {
		mustBreak := false
		switch element.(type) {
		case string:
			if val, ok := value.(map[string]any)[element.(string)]; ok {
				value = val
			} else {
				value = nil
				mustBreak = true
			}
		case int:
			if len(value.([]any)) != 0 {
				value = value.([]any)[element.(int)]
			} else {
				value = nil
				mustBreak = true
			}
		}
		if mustBreak {
			break
		}
	}
	return value
}

func parseSource(response any, newPage bool) ([]map[string]any, string, int) {
	var responseContent []any
	if !newPage {
		content := getValue(response, path{"contents", "twoColumnSearchResultsRenderer", "primaryContents", "sectionListRenderer", "contents"})
		if content != nil {
			responseContent = content.([]any)
		} else {
			return nil, "", 0
		}
	} else {
		content := getValue(response, path{"onResponseReceivedCommands", 0, "appendContinuationItemsAction", "continuationItems"})
		if content != nil {
			responseContent = content.([]any)
		} else {
			return nil, "", 0
		}
	}

	responseContentMaps := make([]map[string]any, len(responseContent))
	for index, value := range responseContent {
		responseContentMaps[index] = value.(map[string]any)
	}
	var responseSource []map[string]any
	var continuationKey string
	if responseContent != nil {
		for _, element := range responseContentMaps {
			if _, ok := element["itemSectionRenderer"]; ok {
				newSource := getValue(element, path{"itemSectionRenderer", "contents"}).([]any)
				responseSource = responseSource[:0]
				for _, value := range newSource {
					responseSource = append(responseSource, value.(map[string]any))
				}
			}
			if _, ok := element["continuationItemRenderer"]; ok {
				continuationKey = getValue(element, path{"continuationItemRenderer", "continuationEndpoint", "continuationCommand", "token"}).(string)
			}
		}
	} else {
		responseSource = getValue(responseContent, path{"contents", "twoColumnSearchResultsRenderer", "primaryContents", "richGridRenderer", "contents"}).([]map[string]any)
		continuationKey = getValue(responseSource, path{"continuationItemRenderer", "continuationEndpoint", "continuationCommand", "token"}).(string)
	}

	estimatedResults, _ := strconv.Atoi(
		getValue(response, path{"estimatedResults"}).(string),
	)
	return responseSource, continuationKey, estimatedResults
}

func parseComponents(responseSource []map[string]any) *searchResult {
	result := &searchResult{}
	for _, element := range responseSource {
		if videoElement, ok := element["videoRenderer"]; ok {
			videoComponent := parseVideoComponent(videoElement.(map[string]any))
			result.Videos = append(result.Videos, videoComponent)
			continue
		}
		if channelElement, ok := element["channelRenderer"]; ok {
			channelComponent := parseChannelComponent(channelElement.(map[string]any))
			result.Channels = append(result.Channels, channelComponent)
			continue
		}
		if playlistElement, ok := element["playlistRenderer"]; ok {
			playlistComponent := parsePlaylistComponent(playlistElement.(map[string]any))
			result.Playlists = append(result.Playlists, playlistComponent)
			continue
		}
		if shelfElement, ok := element["shelfRenderer"]; ok {
			shelfComponent := parseShelfComponent(shelfElement.(map[string]any))
			result.Shelves = append(result.Shelves, shelfComponent)
			continue
		}
		if richItemElement, ok := element["richItemRenderer"]; ok {
			if richItemElementContent, ok := richItemElement.(map[string]any)["content"]; ok {
				if videoElement, ok := richItemElementContent.(map[string]any)["videoRenderer"]; ok {
					videoComponent := parseVideoComponent(videoElement.(map[string]any))
					result.Videos = append(result.Videos, videoComponent)
				}
			}
			continue
		}
	}
	return result
}

func parseVideoComponent(video map[string]any) *videoItem {
	item := &videoItem{}
	if id := getValue(video, path{"videoId"}); id != nil {
		item.ID = id.(string)
		item.URL = "https://www.youtube.com/watch?v=" + item.ID
	}
	if title := getValue(video, path{"title", "runs", 0, "text"}); title != nil {
		item.Title = title.(string)
	}
	if publishedTime := getValue(video, path{"publishedTimeText", "simpleText"}); publishedTime != nil {
		item.PublishedTime = publishedTime.(string)
	}
	if duration := getValue(video, path{"lengthText", "simpleText"}); duration != nil {
		item.Duration = durationToInt(duration)
	}
	if viewCount := getValue(video, path{"viewCountText", "simpleText"}); viewCount != nil {
		item.ViewCount = viewCountToInt(viewCount)
	}
	if thumbnails := getValue(video, path{"thumbnail", "thumbnails"}); thumbnails != nil {
		for _, thumb := range thumbnails.([]any) {
			item.Thumbnails = append(item.Thumbnails, thumbnail{
				URL:    thumb.(map[string]any)["url"].(string),
				Height: int(thumb.(map[string]any)["height"].(float64)), Width: int(thumb.(map[string]any)["width"].(float64)),
			})
		}
	}
	if richThumbnail := getValue(video, path{"richThumbnail", "movingThumbnailRenderer", "movingThumbnailDetails", "thumbnails", 0}); richThumbnail != nil {
		item.RichThumbnail = thumbnail{
			URL:    richThumbnail.(map[string]any)["url"].(string),
			Height: int(richThumbnail.(map[string]any)["height"].(float64)), Width: int(richThumbnail.(map[string]any)["width"].(float64)),
		}
	}
	if description := getValue(video, path{"detailedMetadataSnippets", 0, "snippetText", "runs"}); description != nil {
		item.Description = descriptionToStr(description.([]any))
	}
	item.Channel = channel{}
	if channelTitle := getValue(video, path{"ownerText", "runs", 0, "text"}); channelTitle != nil {
		item.Channel.Title = channelTitle.(string)
	}
	if channelID := getValue(video, path{"ownerText", "runs", 0, "navigationEndpoint", "browseEndpoint", "browseId"}); channelID != nil {
		item.Channel.ID = channelID.(string)
		item.Channel.URL = "https://www.youtube.com/channel/" + item.Channel.ID
	}
	if channelThumbnails := getValue(video, path{"channelThumbnailSupportedRenderers", "channelThumbnailWithLinkRenderer", "thumbnail", "thumbnails"}); channelThumbnails != nil {
		for _, thumb := range channelThumbnails.([]any) {
			item.Channel.Thumbnails = append(item.Channel.Thumbnails, thumbnail{
				URL:    thumb.(map[string]any)["url"].(string),
				Height: int(thumb.(map[string]any)["height"].(float64)), Width: int(thumb.(map[string]any)["width"].(float64)),
			})
		}
	}
	return item
}

func parseChannelComponent(chann map[string]any) *channelItem {
	item := &channelItem{}
	if id := getValue(chann, path{"channelId"}); id != nil {
		item.ID = id.(string)
		item.URL = "https://www.youtube.com/channel/" + item.ID
	}
	if title := getValue(chann, path{"title", "simpleText"}); title != nil {
		item.Title = title.(string)
	}
	if thumbnails := getValue(chann, path{"thumbnail", "thumbnails"}); thumbnails != nil {
		for _, thumb := range thumbnails.([]any) {
			item.Thumbnails = append(item.Thumbnails, thumbnail{
				URL:    "http:" + thumb.(map[string]any)["url"].(string),
				Height: int(thumb.(map[string]any)["height"].(float64)), Width: int(thumb.(map[string]any)["width"].(float64)),
			})
		}
	}
	if videoCount := getValue(chann, path{"videoCountText", "runs", 0, "text"}); videoCount != nil {
		item.VideoCount, _ = strconv.Atoi(videoCount.(string))
	}
	if description := getValue(chann, path{"descriptionSnippet", "runs"}); description != nil {
		item.Description = descriptionToStr(description.([]any))
	}
	if subscribers := getValue(chann, path{"subscriberCountText", "simpleText"}); subscribers != nil {
		item.Subscribers = subscribers.(string)
	}
	return item
}

func parsePlaylistComponent(playlist map[string]any) *playlistItem {
	item := &playlistItem{}
	if id := getValue(playlist, path{"playlistId"}); id != nil {
		item.ID = id.(string)
		item.URL = "https://www.youtube.com/playlist?list=" + item.ID
	}
	if title := getValue(playlist, path{"title", "simpleText"}); title != nil {
		item.Title = title.(string)
	}
	if videoCount := getValue(playlist, path{"videoCount"}); videoCount != nil {
		item.VideoCount, _ = strconv.Atoi(videoCount.(string))
	}
	item.Channel = channel{}
	if channelTitle := getValue(playlist, path{"shortBylineText", "runs", 0, "text"}); channelTitle != nil {
		item.Channel.Title = channelTitle.(string)
	}
	if channelID := getValue(playlist, path{"shortBylineText", "runs", 0, "navigationEndpoint", "browseEndpoint", "browseId"}); channelID != nil {
		item.Channel.ID = channelID.(string)
		item.Channel.URL = "https://www.youtube.com/channel/" + item.Channel.ID
	}
	if thumbnails := getValue(playlist, path{"thumbnailRenderer", "playlistVideoThumbnailRenderer", "thumbnail", "thumbnails"}); thumbnails != nil {
		for _, thumb := range thumbnails.([]any) {
			item.Thumbnails = append(item.Thumbnails, thumbnail{
				URL:    "http:" + thumb.(map[string]any)["url"].(string),
				Height: int(thumb.(map[string]any)["height"].(float64)), Width: int(thumb.(map[string]any)["width"].(float64)),
			})
		}
	}
	return item
}

func parseShelfComponent(shelf map[string]any) *shelfItem {
	item := &shelfItem{}
	if title := getValue(shelf, path{"title", "simpleText"}); title != nil {
		item.Title = title.(string)
	}
	items := getValue(shelf, path{"content", "verticalListRenderer", "items"})
	for _, shelfItm := range items.([]any) {
		if videoElement, ok := shelfItm.(map[string]any)["videoRenderer"]; ok {
			videoComponent := parseVideoComponent(videoElement.(map[string]any))
			item.Items = append(item.Items, videoComponent)
		}
	}
	return item
}

func parseSuggestions(response map[string]any) []string {
	suggestions := response["refinements"]
	if suggestions == nil {
		return []string{}
	}
	result := make([]string, len(suggestions.([]any)))
	for index, item := range suggestions.([]any) {
		result[index] = item.(string)
	}
	return result
}
