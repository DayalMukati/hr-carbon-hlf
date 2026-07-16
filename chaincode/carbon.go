package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for a carbon credit trading registry
type SmartContract struct {
	contractapi.Contract
}

// Credit represents a tradable carbon credit certificate.
type Credit struct {
	CreditID string `json:"CreditID"`
	Owner    string `json:"Owner"`
	Tonnes   int    `json:"Tonnes"`
	Status   string `json:"Status"`
}

// HistoryEntry represents one revision of a credit from the ledger history.
type HistoryEntry struct {
	TxID      string  `json:"TxID"`
	Value     *Credit `json:"Value"`
	Timestamp string  `json:"Timestamp"`
	IsDelete  bool    `json:"IsDelete"`
}

const (
	statusActive  = "ACTIVE"
	statusRetired = "RETIRED"
)

// IssueCredit creates a new carbon credit with status "ACTIVE".
func (s *SmartContract) IssueCredit(ctx contractapi.TransactionContextInterface, creditID string, owner string, tonnes int) error {
	if tonnes <= 0 {
		return fmt.Errorf("tonnes must be positive")
	}

	existing, err := ctx.GetStub().GetState(creditID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("credit %s already exists", creditID)
	}

	credit := Credit{
		CreditID: creditID,
		Owner:    owner,
		Tonnes:   tonnes,
		Status:   statusActive,
	}
	return putCredit(ctx, &credit)
}

// GetCredit returns the credit identified by creditID.
func (s *SmartContract) GetCredit(ctx contractapi.TransactionContextInterface, creditID string) (*Credit, error) {
	data, err := ctx.GetStub().GetState(creditID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("credit %s does not exist", creditID)
	}

	var credit Credit
	if err := json.Unmarshal(data, &credit); err != nil {
		return nil, err
	}
	return &credit, nil
}

// TransferCredit changes the owner of an ACTIVE credit to newOwner.
func (s *SmartContract) TransferCredit(ctx contractapi.TransactionContextInterface, creditID string, newOwner string) error {
	if newOwner == "" {
		return fmt.Errorf("new owner cannot be empty")
	}

	credit, err := s.GetCredit(ctx, creditID)
	if err != nil {
		return err
	}
	if credit.Status != statusActive {
		return fmt.Errorf("credit %s is not ACTIVE (current status: %s)", creditID, credit.Status)
	}

	credit.Owner = newOwner
	return putCredit(ctx, credit)
}

// UpdateTonnes corrects the tonnage of an ACTIVE credit.
func (s *SmartContract) UpdateTonnes(ctx contractapi.TransactionContextInterface, creditID string, newTonnes int) error {
	if newTonnes <= 0 {
		return fmt.Errorf("tonnes must be positive")
	}

	credit, err := s.GetCredit(ctx, creditID)
	if err != nil {
		return err
	}
	if credit.Status != statusActive {
		return fmt.Errorf("credit %s is not ACTIVE (current status: %s)", creditID, credit.Status)
	}

	credit.Tonnes = newTonnes
	return putCredit(ctx, credit)
}

// RetireCredit marks a credit as "RETIRED".
func (s *SmartContract) RetireCredit(ctx contractapi.TransactionContextInterface, creditID string) error {
	credit, err := s.GetCredit(ctx, creditID)
	if err != nil {
		return err
	}
	if credit.Status == statusRetired {
		return fmt.Errorf("credit %s is already RETIRED", creditID)
	}

	credit.Status = statusRetired
	return putCredit(ctx, credit)
}

// GetCreditHistory returns the full revision history of a credit, newest first.
func (s *SmartContract) GetCreditHistory(ctx contractapi.TransactionContextInterface, creditID string) ([]HistoryEntry, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(creditID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for %s: %v", creditID, err)
	}
	defer resultsIterator.Close()

	var history []HistoryEntry
	for resultsIterator.HasNext() {
		modification, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		entry := HistoryEntry{
			TxID:      modification.TxId,
			Timestamp: time.Unix(modification.Timestamp.Seconds, int64(modification.Timestamp.Nanos)).UTC().Format(time.RFC3339),
			IsDelete:  modification.IsDelete,
		}
		if !modification.IsDelete {
			var credit Credit
			if err := json.Unmarshal(modification.Value, &credit); err != nil {
				return nil, err
			}
			entry.Value = &credit
		}
		history = append(history, entry)
	}
	return history, nil
}

// --- helpers ---

func putCredit(ctx contractapi.TransactionContextInterface, credit *Credit) error {
	bytes, err := json.Marshal(credit)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(credit.CreditID, bytes)
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
