docker:
	docker run -it --net=host -v /e/Worker/GoCode/goskynet:/data/goskynet -e GOPROXY=goproxy.io -w /data/goskynet --name gobuild golang:1.17 /bin/bash