rm blockchain
rm *.db

go build -o blockchain *.go
./blockchain