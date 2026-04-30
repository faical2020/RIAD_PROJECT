package logic

import "errors"

func ValidateReservation(r Reservation, c Chambre) error {
    if r.DateDebut >= r.DateFin {
        return errors.New("dates invalides")
    }
    if !c.CanBook() {
        return errors.New("chambre non disponible")
    }
    return nil
}

func (r *Reservation) Checkin(c *Chambre) error {
    if r.Statut != "confirmée" {
        return errors.New("réservation non confirmée")
    }
    r.Statut = "checkin"
    c.Statut = "occupee"
    return nil
}

func (r *Reservation) Checkout(c *Chambre) error {
    if r.Statut != "checkin" {
        return errors.New("réservation n'est pas en checkin")
    }
    r.Statut = "checkout"
    c.Statut = "libre"
    return nil
}