package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter domains (one per line): ")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error: could not read from input %v\n", err)
	}
}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord string

	// Lookup MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: in Lookup MX %v\n", err)
		return
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	// Lookup SPF records
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: in Lookup TXT %v\n", err)
		return
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	// Lookup DMARC records
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: in Lookup TXT %v\n", err)
		return
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			break
		}
	}

	fmt.Printf("Domain: %s, hasMX: %v, hasSPF: %v, SPF Record: %s, hasDMARC: %v, DMARC Records: %v\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecords)
}

/*
Enter domains (one per line):
mailchimp.com
Domain: mailchimp.com, hasMX: true, hasSPF: true, SPF Record: v=spf1 ip4:205.201.128.0/20 ip4:198.2.128.0/18 ip4:148.105.0.0/16 ip4:129.145.74.12 include:_spf.google.com include:mailsenders.netsuite.com include:_spf2.intuit.com include:_spf.qualtrics.com ip4:199.33.145.1 ip4:199.33.145.32 ip4:35.176.132.251 ip4:52.60.115.116 ~all, hasDMARC: true, DMARC Records: [v=DMARC1; p=reject; rua=mailto:19ezfriw@ag.dmarcian.com,mailto:dmarc_rua@emaildefense.proofpoint.com; ruf=mailto:19ezfriw@fr.dmarcian.com,mailto:dmarc_ruf@emaildefense.proofpoint.com;]
*/
