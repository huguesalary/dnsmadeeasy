package dnsmadeeasy

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

var (
	multiDomains []uint32
	domains      []*Domain
	domain       *Domain
	record       *Record
	client       *Client
	domainPrefix string
)

func init() {
	client = &Client{}
	flag.StringVar(&domainPrefix, "prefix", "testcases", "Which domain-prefix should the test use. (E.g. 'mydomain')")
	flag.StringVar(&client.APIKey, "api-key", "", "dnsmadeeasy api key")
	flag.StringVar(&client.APISecret, "api-secret", "", "dnsmadeeasy api secret")
	client.Version = "2.0"
	client.Url = fmt.Sprintf(sandboxApiUrlFmt, "2.0")
	flag.Parse()
	if client.APIKey == "" || client.APISecret == "" {
		flag.PrintDefaults()
		panic("\n=================\nYou need to specify -api-key and -api-secret flags\n=================\n\n")
	}
}

func TestCreateDomains(t *testing.T) {

	i := 1
	var err error
	multiDomains, err = client.CreateDomains([]string{domainPrefix + "-multi1.com", domainPrefix + "-multi2.com"})

	for ; err != nil; multiDomains, err = client.CreateDomains([]string{domainPrefix + "-multi1.com", domainPrefix + "-multi2.com"}) {
		if err, ok := err.(*APIError); ok && err.Code == 400 {
			fmt.Println(err)
			time.Sleep(10 * time.Second)
			if i++; i < 7 {
				continue
			}
		}
		t.Error("Failed creating", "testcases1.com and testcases2.com", err)
		return
	}

}

func TestCreateDomain(t *testing.T) {
	i := 1
	var err error
	for domain, err = client.CreateDomain(domainPrefix + ".com"); err != nil; domain, err = client.CreateDomain(domainPrefix + ".com") {
		if err, ok := err.(*APIError); ok && err.Code == 400 {
			fmt.Println(err)
			time.Sleep(10 * time.Second)
			if i++; i < 7 {
				continue
			}
		}
		t.Error("Failed creating", domainPrefix+".com", err)
		return
	}

}

func TestGetDomains(t *testing.T) {

	doms, err := client.GetDomains()

	if err != nil {
		t.Error("GetDomains returned an error", err)
		return
	}

	for _, dom := range doms.Domains {
		domains = append(domains, &dom)
		if dom.Name != domainPrefix+".com" {
			continue
		} else {
			domain = &dom
			return
		}
	}
	t.Error("Couldn't find the domain testcases.com")
}

func TestUpdDomains(t *testing.T) {
	fields := make(map[string]interface{})
	fields["gtdEnabled"] = false

	i := 0
	for err := client.UpdDomains(multiDomains, fields); err != nil; err = client.UpdDomains(multiDomains, fields) {
		if err, ok := err.(*APIError); ok && err.Code == 400 {
			fmt.Println(err)
			time.Sleep(10 * time.Second)
			if i++; i < 7 {
				continue
			}

		}
		t.Error("Failed creating", domainPrefix+".com", err)
		return
	}
}

func TestGetDomainById(t *testing.T) {
	dom, err := client.GetDomainById(domain.Id)

	if err != nil || dom.Id != domain.Id {
		t.Fail()
	}
}

func TestGetDomainByName(t *testing.T) {
	dom, err := client.GetDomainByName(domain.Name)

	if err != nil || dom.Name != domain.Name {
		t.Fail()
	}
}

func TestAddRecord(t *testing.T) {
	record = &Record{
		Name:        "a-record",
		Value:       "127.0.0.1",
		Type:        "A",
		Ttl:         60,
		GtdLocation: "DEFAULT",
	}

	if err := client.AddRecord(domain.Id, record); err != nil {
		t.Error(err)
	}
}

func TestGetDomainRecords(t *testing.T) {
	if recs, err := client.GetDomainRecords(domain.Id); err != nil {
		t.Fail()
	} else {
		for _, rec := range recs {
			if rec.Name == "a-record" {
				return
			}
		}
	}

	t.Fail()
}

func TestUpdRecord(t *testing.T) {

	if record == nil {
		t.Error("No record to update.")
		return
	}

	record = &Record{
		Name:        "a-updated-record",
		Value:       "127.0.0.1",
		Type:        "A",
		Ttl:         60,
		GtdLocation: "DEFAULT",
		Id:          record.Id,
	}

	if err := client.UpdRecord(domain.Id, record); err != nil {
		t.Error(err)
	}
}

func TestDelRecord(t *testing.T) {

	if record == nil {
		t.Error("No record to delete.")
		return
	}

	if err := client.DelRecord(domain.Id, record.Id); err != nil {
		t.Error(err)
	}
}

func TestDelDomain(t *testing.T) {

	i := 0
	for err := client.DeleteDomain(domain.Id); err != nil; err = client.DeleteDomain(domain.Id) {
		if err, ok := err.(*APIError); ok && err.Code == 400 {
			fmt.Println(err)
			time.Sleep(10 * time.Second)
			if i++; i < 7 {
				continue
			}
		}
		t.Error("Failed deleting", domain.Name, err)
		return
	}
}

func TestDelDomains(t *testing.T) {

	i := 0
	for err := client.DeleteDomains(multiDomains); err != nil; err = client.DeleteDomains(multiDomains) {
		if err, ok := err.(*APIError); ok && err.Code == 400 {
			fmt.Println(err)
			time.Sleep(10 * time.Second)
			if i++; i < 7 {
				continue
			}
		}
		t.Error("Failed deleting", multiDomains, err)
		return
	}
}
