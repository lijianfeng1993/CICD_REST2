FROM golang
MAINTAINER lijianfeng "1205960475@qq.com"
ADD main /go/src
RUN cd /go/src && mkdir template
ADD template/ /go/src/template
ADD start_main.sh /go/src
WORKDIR "/go/src"
CMD ["sh","start_main.sh"]
