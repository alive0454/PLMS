package models

type PersonInfo struct {
	Person   Person            `json:"person"`
	Bicycles []ElectricBicycle `json:"bicycles"`
}
