package entity
import (
	"tutorials"
)

type Lecture struct {
	Title      string	`valid:"Required;"`
	Provider   string    `valid:"Required;"`        //youtube/vimeo/infoq etc
	ImageUrls  map[string]string `valid:"Required;MinSize(1);"`//image links of various resolutions (mobile, tablet, desktop etc)
	Link       string `valid:"Required;"`
	Embeddable bool              //can the video be embedded in the same page
	EmbedLink  string
}

func (f *Lecture) String() string {
	return tutorials.Jsonify(f)
}