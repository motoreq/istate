export GOPATH=/media/sf_Blockchain/medisot/gitrepo/gopath
rm vendor/ -rf
mkdir -p vendor/github.com/prasanths96/iState/
# cp ../../vendor . -rf
cp ../../*.go vendor/github.com/prasanths96/iState
rm $GOPATH/src/github.com/prasanths96/iState -rf
mkdir -p $GOPATH/src/github.com/prasanths96/iState
cp vendor/github.com/prasanths96/iState/* $GOPATH/src/github.com/prasanths96/iState/. -rf
cp $GOPATH/src/github.com/emirpasic/ vendor/github.com/. -rf
# rm $GOPATH/src/github.com/prasanths96/hyperledger/easycompositestate -rf
# mkdir -p $GOPATH/src/github.com/prasanths96/hyperledger/easycompositestate
# cp vendor/github.com/prasanths96/hyperledger/easycompositestate/* $GOPATH/src/github.com/prasanths96/hyperledger/easycompositestate/. -rf
# rm $GOPATH/src/github.com/prasanths96/hyperledger/querylib -rf
# mkdir -p $GOPATH/src/github.com/prasanths96/hyperledger/querylib
# cp vendor/github.com/prasanths96/hyperledger/querylib/* $GOPATH/src/github.com/prasanths96/hyperledger/querylib/. -rf
go build