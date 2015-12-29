package tutorials
import (
	"time"
)

type Course struct {
	ImageUrl   string `valid:"Required;"`                           //Image Url to be displayed when the course is shown on the Home page or search results
	Title      string `valid:"Required;"`
	Rating     int32
	Technology string  `valid:"Required;"`                          //Hibernate, Wicket, JSP , Java etc
	Keywords   []string `valid:"Required;MinSize(1);MaxSize(5)"`    //java/.Net/SQL/python etc
	User       string   `valid:"Required;"`                         //uploaded by
	Lectures   []Lecture `valid:"Required;MinSize(1);MaxSize(100)"` //Limiting to max 100 videos per course. Who would want to watch more than 100 videos?? hah?
	Date       time.Time
}

//Datatore does not support map[string]string struct member. Hence creating a struct to refer to each platform specific image url
type PlatformImageUrls struct {
	Desktop string `valid:"Required;"`
	Mobile  string
	Tablet  string
}

type Lecture struct {
	Title      string    `valid:"Required;"`
	Provider   string    `valid:"Required;"`         //youtube/vimeo/infoq etc
	ImageUrls  PlatformImageUrls `valid:"Required;"` //image links of various resolutions (mobile, tablet, desktop etc)
	Link       string `valid:"Required;"`
	Embeddable bool                                  //can the video be embedded in the same page?
	EmbedLink  string
}

func (f *Lecture) String() string {
	return Jsonify(f)
}

func (f *Course) String() string {
	return Jsonify(f)
}
