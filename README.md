# Rate Limiter

## Como executar

1. Inicie os contêineres com o comando:

    ```bash
    docker-compose up -d --build
    ```

    Este comando iniciará o contêiner do Redis e o contêiner da API. A API estará disponível na porta `8080` e o Redis na porta `6379`.

2. Teste a API com o comando:

    ```bash
    curl -X GET http://localhost:8080
    ```

    ou com requisições utilizando token:

    ```bash
    curl -X GET http://localhost:8080/ -H "Content-Type: application/json" -H "API_KEY: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
    ```

3. Configure as variáveis do limitador de taxa no arquivo `.env` na pasta raiz. As variáveis disponíveis são:
    - `MAX_REQUESTS_WITHOUT_TOKEN_PER_SECOND`: Número máximo de requisições por segundo sem o cabeçalho `API_KEY`.
    - `MAX_REQUESTS_WITH_TOKEN_PER_SECOND`: Número máximo de requisições por segundo com o cabeçalho `API_KEY`.
    - `TIME_BLOCK_IN_SECOND`: Tempo em segundos que o IP ou token será bloqueado.

    **Nota:** Após alterar o arquivo `.env`, reinicie o contêiner da API para aplicar as mudanças usando o comando:

    ```bash
    docker-compose up -d
    ```

4. Execute os testes automatizados com o comando:

    ```bash
    go test ./...
    ```
