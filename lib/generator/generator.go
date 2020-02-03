package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/compiler"
)

func Generate(file string, abi bool, bin bool, bindgo bool, dest string) error {
	pkgNameParts := strings.Split(file, string(os.PathSeparator))
	pkg := strings.Split(pkgNameParts[len(pkgNameParts)-1], ".")[0]

	if dest == "" {
		dest = "generated"
	}
	c, err := compiler.CompileSolidity("solc", file)
	if err != nil {
		fmt.Println("Unable to compile : ", err.Error())
		return err
	}
	for i, v := range c {
		b, err := json.Marshal(v.Info.AbiDefinition) // Info.AbiDefinition is abi
		if err != nil {
			fmt.Println("Unable to marshal : ", err.Error())
			return err
		}
		nameParts := strings.Split(i, ":")
		typeName := nameParts[len(nameParts)-1]

		err = os.MkdirAll(dest+string(os.PathSeparator)+pkg, os.ModePerm)
		if err != nil {
			fmt.Println("Unable to create directories : ", err.Error())
			return err
		}

		// generate abi
		abiBytes, _ := json.Marshal(v.Info.AbiDefinition)
		if abi == true {
			if err := ioutil.WriteFile(dest+string(os.PathSeparator)+pkg+string(os.PathSeparator)+typeName+".abi", abiBytes, 0600); err != nil {
				fmt.Println("Failed to write ABI : ", err.Error())
				return err
			}
		}

		// generate bin
		binArray := strings.Split(v.Code, "x")
		if bin == true {
			if err := ioutil.WriteFile(dest+string(os.PathSeparator)+pkg+string(os.PathSeparator)+typeName+".bin", []byte(binArray[1]), 0600); err != nil {
				fmt.Println("Failed to write bin : ", err.Error())
				return err
			}
		}
		if bindgo == true {
			var (
				abis    []string
				bins    []string
				types   []string
				sigs    []map[string]string
				libs    = make(map[string]string)
				aliases = make(map[string]string)
			)
			abis = append(abis, string(b))
			bins = append(bins, v.Code)
			sigs = append(sigs, v.Hashes)
			types = append(types, typeName)
			lang := bind.LangGo
			code, err := bind.Bind(types, abis, bins, sigs, strings.ToLower(pkg), lang, libs, aliases)
			if err != nil {
				fmt.Println("Unable to bind : ", err.Error())
				return err
			}

			if err := ioutil.WriteFile(dest+string(os.PathSeparator)+pkg+string(os.PathSeparator)+strings.ToLower(nameParts[len(nameParts)-1])+".go", []byte(code), 0600); err != nil {
				fmt.Println("Failed to write ABI binding : ", err.Error())
				return err
			}
		}
	}
	return nil
}
