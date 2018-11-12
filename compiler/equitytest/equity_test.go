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
			"./LockWithHeight",
			"537a64170000007bcda06900c3c251557ac163220000007bcd9f6900c3c251547ac1",
		},
		{
			"./ImportWithHeight",
			"557a6418000000537acda06900c3c251547ac16358000000537acd9f6900c3c25100587a89577a89567a8901747e22537a64170000007bcda06900c3c251557ac163220000007bcd9f6900c3c251547ac189008901c07ec1",
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
