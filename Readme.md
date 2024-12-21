## Rate Limiter - Go Project ##
Este projeto é um Rate Limiter implementado em Go. Ele limita o número de requisições que podem ser feitas a um servidor web com base no endereço IP ou em um token de acesso. A lógica de limitação utiliza o Redis como mecanismo de persistência, mas pode ser estendida para outros armazenamentos com uma abordagem de estratégia.

# Recursos
 - Limitação por IP: Limita o número de requisições com base no endereço IP do cliente.
 - Limitação por Token: Limita o número de requisições com base em um token de acesso fornecido no header API_KEY.
 - Prioridade do Token: Configurações de limite baseadas no token de acesso substituem as configurações de limite baseadas no IP.
 - Persistência: Suporte para armazenamento Redis (configurável para suportar outros armazenamentos).
 - Bloqueio Temporário: Possibilidade de configurar um tempo de bloqueio após exceder o limite.
 - Middleware para Servidor Web: Fácil integração com o servidor web Go.
 - Configuração via .env: Parâmetros de limite configurados por variáveis de ambiente.
 - Resposta Padrão: Retorna HTTP 429 com a mensagem:
   `You have reached the maximum number of requests or actions allowed within a certain time frame.`

# Pré-requisitos
Docker e Docker Compose instalados.
Go instalado (caso queira rodar o servidor localmente sem Docker).

# Configuração
.env
Crie um arquivo .env na raiz do projeto com as seguintes configurações:

# env
  - REDIS_HOST=redis
  - REDIS_PORT=6379
  - REQUEST_LIMIT_IP=5   
  - BLOCK_TIME=300      

  ## Limitações por Token (formato: token=limite)
  - TOKEN_LIMITS=api_key_1=10,api_key_2=20

  ## Mecanismo de Armazenamento
  - STORAGE_BACKEND=redis

# Executando o Projeto
 - Com Docker Compose
    Certifique-se de que o Docker Compose está instalado.
    No diretório do projeto, execute:
    `docker-compose up --build`
    O servidor estará disponível em `http://localhost:8080`.

 - Sem Docker
    Certifique-se de que o Redis está rodando na máquina local (porta padrão 6379).
    Compile e execute o servidor:
    `go build -o main .`
    `./main`
    O servidor estará disponível em `http://localhost:8080`.

# Estrutura do Projeto
├── Dockerfile              # Configuração para build do Docker
├── docker-compose.yml      # Configuração do Docker Compose
├── utils/
│   └── config.go           # Funções para leitura de variáveis de ambiente
├── limiter/
│   ├── persistence_factory.go  # Estratégia para persistência (Redis ou outros)
│   ├── redis_limiter.go        # Implementação do Rate Limiter com Redis
│   ├── rate_limiter.go         # Lógica principal do Rate Limiter
├── middleware/
│   └── middleware.go       # Middleware para integração com o servidor HTTP
└── main.go                 # Ponto de entrada da aplicação

# Endpoints
Qualquer endpoint configurado no servidor será protegido pelo Rate Limiter.
Exemplo de Requisição
 - Header:
    http
    API_KEY: api_key_1
    Resposta em Caso de Exceder o Limite:
    ```json
    {
      "error": "You have reached the maximum number of requests or actions allowed within a certain time frame"
    }```

# Testando
Utilize ferramentas como Postman ou cURL para enviar requisições rápidas ao servidor e observar o comportamento do Rate Limiter.
Ajuste os limites no .env para simular diferentes cenários.

# Extensibilidade
Este projeto foi projetado para ser extensível
 - Novo Armazenamento: Adicione novas implementações ao diretório limiter/ e registre-as em persistence_factory.go.
 - Configurações Adicionais: Amplie as funções em utils/config.go para novos parâmetros no .env.
