module hacksocnotts.co.uk/voting/admin

go 1.15

replace hacksocnotts.co.uk/voting/common => ../common

require (
	github.com/gorilla/mux v1.8.0
	go.mongodb.org/mongo-driver v1.5.4
	hacksocnotts.co.uk/voting/common v0.0.0-00010101000000-000000000000
)
