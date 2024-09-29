package binding

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/unix-world/wails-app/internal/fs"

	"github.com/leaanthony/slicer"
)

var (
	mapRegex          *regexp.Regexp
	keyPackageIndex   int
	keyTypeIndex      int
	valueArrayIndex   int
	valuePackageIndex int
	valueTypeIndex    int
)

func init() {
	mapRegex = regexp.MustCompile(`(?:map\[(?:(?P<keyPackage>\w+)\.)?(?P<keyType>\w+)])?(?P<valueArray>\[])?(?:\*?(?P<valuePackage>\w+)\.)?(?P<valueType>.+)`)
	keyPackageIndex = mapRegex.SubexpIndex("keyPackage")
	keyTypeIndex = mapRegex.SubexpIndex("keyType")
	valueArrayIndex = mapRegex.SubexpIndex("valueArray")
	valuePackageIndex = mapRegex.SubexpIndex("valuePackage")
	valueTypeIndex = mapRegex.SubexpIndex("valueType")
}

func (b *Bindings) GenerateGoBindings(baseDir string) error {
	store := b.db.store
	var obfuscatedBindings map[string]int
	if b.obfuscate {
		obfuscatedBindings = b.db.UpdateObfuscatedCallMap()
	}
	for packageName, structs := range store {
		packageDir := filepath.Join(baseDir, packageName)
		err := fs.Mkdir(packageDir)
		if err != nil {
			return err
		}
		for structName, methods := range structs {
			var jsoutput bytes.Buffer
			jsoutput.WriteString(`// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
`)
			var tsBody bytes.Buffer
			var tsContent bytes.Buffer
			tsContent.WriteString(`// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
`)
			// Sort the method names alphabetically
			methodNames := make([]string, 0, len(methods))
			for methodName := range methods {
				methodNames = append(methodNames, methodName)
			}
			sort.Strings(methodNames)

			var importNamespaces slicer.StringSlicer
			for _, methodName := range methodNames {
				// Get the method details
				methodDetails := methods[methodName]

				// Generate JS
				var args slicer.StringSlicer
				for count := range methodDetails.Inputs {
					arg := fmt.Sprintf("arg%d", count+1)
					args.Add(arg)
				}
				argsString := args.Join(", ")
				jsoutput.WriteString(fmt.Sprintf("\nexport function %s(%s) {", methodName, argsString))
				jsoutput.WriteString("\n")
				if b.obfuscate {
					id := obfuscatedBindings[strings.Join([]string{packageName, structName, methodName}, ".")]
					jsoutput.WriteString(fmt.Sprintf("  return ObfuscatedCall(%d, [%s]);", id, argsString))
				} else {
					jsoutput.WriteString(fmt.Sprintf("  return window['go']['%s']['%s']['%s'](%s);", packageName, structName, methodName, argsString))
				}
				jsoutput.WriteString("\n}\n")

				// Generate TS
				tsBody.WriteString(fmt.Sprintf("\nexport function %s(", methodName))

				args.Clear()
				for count, input := range methodDetails.Inputs {
					arg := fmt.Sprintf("arg%d", count+1)
					entityName := entityFullReturnType(input.TypeName, b.tsPrefix, b.tsSuffix, &importNamespaces)
					args.Add(arg + ":" + goTypeToTypescriptType(entityName, &importNamespaces))
				}
				tsBody.WriteString(args.Join(",") + "):")
				// now build Typescript return types
				// If there is no return value or only returning error, TS returns Promise<void>
				// If returning single value, TS returns Promise<type>
				// If returning single value or error, TS returns Promise<type>
				// If returning two values, TS returns Promise<type1|type2>
				// Otherwise, TS returns Promise<type1> (instead of throwing Go error?)
				var returnType string
				if methodDetails.OutputCount() == 0 {
					returnType = "Promise<void>"
				} else if methodDetails.OutputCount() == 1 && methodDetails.Outputs[0].TypeName == "error" {
					returnType = "Promise<void>"
				} else {
					outputTypeName := entityFullReturnType(methodDetails.Outputs[0].TypeName, b.tsPrefix, b.tsSuffix, &importNamespaces)
					firstType := goTypeToTypescriptType(outputTypeName, &importNamespaces)
					returnType = "Promise<" + firstType
					if methodDetails.OutputCount() == 2 && methodDetails.Outputs[1].TypeName != "error" {
						outputTypeName = entityFullReturnType(methodDetails.Outputs[1].TypeName, b.tsPrefix, b.tsSuffix, &importNamespaces)
						secondType := goTypeToTypescriptType(outputTypeName, &importNamespaces)
						returnType += "|" + secondType
					}
					returnType += ">"
				}
				tsBody.WriteString(returnType + ";\n")
			}

			importNamespaces.Deduplicate()
			importNamespaces.Each(func(namespace string) {
				tsContent.WriteString("import {" + namespace + "} from '../models';\n")
			})
			tsContent.WriteString(tsBody.String())

			jsfilename := filepath.Join(packageDir, structName+".js")
			err = os.WriteFile(jsfilename, jsoutput.Bytes(), 0o755)
			if err != nil {
				return err
			}
			tsfilename := filepath.Join(packageDir, structName+".d.ts")
			err = os.WriteFile(tsfilename, tsContent.Bytes(), 0o755)
			if err != nil {
				return err
			}
		}
	}
	err := b.WriteModels(baseDir)
	if err != nil {
		return err
	}
	return nil
}

func fullyQualifiedName(packageName string, typeName string) string {
	if len(packageName) > 0 {
		return packageName + "." + typeName
	}

	switch true {
	case len(typeName) == 0:
		return ""
	case typeName == "interface{}" || typeName == "interface {}":
		return "any"
	case typeName == "string":
		return "string"
	case typeName == "error":
		return "Error"
	case
		strings.HasPrefix(typeName, "int"),
		strings.HasPrefix(typeName, "uint"),
		strings.HasPrefix(typeName, "float"):
		return "number"
	case typeName == "bool":
		return "boolean"
	default:
		return "any"
	}
}

func arrayifyValue(valueArray string, valueType string) string {
	if len(valueArray) == 0 {
		return valueType
	}

	return "Array<" + valueType + ">"
}

func goTypeToJSDocType(input string, importNamespaces *slicer.StringSlicer) string {
	matches := mapRegex.FindStringSubmatch(input)
	keyPackage := matches[keyPackageIndex]
	keyType := matches[keyTypeIndex]
	valueArray := matches[valueArrayIndex]
	valuePackage := matches[valuePackageIndex]
	valueType := matches[valueTypeIndex]
	// fmt.Printf("input=%s, keyPackage=%s, keyType=%s, valueArray=%s, valuePackage=%s, valueType=%s\n",
	//	input,
	//	keyPackage,
	//	keyType,
	//	valueArray,
	//	valuePackage,
	//	valueType)

	// byte array is special case
	if valueArray == "[]" && valueType == "byte" {
		return "string"
	}

	// if any packages, make sure they're saved
	if len(keyPackage) > 0 {
		importNamespaces.Add(keyPackage)
	}

	if len(valuePackage) > 0 {
		importNamespaces.Add(valuePackage)
	}

	key := fullyQualifiedName(keyPackage, keyType)
	var value string
	if strings.HasPrefix(valueType, "map") {
		value = goTypeToJSDocType(valueType, importNamespaces)
	} else {
		value = fullyQualifiedName(valuePackage, valueType)
	}

	if len(key) > 0 {
		return fmt.Sprintf("{[key: %s]: %s}", key, arrayifyValue(valueArray, value))
	}

	return arrayifyValue(valueArray, value)
}

func goTypeToTypescriptType(input string, importNamespaces *slicer.StringSlicer) string {
	return goTypeToJSDocType(input, importNamespaces)
}

func entityFullReturnType(input, prefix, suffix string, importNamespaces *slicer.StringSlicer) string {
	if strings.ContainsRune(input, '.') {
		nameSpace, returnType := getSplitReturn(input)
		return nameSpace + "." + prefix + returnType + suffix
	}

	return input
}
