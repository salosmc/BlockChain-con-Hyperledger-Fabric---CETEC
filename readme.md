# BlockChain con Hyperledger Fabric

### En colaboracion con el CETEC en marco de la Universidad de Buenos Aires de Ingenieria - Desarrollado en solución a la persistencia de datos en una red descentralizada para la toma de asistencia en instituciones educactivas.



## Arquitectura 

//imagen

//describimos que seria cada uno de los componentes.

## Generamos configuraciones

### Creamos material criptografico //identidad de los "participantes"

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


### Creamos configuraciones de bloque y transacciones. // como se van a comunicar los participantes.

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





