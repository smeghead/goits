FROM golang:1.18

RUN apt-get update && apt-get install -y gettext
