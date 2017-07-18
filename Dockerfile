FROM scratch

COPY deploy /deploy
EXPOSE 3008

CMD ["/deploy"]
