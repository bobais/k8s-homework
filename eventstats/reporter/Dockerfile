FROM python:3.7-alpine

WORKDIR /usr/src/app

COPY requirements.txt ./

RUN \
    apk update \
    && apk add \
        libxml2-dev \
        libxslt-dev \
        gcc \
        libc-dev \
    && rm -rf /var/cache/apk/* \
    && pip install --no-cache-dir -r requirements.txt

RUN addgroup -S noadmin && adduser -S noadmin -G noadmin

# 65534=nobody
COPY --chown=65534 templates/*.jinja2 ./templates/
COPY --chown=65534 reporter.py ./

USER 65534

CMD [ "python", "-u", "./reporter.py" ]
