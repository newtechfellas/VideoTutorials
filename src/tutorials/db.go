package tutorials
import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
)

func CreateOrUpdate(ctx context.Context, obj interface{}, kind string) error {
	_, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, kind, nil), obj)
	if err != nil {
		log.Println("Failed to save object to datastore for kind:", kind)
		return err
	}
	return nil
}


func GetAllCourses(ctx context.Context, course *[]Course) error {
	q := datastore.NewQuery("Course").Order("-Date")
	if _, err := q.GetAll(ctx, course); err != nil {
		return err
	}
	return nil
}
