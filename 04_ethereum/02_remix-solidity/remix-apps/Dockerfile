# Node image
FROM node:latest

USER root

# Create code directory
RUN mkdir /source

# Set working directory
WORKDIR /source

RUN npm install remix-ide@0.8.0 && remix-ide

EXPOSE 8080

