package handlers

import (
	"net/http"

	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/logic"
	"github.com/gin-gonic/gin"
)

var uuidNil = "00000000-0000-0000-0000-000000000000"

// ── Services ──

func GetServices(c *gin.Context) {
	services, err := logic.GetServices(db.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

func CreateService(c *gin.Context) {
	var s logic.Service
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := logic.CreateService(db.GetDB(), &s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

func UpdateService(c *gin.Context) {
	var s logic.Service
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.ID = c.Param("id")
	if err := logic.UpdateService(db.GetDB(), &s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func DeleteService(c *gin.Context) {
	id := c.Param("id")
	if err := logic.DeleteService(db.GetDB(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "service supprimé"})
}

// ── Consommations ──

func GetConsommations(c *gin.Context) {
	reservationID := c.Param("id")
	consommations, err := logic.GetConsommationsByReservation(db.GetDB(), reservationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, consommations)
}

func AddConsommation(c *gin.Context) {
	reservationID := c.Param("id")
	userID, _ := c.Get("user_id")

	var input struct {
		ServiceID    string  `json:"service_id"`
		Libelle      string  `json:"libelle"`
		Quantite     int     `json:"quantite"`
		PrixUnitaire float64 `json:"prix_unitaire"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serviceID := input.ServiceID
	if serviceID == "" {
		serviceID = uuidNil
	}

	conso := logic.Consommation{
		ReservationID: reservationID,
		ServiceID:     serviceID,
		Libelle:       input.Libelle,
		Quantite:      input.Quantite,
		PrixUnitaire:  input.PrixUnitaire,
		AjoutePar:     userID.(string),
	}

	if err := logic.AddConsommation(db.GetDB(), &conso); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, conso)
}

func DeleteConsommation(c *gin.Context) {
	id := c.Param("id")
	if err := logic.DeleteConsommation(db.GetDB(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "consommation supprimée"})
}

// ── Facture ──

func GetFacture(c *gin.Context) {
	reservationID := c.Param("id")
	facture, err := logic.GetFacture(db.GetDB(), reservationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, facture)
}

func AddPaiement(c *gin.Context) {
	reservationID := c.Param("id")
	userID, _ := c.Get("user_id")

	var input struct {
		Montant      float64 `json:"montant"`
		ModePaiement string  `json:"mode_paiement"`
		Reference    string  `json:"reference"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paiement := logic.Paiement{
		ReservationID: reservationID,
		Montant:       input.Montant,
		ModePaiement:  input.ModePaiement,
		Reference:     input.Reference,
		EncaissePar:   userID.(string),
	}

	if err := logic.AddPaiement(db.GetDB(), &paiement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, paiement)
}
