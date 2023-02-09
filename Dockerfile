
FROM scratch

COPY ./build/go-fetch /usr/local/bin/go-fetch

ENTRYPOINT ["go-fetch"]
