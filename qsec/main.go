package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"time"
)

const (
	DefaulTimeout   = 2 * time.Second
	GoogleDNSServer = "8.8.8.8:53"
	CloudDNSServer  = "1.1.1.1:53"
)

type DnsQuery struct {
	QueryType    uint16
	QueryClass   uint16
	QName        string
	NameServer   string
	EnableDnssec bool
	UseRecursive bool
}

func (dqr *DnsQuery) String() string {
	return fmt.Sprintf("ns: %s q: %s type: %d class: %d dnssec: %v recursive: %v",
		dqr.NameServer,
		dqr.QName,
		dqr.QueryClass,
		dqr.QueryType,
		dqr.EnableDnssec,
		dqr.UseRecursive)
}

func NewDnsQueryReq(name, server string, queryType uint16, enableDnssec, useRecursive bool) *DnsQuery {
	return &DnsQuery{
		QueryType:    queryType,
		QueryClass:   dns.ClassINET,
		QName:        name,
		NameServer:   server,
		EnableDnssec: enableDnssec,
		UseRecursive: useRecursive,
	}
}

func (dqr *DnsQuery) Execute() (*dns.Msg, error) {
	fmt.Printf("  Nameserver : %s\n", dqr.NameServer)
	fmt.Printf("      QClass : %d\n", dqr.QueryClass)
	fmt.Printf("       QType : %d\n", dqr.QueryType)
	fmt.Printf("       QName : %s\n", dqr.QName)
	fmt.Printf("EnableDNSSEC : %v\n", dqr.EnableDnssec)
	fmt.Printf("UseRecursive : %v\n", dqr.UseRecursive)

	c := new(dns.Client)
	c.Net = "udp"
	c.ReadTimeout = DefaulTimeout
	c.DialTimeout = DefaulTimeout
	c.WriteTimeout = DefaulTimeout

	m := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Authoritative:     false,
			AuthenticatedData: false,
			CheckingDisabled:  false,
			RecursionDesired:  dqr.UseRecursive,
			Opcode:            dns.OpcodeQuery,
			Rcode:             dns.RcodeSuccess,
			Id:                dns.Id(),
		},
		Question: []dns.Question{
			dns.Question{
				Name:   dns.Fqdn(dqr.QName),
				Qtype:  dqr.QueryType,
				Qclass: dqr.QueryClass,
			},
		},
	}

	if dqr.EnableDnssec {
		o := &dns.OPT{
			Hdr: dns.RR_Header{
				Name:   ".",
				Rrtype: dns.TypeOPT,
			},
		}
		o.SetDo()
		o.SetUDPSize(dns.DefaultMsgSize)
		m.Extra = append(m.Extra, o)
	}

	fmt.Printf(" Query: %s\n", m.String())
	fmt.Printf(" Size: %d\n", m.Len())
	r, rtt, err := c.Exchange(m, dqr.NameServer)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, fmt.Errorf("execute with %v failed: %v", dqr.String(), err)
	}

	fmt.Printf("Rtt: %d\n", rtt)
	return r, nil
}

func printMessage(m *dns.Msg) {
	if m.Answer == nil {
		fmt.Printf("null answer\n")
		return
	}
	if len(m.Answer) == 0 {
		fmt.Printf("empty answer\n")
	}
	// fmt.Printf("RR length of answer %d\n", len(m.Answer))
	for _, ans := range m.Answer {
		switch v := ans.(type) {
		case *dns.A:
			fmt.Printf("IP address: %v\n", v.A)
		case *dns.RRSIG:
			fmt.Printf("RRSig: %v\n", v.String())
		case *dns.NS:
			fmt.Printf("NS: %v\n", v.Ns)
		case *dns.DNSKEY:
			fmt.Printf("DNSKEY: %v\n", v.String())
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}

		hdr := ans.Header()
		if hdr == nil {
			continue
		}

		//rrType := dns.Type(hdr.Rrtype).String()
		fmt.Printf("hdr %s  \n", hdr.String())
		fmt.Println()
	}

	//fmt.Printf("m %v\n", m.String())
}

func verifiedQuery() {
	var dnsNameserver = CloudDNSServer
	var queryRecordType = dns.TypeA
	var query = "cloudflare.com."
	recordQuery := NewDnsQueryReq(query,
		dnsNameserver, queryRecordType, true, true)
	recordResponse, err := recordQuery.Execute()
	if err != nil {
		log.Fatalf("%v", err)
	}
	printMessage(recordResponse)

	rrSet := make([]dns.RR, 0)
	var rrSig *dns.RRSIG = nil
	for _, entry := range recordResponse.Answer {
		hdr := entry.Header()
		if hdr.Rrtype == queryRecordType {
			rrSet = append(rrSet, entry)
			continue
		}

		if hdr.Rrtype == dns.TypeRRSIG {
			rrSigRecord, ok := entry.(*dns.RRSIG)
			if !ok {
				log.Fatalf("cast-err: cannot cast entry of type: %T to *dns.RRSIG", entry)
			}

			rrSig = rrSigRecord
			continue
		}
	}

	if len(rrSet) == 0 || rrSig == nil {
		log.Fatalf("could not extract rrset and/or rrsig")
	}

	dnskeyQuery := NewDnsQueryReq(query,
		dnsNameserver, dns.TypeDNSKEY, true, true)
	dnskeyResponse, err := dnskeyQuery.Execute()
	if err != nil {
		log.Fatalf("%v", err)
	}
	printMessage(dnskeyResponse)

	var (
		zsk       *dns.DNSKEY = nil
		ksk       *dns.DNSKEY = nil
		zskKeyTag uint16
		kskKeyTag uint16
		zskRRSig  *dns.RRSIG = nil
		kskRRSig  *dns.RRSIG = nil
	)
	for _, entry := range dnskeyResponse.Answer {
		hdr := entry.Header()
		if hdr.Rrtype != dns.TypeDNSKEY {
			continue
		}
		dnskeyRR, ok := entry.(*dns.DNSKEY)
		if !ok {
			log.Printf("typecast failed to *dns.DNSKEY\n %T", entry)
			continue
		}
		log.Printf("flags = %v, key-tag = %v\n",
			dnskeyRR.Flags,
			dnskeyRR.KeyTag())
		if dnskeyRR.Flags == 256 {
			zsk = dnskeyRR
			continue
		}
		if dnskeyRR.Flags == 257 {
			ksk = dnskeyRR
			continue
		}
	}

	if zsk == nil || ksk == nil {
		log.Fatalf("could not extract zsk and/or ksk")
	}

	zskKeyTag, kskKeyTag = zsk.KeyTag(), ksk.KeyTag()
	for _, entry := range dnskeyResponse.Answer {
		hdr := entry.Header()
		if hdr.Rrtype != dns.TypeRRSIG {
			continue
		}
		dnskeyRRSig, ok := entry.(*dns.RRSIG)
		if !ok {
			log.Printf("typecast failed to *dns.DNSKEY\n %T", entry)
			continue
		}

		if dnskeyRRSig.KeyTag == zskKeyTag {
			zskRRSig = dnskeyRRSig
			continue
		}

		if dnskeyRRSig.KeyTag == kskKeyTag {
			kskRRSig = dnskeyRRSig
			continue
		}
	}

	if zskRRSig == nil {
		log.Printf("could not extract zsk  RRSIG\n")
	}

	if kskRRSig == nil {
		log.Printf("could not extract ksk RRSIG\n")
	}

	// Verify the  record
	if err := rrSig.Verify(zsk, rrSet); err != nil {
		log.Fatalf("(0) FAILED!!!! verification of rrSet: %s with rrSig: %s and zsk: %s",
			rrSet, rrSet, zsk)
	} else {
		log.Printf("Verified the %v record\n", dns.Type(queryRecordType).String())
	}

	if kskRRSig != nil {
		// Verify the DNSKEY (ZSK)
		if err := kskRRSig.Verify(ksk, []dns.RR{ksk, zsk}); err != nil {
			log.Fatalf("(1) FAILED!!!! verification of zsk: %s rrSig: %s and ksk: %s",
				zsk, zskRRSig, ksk)
		}
		log.Printf("Verified using KSK\n")
	}

	if zskRRSig != nil {
		// Verify the DNSKEY (ZSK)
		if err := zskRRSig.Verify(zsk, []dns.RR{ksk, zsk}); err != nil {
			log.Fatalf("(2) FAILED!!!! verification of zsk: %s rrSig: %s and ksk: %s",
				zsk, zskRRSig, ksk)
		}
		log.Printf("Verified using ZSK\n")
	}

	log.Printf("all verified")
}

func simpleQuery() {
	var query *DnsQuery = nil
	var r *dns.Msg = nil
	var err error = nil

	query = NewDnsQueryReq("ietf.org.",
		GoogleDNSServer, dns.TypeDNSKEY, true, true)
	r, err = query.Execute()
	if err != nil {
		log.Fatalf("%v", err)
	}
	printMessage(r)

	query = NewDnsQueryReq("ietf.org.",
		GoogleDNSServer, dns.TypeA, true, true)
	r, err = query.Execute()
	if err != nil {
		log.Fatalf("%v", err)
	}
	printMessage(r)
}

func main() {
	verifiedQuery()
}
