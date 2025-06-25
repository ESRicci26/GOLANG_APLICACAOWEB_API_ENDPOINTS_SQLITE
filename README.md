# Sistema de Gestão de Produtos WEB API Endpoints - Golang

Aplicação WEB usando Golang referente um CRUD Completo de produtos usando banco de dados SQLITE "products.db" com layout responsivo.

## 🚀 Funcionalidades

- ✅ **CRUD Completo**: Criar, Ler, Atualizar e Deletar produtos
- ✅ **Interface Responsiva**: Layout adaptável para desktop e mobile
- ✅ **Banco de Dados SQLite**: Persistência de dados local
- ✅ **API REST**: Endpoints RESTful para todas as operações
- ✅ **Validação de Dados**: Validação tanto no frontend quanto no backend
- ✅ **Mensagens de Feedback**: Notificações visuais para ações do usuário

## 📋 Pré-requisitos

- Go 1.21 ou superior
- GCC (necessário para compilar o driver SQLite)

### Instalação do GCC:

**Windows:**
```bash
# Instalar através do Chocolatey
choco install mingw

# Ou baixar manualmente: https://www.mingw-w64.org/
```

**macOS:**
```bash
# Instalar através do Homebrew
brew install gcc
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install gcc
```

## 🛠️ Instalação

1. **Clone ou crie o projeto:**
```bash
mkdir products-crud
cd products-crud
```

2. **Crie os arquivos:**
   - Copie o conteúdo do `main.go`
   - Copie o conteúdo do `go.mod`

3. **Instale as dependências:**
```bash
go mod tidy
```

4. **Execute o servidor:**
```bash
go run main.go
```

5. **Acesse a aplicação:**
   - Abra seu navegador em: `http://localhost:8080`

## 📁 Estrutura do Projeto

```
products-crud/
├── main.go          # Servidor principal com handlers
├── go.mod           # Dependências do Go
├── go.sum           # Checksums das dependências (gerado automaticamente)
└── products.db      # Banco SQLite (criado automaticamente)
```

## 🔧 API Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/products` | Lista todos os produtos |
| GET | `/api/products/{id}` | Busca produto por ID |
| POST | `/api/products` | Cria novo produto |
| PUT | `/api/products/{id}` | Atualiza produto |
| DELETE | `/api/products/{id}` | Deleta produto específico |
| DELETE | `/api/products` | Deleta todos os produtos |

### Exemplo de Payload (JSON):
```json
{
  "name": "Smartphone Samsung",
  "seller": "Samsung",
  "price": 1299.99
}
```

## 🎨 Recursos Mantidos do Original

- **Layout Responsivo**: Funciona em desktop e mobile
- **Validação de Formulário**: Campos obrigatórios
- **Feedback Visual**: Mensagens de sucesso/erro com animações
- **Operações CRUD**: Todas as funcionalidades originais
- **Design Bootstrap**: Interface moderna e responsiva

## 🔄 Diferenças do Original

| Aspecto | Original (JavaScript) | Golang |
|---------|----------------------|---------|
| **Armazenamento** | IndexedDB (browser) | SQLite (servidor) |
| **Arquitetura** | Client-side | Server-side |
| **Persistência** | Local browser | Arquivo de banco |
| **Escalabilidade** | Limitada | Maior |
| **Multiplataforma** | Apenas web | Web + API |

## 🚀 Melhorias Implementadas

1. **Banco de Dados Persistente**: SQLite ao invés de IndexedDB
2. **API REST**: Endpoints para integração com outras aplicações
3. **Validação Robusta**: Validação no frontend e backend
4. **Código Organizado**: Estrutura clara e modular
5. **Error Handling**: Tratamento adequado de erros
6. **CORS Support**: Suporte para requisições cross-origin

## 📱 Responsividade

O layout se adapta automaticamente para diferentes tamanhos de tela:

- **Desktop**: Layout completo com tabela expandida
- **Tablet**: Tabela com scroll horizontal
- **Mobile**: Formulário e botões empilhados verticalmente

## 🛡️ Segurança

- Validação de entrada no servidor
- Prepared statements para evitar SQL injection
- Tratamento adequado de erros
- Headers CORS configurados

## 🎯 Próximos Passos

Para expandir o projeto, você pode:

1. **Adicionar Autenticação**: JWT tokens
2. **Implementar Paginação**: Para muitos produtos
3. **Adicionar Filtros**: Busca por nome, fabricante, preço
4. **Upload de Imagens**: Para produtos
5. **Docker**: Containerização da aplicação
6. **Testes**: Unit tests e integration tests

## 📞 Suporte

Se encontrar problemas:

1. Verifique se o Go está instalado: `go version`
2. Verifique se o GCC está instalado: `gcc --version`
3. Execute `go mod tidy` para resolver dependências
4. Verifique se a porta 8080 está livre

## 📄 Licença

Este projeto é uma conversão do meu projeto original IndexedDB para Golang.
