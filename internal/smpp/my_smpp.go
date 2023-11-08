package my_smpp

import (
	"log"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

type SMPPConnection struct {
    transmitter *smpp.Transmitter
    statusChan  <-chan smpp.ConnStatus
}

func NewSMPPConnection() (*SMPPConnection, error) {
    // Create an SMPP transmitter
    tx := &smpp.Transmitter{
        Addr:   "119.235.115.204:15019",
        User:   "Televid",
        Passwd: "Te17evid",
    }

    statusChan := tx.Bind()

    // Wait for the connection to be fully established
    for status := range statusChan {
        if status.Status() == smpp.Connected {
            log.Println("SMPP connection established")
            return &SMPPConnection{transmitter: tx, statusChan: statusChan}, nil
        }
    }

    return nil, smpp.ErrNotBound
}

func (conn *SMPPConnection) SendSMS(phoneNumber string, message string) error {
    // Create an SMS
    sms := &smpp.ShortMessage{
        Src:      "+99362008971", // Replace with your source number
        Dst:      phoneNumber,
        Text:     pdutext.Raw(message), // Use pdutext.Raw to encode the message
        Register: pdufield.NoDeliveryReceipt,
    }

    _, err := conn.transmitter.Submit(sms)
    if err != nil {
        log.Printf("Failed to send SMS: %v", err)
        return err
    }

    log.Printf("SMS sent successfully to %s", phoneNumber)
    return nil
}
