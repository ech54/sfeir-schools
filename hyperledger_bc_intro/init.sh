#!/bin/bash
export VERSION=1.2.1
export CA_VERSION=$VERSION
export THIRDPARTY_IMAGE_VERSION=0.4.12
export ARCH=$(echo "$(uname -s|tr '[:upper:]' '[:lower:]'|sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')")
export MARCH=$(uname -m)
BINARY_FILE=hyperledger-fabric-${ARCH}-${VERSION}.tar.gz
CA_BINARY_FILE=hyperledger-fabric-ca-${ARCH}-${CA_VERSION}.tar.gz


binaryIncrementalDownload() {
      local BINARY_FILE=$1
      local URL=$2
      curl -f -s -C - ${URL} -o ${BINARY_FILE} || rc=$?
      if [ "$rc" = 22 ]; then
	  return 22
      fi
      if [ -z "$rc" ] || [ $rc -eq 33 ] || [ $rc -eq 2 ]; then
          echo "==> File downloaded. Verifying the md5sum..."
          localMd5sum=$(md5sum ${BINARY_FILE} | awk '{print $1}')
          remoteMd5sum=$(curl -s ${URL}.md5)
          if [ "$localMd5sum" == "$remoteMd5sum" ]; then
              echo "==> Extracting ${BINARY_FILE}..."
              tar xzf ./${BINARY_FILE} --overwrite
	      echo "==> Done."
              rm -f ${BINARY_FILE} ${BINARY_FILE}.md5
          else
              echo "Download failed: the local md5sum is different from the remote md5sum. Please try again."
              rm -f ${BINARY_FILE} ${BINARY_FILE}.md5
              exit 1
          fi
      else
          echo "Failure downloading binaries (curl RC=$rc). Please try again and the download will resume from where it stopped."
          exit 1
      fi
}

binaryDownload() {
      local BINARY_FILE=$1
      local URL=$2
      echo "===> Downloading: " ${URL}
      if [ -e ${BINARY_FILE} ]; then
          echo "==> Partial binary file found. Resuming download..."
          binaryIncrementalDownload ${BINARY_FILE} ${URL}
      else
          curl ${URL} | tar xz || rc=$?
          if [ ! -z "$rc" ]; then
              echo "==> There was an error downloading the binary file. Switching to incremental download."
              echo "==> Downloading file..."
              binaryIncrementalDownload ${BINARY_FILE} ${URL}
	  else
	      echo "==> Done."
          fi
      fi
}

binariesInstall() {
  echo "===> Downloading version ${FABRIC_TAG} platform specific fabric binaries"
  binaryDownload ${BINARY_FILE} https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/${ARCH}-${VERSION}/${BINARY_FILE}
  if [ $? -eq 22 ]; then
     echo
     echo "------> ${FABRIC_TAG} platform specific fabric binary is not available to download <----"
     echo
   fi

  echo "===> Downloading version ${CA_TAG} platform specific fabric-ca-client binary"
  binaryDownload ${CA_BINARY_FILE} https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric-ca/hyperledger-fabric-ca/${ARCH}-${CA_VERSION}/${CA_BINARY_FILE}
  if [ $? -eq 22 ]; then
     echo
     echo "------> ${CA_TAG} fabric-ca-client binary is not available to download  (Available from 1.1.0-rc1) <----"
     echo
   fi
}

cleanDirectory() {
  rm .env
  rm -rf ./bin; rm -rf ./config
}

dockerFabricPull() {
  local FABRIC_TAG=$1
  for IMAGES in peer orderer ccenv javaenv tools; do
      echo "==> FABRIC IMAGE: $IMAGES"
      echo
      docker pull hyperledger/fabric-$IMAGES:$FABRIC_TAG
      docker tag hyperledger/fabric-$IMAGES:$FABRIC_TAG hyperledger/fabric-$IMAGES
  done
}

dockerThirdPartyImagesPull() {
  local THIRDPARTY_TAG=$1
  for IMAGES in couchdb kafka zookeeper; do
      echo "==> THIRDPARTY DOCKER IMAGE: $IMAGES"
      echo
      docker pull hyperledger/fabric-$IMAGES:$THIRDPARTY_TAG
      docker tag hyperledger/fabric-$IMAGES:$THIRDPARTY_TAG hyperledger/fabric-$IMAGES
  done
}

dockerCaPull() {
      local CA_TAG=$1
      echo "==> FABRIC CA IMAGE"
      echo
      docker pull hyperledger/fabric-ca:$CA_TAG
      docker tag hyperledger/fabric-ca:$CA_TAG hyperledger/fabric-ca
}


dockerInstallStd() {
  which docker >& /dev/null
  NODOCKER=$?
  if [ "${NODOCKER}" == 0 ]; then
	  echo "===> Pulling fabric Images"
	  dockerFabricPull ${FABRIC_TAG}
	  echo "===> Pulling fabric ca Image"
	  dockerCaPull ${CA_TAG}
	  echo
	  echo "===> List out hyperledger docker images"
	  docker images | grep hyperledger*
  else
    echo "========================================================="
    echo "Docker not installed, bypassing download of Fabric images"
    echo "========================================================="
  fi
}

dockerInstallTp() {
  which docker >& /dev/null
  NODOCKER=$?
  if [ "${NODOCKER}" == 0 ]; then
	  echo "===> Pulling thirdparty docker images"
	  dockerThirdPartyImagesPull ${THIRDPARTY_TAG}
	  echo
	  echo "===> List out hyperledger docker images"
	  docker images | grep hyperledger*
  else
    echo "========================================================="
    echo "Docker not installed, bypassing download of Fabric images"
    echo "========================================================="
  fi
}

function replacePrivateKey() {
  ARCH=$(uname -s | grep Darwin)
  if [ "$ARCH" == "Darwin" ]; then
    OPTS="-it"
  else
    OPTS="-i"
  fi
  cp docker-compose-template.yaml docker-sfeir.yaml
  CURRENT_DIR=$PWD
  cd crypto-config/peerOrganizations/shop.sfeir.lu/ca/
  PRIV_KEY=$(ls *_sk)
  cd "$CURRENT_DIR"
  sed $OPTS "s/CA1_PRIVATE_KEY/${PRIV_KEY}/g" docker-sfeir.yaml
  cd crypto-config/peerOrganizations/warehouse.sfeir.lu/ca/
  PRIV_KEY=$(ls *_sk)
  cd "$CURRENT_DIR"
  sed $OPTS "s/CA2_PRIVATE_KEY/${PRIV_KEY}/g" docker-sfeir.yaml
}


if [[ $VERSION =~ ^1\.[0-1]\.* ]]; then
  export FABRIC_TAG=${MARCH}-${VERSION}
  export CA_TAG=${MARCH}-${CA_VERSION}
  export THIRDPARTY_TAG=${MARCH}-${THIRDPARTY_IMAGE_VERSION}
else
  : ${CA_TAG:="$CA_VERSION"}
  : ${FABRIC_TAG:="$VERSION"}
  : ${THIRDPARTY_TAG:="$THIRDPARTY_IMAGE_VERSION"}
fi

if [ "${1}" == "setup" ]; then
	echo "setup environment"
	export PATH=./bin:${PWD}:$PATH
	export COMPOSE_PROJECT_NAME="logistic"
	echo $COMPOSE_PROJECT_NAME
	export FABRIC_CFG_PATH=$PWD
	cp ./tp/*.yaml .
elif [ "${1}" == "cleanup" ]; then
	rm -rf ./crypto-config/* ./config/*
	rm ./*.yaml
	docker volume rm $(docker volume ls)
	docker network rm $(docker network ls  | grep "sfeir")
	docker rmi $(docker images | grep "hyperledger")
elif [ "${1}" == "replace-pki" ]; then
	echo "replace pki ..."
	replacePrivateKey
elif [ "${1}" == "docker-images" ]; then
	dockerInstallStd
	dockerInstallTp
elif [ "${1}" == "binaries" ]; then
	cleanDirectory
	binariesInstall
	echo "Add to path: $PWD/bin"
	export PATH=$PWD/bin:$PATH
	exec /bin/bash
fi

