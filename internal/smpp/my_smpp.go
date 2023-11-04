package my_smpp

import (
	"log"

	"github.com/fiorix/go-smpp/smpp"
)

type SMPPClient struct {
	transmitter *smpp.Transmitter
}

// NewSMPPClient creates a new instance of the SMPPClient and establishes an SMPP connection.
func NewSMPPClient() (*SMPPClient, error) {
	// Create an SMPP transmitter
	tx := &smpp.Transmitter{
		Addr:   "119.235.115.204:15019",
		User:   "Televid",
		Passwd: "Te17evid",
	}

	if err := tx.Bind(); err != nil {
		log.Fatalf("Failed to bind to SMPP server: %v", err)
		return nil, smpp.ErrNotBound
	}

	log.Println("SMPP connection established")

	return &SMPPClient{transmitter: tx}, nil
}
