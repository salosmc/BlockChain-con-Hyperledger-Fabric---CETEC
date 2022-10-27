# BlockChain con Hyperledger Fabric (HyL-f)

### En colaboracion con el CETEC en marco de la Universidad de Buenos Aires de Ingenieria - Desarrollado en solución a la persistencia de datos en una red descentralizada para la toma de asistencia en instituciones educactivas.



## Arquitectura 

![Imagen Diseño BlockChain](./img/dise%C3%B1oSoluci%C3%B3nBlockchain.png)

//describimos que seria cada uno de los componentes.

## Generamos configuraciones

### Creamos material criptografico // identidad de los "participantes"

* Creamos el archivo crypto-config.yml

    > ### crypto-config.yml
    > Tipo: Configuración
    > Previsualización del archivo:
    >
    >     OrdererOrgs:
    >       - Name: Orderer
    >         Domain: fiuba.com
    >         EnableNodeOUs: true
    >         Specs:
    >           - Hostname: orderer
    >           SANS:
    >               - localhost
    >     PeerOrgs:
    >       - Name: Org1
    >         Domain: org1.fiuba.com
    >         EnableNodeOUs: true
    >         Template:
    >           Count: 1
    >           SANS:
    >           - localhost
    >         Users:
    >           Count: 1

    [ir a archivo crypto-config.yml](./fiuba-network/crypto-config.yml)

* Ejecutamos en la terminal el comando

    >     cryptogen generate --config=./crypto-config.yml

### Resultados

Se crea una carpeta crypto-config que incluye las configuraciones de las organizaciones que conforman la red, incluyendo sus indentidades digitales.

// previsualizacion de la esctrucutura de directorios y certificaciones que deberia crearse por Hyperledger segun lo configurado por crypto-config.yml


### Creamos configuraciones de bloque y transacciones. // como se van a comunicar los participantes?

* Creamos el archivo configtx.yml

    > ### configtx.yml
    > Tipo: Configuración
    > Previsualización del archivo:
    >
    >     ...
    >     Profiles:
    >       ThreeOrgsOrdererGenesis:
    >           <<: *ChannelDefaults
    >           Orderer:
    >               <<: *OrdererDefaults
    >               Organizations:
    >                   - *OrdererOrg
    >               Capabilities:
    >                   <<: *OrdererCapabilities
    >           Consortiums:
    >               SampleConsortium:
    >               Organizations:
    >                   - *Org1
    >                   - *Org2
    >                   - *Org3
    >       ThreeOrgsChannel:
    >           Consortium: SampleConsortium
    >           <<: *ChannelDefaults
    >           Application:
    >               <<: *ApplicationDefaults
    >               Organizations:
    >                   - *Org1
    >                   - *Org2
    >                   - *Org3
    >           Capabilities:
    >               <<: *ApplicationCapabilities

    [ir a archivo configtx.yml](./fiuba-network/configtx.yml)

* Ejecutamos los siguientes comandos

    >### Generamos el genesis.block
    >     configtxgen -profile ThreeOrgsOrdererGenesis -channelID system-channel -outputBlock ./channel-artifacts/genesis.block
    >### Generamos channel.tx 
    >     configtxgen -profile ThreeOrgsChannel -channelID attendance -outputCreateChannelTx ./channel-artifacts/channel.tx
    >### Generamos Org1MSPanchors.tx
    >     configtxgen -profile ThreeOrgsChannel -channelID attendance -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -asOrg Org1MSP
    >### Generamos Org2MSPanchors.tx
    >     configtxgen -profile ThreeOrgsChannel -channelID attendance -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -asOrg Org2MSP
    >### Generamos Org3MSPanchors.tx
    >     configtxgen -profile ThreeOrgsChannel -channelID attendance -outputAnchorPeersUpdate ./channel-artifacts/Org3MSPanchors.tx -asOrg Org3MSP


### Resultados

Verificamos que se haya creado los archivos en channel-artifacts.

//previsualización diagrama de rama de archivos generados.


## Levantamos la red

Con docker compose vamos a levantar la red que utiliza todas las configuraciones previamente creadas.

### Contenedor base de los peers
Creamos el directorio de trabajo /base.

* En el directorio /base creamos el archivo peer-base.yaml, que se se encarga de crear contenedor base para los peers de (HyL-f)

    > ### peer-base.yaml
    >
    > Tipo: Levantar Red
    >
    > Previsualización del archivo:
    >
    >     version: '2'
    >     services:
    >         peer-base:
    >             image: hyperledger/fabric-peer:2.2.0
    >             environment:
    >                 - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
    >                 - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=fiuba-network_basic
    >                 - FABRIC_LOGGING_SPEC=INFO
    >                 - CORE_PEER_TLS_ENABLED=true
    >                 - CORE_PEER_GOSSIP_USELEADERELECTION=true
    >                 - CORE_PEER_GOSSIP_ORGLEADER=false
    >                 - CORE_PEER_PROFILE_ENABLED=true
    >                 - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
    >                 - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
    >                 - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
    >     working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    >     command: peer node start


    [ir a archivo peer-base.yaml](./fiuba-network/base/peer-base.yaml)

### Contenedores de los participantes y del servicio de ordenamiento

* Creamos el archivo docker-compose-base.yaml que se va a encargar de levantar los contenedores para los participantes y del servicio de ordenamiento.

    > ### docker-compose-base.yaml
    > 
    > Tipo : Levantar Red
    > 
    > Previsualización del archivo :
    > 
    >     services:
    >       orderer.fiuba.com:
    >           container_name: orderer.fiuba.com
    >           image: hyperledger/fabric-orderer:2.2.0
    >           environment:
    >               ...
    >           working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    >           command: orderer
    >           volumes:
    >               ...
    >           ports:
    >               - 7050:7050
    >
    >       peer0.org1.fiuba.com:
    >           container_name: peer0.org1.fiuba.com
    >           extends:
    >               file: peer-base.yaml
    >               service: peer-base
    >           environment:
    >               ...
    >           volumes:
    >               ...
    >           ports:
    >               - 7051:7051
    >               - 7053:7053
    >
    >        peer0.org2.fiuba.com:
    >           ...

    [ir a archivo docker-compose-base.yml](./fiuba-network/base/docker-compose-base.yaml)

### Orquestamos los contenedores creados y ademas se agregan las bases de datos y servicios CA y CLI

* Creamos el archivo docker-compose-cli-couchdb.yaml que se va a encargar de orquestar todos los contenedores en una sola red y ademas agrega las bases de datos de cada participante y servicios como CA y CLI necesarios para la arquitectura (HyL-f).

    > ### docker-compose-cli-couchdb.yaml
    > 
    > Tipo : Levantar Red
    > 
    > Previsualización del archivo :
    > 
    >     networks:
    >       basic:
    >     services:
    >       orderer.fiuba.com:
    >           extends:
    >               file:   base/docker-compose-base.yaml
    >               service: orderer.fiuba.com
    >           container_name: orderer.fiuba.com
    >           networks:
    >               - basic
    >
    >       peer0.org1.fiuba.com:
    >           container_name: peer0.org1.fiuba.com
    >           extends:
    >               file:  base/docker-compose-base.yaml
    >               service: peer0.org1.fiuba.com
    >           environment:
    >               ...
    >           depends_on:
    >               - orderer.fiuba.com
    >               - couchdb0
    >           networks:
    >               - basic
    >
    >       peer0.org2.fiuba.com:
    >           ...
    >       
    >       peer0.org3.fiuba.com:
    >           ...
    >
    >       ca.org1.fiuba.com:
    >           image: hyperledger/fabric-ca:1.4.8
    >           environment:
    >               ...
    >           ports:
    >               - "7054:7054"
    >           command: sh -c 'fabric-ca-server start -b admin:adminpw'
    >           volumes:
    >               ...
    >           container_name: ca.org1.fiuba.com
    >           networks:
    >               - basic
    >       cli:
    >           container_name: cli
    >           image: hyperledger/fabric-tools:2.2
    >           tty: true
    >           stdin_open: true
    >           environment:
    >               ...
    >           working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    >           command: /bin/bash
    >           volumes:
    >               ...
    >           depends_on:
    >               - orderer.fiuba.com
    >               - peer0.org1.fiuba.com
    >               - peer0.org2.fiuba.com
    >               - peer0.org3.fiuba.com
    >           networks:
    >               - basic
    >   
    >           couchdb0:
    >               image: couchdb:3.1
    >           environment:
    >               ...
    >           ports: 
    >               - 5984:5984
    >           container_name: couchdb0
    >           networks:
    >               - basic
    >
    >           couchdb1:
    >               ...
    >
    >           couchdb2:
    >               ...

    [ir a archivo docker-compose-cli-couchdb.yaml](./fiuba-network/docker-compose-cli-couchdb.yaml)

