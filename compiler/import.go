package compiler

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
)

func parseImportDirectives(p *parser) []*Contract {
	var result []*Contract
	for peekKeyword(p) == "import" {
		contracts := parseImportDirective(p)
		for _, contract := range contracts {
			result = append(result, contract)
		}
	}
	return result
}

func parseImportDirective(p *parser) []*Contract {
	pathFile := parseImport(p)
	if len(pathFile) == 0 {
		p.errorf("Import path is empty")
	}

	// acquire absolute path and check the file status
	importFile, err := absolutePath(string(pathFile))
	if err != nil {
		p.errorf("Check absolute path error: %v", err)
	}

	inputFile, err := os.Open(importFile)
	if err != nil {
		p.errorf("Open the import contract file \"%s\" error: %v", importFile, err)
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	importContract, err := ioutil.ReadAll(inputReader)
	if err != nil {
		p.errorf("Read the import contract file \"%s\" error: %v", inputFile.Name(), err)
	}

	// parse the import contract
	contracts, err := parse(importContract)
	if err != nil {
		p.errorf("Parse the import contract file \"%s\" error: %v", inputFile.Name(), err)
	}
	return contracts
}

func parseImport(p *parser) []byte {
	consumeKeyword(p, "import")
	importPathFile, newOffset := scanStrLiteral(p.buf, p.pos)
	if newOffset < 0 {
		p.errorf("Invalid import character format")
	}
	p.pos = newOffset

	return importPathFile
}

func absolutePath(pathFile string) (string, error) {
	absPathFile, err := filepath.Abs(pathFile)
	if err != nil {
		return "", err
	}

	// check the status of absolute path file
	if _, err := os.Stat(absPathFile); err != nil {
		return "", err
	}
	return absPathFile, nil
}
