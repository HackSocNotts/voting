module hacksocnotts.co.uk/voting/registration

go 1.15

require (
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	go.mongodb.org/mongo-driver v1.5.4
	hacksocnotts.co.uk/voting/common v0.0.0-00010101000000-000000000000
)

replace hacksocnotts.co.uk/voting/common => ../common
