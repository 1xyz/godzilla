# Understanding Dnssec

Dnssec is used to validate the content served by the DNS server to prevent DNS spoofing and cache poisoning.

In Dnssec, the parent zone vets the response from the child zone, this is called the chain of trust.

Dnssec introduces the following types of records:
1) RRSIG; a digital signature of a Resource Record (RR) set; also known as a RRSET.
2) DNSKEY; a resource record that holds the public key information
3) DS; hash of the DNSKEY record
4) NSEC or NSEC3; records for negative (i.e lack of record) response.
5) CDNSKEY and CDS key; used for key rotation. This is rarely used; so I skipped understanding this.

## RRSET

RRSET is a wire-format concatenation of one or more resource records (RR) of the same type and label for a zone.

**Example**:
For a zone foo.bar.com
there are two A resource records for foo.bar.com say to 10.12.12.1 and 10.12.12.2; then the two resource records form a RRSET.

## Zone Signing Key (ZSK)

* Every zone gets a supposedly unique public/private key pair called the Zone Signing Key Pair (ZSK)
* For every zone; and every RRSET; a RRSIG record is computed. The RRSIG is a digital signature of the RRSET.
* In other words; the RRSET is hashed and then encrypted with the private part of the ZSK.
  * [IANA](https://www.iana.org/assignments/dns-sec-alg-numbers/dns-sec-alg-numbers.xhtml#dns-sec-alg-numbers-1) has a list of algorithm numbers and Mnemonics
  * Example: ECDSAP256SHA256 (number: 13) represents an Elliptic Curve Digital Signature Algorithm (ECDSA) that uses a 256-bit public key and a Secure Hashing Algorithm (SHA) to hash the message contents before signing. The length of the hash is 256 bits.
* The public part of the ZSK is placed in a DNSKEY resource record.
  * More on the RRSIG part of the DNSKEY record.

### Problem
In essence, we can validate the Resource Records provided by the DNS server (using the digital signature) for a zone as long as we can trust the ZSK.
1) How can we trust the DNSKEY : (the public ZSK)?
2) How do you know that DNSKEY is not spoofed.

## Key Signing Key (KSK)

### Problem:
I need to know:
1) If the ZSK is valid and not compromised
2) Validate that ZSK is indeed from the zone - this is called the chain of trust (more on this in the next section).

Key signing key pair (KSK) is a supposedly unique public/private key pair for a zone; which is used to generate a RRSIG of the DNSKEY resource records.

Note:
1) A KSK is also published as a DNSKEY record
2) The RRSIG of the DNSKEY is a digital signature computed by hashing the concatenated wire-formats of public KSK and the ZSK (the DNSKEY RRSET) and encrypted using the private KSK.

Note:
1) Now we use the KSK to digitally sign the ZSK and the KSK; however we have the same problem;
2) How do we validate the KSK is not compromised and actually belongs to this zone.

## Chain of trust and the DS record

The Delegation signer (DS) record is the hash of the KSK using a hashing algorithm, such as SHA-256 for example. This DS record is manually entered into the parent's zone by an authorized entity. So when a child zone is referred to by the parent, the referral contains the DS record; which can be independently verified by taking the digital signature approach.



## References