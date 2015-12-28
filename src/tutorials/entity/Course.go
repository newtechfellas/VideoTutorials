package entity
import (
	"encoding/json"
)

type Course struct {
	Title      string `valid:"Required;"`
	Rating     int32
	Technology string  `valid:"Required;"` //Hibernate, Wicket, JSP , Java etc
	Keywords   []string `valid:"Required;MinSize(1);"`//java/.Net/SQL/python etc
	User       string   `valid:"Required;"` //uploaded by
	Lectures   []*Lecture `valid:"Required;MinSize(1);"`
}

func (f *Course) String() string {
	b, err := json.Marshal(f)
	if err != nil {
		return ""
	}
	return string(b)
}
