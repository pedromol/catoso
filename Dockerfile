FROM pedromol/catoso:base
    
COPY cmd /go/src/app/cmd
COPY pkg /go/src/app/pkg
COPY go.mod /go/src/app/
COPY go.sum /go/src/app/
COPY data/haarcascade_frontalcatface_extended.xml /
COPY data/haarcascade_frontalcatface.xml /

RUN . "$HOME/.asdf/asdf.sh" && \
    go build -o catoso cmd/catoso/main.go && \
    cp catoso /usr/local/bin

WORKDIR /
RUN rm -Rf /go

ENTRYPOINT ["catoso"]
