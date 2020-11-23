FROM golang:1.15

WORKDIR /go/src/github.com/inaccel/protector4j

ADD go.mod .
RUN go mod download

ADD . .

RUN CGO_ENABLED=0 go build -o /go/bin/vlinx-protector4j

FROM buildpack-deps:stretch
LABEL maintainer=InAccel

ARG PROTECTOR4J_VERSION
RUN wget -q https://protector4j.com/resources/pub/protector4j-${PROTECTOR4J_VERSION}.linux64.tar.gz && \
	tar -xzf protector4j-${PROTECTOR4J_VERSION}.linux64.tar.gz -C /usr/share && \
	rm -f protector4j-${PROTECTOR4J_VERSION}.linux64.tar.gz && \
	ln -s /usr/share/protector4j/protector4j /usr/bin/protector4j

COPY --from=0 /go/bin/vlinx-protector4j /bin/vlinx-protector4j

ENTRYPOINT ["vlinx-protector4j"]
