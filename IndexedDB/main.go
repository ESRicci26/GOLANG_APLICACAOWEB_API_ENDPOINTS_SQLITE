package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	ID     int     `json:"id" db:"id"`
	Name   string  `json:"name" db:"name"`
	Seller string  `json:"seller" db:"seller"`
	Price  float64 `json:"price" db:"price"`
}

type Server struct {
	db *sql.DB
}

func NewServer() *Server {
	db, err := sql.Open("sqlite3", "./products.db")
	if err != nil {
		log.Fatal(err)
	}

	// Criar tabela se não existir
	createTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		seller TEXT NOT NULL,
		price REAL NOT NULL
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	return &Server{db: db}
}

func (s *Server) Close() {
	s.db.Close()
}

// Handlers
func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <meta http-equiv="X-UA-Compatible" content="ie=edge" />
  <title>Golang Products CRUD</title>
  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"/>
  <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.8.1/css/all.css"/>
  <style>
    body { box-sizing: border-box; }
    .d-flex input { margin: .9em 0em; }
    .d-flex button { margin: 1.5em .6em; padding: .3em 2.4em; }
    .d-flex table { margin: 1em 10em; }
    table .btnedit { color: lightgreen; cursor: pointer; }
    table .btndelete { color: tomato; cursor: pointer; }
    .updatemsg, .deletemsg, .insertmsg {
      position: absolute; top: -40px; z-index: 1000;
    }
    .movedown { animation: slideup 3.4s ease; }
    @keyframes slideup {
      50% { top: 0 }
      100% { top: -40px; }
    }
    @media (max-width: 768px) {
      .d-flex table { margin: 1em 1em; }
      .w-50 { width: 90% !important; }
    }
  </style>
</head>
<body>
  <main>
    <div class="container text-center">
      <h1 class="bg-light py-4 text-info">
        <i class="fas fa-plug"></i> Ricci Amazon Eletrônicos S/A
      </h1>
      <div class="d-flex justify-content-center">
        <form class="w-50" id="productForm">
          <input type="hidden" id="userid" />
          <input type="text" id="proname" class="form-control" autocomplete="off" placeholder="Nome Produto" required />
          <div class="row">
            <div class="col">
              <input type="text" id="seller" class="form-control m-0" autocomplete="off" placeholder="Fabricante" required />
            </div>
            <div class="col">
              <input type="number" step="0.01" id="price" class="form-control m-0" autocomplete="off" placeholder="Valor" required />
            </div>
          </div>
        </form>
      </div>
      <div class="d-flex justify-content-center flex-wrap">
        <button class="btn btn-success" id="btn-create">Inserir <i class="fas fa-plus-square"></i></button>
        <button class="btn btn-primary" id="btn-read">Lista <i class="fas fa-list-ol"></i></button>
        <button class="btn btn-warning" id="btn-update">Alterar <i class="fas fa-edit"></i></button>
        <button class="btn btn-danger" id="btn-delete">Deletar Todos <i class="fas fa-trash-alt"></i></button>
      </div>

      <div class="d-flex table-data table-responsive">
        <table class="table table-striped w-100">
          <thead>
            <tr>
              <th scope="col">ID</th>
              <th scope="col">Nome Produto</th>
              <th scope="col">Fabricante</th>
              <th scope="col">Preço</th>
              <th scope="col">Alterar</th>
              <th scope="col">Deletar</th>
            </tr>
          </thead>
          <tbody id="tbody">
          </tbody>
        </table>
      </div>
      <div id="notfound"></div>
    </div>
  </main>

  <div class="w-100 btn btn-success insertmsg">Dados INSERIDOS com sucesso...!</div>
  <div class="w-100 btn btn-warning updatemsg">Dados ALTERADOS com sucesso..!</div>
  <div class="w-100 btn btn-danger deletemsg">Dados EXCLUIDOS com sucesso...!</div>

  <script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
  <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js"></script>

  <script>
    const API_BASE = '/api/products';
    
    // Elementos do DOM
    const userid = document.getElementById('userid');
    const proname = document.getElementById('proname');
    const seller = document.getElementById('seller');
    const price = document.getElementById('price');
    const tbody = document.getElementById('tbody');
    const notfound = document.getElementById('notfound');

    // Funções utilitárias
    function showMessage(messageClass) {
      const element = document.querySelector('.' + messageClass);
      element.classList.add('movedown');
      setTimeout(() => {
        element.classList.remove('movedown');
      }, 4000);
    }

    function clearForm() {
      userid.value = '';
      proname.value = '';
      seller.value = '';
      price.value = '';
    }

    // Função para buscar produtos
    async function fetchProducts() {
      try {
        const response = await fetch(API_BASE);
        const products = await response.json();
        return products || [];
      } catch (error) {
        console.error('Erro ao buscar produtos:', error);
        return [];
      }
    }

    // Função para renderizar tabela
    async function renderTable() {
      tbody.innerHTML = '';
      notfound.textContent = '';
      
      const products = await fetchProducts();
      
      if (products.length === 0) {
        notfound.textContent = 'Nenhum registro encontrado no banco de dados...!';
        return;
      }
      
      products.forEach(product => {
        const row = document.createElement('tr');
        row.innerHTML = ` + "`" + `
          <td>${product.id}</td>
          <td>${product.name}</td>
          <td>${product.seller}</td>
          <td>R$ ${product.price.toFixed(2)}</td>
          <td><i class="fas fa-edit btnedit" data-id="${product.id}"></i></td>
          <td><i class="fas fa-trash-alt btndelete" data-id="${product.id}"></i></td>
        ` + "`" + `;
        tbody.appendChild(row);
      });
      
      // Adicionar event listeners
      document.querySelectorAll('.btnedit').forEach(btn => {
        btn.addEventListener('click', editProduct);
      });
      
      document.querySelectorAll('.btndelete').forEach(btn => {
        btn.addEventListener('click', deleteProduct);
      });
    }

    // Função para criar produto
    async function createProduct() {
      if (!proname.value || !seller.value || !price.value) {
        alert('Preencha todos os campos!');
        return;
      }

      try {
        const response = await fetch(API_BASE, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            name: proname.value,
            seller: seller.value,
            price: parseFloat(price.value)
          })
        });

        if (response.ok) {
          showMessage('insertmsg');
          clearForm();
          renderTable();
        }
      } catch (error) {
        console.error('Erro ao criar produto:', error);
      }
    }

    // Função para atualizar produto
    async function updateProduct() {
      if (!userid.value) {
        alert('Selecione um produto para alterar!');
        return;
      }

      try {
        const response = await fetch(` + "`" + `${API_BASE}/${userid.value}` + "`" + `, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            name: proname.value,
            seller: seller.value,
            price: parseFloat(price.value)
          })
        });

        if (response.ok) {
          showMessage('updatemsg');
          clearForm();
          renderTable();
        }
      } catch (error) {
        console.error('Erro ao atualizar produto:', error);
      }
    }

    // Função para editar produto
    async function editProduct(event) {
      const id = event.target.dataset.id;
      try {
        const response = await fetch(` + "`" + `${API_BASE}/${id}` + "`" + `);
        const product = await response.json();
        
        userid.value = product.id;
        proname.value = product.name;
        seller.value = product.seller;
        price.value = product.price;
      } catch (error) {
        console.error('Erro ao buscar produto:', error);
      }
    }

    // Função para deletar produto individual
    async function deleteProduct(event) {
      const id = event.target.dataset.id;
      if (confirm('Tem certeza que deseja excluir este produto?')) {
        try {
          const response = await fetch(` + "`" + `${API_BASE}/${id}` + "`" + `, {
            method: 'DELETE'
          });

          if (response.ok) {
            renderTable();
          }
        } catch (error) {
          console.error('Erro ao deletar produto:', error);
        }
      }
    }

    // Função para deletar todos os produtos
    async function deleteAllProducts() {
      if (confirm('Tem certeza que deseja excluir TODOS os produtos?')) {
        try {
          const response = await fetch(API_BASE, {
            method: 'DELETE'
          });

          if (response.ok) {
            showMessage('deletemsg');
            clearForm();
            renderTable();
          }
        } catch (error) {
          console.error('Erro ao deletar produtos:', error);
        }
      }
    }

    // Event Listeners
    document.getElementById('btn-create').addEventListener('click', createProduct);
    document.getElementById('btn-read').addEventListener('click', renderTable);
    document.getElementById('btn-update').addEventListener('click', updateProduct);
    document.getElementById('btn-delete').addEventListener('click', deleteAllProducts);

    // Carregar dados na inicialização
    window.addEventListener('load', renderTable);
  </script>
</body>
</html>`

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, nil)
}

// API Handlers
func (s *Server) getProductsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("SELECT id, name, seller, price FROM products ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Seller, &p.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (s *Server) getProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p Product
	err = s.db.QueryRow("SELECT id, name, seller, price FROM products WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.Seller, &p.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := s.db.Exec("INSERT INTO products (name, seller, price) VALUES (?, ?, ?)",
		p.Name, p.Seller, p.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (s *Server) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = s.db.Exec("UPDATE products SET name = ?, seller = ?, price = ? WHERE id = ?",
		p.Name, p.Seller, p.Price, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	_, err = s.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) deleteAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := s.db.Exec("DELETE FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reset auto increment
	_, err = s.db.Exec("DELETE FROM sqlite_sequence WHERE name='products'")
	if err != nil {
		log.Printf("Erro ao resetar sequência: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	path := r.URL.Path
	method := r.Method

	switch {
	case path == "/api/products" && method == "GET":
		s.getProductsHandler(w, r)
	case path == "/api/products" && method == "POST":
		s.createProductHandler(w, r)
	case path == "/api/products" && method == "DELETE":
		s.deleteAllProductsHandler(w, r)
	case strings.HasPrefix(path, "/api/products/") && method == "GET":
		s.getProductHandler(w, r)
	case strings.HasPrefix(path, "/api/products/") && method == "PUT":
		s.updateProductHandler(w, r)
	case strings.HasPrefix(path, "/api/products/") && method == "DELETE":
		s.deleteProductHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func main() {
	server := NewServer()
	defer server.Close()

	http.HandleFunc("/", server.homeHandler)
	http.HandleFunc("/api/products", server.apiHandler)
	http.HandleFunc("/api/products/", server.apiHandler)

	fmt.Println("🚀 Servidor rodando em http://localhost:8080")
	fmt.Println("📦 Banco de dados SQLite: products.db")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
