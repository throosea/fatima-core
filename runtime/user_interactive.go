//
// Copyright (c) 2017 SK TECHX.
// All right reserved.
//
// This software is the confidential and proprietary information of SK TECHX.
// You shall not disclose such Confidential Information and
// shall use it only in accordance with the terms of the license agreement
// you entered into with SK TECHX.
//
//
// @project fatima
// @author 1100282
// @date 2017. 8. 28. AM 10:25
//

package runtime

import (
	"throosea.com/fatima"
	"path/filepath"
	"os"
	"fmt"
	"encoding/xml"
)

type UserInteractive struct {

}

func newUserInteractive() *UserInteractive {
	return &UserInteractive{}
}

func (ui *UserInteractive) Initialize() bool {
	inputFile := filepath.Join(
		process.GetEnv().GetFolderGuide().GetAppFolder(),
		process.GetEnv().GetSystemProc().GetProgramName(),
		process.GetEnv().GetSystemProc().GetProgramName() + ".ui.xml")

	xmlFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("fail to load user interactive xml file : %s", err.Error())
		return false
	}

	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	var inElement string
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...and its name is "page"
			if inElement == "prompt" {
				fmt.Printf("found inElement - prompt\n")
			}
		default:
		}

	}

	return true
}

func (ui *UserInteractive) Bootup() {
	go func() {
		startUserInteractive()
	}()
}

func (ui *UserInteractive) Shutdown() {
}

func (ui *UserInteractive) GetType() fatima.FatimaComponentType {
	return fatima.COMP_GENERAL
}

func startUserInteractive() {
	fmt.Printf("startUserInteractive....\n")
}

