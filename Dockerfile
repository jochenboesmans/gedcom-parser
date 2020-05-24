FROM golang:1.14

RUN echo "[url \"git@github.com:\"]\n\tinsteadOf = https://github.com/" >> /root/.gitconfig
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

ENV WDIR /go/src/github.com/jochenboesmans/gedcom-parser
WORKDIR ${WDIR}
COPY . ${WDIR}
RUN mkdir -p ./artifacts

RUN go get github.com/jochenboesmans/gedcom-parser
RUN go build

CMD ["./gedcom-parser"]
