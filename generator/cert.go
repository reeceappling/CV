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

func NewDate(month int, year int) *monthYr {
	return &monthYr{
		Month: month,
		Year:  year,
	}
}

var (
	certAwsSaa = NewCert("Aws Solutions Architect - Associate", 04, 2023, nil, nil) // TODO: ENSURE DATES OK
	certOWASP  = NewCert("Owasp top 10 course", 9, 2022, nil, nil)                  // TODO: ENSURE DATES OK
	// TODO; RF SAFETY COURSE // TODO: link
	certOSHA      = NewCert("OSHA 10hr", 9, 2019, nil, nil)                                              // TODO: link
	certComtrain2 = NewCert("Comtrain Certified Tower Climber/Rescuer", 9, 2019, NewDate(9, 2022), nil)  // TODO: link
	certComtrain  = NewCert("Comtrain Certified Tower Climber/Rescuer", 06, 2017, NewDate(6, 2020), nil) // TODO: link
	certLifeguard = NewCert("Red Cross certified Lifeguard", 01, 2013, NewDate(6, 2018), nil)            // TODO: link
	certCPR       = NewCert("CPR", 01, 2013, NewDate(1, 2018), nil)                                      // TODO: link
)
