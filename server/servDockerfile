FROM ubuntu:24.04

WORKDIR /opt/serv

RUN apt-get -yqq update
RUN apt-get -yqq install dbus

COPY ./start.sh /usr/local/bin
RUN chmod +x /usr/local/bin/start.sh

RUN dbus-uuidgen > /var/lib/dbus/machine-id
RUN mkdir -p /var/run/dbus

COPY ./session.conf /usr/share/dbus-1/
COPY ./session.conf /etc/dbus-1/session.d/
COPY ./session.conf /etc/dbus-1/

EXPOSE 12345

CMD [ "/bin/bash", "/usr/local/bin/start.sh" ]
