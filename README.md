
# Caching Proxy en Go

Este proyecto es un servidor proxy con capacidades de caching, desarrollado en Go. Utiliza Redis para almacenar en caché las respuestas de un servidor de origen.

## Estructura del Proyecto

## Requisitos

- Go 1.22.2 o superior
- Redis

## Instalación

1. Clona el repositorio:

    ```sh
    git clone https://github.com/Darcoprogramador/caching-proxy-go.git
    cd caching-proxy-go
    ```

2. Instala las dependencias:

    ```sh
    go mod tidy
    ```

3. Asegúrate de tener un servidor Redis corriendo en `localhost:6379`.

## Uso

Para iniciar el servidor proxy, ejecuta el siguiente comando:

```sh
go run main.go -port <PUERTO> -origin <SERVIDOR_ORIGEN>
```

Para borrar la cache del servidor proxy puedes hacer uso de el comando
```sh
go run main.go -port <PUERTO> -origin <SERVIDOR_ORIGEN> --clear-cache
```

****
Proyecto de: [roadmap.sh projects](https://roadmap.sh/projects/caching-server)