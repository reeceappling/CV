package main

type Certification struct {
	Name       string
	Link       *string
	CertDate   monthYr
	ExpiryDate *monthYr
}

func NewCert(name string, month int, year int, end *monthYr, link *string) *Certification {
	return &Certification{
		Name: name,
		Link: link,
		CertDate: monthYr{
			Month: month,
			Year:  year,
		},
		ExpiryDate: end,
	}
}

var (
	certAwsSaa    = NewCert("Aws Solutions Architect - Associate", 04, 2023, nil, nil)      // TODO: ENSURE DATES OK
	certOWASP     = NewCert("Owasp top 10 course", 9, 2022, nil, nil)                       // TODO: ENSURE DATES OK
	certComtrain2 = NewCert("Comtrain Certified Tower Climber/Rescuer", 9, 2019, nil, nil)  // TODO: ENSURE DATES OK
	certComtrain  = NewCert("Comtrain Certified Tower Climber/Rescuer", 06, 2017, nil, nil) // TODO: ENSURE DATES OK
	certLifeguard = NewCert("Red Cross certified LifeGuard", 01, 2013, nil, nil)            // TODO: ENSURE DATES OK
	certCPR       = NewCert("CPR", 01, 2013, nil, nil)                                      // TODO: ENSURE DATES OK
)
