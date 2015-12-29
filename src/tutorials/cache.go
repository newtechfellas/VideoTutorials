package tutorials
import (
	"time"
	"os"
	"log"
	"golang.org/x/net/context"
	"strconv"
)

// Although packages exist for cache in Go, below cache implementation is purely local and simple usage
// for app engine usage

type CourseCache struct {
	courses    []Course
	cachedTime time.Time
}

var cacheExpiryInMin int
func loadCacheAllowedMinutes() (int) {
	if ( cacheExpiryInMin == 0 ) {
		cacheExpiryInMin = 24 * 60 //24 hours default
		if v := os.Getenv("CACHE_EXPIRY_MINUTES"); v != "" {
			log.Println("CACHE_EXPIRY_MINUTES = ", v)
			if i, err := strconv.Atoi(v); err == nil {
				cacheExpiryInMin = i
			}
		}
	}
	return cacheExpiryInMin
}

var CACHE CourseCache = CourseCache{cachedTime : time.Now() }//global cache for courses

//return the courses from cache if available and within allowed expiry.
func GetCoursesFromCache(ctx context.Context) []Course {
	hoursSinceCached := int(time.Now().Sub(CACHE.cachedTime).Hours())
	if ( len(CACHE.courses) > 0 &&  hoursSinceCached <= loadCacheAllowedMinutes()) {
		log.Println("Returning cached courses count = ", len(CACHE.courses))
		return CACHE.courses
	}
	log.Println("Course cache does not have data. Loading ...")
	LoadCoursesToCache(ctx)
	return CACHE.courses
}

func LoadCoursesToCache(ctx context.Context) {
	var courses []Course
	if err := GetAllCourses(ctx, &courses); err != nil {
		log.Println("ERROR in loading courses from datastore. Can not cache the values")
	}
	CACHE.cachedTime = time.Now()
	CACHE.courses = courses
}
