# Node image
FROM node:12.4.0

# Create code directory
RUN mkdir /source

# Set working directory
WORKDIR /source

# Install Truffle
#RUN npm install solc && npm install truffle@4.1.14 --unsafe-perm && npm config set bin-links false && npm install -g express-basic-auth

RUN npm install -g truffle@4.1.15 && npm config set bin-links false

WORKDIR MetaCoin

