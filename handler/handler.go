package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"github.com/mager/bouncer/config"
	"github.com/mager/bouncer/premint"
	"go.uber.org/zap"
)

var (
	// roleID is the role for Premint users
	roleID = "901476988349468702"
	// guildID is the Discord guild ID
	guildID = "892132883844726814"
)

// Handler struct for HTTP requests
type Handler struct {
	logger  *zap.SugaredLogger
	router  *mux.Router
	premint premint.PremintClient
	cfg     config.Config
	discord *discordgo.Session
}

// New creates a Handler struct
func New(
	logger *zap.SugaredLogger,
	router *mux.Router,
	premint premint.PremintClient,
	cfg config.Config,
	discord *discordgo.Session,
) *Handler {
	h := Handler{logger, router, premint, cfg, discord}
	h.registerRoutes()

	return &h
}

// RegisterRoutes registers all the routes for the route handler
func (h *Handler) registerRoutes() {
	h.router.HandleFunc("/getStatus", h.getStatus).Methods("POST")
	h.router.HandleFunc("/allowEntry", h.allowEntry).Methods("POST")
}

// Req is the request
type Req struct {
	Snowflake string `json:"snowflake"`
	Address   string `json:"address"`
}

// Resp is the response
type Resp struct {
	AddressInList  bool `json:"address_in_list"`
	DiscordRoleSet bool `json:"discord_role_set"`
}

// allowEntry is the route handler for the POST /allowEntry endpoint
func (h *Handler) allowEntry(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req Req
	)

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Make sure there is an address in the request
	if req.Address == "" {
		h.logger.Error("No address in request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make sure there is a Discord snowflake in the request
	if req.Snowflake == "" {
		h.logger.Error("No discord in request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		address   = req.Address
		snowflake = req.Snowflake
		resp      Resp
	)

	// Call the Premint API to make sure the address is in the list
	addresses, err := h.premint.GetWalletAddresses(address)
	if err != nil {
		h.logger.Errorf("Error getting wallet addresses: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If the address is not in the list, return an error
	if !stringInSlice(address, addresses) {
		h.logger.Errorf("Address %s not in list", address)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp.AddressInList = true

	// Make a request to Discord to add role
	var (
		roleID  = "901476988349468702"
		guildID = "892132883844726814"
	)

	err = h.discord.GuildMemberRoleAdd(guildID, snowflake, roleID)
	if err != nil {
		h.logger.Errorf("Error adding role to user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.DiscordRoleSet = true

	json.NewEncoder(w).Encode(resp)
}

// getStatus is the route handler for the POST /allowEntry endpoint
func (h *Handler) getStatus(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req Req
	)

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Make sure there is an address in the request
	if req.Address == "" {
		h.logger.Error("No address in request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make sure there is a Discord snowflake in the request
	if req.Snowflake == "" {
		h.logger.Error("No discord in request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		address   = req.Address
		snowflake = req.Snowflake
		resp      Resp
	)

	// Call the Premint API to make sure the address is in the list
	addresses, err := h.premint.GetWalletAddresses(address)
	if err != nil {
		h.logger.Errorf("Error getting wallet addresses: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If the address is not in the list, return an error
	if !stringInSlice(address, addresses) {
		h.logger.Errorf("Address %s not in list", address)
	} else {
		resp.AddressInList = true
	}

	// Make a request to Discord to fetch user
	m, err := h.discord.GuildMember(guildID, snowflake)
	if err != nil {
		h.logger.Errorf("Error fetching member: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if user has role
	resp.DiscordRoleSet = stringInSlice(roleID, m.Roles)

	json.NewEncoder(w).Encode(resp)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.EqualFold(a, b) {
			return true
		}
	}
	return false
}
