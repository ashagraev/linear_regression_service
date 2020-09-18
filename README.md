# linear_regression_service

## 1. Simple linear regression problem

A simple linear regression problem states as follows: given two-dimensional sample points having one independent variable and one target value, build a linear model to minimize the residual sum of squared errors. See [1]   for further details.

![](https://user-images.githubusercontent.com/6789687/93579011-a5d04a80-f9a6-11ea-975c-1f69443bcf0c.png)

This problem is relatively simple in terms of computation costs. However, numerical errors could potentially lead to unstable and improper results. To deal with that problem, we use Welford's method [2] for calculating means and covariations as well as Kahan's summation algorithm [3].

Links:
1. https://en.wikipedia.org/wiki/Simple_linear_regression
2. https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance#Welford's_online_algorithm
3. https://en.wikipedia.org/wiki/Kahan_summation_algorithm

## 2. Training, storing and applying models

Models are trained on the server and then stored in the Spanner database, so that in server mode one needs to have the corresponding credentials. Models are calculated on the server, too.

To access the compute server, run the program in one of the client modes:
- ```---http-apply``` for calculating model values using HTTP calls;
- ```---http-train``` for training models using HTTP calls;
- ```---grpc-apply``` for calculating model values using gRPC calls;
- ```---grpc-train``` for training models using gRPC calls;

See the following sections for details.

## 3. Install dependencies

```
sudo apt-get update
sudo apt-get -y upgrade

sudo apt-get install -y git
sudo apt-get install -y wget
sudo apt-get install -y zip

wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz
tar -xvf go1.13.3.linux-amd64.tar.gz
sudo mv go /usr/local

export GOROOT=/usr/local/go
export GOPATH=~/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

go get cloud.google.com/go/spanner
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
sudo unzip -o protoc-3.7.1-linux-x86_64.zip -d /usr/local bin/protoc
sudo unzip -o protoc-3.7.1-linux-x86_64.zip -d /usr/local 'include/*'
```

## 4. Build the app

```
cd $GOPATH/src

git clone https://github.com/ashagraev/linear_regression_service
cd linear_regression_service
protoc regression.proto --go_out=.
protoc regression.proto --go-grpc_out=.
go build .
```

## 5. Run the HTTP and gRPC services

```
GOOGLE_APPLICATION_CREDENTIALS=/home/user/token.json ./linear_regression_service --http-server --port 8080
GOOGLE_APPLICATION_CREDENTIALS=/home/user/token.json ./linear_regression_service --grpc-server --address localhost:8081
```

## 6. Train and apply the model via HTTP API

```
./linear_regression_service --http-train --server http://localhost:8080 < ./sample_instances.tsv
{
    "Model": {
        "Coefficient": 0.8727272727272727,
        "Intercept": 0.9400000000000004
    },
    "SumSquaredErrors": 4.187636363636379,
    "Name": "RGtx-35CXkm5Kw==",
    "CreationTime": "2020-09-18T09:49:00.67646Z"
}

./linear_regression_service --http-apply --server http://localhost:8080 --model RGtx-35CXkm5Kw==
1
{
    "Value": 1.812727272727273,
    "Argument": 1,
    "Model": {
        "Name": "RGtx-35CXkm5Kw==",
        "Coefficient": 0.8727272727272727,
        "Intercept": 0.9400000000000004
    },
    "FromCache": false,
    "CalculationTime": "2020-09-18T09:49:42.278102358Z"
}
5
{
    "Value": 5.303636363636364,
    "Argument": 5,
    "Model": {
        "Name": "RGtx-35CXkm5Kw==",
        "Coefficient": 0.8727272727272727,
        "Intercept": 0.9400000000000004
    },
    "FromCache": true,
    "CalculationTime": "2020-09-18T09:49:43.108535903Z"
}
```

## 7. Train and apply the model via gRPC API

```
./linear_regression_service --grpc-train --server localhost:8081 < ./sample_instances.tsv
{
    "Model": {
        "Coefficient": 0.8727272727272727,
        "Intercept": 0.9400000000000004
    },
    "SumSquaredErrors": 4.187636363636379,
    "Name": "JeNnhLsrkEK0TQ==",
    "CreationTime": "2020-09-18 09:50:20.522453 +0000 UTC"
}

./linear_regression_service --grpc-apply --server localhost:8081 --model JeNnhLsrkEK0TQ==
1
{
    "Value": 1.812727272727273,
    "Argument": 1,
    "Model": {
        "Name": "JeNnhLsrkEK0TQ==",
        "Coefficient": 0.8727272727272727,
        "Intercept": 0.9400000000000004
    },
    "CalculationTime": "2020-09-18 09:50:50.922917428 +0000 UTC m=+302.800959414"
}
5
{
    "Value": 5.303636363636364,
    "Argument": 5,
    "Model": {
        "Name": "JeNnhLsrkEK0TQ==",
        "Coefficient": 0.8727272727272727,
        "Intercept": 0.9400000000000004
    },
    "FromCache": true,
    "CalculationTime": "2020-09-18 09:50:51.551225717 +0000 UTC m=+303.429267694"
}

```
