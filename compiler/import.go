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
	path := parseImport(p)
	if len(path) == 0 {
		p.errorf("Import path is empty")
	}

	// acquire absolute path and check the
	filename, err := absolutePath(string(path))
	if err != nil {
		p.errorf("Check absolute path error: %v", err)
	}

	inputFile, err := os.Open(filename)
	if err != nil {
		p.errorf("Open the import contract file \"%s\" error: %v", filename, err)
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
	importPath, newOffset := scanStrLiteral(p.buf, p.pos)
	if newOffset < 0 {
		p.errorf("Invalid import character format")
	}
	p.pos = newOffset

	return importPath
}

func absolutePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// check the status of absolute path
	if _, err := os.Stat(absPath); err != nil {
		return "", err
	}
	return absPath, nil
}
