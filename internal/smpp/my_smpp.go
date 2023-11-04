package my_smpp

import (
	"log"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

type SMPPConnection struct {
    transmitter *smpp.Transmitter
}

func NewSMPPConnection() (*SMPPConnection, error) {
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

    return &SMPPConnection{transmitter: tx}, nil
}

func (conn *SMPPConnection) SendSMS(phoneNumber string, message string) error {
    // Create an SMS
    sms := &smpp.ShortMessage{
        Src:      "+99362008971", // Replace with your source number
        Dst:      phoneNumber,
        Text:     pdutext.Raw(message), // Use pdutext.Raw to encode the message
        Register: pdufield.NoDeliveryReceipt,
        // Add any other required parameters
    }

    // Send the SMS using the established SMPP connection
    if _, err := conn.transmitter.Submit(sms); err != nil {
        log.Printf("Failed to send SMS: %v", err)
        return err
    }

    log.Printf("SMS sent successfully to %s", phoneNumber)
    return nil
}