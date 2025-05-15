# Sistema de Teste de Carga em Go (loadtester)

Este projeto contém uma aplicação CLI em Go para realizar testes de carga em um serviço web.

## Funcionalidades

- Permite especificar a URL do serviço, o número total de requests e a quantidade de chamadas simultâneas via parâmetros de linha de comando.
- Realiza requests HTTP para a URL especificada, distribuindo os requests de acordo com o nível de concorrência.
- Gera um relatório ao final dos testes contendo:
    - Tempo total gasto na execução.
    - Quantidade total de requests realizados.
    - Quantidade de requests com status HTTP 200.
    - Distribuição de outros códigos de status HTTP.


## Como Compilar e Executar Localmente

1.  **Navegue até o diretório do projeto:**
    ```bash
    cd /home/ubuntu/loadtester
    ```

2.  **Compile a aplicação:**
    ```bash
    go build -o loadtester_local .
    ```
    Isso gerará um executável chamado `loadtester_local` no diretório atual.

3.  **Execute a aplicação:**
    ```bash
    ./loadtester_local --url=<URL_DO_SERVICO> --requests=<NUMERO_DE_REQUESTS> --concurrency=<NUMERO_DE_CHAMADAS_SIMULTANEAS>
    ```
    **Exemplo:**
    ```bash
    ./loadtester_local --url=http://httpbin.org/status/200 --requests=100 --concurrency=10
    ```

## Como Construir e Executar com Docker

1.  **Navegue até o diretório do projeto onde o `Dockerfile` está localizado:**
    ```bash
    cd /home/ubuntu/loadtester
    ```

2.  **Construa a imagem Docker:**
    ```bash
    docker build -t loadtester_app .
    ```

3.  **Execute a aplicação usando o contêiner Docker:**
    ```bash
    docker run loadtester_app --url=<URL_DO_SERVICO> --requests=<NUMERO_DE_REQUESTS> --concurrency=<NUMERO_DE_CHAMADAS_SIMULTANEAS>
    ```
    **Exemplo:**
    ```bash
    docker run loadtester_app --url=http://httpbin.org/status/200 --requests=1000 --concurrency=10
    ```

## Parâmetros da CLI

-   `--url`: (Obrigatório) URL do serviço a ser testado.
-   `--requests`: (Opcional, padrão: 100) Número total de requests a serem realizados.
-   `--concurrency`: (Opcional, padrão: 10) Número de chamadas simultâneas.


