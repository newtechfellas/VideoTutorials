package tutorials

import (
	"code.google.com/p/google-api-go-client/googleapi/transport"
	"code.google.com/p/google-api-go-client/youtube/v3"
	"log"
	"net/http"
)

var service *youtube.Service

func init() {
	var err error
	log.Println("Apikey = ", apiKey)
	client := &http.Client{Transport: &transport.APIKey{Key: apiKey}}
	service, err = youtube.New(client)
	if err != nil {
		log.Println("ERROR in creating youtube New client ", err)
	}
	var items *youtube.PlaylistItemListResponse
	if items, err = service.PlaylistItems.List("snippet").PlaylistId("PLFE2CE09D83EE3E28").Do(); err != nil {
		log.Println("Error in fetching playlist items ", err)
	}
	log.Println(Jsonify(items))
}
