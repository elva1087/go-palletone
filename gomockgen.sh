go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen
mockgen -source=./dag/interface.go -destination=./dag/dag_mock.go -package=dag -self_package="github.com/palletone/go-palletone/dag"
mockgen -source=./mediator_connection.go  -destination=./mediator_connection_mock.go -package=ptn
