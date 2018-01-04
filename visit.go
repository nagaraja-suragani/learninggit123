package doctorservice

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

/*
//SuccessResponse store response
type SuccessResponse struct {
	//	visit   VisitEntity `json:"entity"`
	ID      int64  `json:"Id"`
	Message string `json:"message"`
}
*/

func getallvisitshandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}

	// get all
	visitList := []*VisitEntity{}
	q := datastore.NewQuery(KINDVISIT)
	keys, err := q.GetAll(ctx, &visitList)
	if err != nil {
		panic(err)
	}

	for i, v := range visitList {
		v.ID = keys[i].IntID()
	}
	if err := json.NewEncoder(w).Encode(visitList); err != nil {
		panic(err)
	}
}

func visitsbypatienthandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}
	//var patient PatientEntity
	x := ps.ByName("id")
	y, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	visitList := []*VisitEntity{}
	q := datastore.NewQuery(KINDVISIT)
	patientkey := datastore.NewKey(ctx, KINDPATIENT, "", int64(y), nil)
	log.Printf("%#v Getting values url - key ", patientkey)
	q = q.Ancestor(patientkey)
	log.Printf("%#v Getting values from  ancestor ", q)

	//get Patient ID From URL  and get the key pkey
	//q = q.Ancestor(pkey)

	keys, err := q.GetAll(ctx, &visitList)
	if err != nil {
		panic(err)
	}

	for i, v := range visitList {
		v.ID = keys[i].IntID()
	}
	if err := json.NewEncoder(w).Encode(visitList); err != nil {
		panic(err)
	}
}

func getavisithandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}
	x := ps.ByName("id")
	y, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	// get one
	visit := &VisitEntity{}

	key := datastore.NewKey(ctx, KINDVISIT, "", int64(y), nil)
	log.Printf("%#v Getting values key ", key)
	if err := datastore.Get(ctx, key, visit); err != nil {
		panic(err)
	}
	log.Printf("%#v Getting values visit ", visit)

	visit.ID = key.IntID()
	if err := json.NewEncoder(w).Encode(visit); err != nil {
		panic(err)
	}
	return
}

func putavisthandler(w http.ResponseWriter, r *http.Request) {

	visit := &VisitEntity{}
	log.Printf("%#v Getting values visit ", visit)
	if err := json.NewDecoder(r.Body).Decode(visit); err != nil {
		panic(err)
	}
	keys := make([]*datastore.Key, 1)
	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}
	// Send a visit.PatientID in Newkey and get the key pkey
	//in below line instead of sending nil send pkey

	if err != nil {
		panic(err)
	}
	x := visit.PatientID
	log.Printf("%#v Getting values xxxxxxx ", x)
	//	y, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	patientkey := datastore.NewKey(ctx, KINDPATIENT, "", int64(x), nil)

	keys[0] = datastore.NewIncompleteKey(ctx, KINDVISIT, patientkey)

	k, err := datastore.Put(ctx, keys[0], visit)
	if err != nil {
		panic(err)
	}

	visit.ID = k.IntID()
	json.NewEncoder(w).Encode(visit)
}

func updateavisithandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	visit := &VisitEntity{}
	if err := json.NewDecoder(r.Body).Decode(visit); err != nil {
		panic(err)
	}

	keys := make([]*datastore.Key, 1)
	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}
	x := ps.ByName("id")

	log.Printf("%#v Getting values url - x ", x)
	y, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v Getting values url - y ", y)

	keys[0] = datastore.NewKey(ctx, KINDVISIT, "", int64(y), nil)

	k, err := datastore.Put(ctx, keys[0], visit)
	if err != nil {
		panic(err)
	}

	visit.ID = k.IntID()
	json.NewEncoder(w).Encode(visit)
}

func deleteallvisithandler(w http.ResponseWriter, r *http.Request) {

	var visitslist []VisitEntity
	var ctx context.Context
	ctx = appengine.NewContext(r)
	keys := make([]*datastore.Key, 1)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		return
	}

	q := datastore.NewQuery(KINDVISIT)
	keys, _ = q.GetAll(ctx, &visitslist)
	option := &datastore.TransactionOptions{XG: true}
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		return datastore.DeleteMulti(c, keys)
	}, option)
	json.NewEncoder(w).Encode("Deleted All Visit Records!")
}

func deleteavisithandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var ctx context.Context
	ctx = appengine.NewContext(r)
	keys := make([]*datastore.Key, 1)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		return
	}
	x := ps.ByName("id")
	log.Printf("%#v Getting values url - x ", x)
	y, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v Getting values url - y ", y)
	if y > 0 {

		keys[0] = datastore.NewKey(ctx, KINDVISIT, "", int64(y), nil)
		option := &datastore.TransactionOptions{XG: true}
		err = datastore.RunInTransaction(ctx, func(c context.Context) error {
			return datastore.DeleteMulti(c, keys)
		}, option)
		json.NewEncoder(w).Encode("Deleted A Visit Records!")
	}
}
