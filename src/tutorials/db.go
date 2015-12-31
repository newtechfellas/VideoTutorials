package tutorials
import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
)

func CreateOrUpdate(ctx context.Context, obj interface{}, kind string, numericID int64) error {
	_, err := datastore.Put(ctx, datastore.NewKey(ctx, kind, "", numericID, nil), obj)
	if err != nil {
		log.Println("Failed to save object to datastore for kind:", kind, err)
		return err
	}
	return nil
}


func GetEntity(ctx context.Context, id int64, kind string, entity interface{}) (err error) {
	if err = datastore.Get(ctx, datastore.NewKey(ctx, kind, "", id, nil), entity); err != nil {
		log.Println("Did not find the entity with id ", id, " for kind = ", kind)
	}
	return
}

func GetAllCourses(ctx context.Context, course *[]Course) error {
	q := datastore.NewQuery("Course").Order("-Views") //sort to have the most viewed courses at the top
	if _, err := q.GetAll(ctx, course); err != nil {
		return err
	}
	return nil
}
