FROM scratch
MAINTAINER Jana Bollineni (jana.bollineni@gmail.com)
LABEL version="1.2"
LABEL org.label-schema.name="HTTP Echo" \
	  org.label-schema.description="Webservice written in GO to echo HTTP header values" \
	  org.label-schema.schema-version="1.0"

ENV PORT 5000
EXPOSE 5000

COPY http-echo /
COPY template.html /

ENTRYPOINT ["/http-echo"]