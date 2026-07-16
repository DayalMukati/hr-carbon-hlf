package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for a carbon credit trading registry
type SmartContract struct {
	contractapi.Contract
}

// Credit represents a tradable carbon credit certificate.
type Credit struct {
	CreditID string `json:"CreditID"` // unique credit id, e.g. "cc1"
	Owner    string `json:"Owner"`    // current owner name
	Tonnes   int    `json:"Tonnes"`   // tonnes of CO2 represented
	Status   string `json:"Status"`   // ACTIVE | RETIRED
}

// HistoryEntry represents one revision of a credit from the ledger history.
type HistoryEntry struct {
	TxID      string  `json:"TxID"`
	Value     *Credit `json:"Value"`
	Timestamp string  `json:"Timestamp"`
	IsDelete  bool    `json:"IsDelete"`
}

// IssueCredit creates a new carbon credit with status "ACTIVE".
// It must fail if the credit already exists or tonnes is not positive.
func (s *SmartContract) IssueCredit(ctx contractapi.TransactionContextInterface, creditID string, owner string, tonnes int) error {

	return nil
}

// GetCredit returns the credit identified by creditID.
// It must fail if the credit does not exist.
func (s *SmartContract) GetCredit(ctx contractapi.TransactionContextInterface, creditID string) (*Credit, error) {

	return nil, nil
}

// TransferCredit changes the owner of an ACTIVE credit to newOwner.
// It must fail if the credit does not exist, is RETIRED, or newOwner is empty.
func (s *SmartContract) TransferCredit(ctx contractapi.TransactionContextInterface, creditID string, newOwner string) error {

	return nil
}

// UpdateTonnes corrects the tonnage of an ACTIVE credit (for example after a
// re-audit of the underlying project).
// It must fail if the credit does not exist, is RETIRED, or newTonnes is not
// positive.
func (s *SmartContract) UpdateTonnes(ctx contractapi.TransactionContextInterface, creditID string, newTonnes int) error {

	return nil
}

// RetireCredit marks a credit as "RETIRED" (permanently consumed, cannot be
// transferred or updated again). It must fail if the credit does not exist or
// is already RETIRED.
func (s *SmartContract) RetireCredit(ctx contractapi.TransactionContextInterface, creditID string) error {

	return nil
}

// GetCreditHistory returns the full revision history of a credit, newest first,
// using GetHistoryForKey.
func (s *SmartContract) GetCreditHistory(ctx contractapi.TransactionContextInterface, creditID string) ([]HistoryEntry, error) {

	return nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		panic("Error creating carbon chaincode: " + err.Error())
	}

	if err := chaincode.Start(); err != nil {
		panic("Error starting carbon chaincode: " + err.Error())
	}
}
