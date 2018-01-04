package doctorservice

const NAMESPACENAME = "-kashyak-"
const KINDVISIT = "Visit"
const KINDPATIENT = "Patient"
const KINDPRESCRIPTION = "Prescription"

//User is struct to hold the User details
type User struct {
	Username string `datastore:"-"`
	Password string `json:"password"`
}

//PatientEntity is struct to hold the Patient details
type PatientEntity struct {
	ID        int64  `datastore:"-"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

//VisitEntity is struct to hold the vist details
type VisitEntity struct {
	ID            int64  `datastore:"-"`
	PatientID     int64  `json:"patientid,omitempty"`
	UserName      string `json:"username,omitempty"`
	Height        string `json:"height,omitempty"`
	Weight        string `json:"weight,omitempty"`
	Temperature   string `json:"temperature,omitempty"`
	BloodPressure string `json:"bloodpressure,omitempty"`
	DoctorNote    string `json:"doctornote,omitempty"`
	PatientNote   string `json:"patientnote,omitempty"`
	NurseNote     string `json:"nursenote,omitempty"`
}

//PrescriptionEntity is struct to hold the prescription details
type PrescriptionEntity struct {
	ID             int64  `datastore:"-"`
	PatientName    string `json:"patientname,omitempty"`
	Prescription   string `json:"prescription,omitempty"`
	Drug           string `json:"drug,omitempty"`
	Administration string `json:"administration,omitempty"`
	Duration       string `json:"duration,omitempty"`
}
