FROM scratch
COPY mapcleaner /bin/mapcleaner
ENTRYPOINT ["/bin/mapcleaner"]