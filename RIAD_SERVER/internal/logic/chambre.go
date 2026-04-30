package logic

import "errors"

func ValidateChambre(c Chambre) error {
    if c.Numero <= 0 {
        return errors.New("numéro de chambre invalide")
    }
    if c.Prix <= 0 {
        return errors.New("prix invalide")
    }
    return nil
}

func (c *Chambre) CanBook() bool {
    return c.Statut == "libre"
}