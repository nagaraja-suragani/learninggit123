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

func getallpatientshandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}

	// get all
	patientList := []*PatientEntity{}
	q := datastore.NewQuery(KINDPATIENT)
	keys, err := q.GetAll(ctx, &patientList)
	if err != nil {
		panic(err)
	}

	for i, v := range patientList {
		v.ID = keys[i].IntID()
	}
	if err := json.NewEncoder(w).Encode(patientList); err != nil {
		panic(err)
	}
}

func putapatienthandler(res http.ResponseWriter, req *http.Request) {
	patient := &PatientEntity{}
	if err := json.NewDecoder(req.Body).Decode(patient); err != nil {
		panic(err)
	}

	keys := make([]*datastore.Key, 1)
	ctx := appengine.NewContext(req)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}

	keys[0] = datastore.NewIncompleteKey(ctx, KINDPATIENT, nil)
	//keys[0] = datastore.NewKey(ctx, "Patient",patient.Phone, 0, nil)

	k, err := datastore.Put(ctx, keys[0], patient)
	if err != nil {
		panic(err)
	}
	patient.ID = k.IntID()
	json.NewEncoder(res).Encode(patient)

}
func getapatienthandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	patient := &PatientEntity{}

	key := datastore.NewKey(ctx, KINDPATIENT, "", int64(y), nil)
	log.Printf("%#v Getting values key ", key)
	if err := datastore.Get(ctx, key, patient); err != nil {
		panic(err)
	}
	log.Printf("%#v Getting values visit ", patient)

	patient.ID = key.IntID()
	if err := json.NewEncoder(w).Encode(patient); err != nil {
		panic(err)
	}
	return
}

func updateapatienthandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	patient := &PatientEntity{}
	if err := json.NewDecoder(r.Body).Decode(patient); err != nil {
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

	keys[0] = datastore.NewKey(ctx, KINDPATIENT, "", int64(y), nil)

	k, err := datastore.Put(ctx, keys[0], patient)
	if err != nil {
		panic(err)
	}

	patient.ID = k.IntID()
	json.NewEncoder(w).Encode(patient)
}

func deleteallpatienthandler(w http.ResponseWriter, r *http.Request) {

	var patientslist []PatientEntity
	var ctx context.Context
	ctx = appengine.NewContext(r)
	keys := make([]*datastore.Key, 1)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		return
	}

	q := datastore.NewQuery(KINDPATIENT)
	keys, _ = q.GetAll(ctx, &patientslist)
	option := &datastore.TransactionOptions{XG: true}
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		return datastore.DeleteMulti(c, keys)
	}, option)
	json.NewEncoder(w).Encode("Deleted All patient Records!")
}

func deleteapatienthandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

		keys[0] = datastore.NewKey(ctx, KINDPATIENT, "", int64(y), nil)
		option := &datastore.TransactionOptions{XG: true}
		err = datastore.RunInTransaction(ctx, func(c context.Context) error {
			return datastore.DeleteMulti(c, keys)
		}, option)
		json.NewEncoder(w).Encode("Deleted A Visit Records!")
	}
}

func searchpatientsphoneHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	ctx := appengine.NewContext(req)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}
	x := ps.ByName("phone")
	log.Printf("%#v Getting values q ", x)
	patients := []*PatientEntity{}
	q := datastore.NewQuery(KINDPATIENT).Filter("Phone =", x)
	log.Printf("%#v Getting values q ", q)

	keys, err := q.GetAll(ctx, &patients)
	log.Printf("%#v Getting values m ", keys)

	if err != nil {
		// Handle error
		return
	}
	for i, v := range patients {
		v.ID = keys[i].IntID()

	}
	json.NewEncoder(res).Encode(patients)
}
