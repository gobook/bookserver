GOOS=linux GOARCH=amd64 go build -o bookserver_amd64
echo "uploading ..."
scp ./bookserver_amd64 lunny@xorm.io:/home/lunny/gobook/
scp -r ./themes lunny@xorm.io:/home/lunny/gobook/
scp -r ./public lunny@xorm.io:/home/lunny/gobook/
scp -r ./templates lunny@xorm.io:/home/lunny/gobook/
echo "gobook deployed, please run it"