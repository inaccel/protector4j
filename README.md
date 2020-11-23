# What is Protector4J

[VLINX Protector4J](https://protector4j.com) is a tool to prevent Java applications from decompilation. Protector4J provides a custom native ClassLoader by modifying the JVM. The Java classes are encrypted by AES and decrypted in the native ClassLoader.

# Quick reference

* **Maintained by**:

	[InAccel](https://inaccel.com)

* **Source of this description**:

	[inaccel/protector4j repo](https://github.com/inaccel/protector4j)

# How to use this image

You can run a VLINX Protector4J task by using this Docker image directly, passing the Java task parameters to `docker run`:

```console
$ docker run -it --rm --name my-java-task -u $(id -u):$(id -g) -v $(pwd):/usr/src/myvlinx -w /usr/src/myvlinx inaccel/protector4j \
	--version <jre-version> \
	--email <account-email> \
	--password <md5-of-password> \
	--protect-all \
	jar-path1 jar-path2 ...
```

## Building local Docker image (optional)

This is a base image that you can extend, so it has the bare minimum packages needed. If you add custom package(s) to the `Dockerfile`, then you can build your local Docker image like this:

```console
$ docker-compose build
```

# Multi-stage Builds

You can build your application with Maven, protect it with Protector4J and package everything in an image that does not include Maven nor Protector4J using [multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build).

```dockerfile
# build
FROM maven
WORKDIR /usr/src/app
COPY pom.xml .
RUN mvn -B -e -C -T 1C org.apache.maven.plugins:maven-dependency-plugin:3.1.1:go-offline
COPY . .
RUN mvn -B -e -o -T 1C verify

# protect
FROM inaccel/protector4j
WORKDIR /usr/src/app
COPY --from=0 /usr/src/app .
ARG VLINX_EMAIL
ARG VLINX_PASSWORD
RUN vlinx-protector4j \
	--version 11 \
	--email ${VLINX_EMAIL} \
	--password ${VLINX_PASSWORD} \
	--protect-all \
	target/*.jar

# package without maven, protector4j
FROM debian
ENV JAVA_HOME="/usr/vlinx/jre"
ENV PATH="${JAVA_HOME}/bin:${PATH}"
COPY --from=1 /usr/src/app/jre ${JAVA_HOME}
COPY --from=1 /usr/src/app/target/*.jar ./
```
