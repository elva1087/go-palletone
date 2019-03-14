#!/bin/bash

./deploy.sh

cp ./createaccount.sh ./createaccounts.sh ./node1

cd ./node1

./createaccounts.sh $1

cd ..

./start.sh

./node1/gptn --exec 'mediator.startProduce()' attach ./node1/palletone/gptn.ipc