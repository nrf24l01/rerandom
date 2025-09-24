FROM alpine:3.18
RUN apk add --no-cache dnsmasq
COPY dnsmasq.conf /etc/dnsmasq.conf
EXPOSE 5300 5300/udp
CMD ["dnsmasq", "-k"]
