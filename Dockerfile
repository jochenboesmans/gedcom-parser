FROM golang:1.14

RUN echo "[url \"git@github.com:\"]\n\tinsteadOf = https://github.com/" >> /root/.gitconfig
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config
ADD .  /go/src/github.com/jochenboesmans/gedcom-parser
CMD cd /go/src/github.com/jochenboesmans/gedcom-parser && go get github.com/jochenboesmans/gedcom-parser && go build -o /parse
