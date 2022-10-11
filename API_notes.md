



## G API notes



Describe a DNSSEC enabled zone

```shell
gcloud dns  managed-zones describe ramashenai-zone --log-http >> API_notes.md
=======================
==== request start ====
uri: https://dns.googleapis.com/dns/v1/projects/xyz-dev-274318/managedZones/ramashenai-zone?alt=json
method: GET
== headers start ==
b'accept': b'application/json'
b'accept-encoding': b'gzip, deflate'
b'authorization': --- Token Redacted ---
b'content-length': b'0'
b'user-agent': b'google-cloud-sdk gcloud/387.0.0 command/gcloud.dns.managed-zones.describe invocation-id/712d7a57c9cf4403b352e8f0626f8364 environment/None environment-version/None interactive/True from-script/False python/3.9.10 term/xterm-256color (Macintosh; Intel Mac OS X 20.6.0)'
== headers end ==
== body start ==

== body end ==
==== request end ====
---- response start ----
status: 200
-- headers start --
Alt-Svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000,h3-Q050=":443"; ma=2592000,h3-Q046=":443"; ma=2592000,h3-Q043=":443"; ma=2592000,quic=":443"; ma=2592000; v="46,43"
Cache-Control: private
Content-Encoding: gzip
Content-Type: application/json; charset=UTF-8
Date: Mon, 10 Oct 2022 17:37:07 GMT
ETag: 7eeda99fbcb1995f00000183c2e67817
Server: ESF
Transfer-Encoding: chunked
Vary: Origin, X-Origin, Referer
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-XSS-Protection: 0
-- headers end --
-- body start --
{
  "name": "ramashenai-zone",
  "dnsName": "ramashenai.com.",
  "description": "",
  "id": "9146152921789995359",
  "nameServers": [
    "ns-cloud-e1.googledomains.com.",
    "ns-cloud-e2.googledomains.com.",
    "ns-cloud-e3.googledomains.com.",
    "ns-cloud-e4.googledomains.com."
  ],
  "creationTime": "2022-10-10T17:17:07.479Z",
  "dnssecConfig": {
    "state": "on",
    "defaultKeySpecs": [
      {
        "keyType": "keySigning",
        "algorithm": "rsasha256",
        "keyLength": 2048,
        "kind": "dns#dnsKeySpec"
      },
      {
        "keyType": "zoneSigning",
        "algorithm": "rsasha256",
        "keyLength": 1024,
        "kind": "dns#dnsKeySpec"
      }
    ],
    "nonExistence": "nsec3",
    "kind": "dns#managedZoneDnsSecConfig"
  },
  "visibility": "public",
  "cloudLoggingConfig": {
    "enableLogging": false,
    "kind": "dns#managedZoneCloudLoggingConfig"
  },
  "kind": "dns#managedZone"
}

-- body end --
total round trip time (request+response): 0.222 secs
---- response end ----
----------------------
```

Describe KSK

```shell
gcloud dns dns-keys list --zone ramashenai-zone --format='value(ds_record())' --filter='type=keySigning' --log-http
=======================
==== request start ====
uri: https://dns.googleapis.com/dns/v1/projects/xyz-dev-274318/managedZones/ramashenai-zone/dnsKeys?alt=json&maxResults=100
method: GET
== headers start ==
b'accept': b'application/json'
b'accept-encoding': b'gzip, deflate'
b'authorization': --- Token Redacted ---
b'content-length': b'0'
b'user-agent': b'google-cloud-sdk gcloud/387.0.0 command/gcloud.dns.dns-keys.list invocation-id/7de34372556746fe98adf3b7d078032b environment/None environment-version/None interactive/True from-script/False python/3.9.10 term/xterm-256color (Macintosh; Intel Mac OS X 20.6.0)'
== headers end ==
== body start ==

== body end ==
==== request end ====
---- response start ----
status: 200
-- headers start --
Alt-Svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000,h3-Q050=":443"; ma=2592000,h3-Q046=":443"; ma=2592000,h3-Q043=":443"; ma=2592000,quic=":443"; ma=2592000; v="46,43"
Cache-Control: private
Content-Encoding: gzip
Content-Type: application/json; charset=UTF-8
Date: Mon, 10 Oct 2022 17:45:38 GMT
Server: ESF
Transfer-Encoding: chunked
Vary: Origin, X-Origin, Referer
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-XSS-Protection: 0
-- headers end --
-- body start --
{
  "dnsKeys": [
    {
      "id": "0",
      "algorithm": "rsasha256",
      "keyLength": 2048,
      "publicKey": "AwEAAYAjWFVOozRU6IyuQrKAwDG/BH6rl4fCf6oMuYZsTPJ0QnRTWOv6cKymLKxrVIxm24wQjcfOMrwT8Be1etmvhjGFxmIUd02azHluNapNS7K8KI7F/+Gw/fn/RPbPxMfQE1RKViL78ti+HONCqPTSNHLxN/yiYH1RyWUTW6AGZL7CqANUwWUKgTdwEArYAR4zjeplpHpATeu33aVzIR39TlDZuogGv0Xr+aZjVUvUe3W0w+fo+8spxbDUpgpyHCmKN7mJCr5PVIXN96qKOi1ybAjiylHTyV3Z0shf/RSaQIRQiFkWsDPsl/MTBJTwaVYb3rDOkbRKX6BP/2i3XAq00D0=",
      "creationTime": "2022-10-10T17:17:07.497Z",
      "isActive": true,
      "type": "keySigning",
      "keyTag": 34946,
      "digests": [
        {
          "type": "sha256",
          "digest": "2286DD954EF44E7F47B2AD2AC95027AD67B8F247D121D9312AD96FF8A26A263E"
        }
      ],
      "kind": "dns#dnsKey"
    },
    {
      "id": "1",
      "algorithm": "rsasha256",
      "keyLength": 1024,
      "publicKey": "AwEAAbp+JAEsFoUkyVIxsLTBrdbZfimm4D/QIO7jg3M85UiBd0fqiWMpvWE3iDiyvTZjuT9wZZyGZoaDkB8WDJO0e8XcwyOqyRB1WoznUU8VAwyHSfkIHIzg1A2WOEs+BnxRO2ltqyxlgxqrWj0tRo7bfkmqXrbhVmZsFzcAJXlQM+qx",
      "creationTime": "2022-10-10T17:17:07.609Z",
      "isActive": true,
      "type": "zoneSigning",
      "keyTag": 10785,
      "digests": [],
      "kind": "dns#dnsKey"
    }
  ],
  "kind": "dns#dnsKeysListResponse"
}

-- body end --
total round trip time (request+response): 3.689 secs
---- response end ----
----------------------
34946 8 2 2286DD954EF44E7F47B2AD2AC95027AD67B8F247D121D9312AD96FF8A26A263E
```