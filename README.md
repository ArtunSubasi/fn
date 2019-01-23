<a id="top"></a>
![Fn Project](http://fnproject.io/images/fn-300x125.png)

# Fn Zeebe Extension Prototyp
* Dies ist privates Fork von dem GitHub-Projekt "fnproject/fn":  https://github.com/fnproject/fn
* Der Code von dem fnproject wurde nicht angepasst, sondern nur um die Zeebe-Extension erweitert. Die Zeebe-Extension befindet sich in dem Verzeichnis vendor/git.esentri.com/fn/ext-zeebe
* Nachdem die Zeebe-Extension in einem öffentlichen Repository wie GitHub angelegt wird, kann die Extension über den vorgesehenen Mechanismus mittels der ext.yaml Datei gebaut werden. Dann wird dieses Repository nicht mehr gebraucht.

Voraussetzungen:
* docker
* docker-compose

Build mit Docker:
```sh
./buildFnserver
```
das Skript baut ein Docker-Image `artunsubasi/fnserver`
Da das Bauen in einem Docker Container stattfindet, ist eine lokale Go-Installation nicht unbedingt notwendig.

# Zusammen mit Zeebe Broker und Simple Monitor starten
* https://github.com/ArtunSubasi/zeebe-simple-monitor auschecken
* Mit docker-compose starten:
```sh
docker/run
```
Die Docker Compose-Datei ist in dem obigen Repository so konfiguriert, dass der Docker Image `artunsubasi/fnserver` verwendet wird. D.h. der fnserver mit der Zeebe-Extension muss vorher gebaut werden. Der Fn Server wird in All-in-One-Modus gestartet und beinhaltet somit den API Server, den Load Balancer sowie einen Fn Runner Server im selben Container.

## Fn Server ohne Docker Compose starten:
fnserver Starten: 
```sh
docker run --rm -i --name fnserver \
    -e FN_LB_URL=http://localhost:8080 \ # test
    -e FN_API_SERVER_URL=http://localhost:8080 \
    -e FN_ZEEBE_GATEWAY_URL=http://localhost:26500 \
    -v ./fn/data:/app/data  \
    -v /var/run/docker.sock:/var/run/docker.sock  \
    -p 8080:8080  \
    artunsubasi/fnserver
```

## Umgebungsvariablen
* FN_LB_URL: Die URL des Fn Load Balancers
* FN_API_SERVER_URL: Die URL des Fn API Servers
* FN_ZEEBE_GATEWAY_URL= Die URL des Zeebe Gateways (gRPC-Port)

## Docker-Volumes
* /app/data ist die Datenbank mit den Fn Apps und Funktionen
* /var/run/docker.sock ist das Unix-Socket des Docker-Daemons, damit der Fn Server über den Docker-Deamon intern Docker Container verwalten zu können, da die Funktionen auch als Docker Container laufen

# Weitere Erkenntnisse
* fnserver lässt sich unter MacOS nicht nativ starten. Siehe Slack-Chat: https://fnproject.slack.com/archives/C6FL02Q02/p1547646705082200

