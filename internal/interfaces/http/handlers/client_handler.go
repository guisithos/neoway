package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/guisithos/neoway/internal/application/services"
	"github.com/guisithos/neoway/internal/domain/models"
)

type ClientHandler struct {
	clientService services.ClientService
}

func NewClientHandler(clientService services.ClientService) *ClientHandler {
	return &ClientHandler{clientService: clientService}
}

func normalizeDocument(document string) string {
	reg := regexp.MustCompile(`[^0-9]`)
	return reg.ReplaceAllString(document, "")
}

// Funcoes de validacao
func validateDocument(document string, clientType models.ClientType) error {
	doc := normalizeDocument(document)

	switch clientType {
	case models.PersonType:
		if len(doc) != 11 {
			return fmt.Errorf("CPF precisa ter 11 dígitos para o tipo PERSON")
		}
	case models.BusinessType:
		if len(doc) != 14 {
			return fmt.Errorf("CNPJ precisa ter 14 dígitos para o tipo BUSINESS")
		}
	default:
		return fmt.Errorf("tipo de cliente inválido: %s", clientType)
	}

	return nil
}

// Essa função checa se o documento (CPF ou CNPJ) está correto pro tipo de cliente.
// Se for pessoa física, o CPF tem q ter 11 números.
// Se for empresa, o CNPJ preciza ter 14 números.
// Primeiro a gente limpa o documento tirando tudo q não é número, após checa o tamanho.
// Se der ruim vai voltar erro.

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var client models.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	client.Document = normalizeDocument(client.Document)

	// Validar se o documento corresponde ao tipo de cliente
	if err := validateDocument(client.Document, client.Type); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.clientService.CreateClient(&client); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	log.Println("ListClients chamado")

	clients, err := h.clientService.ListClients()
	if err != nil {
		log.Printf("erro ao buscar clientes: %v", err)
		http.Error(w, "erro ao listar clientes", http.StatusInternalServerError)
		return
	}

	log.Printf("encontrado %d clientes", len(clients))

	w.Header().Set("Content-Type", "application/json")

	if len(clients) == 0 {
		log.Println("nenhum cliente encontrado, retornando array vazio")
		w.Write([]byte("[]"))
		return
	}

	if err := json.NewEncoder(w).Encode(clients); err != nil {
		log.Printf("erro ao codificar resposta: %v", err)
		http.Error(w, "erro ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("resposta enviada")
}

func (h *ClientHandler) BlockClient(w http.ResponseWriter, r *http.Request) {
	document := chi.URLParam(r, "document")
	log.Printf("recebendo requisição de bloqueio para o documento: %s", document)

	if document == "" {
		http.Error(w, "documento é obrigatório", http.StatusBadRequest)
		return
	}

	if err := h.clientService.BlockClient(document); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "cliente bloqueado"})
}

func (h *ClientHandler) UnblockClient(w http.ResponseWriter, r *http.Request) {
	document := chi.URLParam(r, "document")
	log.Printf("recebendo requisição de desbloqueio para o documento: %s", document)

	if document == "" {
		http.Error(w, "documento é obrigatório", http.StatusBadRequest)
		return
	}

	if err := h.clientService.UnblockClient(document); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "cliente desbloqueado"})
}

func (h *ClientHandler) GetClientByDocument(w http.ResponseWriter, r *http.Request) {
	document := chi.URLParam(r, "document")
	normalizedDoc := normalizeDocument(document)

	client, err := h.clientService.GetClientByDocument(normalizedDoc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}
