package main

import (
	"TemplateParser/TemplateParser"
	"fmt"
)

func Decode(ts string, r []TemplateParser.ObjectType, ok bool, errmsg string) {
	if ok {
		fmt.Printf("Successful parse of %s\n", ts)
		for idx, obj := range r {
			fmt.Printf("%d %v %v\n", idx, obj.ObjectTypeId, obj.ObjectValue)
		}
	} else {
		fmt.Printf("\nFailed parse of %s\n", ts)
		fmt.Println(errmsg)
	}
}

func main() {
	// Create the template of what we want
	fmt.Printf("Our objects expect an identifer")
	var templateList = []TemplateParser.TemplateObject{
		{
			TemplateType:  TemplateParser.TokenIdentifier,
			TemplateError: "Expected an identifier",
		},
		{
			TemplateType:  TemplateParser.TokenRegister,
			TemplateError: "Expected a destination register",
		},
		{
			TemplateType:  TemplateParser.TokenRegister,
			TemplateError: "Expected a source register",
		},
	}
	// Parse a correct line -- returns the objects filled im
	// A boolean if we were successful and a possible error message
	// Will succeed
	testText := "mov64 r10 r11"
	returnedObjs, ok, errmsg := TemplateParser.ParseLine(testText, templateList)
	Decode(testText, returnedObjs, ok, errmsg)
	// Will fail
	testText = "mov64 bob alice"
	ret, ok, errmsg := TemplateParser.ParseLine(testText, templateList)
	Decode(testText, ret, ok, errmsg)
	fmt.Println("Done")
}
