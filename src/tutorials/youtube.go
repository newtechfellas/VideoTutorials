package tutorials

import (
	"code.google.com/p/google-api-go-client/googleapi/transport"
	"code.google.com/p/google-api-go-client/youtube/v3"
	"encoding/json"
	"errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"log"
	"net/http"
	"time"
)

func AddPlaylistVideosAsCourse(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = isTrustedReq(w, r); err != nil {
		return //response already handled by isTrustedReq
	}

	decoder := json.NewDecoder(r.Body)
	now := time.Now()
	c := Course{Id: now.Unix(), Date: now}
	if err = decoder.Decode(&c); err != nil {
		log.Println("Error in decoding request body. Error is ", err)
		ErrorResponse(w, errors.New("Invalid json details in request body"), http.StatusBadRequest)
		return
	}
	log.Println("decoded course value:", c)

	var service *youtube.Service
	ctx := appengine.NewContext(r)

	transport := &transport.APIKey{
		Key:       apiKey,
		Transport: &urlfetch.Transport{Context: ctx}}
	client := &http.Client{Transport: transport}

	service, err = youtube.New(client)
	if err != nil {
		log.Println("ERROR in creating youtube New client ", err)
	}
	var response *youtube.PlaylistItemListResponse
	var playListItems []*youtube.PlaylistItem
	var nextPageToken string
	for {
		playlistItemCall := service.PlaylistItems.List("snippet").PlaylistId("PLFE2CE09D83EE3E28")
		if len(nextPageToken) > 0 {
			playlistItemCall.PageToken(nextPageToken)
		}
		if response, err = playlistItemCall.Do(); err != nil {
			log.Println("Error in fetching playlist items ", err)
		}
		if len(response.Items) > 0 {
			playListItems = append(playListItems, response.Items...)
		}
		nextPageToken = response.NextPageToken
		if len(nextPageToken) <= 0 {
			break
		}
	}
	for index, item := range playListItems {
		c.Lectures = append(c.Lectures, Lecture{ImageUrl: item.Snippet.Thumbnails.Default.Url,
			Title:    item.Snippet.Title,
			Provider: "youtube",
			Order:    index,
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
