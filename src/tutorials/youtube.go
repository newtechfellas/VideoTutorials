package tutorials

import (
	"code.google.com/p/google-api-go-client/googleapi/transport"
	"code.google.com/p/google-api-go-client/youtube/v3"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"log"
	"net/http"
	"strings"
	"time"
)

func AddPlaylistVideosAsCourse(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = isTrustedReq(w, r); err != nil {
		return //response already handled by isTrustedReq
	}
	var c Course
	if c, err = decodeCourse(w, r); err != nil {
		return
	}
	ctx := appengine.NewContext(r)
	var service *youtube.Service
	if service, err = getYoutubeService(w, r, ctx); err != nil {
		return
	}
	var playListItems []*youtube.PlaylistItem
	if playListItems, err = getPlaylistItems(service, w,r); err != nil {
		return
	}
	var videoDurations map[string]string
	log.Println("playListItems= ", Jsonify(playListItems))
	videoDurations, err = getDurations(service, playListItems)
	for index, item := range playListItems {
		log.Println("iterating over ",index," for ", Jsonify(item))
		c.Lectures = append(c.Lectures, Lecture{ImageUrl: item.Snippet.Thumbnails.Default.Url,
			Title:    item.Snippet.Title,
			Provider: "youtube",
			Order:    index,
			Duration: videoDurations[item.Snippet.ResourceId.VideoId],
			Link:     "https://youtube.com/watch?v=" + item.Snippet.ResourceId.VideoId})
	}
	if err = CreateOrUpdate(ctx, &c, "Course", c.Id); err != nil {
		log.Println("Error in saving course created from playlist items ", c)
		ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	AddToCache(c)
	JsonResponse(w, c, nil, http.StatusOK)
}

func decodeCourse(w http.ResponseWriter, r *http.Request) (c Course, err error) {
	decoder := json.NewDecoder(r.Body)
	now := time.Now()
	c = Course{Id: now.Unix(), Date: now}
	if err = decoder.Decode(&c); err != nil {
		log.Println("Error in decoding request body. Error is ", err)
		ErrorResponse(w, errors.New("Invalid json details in request body"), http.StatusBadRequest)
		return
	}
	log.Println("decoded course value:", c)
	return
}

func getPlaylistItems(service *youtube.Service, w http.ResponseWriter, r *http.Request) (playListItems []*youtube.PlaylistItem, err error) {
	var response *youtube.PlaylistItemListResponse
	playlistId := r.URL.Query().Get("id")
	var nextPageToken string
	for {
		playlistItemCall := service.PlaylistItems.List("snippet").PlaylistId(playlistId)
		if len(nextPageToken) > 0 {
			playlistItemCall.PageToken(nextPageToken)
		}
		if response, err = playlistItemCall.Do(); err != nil {
			log.Println("Error in fetching playlist items ", err)
			ErrorResponse(w, err, http.StatusInternalServerError)
			break
		}
		if len(response.Items) > 0 {
			//filter deleted videos
			for _, item := range response.Items {
				if item.Snippet.Description != "This video is unavailable." {
					playListItems = append(playListItems, item)
				}
			}
		}
		nextPageToken = response.NextPageToken
		if len(nextPageToken) <= 0 {
			break
		}
	}
	return
}

func getYoutubeService(w http.ResponseWriter, r *http.Request, ctx context.Context) (service *youtube.Service, err error) {
	transport := &transport.APIKey{
		Key:       apiKey,
		Transport: &urlfetch.Transport{Context: ctx}}
	client := &http.Client{Transport: transport}

	service, err = youtube.New(client)
	if err != nil {
		log.Println("ERROR in creating youtube New client ", err)
		ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	return
}

func getDurations(service *youtube.Service, playlistItems []*youtube.PlaylistItem) (data map[string]string, err error) {
	var videoIds []string
	data = make(map[string]string)
	for _, listItem := range playlistItems {
		videoIds = append(videoIds, listItem.Snippet.ResourceId.VideoId)
	}
	response, err := service.Videos.List("contentDetails").Id(strings.Join(videoIds, ",")).Do()
	if err != nil {
		log.Println("Error in fetching duration of the video ", err)
		return
	}

	for _, item := range response.Items {
		data[item.Id] = item.ContentDetails.Duration
	}
	return
}
