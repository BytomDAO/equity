package equitytest

import (
	"bufio"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/equity/compiler"
)

func TestCompileContract(t *testing.T) {
	cases := []struct {
		pathFile string
		want     string
	}{
		{
			"./RepayCollateral",
			"557a641f0000007bcda069007b7b51547ac16951c3c251547ac1632a0000007bcd9f6900c3c251567ac1",
		},
		{
			"./LoanCollateral",
			"567a64650000007bcda06900c3537ac2547a5100597989587a89577a89557a89537a8901747e2a557a641f0000007bcda069007b7b51547ac16951c3c251547ac1632a0000007bcd9f6900c3c251567ac189008901c07ec16951c3c251547ac163700000007bcd9f6900c3c251577ac1",
		},
	}

	for _, c := range cases {
		contractName := filepath.Base(c.pathFile)
		t.Run(contractName, func(t *testing.T) {
			absPathFile, err := filepath.Abs(c.pathFile)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := os.Stat(absPathFile); err != nil {
				t.Fatal(err)
			}

			inputFile, err := os.Open(absPathFile)
			if err != nil {
				t.Fatal(err)
			}
			defer inputFile.Close()

			inputReader := bufio.NewReader(inputFile)
			contracts, err := compiler.Compile(inputReader)
			if err != nil {
				t.Fatal(err)
			}

			contract := contracts[len(contracts)-1]
			got := hex.EncodeToString(contract.Body)
			if got != c.want {
				t.Errorf("%s got %s\nwant %s", contractName, got, c.want)
			}
		})
	}
}
