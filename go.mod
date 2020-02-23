module github.com/710leo/urlooker

go 1.12

require (
	github.com/astaxie/beego v1.12.0
	github.com/codegangsta/negroni v1.0.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.9
	github.com/gorilla/context v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/securecookie v1.1.1
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mmitton/ldap v0.0.0-20170928125358-2890b18948f2
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/quipo/statsd v0.0.0-20180118161217-3d6a5565f314
	github.com/sirupsen/logrus v1.4.2
	github.com/toolkits/container v0.0.0-20151219225805-ba7d73adeaca
	github.com/toolkits/core v0.0.0-20141116054942-0ebf14900fe2 // indirect
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	github.com/toolkits/net v0.0.0-20160910085801-3f39ab6fe3ce
	github.com/toolkits/smtp v0.0.0-20190110072832-af41f29c3d89
	github.com/toolkits/str v0.0.0-20160913030958-f82e0f0498cb
	github.com/toolkits/sys v0.0.0-20170615103026-1f33b217ffaf
	github.com/toolkits/web v0.0.0-20160312232617-dc0f03327e1d
	github.com/unrolled/render v1.0.1
	golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392 // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58 // indirect
	golang.org/x/sys v0.0.0-20190924154521-2837fb4f24fe // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-00010101000000-000000000000 // indirect
	gopkg.in/yaml.v2 v2.2.7
	stathat.com/c/consistent v1.0.0
)

replace gopkg.in/asn1-ber.v1 => github.com/go-asn1-ber/asn1-ber v1.3.1
