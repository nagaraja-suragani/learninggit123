package doctorservice

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

func init() {

	initKeys()
	//Testing
	http.Handle("/resource/", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))

	router := httprouter.New()
	router.POST("/register", RegisterHandler)
	router.POST("/log", LoginHandler)
	router.POST("/forgot", Forgotpassword)
	router.NotFound = http.FileServer(http.Dir("./Angular/dist"))

	//START OF PATIENT HANDLERS
	nputapatient := negroni.New()
	nputapatient.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nputapatient.UseHandlerFunc(putapatienthandler)
	router.Handler("POST", "/api/postapatient/", nputapatient)

	// add middleware for a specific route and get params from rout
	nUpdateapatient := negroni.New()
	nUpdateapatient.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nUpdateapatient.UseHandlerFunc(callwithParams(router, updateapatienthandler))
	router.Handler("POST", "/api/postapatient/:id", nUpdateapatient)

	// add middleware for a specific route to protect
	ngetallpatients := negroni.New()
	ngetallpatients.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	ngetallpatients.UseHandlerFunc(getallpatientshandler)
	router.Handler("GET", "/api/getallpatients/", ngetallpatients)

	// add middleware for a specific route and get params from route
	nGetapatient := negroni.New()
	nGetapatient.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nGetapatient.UseHandlerFunc(callwithParams(router, getapatienthandler))
	router.Handler("GET", "/api/getallpatients/:id", nGetapatient)

	// add middleware for a specific route to protect
	nDeleteallpatients := negroni.New()
	nDeleteallpatients.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nDeleteallpatients.UseHandlerFunc(deleteallpatienthandler)
	router.Handler("DELETE", "/api/deleteapatient/", nDeleteallpatients)

	// add middleware for a specific route and get params from rout
	nDeleteapatient := negroni.New()
	nDeleteapatient.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nDeleteapatient.UseHandlerFunc(callwithParams(router, deleteapatienthandler))
	router.Handler("DELETE", "/api/deleteapatient/:id", nDeleteapatient)

	// add middleware for a specific route and get params from rout
	nSearchapatient := negroni.New()
	nSearchapatient.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nSearchapatient.UseHandlerFunc(callwithParams(router, searchpatientsphoneHandler))
	router.Handler("GET", "/api/search/:phone", nSearchapatient)

	//END OF PATIENT HANDLERS

	// START ALL THE ENDPOINTS  OF VISIT ARE PROTECTED

	// add middleware for a specific route to protect
	nPostavisit := negroni.New()
	nPostavisit.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	//nPostavisit.UseHandlerFunc(callwithParams(router, putavisthandler))
	nPostavisit.UseHandlerFunc(putavisthandler)
	router.Handler("POST", "/api/postavisit/", nPostavisit)

	// add middleware for a specific route and get params from rout
	nUpdateavisit := negroni.New()
	nUpdateavisit.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nUpdateavisit.UseHandlerFunc(callwithParams(router, updateavisithandler))
	router.Handler("POST", "/api/postavisit/:id", nUpdateavisit)

	// add middleware for a specific route to protect
	ngetallvisits := negroni.New()
	ngetallvisits.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	ngetallvisits.UseHandlerFunc(getallvisitshandler)
	router.Handler("GET", "/api/getallvisits/", ngetallvisits)

	// add middleware for a specific route and get params from route
	nHello := negroni.New()
	nHello.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nHello.UseHandlerFunc(callwithParams(router, getavisithandler))
	router.Handler("GET", "/api/getallvisits/:id", nHello)

	// add middleware for a specific route to protect
	nDeleteallvisit := negroni.New()
	nDeleteallvisit.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nDeleteallvisit.UseHandlerFunc(deleteallvisithandler)
	router.Handler("DELETE", "/api/deleteavisit/", nDeleteallvisit)

	// add middleware for a specific route and get params from rout
	nDeleteavisit := negroni.New()
	nDeleteavisit.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nDeleteavisit.UseHandlerFunc(callwithParams(router, deleteavisithandler))
	router.Handler("DELETE", "/api/deleteavisit/:id", nDeleteavisit)

	nVisitsbypatient := negroni.New()
	nVisitsbypatient.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nVisitsbypatient.UseHandlerFunc(callwithParams(router, visitsbypatienthandler))
	router.Handler("GET", "/api/patientvisits/:id", nVisitsbypatient)

	//END OF VISIT HANDLERS

	//START OF PRESCRIPTION HANDLERS

	nputaprescription := negroni.New()
	nputaprescription.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nputaprescription.UseHandlerFunc(putaprescriptionhandler)
	router.Handler("POST", "/api/postaprescription/", nputaprescription)

	// add middleware for a specific route and get params from rout
	nUpdateaprescription := negroni.New()
	nUpdateaprescription.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nUpdateaprescription.UseHandlerFunc(callwithParams(router, updateprescriptionahandler))
	router.Handler("POST", "/api/postaprescription/:id", nUpdateaprescription)

	// add middleware for a specific route to protect
	ngetallprescriptions := negroni.New()
	ngetallprescriptions.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	ngetallprescriptions.UseHandlerFunc(getallprescriptionhandler)
	router.Handler("GET", "/api/getallprescriptions/", ngetallprescriptions)

	// add middleware for a specific route and get params from route
	nGetaprescription := negroni.New()
	nGetaprescription.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nGetaprescription.UseHandlerFunc(callwithParams(router, getaprescriptionhandler))
	router.Handler("GET", "/api/getallprescriptions/:id", nGetaprescription)

	// add middleware for a specific route to protect
	nDeleteallprescriptions := negroni.New()
	nDeleteallprescriptions.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nDeleteallprescriptions.UseHandlerFunc(deleteallprescriptionhandler)
	router.Handler("DELETE", "/api/deleteaprescription/", nDeleteallprescriptions)

	// add middleware for a specific route and get params from rout
	nDeleteaprescription := negroni.New()
	nDeleteaprescription.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	nDeleteaprescription.UseHandlerFunc(callwithParams(router, deleteaprescriptionhandler))
	router.Handler("DELETE", "/api/deleteaprescription/:id", nDeleteaprescription)

	//http.HandleFunc("/api/patient/", getallvisitsbasedonpatienthandler)
	http.Handle("/", router)

}

// callwithParams function is helping us to call controller from middleware having access to URL params
func callwithParams(router *httprouter.Router, handler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := getURLParams(router, r)
		handler(w, r, params)
	}
}

// getUrlParams function is extracting URL parameters
func getURLParams(router *httprouter.Router, req *http.Request) httprouter.Params {

	_, params, _ := router.Lookup(req.Method, req.URL.Path)

	return params
}
