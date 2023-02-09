
FROM scratch

COPY ./build/artifact-manager /usr/local/bin/artifact-manager

ENTRYPOINT ["artifact-manager"]
