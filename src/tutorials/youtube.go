package tutorials

import (
	"code.google.com/p/google-api-go-client/googleapi/transport"
	"code.google.com/p/google-api-go-client/youtube/v3"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"log"
	"net/http"
)

func AddPlaylistVideosAsCourse(w http.ResponseWriter, r *http.Request) {
	//TODO: Add security check. This API must be secured using a key
	if err := isTrustedReq(w, r); err != nil {
		return //response already handled by isTrustedReq
	}

	var service *youtube.Service
	ctx := appengine.NewContext(r)

	transport := &transport.APIKey{
		Key:       apiKey,
		Transport: &urlfetch.Transport{Context: ctx}}
	client := &http.Client{Transport: transport}

	var err error
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
