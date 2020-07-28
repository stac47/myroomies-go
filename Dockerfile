FROM golang AS builder
WORKDIR /usr/src/myroomies
COPY . .
RUN make all

FROM debian
ENV MYROOMIES_ROOT_LOGIN="root"
ENV MYROOMIES_ROOT_PASSWORD="password"
ENV MYROOMIES_DATA_STORAGE="memory"
EXPOSE 8080
COPY --from=builder /usr/src/myroomies/myroomies-server /
CMD /myroomies-server --storage ${MYROOMIES_DATA_STORAGE}
