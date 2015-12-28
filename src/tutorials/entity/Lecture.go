package entity
import (
	"tutorials"
)

//Datatore does not support map[string]string struct member. Hence creating a struct to refer to each platform specific image url
type PlatformImageUrls struct {
	Desktop string `valid:"Required;"`
	Mobile string
	Tablet string
}

type Lecture struct {
	Title      string	`valid:"Required;"`
	Provider   string    `valid:"Required;"`        //youtube/vimeo/infoq etc
	ImageUrls PlatformImageUrls `valid:"Required;"`//image links of various resolutions (mobile, tablet, desktop etc)
	Link       string `valid:"Required;"`
	Embeddable bool              //can the video be embedded in the same page
	EmbedLink  string
}

func (f *Lecture) String() string {
	return tutorials.Jsonify(f)
}