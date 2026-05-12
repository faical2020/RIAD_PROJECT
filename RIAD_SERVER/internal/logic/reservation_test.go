package logic

import (
	"testing"
)

func TestValidateReservation_Valid(t *testing.T) {
	chambre := Chambre{Statut: "libre"}
	res := Reservation{
		DateDebut: "2026-06-01",
		DateFin:   "2026-06-05",
		Montant:   1000,
	}

	err := ValidateReservation(res, chambre)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateReservation_InvalidDates(t *testing.T) {
	chambre := Chambre{Statut: "libre"}
	res := Reservation{
		DateDebut: "2026-06-10",
		DateFin:   "2026-06-05",
	}

	err := ValidateReservation(res, chambre)
	if err == nil {
		t.Error("expected error for invalid dates")
	}
}

func TestValidateReservation_UnavailableRoom(t *testing.T) {
	chambre := Chambre{Statut: "occupee"}
	res := Reservation{
		DateDebut: "2026-06-01",
		DateFin:   "2026-06-05",
	}

	err := ValidateReservation(res, chambre)
	if err == nil {
		t.Error("expected error for unavailable room")
	}
}

func TestChambreCanBook(t *testing.T) {
	tests := []struct {
		statut string
		want   bool
	}{
		{"libre", true},
		{"occupee", false},
		{"occupe", false},
		{"nettoyage", false},
	}

	for _, tt := range tests {
		c := Chambre{Statut: tt.statut}
		if got := c.CanBook(); got != tt.want {
			t.Errorf("CanBook(%s) = %v, want %v", tt.statut, got, tt.want)
		}
	}
}

func TestReservationCheckin(t *testing.T) {
	res := Reservation{Statut: "confirmée"}
	chambre := Chambre{Statut: "libre"}

	if err := res.Checkin(&chambre); err != nil {
		t.Fatalf("Checkin failed: %v", err)
	}

	if res.Statut != "checkin" {
		t.Errorf("expected statut checkin, got %s", res.Statut)
	}
	if chambre.Statut != "occupee" {
		t.Errorf("expected chambre occupee, got %s", chambre.Statut)
	}
}

func TestReservationCheckin_NotConfirmed(t *testing.T) {
	res := Reservation{Statut: "en attente"}
	chambre := Chambre{}

	if err := res.Checkin(&chambre); err == nil {
		t.Error("expected error for non-confirmed reservation")
	}
}

func TestReservationCheckout(t *testing.T) {
	res := Reservation{Statut: "checkin"}
	chambre := Chambre{Statut: "occupee"}

	if err := res.Checkout(&chambre); err != nil {
		t.Fatalf("Checkout failed: %v", err)
	}

	if res.Statut != "checkout" {
		t.Errorf("expected statut checkout, got %s", res.Statut)
	}
	if chambre.Statut != "libre" {
		t.Errorf("expected chambre libre, got %s", chambre.Statut)
	}
}

func TestReservationCheckout_NotCheckedIn(t *testing.T) {
	res := Reservation{Statut: "confirmée"}
	chambre := Chambre{}

	if err := res.Checkout(&chambre); err == nil {
		t.Error("expected error for non-checked-in reservation")
	}
}
