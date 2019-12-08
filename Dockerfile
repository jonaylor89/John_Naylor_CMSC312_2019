
######## Start build phase of execution #######
FROM golang:latest as builder

LABEL maintainer="John Naylor <jonaylor89@gmail.com>"
 
WORKDIR /app

COPY . .
 
RUN make

######## Start a new stage from scratch #######
FROM alpine:latest 

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/jose .
COPY --from=builder /app/config.yml .

COPY ProgramFiles/ ./ProgramFiles

CMD ["./jose"]
