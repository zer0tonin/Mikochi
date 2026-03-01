FROM ubuntu:latest

COPY --chmod=755 /backend/mikochi /usr/bin/mikochi
RUN mkdir /usr/share/mikochi
COPY  ./frontend/dist /usr/share/mikochi/static

EXPOSE 8080

CMD ["mikochi"]
