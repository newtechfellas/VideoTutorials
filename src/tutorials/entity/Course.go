package entity
import (
	"tutorials"
	"time"
)

type Course struct {
	Title      string `valid:"Required;"`
	Rating     int32
	Technology string  `valid:"Required;"`             //Hibernate, Wicket, JSP , Java etc
	Keywords   []string `valid:"Required;MinSize(1);"` //java/.Net/SQL/python etc
	User       string   `valid:"Required;"`            //uploaded by
	Lectures   []Lecture `valid:"Required;MinSize(1);"`
	Date       time.Time
}

func (f *Course) String() string {
	return tutorials.Jsonify(f)
}
