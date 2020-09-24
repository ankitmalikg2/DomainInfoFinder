package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/ankitmalikg2/DomainInfoFinder/models"
)

const MAX = 300
const NewDomainsFileName = "domain-names.txt"
const SkipLines = 47543

//ReadFileIntoArray it reads the file and sent Array of strings as output
func ReadFileIntoArray(filePath string) ([]string, error) {
	file, err := os.Open(filePath)

	defer file.Close()

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	return txtlines, nil
}

//GetWhoisInfo provides the whois info for the domain
func GetWhoisInfo(domain string) {

}
func main() {
	fmt.Println("Domain Info Finder")
	sem := make(chan int, MAX)

	domainFileName := NewDomainsFileName
	txtlines, err := ReadFileIntoArray(domainFileName)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	count := 0
	successCount := 0

	for _, eachline := range txtlines {

		// if count >= 5 {
		// 	break
		// }

		count++

		if count < SkipLines {
			fmt.Println(count)
			continue
		}

		sem <- 1 //for limiting the function
		go func() {
			fmt.Println(eachline)
			domain := eachline
			out, err := exec.Command("whois", domain).Output()
			if err != nil {
				fmt.Println(domain, "error occured", err)
			}
			outputArr := strings.Split(string(out), "\n")
			registrantInfo := models.DomainInfo{}
			for _, singleLine := range outputArr {
				ExtractValues(singleLine, &registrantInfo)
			}

			// fmt.Printf("%+v ", registrantInfo)

			//saving the data
			if !reflect.DeepEqual(registrantInfo, models.DomainInfo{}) {
				successCount++
				errMap := SaveDispatch(domain, registrantInfo)
				fmt.Println(errMap)
			}
			<-sem //for completing the function
		}()
	}

	fmt.Printf("Total Count: %d and Success Count: %d", count, successCount)

}

// ExtractValues : it extacts the values from stirng and saves to model
func ExtractValues(str string, val *models.DomainInfo) {
	if strings.Contains(str, "Registrant Name:") {
		val.RegistrantName = strings.TrimSpace(strings.ReplaceAll(str, "Registrant Name:", ""))
	} else if strings.Contains(str, "Registrant Organization:") {
		val.RegistrantOrg = strings.TrimSpace(strings.ReplaceAll(str, "Registrant Organization:", ""))
	} else if strings.Contains(str, "Registrant Street:") {
		val.RegistrantStreet = strings.TrimSpace(strings.ReplaceAll(str, "Registrant Street:", ""))
	} else if strings.Contains(str, "Registrant City:") {
		val.RegistrantCity = strings.TrimSpace(strings.ReplaceAll(str, "Registrant City:", ""))
	} else if strings.Contains(str, "Registrant State/Province:") {
		val.RegistrantProvince = strings.TrimSpace(strings.ReplaceAll(str, "Registrant State/Province:", ""))
	} else if strings.Contains(str, "Registrant Postal Code:") {
		val.RegistrantPostalCode = strings.TrimSpace(strings.ReplaceAll(str, "Registrant Postal Code:", ""))
	} else if strings.Contains(str, "Registrant Country:") {
		val.RegistrantCountry = strings.TrimSpace(strings.ReplaceAll(str, "Registrant Country:", ""))
	} else if strings.Contains(str, "Registrant Phone:") {
		val.RegistrantPhone = strings.TrimSpace(strings.ReplaceAll(str, "Registrant Phone:", ""))
	} else if strings.Contains(str, "Registrant Phone Ext:") {
		val.RegistrantPhoneExt = strings.TrimSpace(strings.ReplaceAll(str, "Registrant Phone Ext:", ""))
	} else if strings.Contains(str, "Registrant Fax:") {
		val.RegistrantFax = strings.TrimSpace(strings.ReplaceAll(str, "Registrant Fax:", ""))
	} else if strings.Contains(str, "Registrant Fax Ext:") {
		val.RegistrantFaxExt = strings.TrimSpace(strings.ReplaceAll(str, "Registrant Fax Ext:", ""))
	} else if strings.Contains(str, "Registrant Email:") {
		val.RegistrantEmail = strings.TrimSpace(strings.ReplaceAll(str, "Registrant Email:", ""))
	}
}

// SaveDispatch : it diptaches the model to save data
func SaveDispatch(domain string, model models.DomainInfo) map[string]error {
	outputErrors := make(map[string]error)
	err := SavetoCSV(domain, model)
	outputErrors["CSV"] = err

	return outputErrors
}

func CreateFile(fileVal string) error {
	if !FileExists(fileVal) {
		fmt.Println("file create mode")
		fl, err := os.Create(fileVal)
		if err != nil {
			return err
		}
		ColumnNames := []string{
			"domain",
			"Name",
			"Organisation",
			"Email",
			"Country",
			"City",
			"Street",
			"Phone",
			"PhoneExt",
			"Fax",
			"FaxExt",
			"PostalCode",
			"Province/State",
		}
		writer := csv.NewWriter(fl)
		err = writer.Write(ColumnNames)
		if err != nil {
			return err
		}
		writer.Flush()
		fl.Close()
	}
	return nil
}

//SavetoCSV : It will save the model data to CSV file
func SavetoCSV(domain string, model models.DomainInfo) error {
	val := time.Now().Format("2006-01-02") + ".csv"
	var file *os.File
	var err error

	if !FileExists(val) {
		fmt.Println("file create mode")
		fl, err := os.Create(val)
		if err != nil {
			return err
		}
		ColumnNames := []string{
			"domain",
			"Name",
			"Organisation",
			"Email",
			"Country",
			"City",
			"Street",
			"Phone",
			"PhoneExt",
			"Fax",
			"FaxExt",
			"PostalCode",
			"Province/State",
		}
		writer := csv.NewWriter(fl)
		err = writer.Write(ColumnNames)
		if err != nil {
			return err
		}
		writer.Flush()
		fl.Close()
	}

	file, err = os.OpenFile(val, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	data := []string{
		domain,
		model.RegistrantName,
		model.RegistrantOrg,
		model.RegistrantEmail,
		model.RegistrantCountry,
		model.RegistrantCity,
		model.RegistrantStreet,
		model.RegistrantPhone,
		model.RegistrantPhoneExt,
		model.RegistrantFax,
		model.RegistrantFaxExt,
		model.RegistrantPostalCode,
		model.RegistrantProvince,
	}
	// fmt.Println("writing data", data)
	err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

//FileExists checks if file exits or not
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
