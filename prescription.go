package doctorservice

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"golang.org/x/net/context"
	//"strings"
	"github.com/julienschmidt/httprouter"
	//"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	//"github.com/julienschmidt/httprouter"
)

//KINDNAME is the table name to store the values
//const KINDPRESCRIPTION = "Prescription"

//NAMESPACENAME is the Namespace
//const NAMESPACENAME = "-kashyak-"

/*
type prescriptiondata struct{

	ID              int64   `datastore:"-"`
	PatientName     string  `json:"patientname,omitempty"`
	Prescription    string  `json:"prescription,omitempty"`
	Drug            string  `json:"drug,omitempty"`
	Administration  string  `json:"administration,omitempty"`
	Duration        string  `json:"duration,omitempty"`
}*/

/*func init() {
	http.HandleFunc("/api/postaprescrition/", putaprescription)
	http.HandleFunc("/api/getaprescrition/", getaprescription)
	http.HandleFunc("/api/deleteaprescription/",deleteaprescription)

}*/


func putaprescriptionhandler(w http.ResponseWriter, r *http.Request) {

	prescription := &PrescriptionEntity{}
	if err := json.NewDecoder(r.Body).Decode(prescription); err != nil {
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

	keys[0] = datastore.NewIncompleteKey(ctx, KINDPRESCRIPTION, nil)

	k, err := datastore.Put(ctx, keys[0], prescription)
	if err != nil {
		panic(err)
	}

	prescription.ID = k.IntID()
	json.NewEncoder(w).Encode(prescription)
}


func updateprescriptionahandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	prescription := &PrescriptionEntity{}
	if err := json.NewDecoder(r.Body).Decode(prescription); err != nil {
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

	keys[0] = datastore.NewKey(ctx, KINDPRESCRIPTION, "", int64(y), nil)

	k, err := datastore.Put(ctx, keys[0], prescription)
	if err != nil {
		panic(err)
	}

	prescription.ID = k.IntID()
	json.NewEncoder(w).Encode(prescription)
}


func getallprescriptionhandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}

	// get all
	
	prescriptionList := []*PrescriptionEntity{}
	q := datastore.NewQuery(KINDPRESCRIPTION)
	keys, err := q.GetAll(ctx, &prescriptionList)
	if err != nil {
		panic(err)
	}

	for i, v := range prescriptionList {
		v.ID = keys[i].IntID()
	}
	if err := json.NewEncoder(w).Encode(prescriptionList); err != nil {
		panic(err)
	}
}


func getaprescriptionhandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	prescriptionList := &PrescriptionEntity{}

	key := datastore.NewKey(ctx, KINDPRESCRIPTION, "", int64(y), nil)
	log.Printf("%#v Getting values key ", key)
	if err := datastore.Get(ctx, key, prescriptionList); err != nil {
		panic(err)
	}
	log.Printf("%#v Getting values visit ", prescriptionList)

	prescriptionList.ID = key.IntID()
	if err := json.NewEncoder(w).Encode(prescriptionList); err != nil {
		panic(err)
	}
	return
}


func deleteallprescriptionhandler(w http.ResponseWriter, r *http.Request) {

	var prescriptionList []PrescriptionEntity
	var ctx context.Context
	ctx = appengine.NewContext(r)
	keys := make([]*datastore.Key, 1)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		return
	}

	q := datastore.NewQuery(KINDPRESCRIPTION)
	keys, _ = q.GetAll(ctx, &prescriptionList)
	option := &datastore.TransactionOptions{XG: true}
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		return datastore.DeleteMulti(c, keys)
	}, option)
	json.NewEncoder(w).Encode("Deleted All Visit Records!")
}

func deleteaprescriptionhandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

		keys[0] = datastore.NewKey(ctx, KINDPRESCRIPTION, "", int64(y), nil)
		option := &datastore.TransactionOptions{XG: true}
		err = datastore.RunInTransaction(ctx, func(c context.Context) error {
			return datastore.DeleteMulti(c, keys)
		}, option)
		json.NewEncoder(w).Encode("Deleted A Visit Records!")
	}
}
