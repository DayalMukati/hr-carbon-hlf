# hr-carbon-hlf

Solution repository for the **Carbon Credit Trading** Hyperledger Fabric chaincode
challenge (NPCI / HackerRank, Hard).

Standard Fabric **test-network** plus a chaincode skeleton at
[`chaincode/carbon.go`](chaincode/carbon.go). Cloned into the candidate's
environment by the HackerRank Setup Script (via [`setup.sh`](setup.sh)).

## Candidate task
1. Implement the functions in `chaincode/carbon.go`, including a range query over
   the world state and an aggregation that excludes retired credits.
2. Deploy: `cd test-network && ./network.sh deployCC -ccn carboncc -ccp ../chaincode -ccl go`
3. Issue cc1/cc2 (GreenCo) and cc3 (EcoLtd), transfer cc1 to EcoLtd, retire cc2.

---

Authored by **Dayal Mukati** — [dayalmukati.com](https://dayalmukati.com)
