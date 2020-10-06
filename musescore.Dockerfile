FROM ubuntu:18.04

RUN apt update \
        && apt install -y software-properties-common \
        && add-apt-repository ppa:mscore-ubuntu/mscore3-stable \
        && apt-get update \
        && apt-get -y install musescore3

ENTRYPOINT ["musescore3"]
